package efun

import (
	"github.com/lycis/kami/privilege"
	"github.com/lycis/kami/script"
	"github.com/robertkrimen/otto"
	"strings"
)

// The efun "log" allows you to write log messages into the driver log.
// syntax:
//   log(level: string, message: string)
//
//   levels: DEBUG, INFO, WARN, FATAL
type efunLog struct {
	script *script.ScriptContext
}

func create_log(i *script.ScriptContext) script.ExposedFunction {
	return &efunLog{
		script: i,
	}
}

func (e *efunLog) Function() func(call otto.FunctionCall) otto.Value {
	return func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) != 2 {
			panic(e.script.Vm().MakeSyntaxError("log(level: string, message: string): expects at least two arguments"))
		}

		if !call.Argument(0).IsString() {
			panic(e.script.Vm().MakeSyntaxError("log(level: string, message: string): level has to be string"))
		}
		level := call.Argument(0).String()

		if !call.Argument(1).IsString() {
			panic(e.script.Vm().MakeSyntaxError("log(level: string, message: string): message has to be string"))
		}
		message := call.Argument(1).String()

		l := e.script.Driver().Logger()
		var fun func(args ...interface{})
		switch strings.ToUpper(level) {
		case "DEBUG":
			fun = l.Debug
		case "INFO":
			fun = l.Info
		case "WARN":
			fun = l.Warn
		case "FATAL":
			fun = l.Fatal
		default:
			panic(e.script.Vm().MakeSyntaxError("log(level: string, message: string): invalid level"))
		}

		fun(message)
		return otto.UndefinedValue()
	}
}

func (e *efunLog) RequiredPrivilegeLevel() privilege.Level {
	return privilege.PrivilegeBasic
}
