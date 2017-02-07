package script

import (
	"github.com/lycis/kami/entity"
	"github.com/Sirupsen/logrus"
)

type DriverAPI interface{
	SpawnExclusive(rpath string) (*entity.Entity, error)
	SpawnEntity(rpath string) (*entity.Entity, error)
	GetEntityById(id string) *entity.Entity
	Logger() *logrus.Logger
}
