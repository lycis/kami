package driver

import (
	"github.com/Sirupsen/logrus"
	"gitlab.com/lycis/kami/script"
)

// Driver represents the overall game driver state and driver base functions
// It takes care of loading and executing the game world
type Driver interface {
	script.DriverAPI

	// SetLogger indicates which logger is to be used by the driver
	SetLogger(*logrus.Logger)

	// Init initialises the driver by calling the given script file
	Init(file string)

	// RunWorld is supposed to do everything to execute the game world.
	// It is called once the driver was initialised successfully.
	RunWorld()
}
