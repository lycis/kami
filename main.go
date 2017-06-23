package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/lycis/kami/driver/local"
	flag "github.com/ogier/pflag"
	"os"
	"strings"

	// import for side effect of registering functions
	_ "github.com/lycis/kami/driver/dfun"
	_ "github.com/lycis/kami/script/efun"
)

func main() {
	printLicenseHint()

	libDir := flag.String("lib", "/usr/lib/kami/", "root directory of the game library")
	initScript := flag.String("init", "/sys/init.js", "name of the script that will initialise the library (run on startup)")
	logLevel := flag.String("log-level", "INFO", "log level of printed messages (DEBUG, INFO, WARNING, ERROR, FATAl, PANIC)")
	flag.Parse()

	l := log.New()
	l.Level = stringToLogLevel(*logLevel)
	l.Out = os.Stdout

	mainDriver := local.New(*libDir)
	mainDriver.SetLogger(l)
	mainDriver.Init(*initScript)
	mainDriver.RunWorld()
}
func printLicenseHint() {
	fmt.Println(`Kami Game Driver v0.0

	Copyright (C) 2017  Ing. Daniel Eder (daniel@deder.at)

	This program is free software: you can redistribute it and/or modify
	it under the terms of the GNU General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	This program is distributed in the hope that it will be useful,
		but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU General Public License for more details.

	You should have received a copy of the GNU General Public License
	along with this program.  If not, see <http://www.gnu.org/licenses/>.)
	`)
}

func stringToLogLevel(s string) log.Level {
	switch strings.ToUpper(s) {
	case "DEBUG":
		return log.DebugLevel
	case "INFO":
		return log.InfoLevel
	case "WARNING":
		return log.WarnLevel
	case "ERROR":
		return log.ErrorLevel
	case "FATAL":
		return log.FatalLevel
	case "PANIC":
		return log.PanicLevel
	default:
		panic("invalid log level")
	}
}
