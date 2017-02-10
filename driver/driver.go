package driver

import (
	"github.com/Sirupsen/logrus"
)

// Driver represents the overall game driver state and driver base functions
// It takes care of loading and executing the game world
type Driver interface {
	SetLogger(*logrus.Logger)
	Init(file string)
	RunWorld()
	Logger() *logrus.Logger
}
