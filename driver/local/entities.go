package local

import (
	"github.com/lycis/kami/entity"
	log "github.com/Sirupsen/logrus"
	"github.com/nu7hatch/gouuid"
	"github.com/lycis/kami/script"
)

func (driver *LocalDriver) SpawnExclusive(rpath string) (*entity.Entity, error) {
	// TODO check if instances exist

	e, err := driver.createEntityInstance(rpath)
	if err != nil {
		return nil, err
	}

	e.SetProp(entity.P_SYS_ID, rpath)
	e.SetProp(entity.P_SYS_EXCLUSIVE, true)

	driver.registerEntity(e)
	return e, nil
}

// SpawnEntity loads and spawns an entity from the given script path
func (driver *LocalDriver) SpawnEntity(rpath string) (*entity.Entity, error) {
	// TODO check if exclusive entity exists

	e, err := driver.createEntityInstance(rpath)
	if err != nil {
		return nil, err
	}

	id, err := uuid.NewV4()
	if err !=nil {
		return nil, err
	}

	e.SetProp(entity.P_SYS_ID, id.String())
	e.SetProp(entity.P_SYS_EXCLUSIVE, false)

	driver.registerEntity(e)
	return e, nil
}

func (driver *LocalDriver) createEntityInstance(rpath string) (*entity.Entity, error) {
	ctx, err := script.ContextForScript(rpath, driver.LibraryDir(), &driver.scriptCache)
	if err != nil {
		return nil, err
	}

	e, err := entity.NewEntity(&ctx)
	if err != nil {
		return nil, err
	}

	e.SetProp(entity.P_SYS_PATH, rpath)

	return e, nil
}

func (driver *LocalDriver) registerEntity(e *entity.Entity) {
	driver.entityListMutex.Lock()
	defer driver.entityListMutex.Unlock()
	id := e.GetProp(entity.P_SYS_ID).(string)
	path := e.GetProp(entity.P_SYS_PATH).(string)

	driver.activeEntities[id] = e

	if driver.entityInstances[path] == nil {
		driver.entityInstances[path] = make([]*entity.Entity, 0)
	}
	driver.entityInstances[path] = append(driver.entityInstances[path], e)

	if e.GetProp(entity.P_SYS_EXCLUSIVE).(bool) {
		driver.Log.WithFields(log.Fields{"path": id}).Info("Exclusive entity spawned.")
	} else {
		driver.Log.WithFields(log.Fields{"uuid": id, "path": path}).Info("Entity instance spawned.")
	}
}
