package local

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/lycis/kami/driver"
	"github.com/lycis/kami/driver/dfun"
	"github.com/lycis/kami/entity"
	"github.com/lycis/kami/script"
	"github.com/lycis/kami/subsystem"
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
						d.Log.WithField("entity", fmt.Sprintf("%s#%s", xerr.Entity.GetProp(entity.P_SYS_PATH), xerr.Entity.GetProp(entity.P_SYS_PATH))).WithError(xerr.Error).Error("Calling heratbeat error hook failed")
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
	hv, ok := d.hooks[dfun.H_NEW_USER]
	if !ok {
		return "", fmt.Errorf("NEW_USER driver hook missing")
	}

	if !hv.IsFunction() {
		return "", fmt.Errorf("NEW_USER driver hook has invalid type (expected: function)")
	}

	token, err := hv.Call(otto.ToValue(d))
	if err != nil {
		return "", err
	}

	strToken, err := token.ToString()
	if err != nil {
		return "", fmt.Errorf("NEW_USER driver hook did not return token string")
	}

	return strToken, nil
}
