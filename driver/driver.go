package driver

import (
	log "github.com/Sirupsen/logrus"
	"github.com/lycis/kami/script"
	"time"
	"github.com/lycis/kami/entity"
)

// Driver represents the overall game driver state and driver base functions
// It takes care of loading and executing the game world
type Driver struct {
	libraryDir string
	Log        *log.Logger
	scriptCache script.ScriptCache
	cacheCleanupTimer *time.Timer
	activeEntities map[string]*entity.Entity
	entityInstances map[string]*entity.Entity
}

// New generates a new driver instance that can be executed.
// The libDir tells the driver where to search for scripts and programs
// to be executed.
func New(libDir string) Driver {
	return Driver{
		libraryDir: libDir,
		Log: log.New(),
		scriptCache: script.NewCache(),
		activeEntities: make(map[string]*entity.Entity),
		entityInstances: make(map[string]*entity.Entity),
	}
}

// SetLogger gives the driver a logger that all output will
// be written to. If this is not set it will by default use
// the default logger (usually stdout)
func (d *Driver) SetLogger(l *log.Logger) {
	if l == nil {
		return
	}

	d.Log = l
}

func (d Driver) LibraryDir() string {
	return d.libraryDir
}

func (driver *Driver) RunWorld() {

}


func (d Driver) Logger() *log.Logger {
	return d.Log
}