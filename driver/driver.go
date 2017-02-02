package driver

import (
	log "github.com/Sirupsen/logrus"
	"fmt"
	"os"
	"github.com/lycis/kami/script"
	"time"
	"github.com/lycis/kami/entity"
	"github.com/nu7hatch/gouuid"
	"github.com/lycis/kami/driver/dfun"
)

// Driver represents the overall game driver state and driver base functions
// It takes care of loading and executing the game world
type Driver struct {
	libraryDir string
	Log        *log.Logger
	scriptCache script.ScriptCache
	cacheCleanupTimer *time.Timer
	activeEntities map[string]*entity.Entity
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

// Init will initialise and start the driver and also game world.
func (d *Driver) Init(file string) {
	d.Log.Info("Starting game driver.")

	d.Log.Info("Registering dfuns.")
	d.spawnTimers()
	d.callInitScript(file)

	d.Log.Info("Driver initialised.")
}

func (driver *Driver) spawnTimers() {
	driver.cacheCleanupTimer = time.AfterFunc(time.Minute*5, driver.cleanupCaches)
	driver.Log.WithField("interval", "5m").Info("Spawned cache cleanup timer.")
}

func (driver *Driver) cleanupCaches() {
	driver.scriptCache.Cleanup(time.Minute*5)
	driver.cacheCleanupTimer = time.AfterFunc(time.Minute*5, driver.cleanupCaches)
}

func (d *Driver) callInitScript(file string) {
	d.Log.WithField("init", fmt.Sprintf("%s%s", d.libraryDir, file)).Info("Loading init script.")
	if _, err := os.Stat(d.libraryDir); os.IsNotExist(err) {
		d.Log.WithField("lib", d.libraryDir).Fatal("Game library directory does not exist")
		return
	}

	ctx := script.NewContext(d.libraryDir, &d.scriptCache)
	ctx.Bind("_driver", dfun.NewProvider(d))

	if err := ctx.RunScript(file); err != nil {
		log.WithError(err).Fatal("Executing the init script failed.")
		return
	}
}

func (driver *Driver) RunWorld() {

}

// SpawnEntity loads and spawns an entity from the given script path
func (driver *Driver) SpawnEntity(rpath string) (*entity.Entity, error) {
	ctx := script.NewContext(driver.libraryDir, &driver.scriptCache)
	if err := ctx.RunScript(rpath); err != nil {
		return nil, err
	}

	instance := ctx.GetInstance()

	e, err := entity.NewEntity(instance)
	if err != nil {
		return nil, err
	}

	id, err := uuid.NewV4()
	if err !=nil {
		return nil, err
	}

	e.SetProp("$uuid", id.String())
	driver.activeEntities[id.String()] = e
	driver.Log.WithFields(log.Fields{"$uuid": id, "path": rpath}).Info("Entity spawned.")
	return e, nil
}

func (d Driver) Logger() *log.Logger {
	return d.Log
}