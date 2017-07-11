package dfun

import (
	"github.com/lycis/kami/privilege"
	"github.com/lycis/kami/script"
	"github.com/robertkrimen/otto"
)

// syntax:
//   enable_subsystem(type: int[, options])
//
// enable_subsystem is used to turn on a driver subsystem. Optionally
// you can pass options to it, if they are supported.

type dfun_enable_subsystem struct {
	script *script.ScriptContext
}

func create_enable_subsystem(i *script.ScriptContext) script.ExposedFunction {
	return &dfun_enable_subsystem{
		script: i,
	}
}

func (dfun_enable_subsystem) RequiredPrivilegeLevel() privilege.Level {
	return privilege.PrivilegeRoot
}

const (
	D_SUBSYSTEM_REST = 0
)

func (df dfun_enable_subsystem) Function() func(call otto.FunctionCall) otto.Value {
	return func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) < 1 {
			panic(df.script.Vm().MakeSyntaxError("enable_subsystem(type: int[, options]): requires at least one parameter"))
		}

		ss_type, err := call.Argument(0).ToInteger()
		if err != nil {
			panic(df.script.Vm().MakeSyntaxError("enable_subsystem(type: int[, options]): requires type to be integer"))
		}

		// TODO option passing

		return change_subsystem(df.script.Driver(), *df.script.Vm(), ss_type, true)

	}
}
func change_subsystem(driver script.DriverAPI, vm otto.Otto, ss_type int64, status bool) otto.Value {
	switch ss_type {
	case D_SUBSYSTEM_REST:
		if err := driver.SetSubsystemState(ss_type, status); err != nil {
			v, e := otto.ToValue(err.Error())
			if e != nil {
				panic(vm.MakeCustomError("conversion_error", "enable subsystem failed. unknown error"))
			}

			return v
		}
	default:
		panic(vm.MakeCustomError("driver_invalid_subsystem", "enable_subsystem: unsupported subsystem type"))
	}

	return otto.TrueValue()
}

// syntax:
//   disable_subsystem(type: int)
//
// enable_subsystem is used to turn on a driver subsystem. Optionally
// you can pass options to it, if they are supported.

type dfun_disable_subsystem struct {
	script *script.ScriptContext
}

func create_disable_subsystem(i *script.ScriptContext) script.ExposedFunction {
	return &dfun_disable_subsystem{
		script: i,
	}
}

func (dfun_disable_subsystem) RequiredPrivilegeLevel() privilege.Level {
	return privilege.PrivilegeRoot
}

func (df dfun_disable_subsystem) Function() func(call otto.FunctionCall) otto.Value {
	return func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) != 1 {
			panic(df.script.Vm().MakeSyntaxError("disable_subsystem(type: int): requires at least one parameter"))
		}

		ss_type, err := call.Argument(0).ToInteger()
		if err != nil {
			panic(df.script.Vm().MakeSyntaxError("disable_subsystem(type: int): requires type to be integer"))
		}

		return change_subsystem(df.script.Driver(), *df.script.Vm(), ss_type, false)

		return otto.TrueValue()
	}
}
