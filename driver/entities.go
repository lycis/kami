package driver

import (
	"github.com/lycis/kami/entity"
	log "github.com/Sirupsen/logrus"
	"github.com/nu7hatch/gouuid"
	"github.com/lycis/kami/script"
)

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
