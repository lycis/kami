package main

import (
	"gitlab.com/lycis/kami/util"
	log "github.com/Sirupsen/logrus"
	"gitlab.com/lycis/kami/driver/local"
	flag "github.com/ogier/pflag"
	"os"
	"strings"

	// import for side effect of registering functions
	_ "gitlab.com/lycis/kami/driver/dfun"
	_ "gitlab.com/lycis/kami/script/efun"
)

func main() {
	util.PrintLicenseHint("Kami Game Driver v0.0")

	libDir := flag.String("lib", "/usr/lib/kami/", "root directory of the game library")
	initScript := flag.String("init", "/sys/init.js", "name of the script that will initialise the library (run on startup)")
	logLevel := flag.String("log-level", "INFO", "log level of printed messages (DEBUG, INFO, WARNING, ERROR, FATAl, PANIC)")
	flag.Parse()

	l := log.New()
	l.Level = stringToLogLevel(*logLevel)
	l.Out = os.Stdout
	log.SetOutput(os.Stdout)
	log.SetLevel(stringToLogLevel(*logLevel))

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
