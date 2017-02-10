package script

import (
	"github.com/Sirupsen/logrus"
	"github.com/lycis/kami/entity"
)

type DriverAPI interface {
	SpawnExclusive(rpath string, creator ContextCreator) (*entity.Entity, error)
	SpawnEntity(rpath string, creator ContextCreator) (*entity.Entity, error)
	GetEntityById(id string) *entity.Entity
	Logger() *logrus.Logger
}
