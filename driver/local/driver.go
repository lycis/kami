package local

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"gitlab.com/lycis/kami/driver"
	"gitlab.com/lycis/kami/driver/dfun"
	"gitlab.com/lycis/kami/entity"
	"gitlab.com/lycis/kami/script"
	"gitlab.com/lycis/kami/subsystem"
	"github.com/robertkrimen/otto"
	"reflect"
	"runtime"
	"sync"
	"time"
)

// Driver represents the overall game driver state and driver base functions
// It takes care of loading and executing the game world
type LocalDriver struct {
	libraryDir        string
	Log               *log.Logger
	scriptCache       script.ScriptCache
	cacheCleanupTimer *time.Timer

	activeEntities  map[string]*entity.Entity
	entityInstances map[string][]*entity.Entity
	entityListMutex sync.Mutex

	lastHeartbeat time.Time

	hooks map[int64]otto.Value

	// indicates whether world is running
	running bool

	restInterface subsystem.NetworkingInterface
}

// New generates a new driver instance that can be executed.
// The libDir tells the driver where to search for scripts and programs
// to be executed.
func New(libDir string) driver.Driver {
	return &LocalDriver{
		libraryDir:      libDir,
		Log:             log.New(),
		scriptCache:     script.NewCache(),
		activeEntities:  make(map[string]*entity.Entity),
		entityInstances: make(map[string][]*entity.Entity),
		hooks:           make(map[int64]otto.Value),
		restInterface:   nil,
	}
}

// SetLogger gives the driver a logger that all output will
// be written to. If this is not set it will by default use
// the default logger (usually stdout)
func (d *LocalDriver) SetLogger(l *log.Logger) {
	if l == nil {
		return
	}

	d.Log = l
}

func (d LocalDriver) LibraryDir() string {
	return d.libraryDir
}

func (d LocalDriver) Logger() *log.Logger {
	return d.Log
}

func (d LocalDriver) GetEntityById(id string) *entity.Entity {
	return d.activeEntities[id]
}

func (d *LocalDriver) SetHook(hook int64, value interface{}) error {
	ov, ok := value.(otto.Value)
	if !ok {
		return fmt.Errorf("local driver only supports javascript function calls")
	}

	if !ov.IsFunction() {
		return fmt.Errorf("local driver only supports javascript function calls")
	}

	d.hooks[hook] = ov
	return nil
}

func (d *LocalDriver) Shutdown(reason string) error {
	if !d.running {
		return fmt.Errorf("world not in running state")
	}

	d.Log.WithField("reason", reason).Info("* * * Driver is going to H A L T * * *")
	d.Log.Debug("Informing entities of shutdown.")
	var wg sync.WaitGroup
	for _, instances := range d.entityInstances {
		wg.Add(len(instances))
	}

	d.forAllInstances(func(shard []*entity.Entity) {
		for _, e := range shard {
			defer func() {
				if err := recover(); err != nil {
					xerr, ok := err.(entity.FunctionInvocationError)
					if ok {
						d.Log.WithField("entity", fmt.Sprintf("%s#%s", xerr.Entity.GetProp(entity.P_SYS_PATH), xerr.Entity.GetProp(entity.P_SYS_PATH))).WithError(xerr.Err).Error("Calling heratbeat error hook failed")
					} else {
						d.Log.WithError(err.(error)).Warn("Error in shutdown processing of unknown entity.")
					}
				}
			}()

			d.Log.WithField("entity", fmt.Sprintf("%s#%s", e.GetProp(entity.P_SYS_PATH), e.GetProp(entity.P_SYS_PATH))).Debug("executing onShutdown for entity")
			e.OnShutdown(reason)
			wg.Done()
		}
	})
	wg.Wait()
	d.Log.Debug("Entities informed.")

	d.running = false
	d.Log.Debug("World stopped.")

	d.Log.Info("Driver stopped.")
	return nil
}

func (d *LocalDriver) forAllInstances(f func([]*entity.Entity)) {
	for path, instances := range d.entityInstances {
		d.Log.WithFields(log.Fields{"path": path, "function": getFunctionName(f)}).Debug("Calling function for instance shard.")
		go func() {
			f(instances)
		}()
	}
}

func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func (d *LocalDriver) SetSubsystemState(stype int64, status bool) error {
	d.Log.WithFields(log.Fields{"subsystem": stype, "status": status}).Debug("Changing subsystem state.")
	switch stype {
	case dfun.D_SUBSYSTEM_REST:
		if status {
			return d.enable_rest()
		} else {
			return d.disable_rest()
		}
	default:
		return fmt.Errorf("driver does not support subsystem")
	}
}

func (d LocalDriver) isRestEnabled() bool {
	return d.restInterface != nil
}

func (d *LocalDriver) enable_rest() error {
	if d.isRestEnabled() {
		return fmt.Errorf("REST subsystem already enabled")
	}

	// TODO make configurable
	d.restInterface = subsystem.CreateNetworkingInterface(subsystem.NWI_REST, "0.0.0.0", 8080)
	d.restInterface.SetHandler(d)
	return d.restInterface.Listen()
}

func (d *LocalDriver) disable_rest() error {
	if !d.isRestEnabled() {
		return nil
	}

	d.restInterface.Close()
	d.restInterface = nil
	return nil
}

// Provide a user token for new logins
func (d *LocalDriver) UserTokenRequested(nwi subsystem.NetworkingInterface) (string, error) {
	d.Log.Info("New user token requested.")
	hv, err := d.getHookFuncCallable(dfun.H_NEW_USER)
	if err != nil {
		log.WithFields(log.Fields{"hook": "H_NEW_USER"}).WithError(err).Error("User token could not be provided. Driver Hook error.")
		return "", err
	}

	token, err := hv.Call(otto.UndefinedValue())
	if err != nil {
		log.WithError(err).Error("User token not provided. Error when calling javascript function for H_NEW_USER.")
		return "", err
	}

	strToken, err := token.ToString()
	if err != nil {
		log.Error("User token not provided. H_NEW_USER hook did not return token string")
		return "", fmt.Errorf("NEW_USER driver hook did not return token string")
	}

	log.WithFields(log.Fields{"token": strToken}).Info("User token provided.")
	return strToken, nil
}

func (d *LocalDriver) UserInputProvided(nwi subsystem.NetworkingInterface, token, input string) error {
	d.Log.WithFields(log.Fields{"token": token}).Info("User input provided.")
	hv, err := d.getHookFuncCallable(dfun.H_USER_INPUT)
	if err != nil {
		log.WithError(err).Error("Processing user input failed.")
		return err
	}

	successV, err := hv.Call(otto.UndefinedValue(), token, input)
	if err != nil {
		return err
	}

	success, err := successV.ToBoolean()
	if err != nil {
		return err
	}

	if !success {
		return fmt.Errorf("H_USER_INPUT failed. See driver log for details.")
	}

	return nil
}

func (d LocalDriver) getHookFuncCallable(hook int64) (otto.Value, error) {
	hv, ok := d.hooks[hook]
	if !ok {
		return otto.UndefinedValue(), fmt.Errorf("driver hook %d not set", hook)
	}

	if !hv.IsFunction() {
		return otto.UndefinedValue(), fmt.Errorf("driver hook %d has invalid type (expected: function)", hook)
	}

	return hv, nil
}
