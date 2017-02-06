package driver

import (
	"github.com/Sirupsen/logrus"
	"github.com/lycis/kami/entity"
)

// Driver represents the overall game driver state and driver base functions
// It takes care of loading and executing the game world
type Driver interface {
	SetLogger(*logrus.Logger)
	Init(file string)
	RunWorld()
	SpawnExclusive(rpath string) (*entity.Entity, error)
	SpawnEntity(rpath string) (*entity.Entity, error)
	Logger() *logrus.Logger
}