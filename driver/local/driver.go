package local

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/lycis/kami/driver"
	"github.com/lycis/kami/entity"
	"github.com/lycis/kami/script"
	"github.com/robertkrimen/otto"
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
