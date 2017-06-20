package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/lycis/kami/driver/local"
	flag "github.com/ogier/pflag"
	"os"
	"strings"
)

func main() {
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
