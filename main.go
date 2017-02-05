package main

import (
	flag "github.com/ogier/pflag"
	log "github.com/Sirupsen/logrus"
	"github.com/lycis/kami/driver"
	"os"
)

func main() {
	libDir := flag.String("lib", "/usr/lib/kami/", "root directory of the game library")
	initScript := flag.String("init", "/sys/init.js", "name of the script that will initialise the library (run on startup)")
	flag.Parse()

	l := log.New()
	l.Out = os.Stdout

	mainDriver := driver.New(*libDir)
	mainDriver.SetLogger(l)
	mainDriver.Init(*initScript)
	mainDriver.RunWorld()
}