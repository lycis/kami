package driver

import (
	log "github.com/Sirupsen/logrus"
	"fmt"
	"os"
	"github.com/lycis/kami/script"
)

// Driver represents the overall game driver state and driver base functions
// It takes care of loading and executing the game world
type Driver struct {
	libraryDir string
	Log        *log.Logger
}

// New generates a new driver instance that can be executed.
// The libDir tells the driver where to search for scripts and programs
// to be executed.
func New(libDir string) Driver {
	return Driver{
		libraryDir: libDir,
		Log: log.New(),
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

	d.Log.WithField("init", fmt.Sprintf("%s%s", d.libraryDir, file)).Info("Loading init script.")
	if _, err := os.Stat(d.libraryDir); os.IsNotExist(err) {
		d.Log.WithField("lib", d.libraryDir).Fatal("Game library directory does not exist")
		return
	}

	ctx := script.NewContext(d.libraryDir)
	ctx.Bind("_driver", d)
	if err := ctx.RunScript(file); err != nil {
		log.WithError(err).Fatal("Executing the init script failed.")
		return
	}

	d.Log.Info("Driver initialised.")
}