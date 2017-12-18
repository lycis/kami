package local

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"gitlab.com/lycis/kami/privilege"
	"gitlab.com/lycis/kami/script"
	"os"
	"time"
)

// Init will initialise and start the driver and also game world.
func (d *Driver) Init(file string) {
	d.Log.Info("Starting game driver.")

	d.spawnTimers()
	d.callInitScript(file)

	d.Log.Info("Driver initialised.")
}

func (d *Driver) spawnTimers() {
	d.cacheCleanupTimer = time.AfterFunc(time.Minute*5, d.cleanupCaches)
	d.Log.WithField("interval", "5m").Info("Spawned cache cleanup timer.")
}

func (d *Driver) cleanupCaches() {
	d.scriptCache.Cleanup(time.Minute * 5)
	d.cacheCleanupTimer = time.AfterFunc(time.Minute*5, d.cleanupCaches)
}

func (d *Driver) callInitScript(file string) {
	d.Log.WithField("init", fmt.Sprintf("%s%s", d.libraryDir, file)).Info("Loading init script.")
	if _, err := os.Stat(d.libraryDir); os.IsNotExist(err) {
		d.Log.WithField("lib", d.libraryDir).Fatal("Game library directory does not exist")
		return
	}

	ctx := script.NewContext(d, d.libraryDir, &d.scriptCache)
	ctx.SetCreator(d)
	ctx.GrantPrivilege(privilege.PrivilegeRoot)

	if err := ctx.RunScript(file); err != nil {
		log.WithError(err).Fatal("Executing the init script failed.")
		return
	}
}
