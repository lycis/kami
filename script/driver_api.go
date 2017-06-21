package script

import (
	"github.com/Sirupsen/logrus"
	"github.com/lycis/kami/entity"
)

type DriverAPI interface {

	// SpawnExcluse creates a new exclusive(!) entity from the given file.
	// it works like SpawnEntity but ensures that the entity exists only once
	// or will return an error
	SpawnExclusive(rpath string, creator ContextCreator) (*entity.Entity, error)

	// SpawnEntity creates a new entity from the given script file
	SpawnEntity(rpath string, creator ContextCreator) (*entity.Entity, error)

	// GetEntityById will return the entity object matching to the given ID or
	// nil in case no entity with this ID was found
	GetEntityById(id string) *entity.Entity

	// Logger() is used to provide a logger for the driver.
	Logger() *logrus.Logger

	// SetHook sets a driver hook to the given value
	SetHook(hook int64, value interface{}) error

	// Shutdown is called when something wants the driver to stop.
	Shutdown(reason string) error
}
