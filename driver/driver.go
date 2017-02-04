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

func (driver *Driver) SpawnExcluseive(rpath string) (*entity.Entity, error) {
	// TODO check if instances exist

	e, err := driver.createEntityInstance(rpath)
	if err != nil {
		return nil, err
	}

	e.SetProp("$id", rpath)
	e.SetProp("$exclusive", true)

	driver.registerEntity(e)
	return e, nil
}

// SpawnEntity loads and spawns an entity from the given script path
func (driver *Driver) SpawnEntity(rpath string) (*entity.Entity, error) {
	// TODO check if exclusive entity exists

	e, err := driver.createEntityInstance(rpath)
	if err != nil {
		return nil, err
	}

	id, err := uuid.NewV4()
	if err !=nil {
		return nil, err
	}

	e.SetProp("$uuid", id.String())
	e.SetProp("$unique", false)

	driver.registerEntity(e)
	return e, nil
}

func (driver *Driver) createEntityInstance(rpath string) (*entity.Entity, error) {
	ctx, err := script.ContextForScript(rpath, driver.LibraryDir(), &driver.scriptCache)
	if err != nil {
		return nil, err
	}

	instance := ctx.GetInstance()
	e, err := entity.NewEntity(instance)
	if err != nil {
		return nil, err
	}

	e.SetProp("$path", rpath)

	return e, nil
}

func (driver *Driver) registerEntity(e *entity.Entity) {
	id := e.GetProp("$uuid").(string)
	path := e.GetProp("$path").(string)

	driver.activeEntities[id] = e
	driver.entityInstances[path] = e
	if e.GetProp("$unique").(bool) {
		driver.Log.WithFields(log.Fields{"path": id}).Info("Exclusive entity spawned.")
	} else {
		driver.Log.WithFields(log.Fields{"$uuid": id, "path": path}).Info("Entity instance spawned.")
	}
}

func (d Driver) Logger() *log.Logger {
	return d.Log
}