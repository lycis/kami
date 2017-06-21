package local

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/lycis/kami/driver"
	"github.com/lycis/kami/entity"
	"github.com/lycis/kami/script"
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
	d.running = false
	d.Log.Debug("World stopped.")
	d.Log.Debug("Informing entities of shutdown.")
	d.forAllInstances(func(shard []*entity.Entity) {
		for _, e := range shard {
			defer func() {
				if err := recover(); err != nil {
					xerr, ok := err.(entity.FunctionInvocationError)
					if ok {
						log.WithField("entity", fmt.Sprintf("%s#%s", xerr.Entity.GetProp(entity.P_SYS_PATH), xerr.Entity.GetProp(entity.P_SYS_PATH))).WithError(xerr.Error).Error("Calling heratbeat error hook failed")
					} else {
						d.Log.WithError(err.(error)).Warn("Error in shutdown processing of unknown entity.")
					}
				}
			}()

			d.Log.WithField("entity", fmt.Sprintf("%s#%s", e.GetProp(entity.P_SYS_PATH), e.GetProp(entity.P_SYS_PATH))).Debug("executing onShutdown for entity")
			e.OnShutdown(reason)
		}
	})
	d.Log.Debug("Entities informed.")

	d.Log.Info("Driver stopped.")
	return nil
}

func (d *LocalDriver) forAllInstances(f func([]*entity.Entity)) {
	for path, instances := range d.entityInstances {
		log.WithFields(log.Fields{"path": path, "function": getFunctionName(f)}).Debug("Calling function for instance shard.")
		go func() {
			f(instances)
		}()
	}
}

func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
