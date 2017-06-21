package dfun

import (
	"github.com/lycis/kami/kerror"
	"github.com/lycis/kami/privilege"
	"github.com/lycis/kami/script"
	"github.com/robertkrimen/otto"
)

// set_driver_hook is used to set specific driver callback functions.
//
// This functions requires the executing script to have ROOT privileges.
//
// Hooks:
//   H_HB_ON_ERROR - function is called when an error occurs in the
//                     heartbeat call
//  H_WHEN_WORLD_RUN - this function is called when the world enters the running state (this = driver!)

const (
	H_HB_ON_ERROR    = 0
	H_WHEN_WORLD_RUN = 1
)

type dfun_set_driver_hook struct {
	script *script.ScriptContext
}

func create_dfun_set_driver_hook(i *script.ScriptContext) script.ExposedFunction {
	return &dfun_set_driver_hook{
		script: i,
	}
}

func (e dfun_set_driver_hook) RequiredPrivilegeLevel() privilege.Level {
	return privilege.PrivilegeRoot
}

func (df dfun_set_driver_hook) Function() func(call otto.FunctionCall) otto.Value {
	return func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) != 2 {
			panic(df.script.Vm().MakeSyntaxError("set_driver_hook(hook: int, value) requires two parameters"))
		}

		hook_id, err := call.Argument(0).ToInteger()
		if err != nil {
			panic(df.script.Vm().MakeSyntaxError("set_driver_hook(hook-id: int, value) requires hook to be int"))
		}

		if err := df.script.Driver().SetHook(hook_id, call.Argument(1)); err != nil {
			panic(df.script.Vm().MakeCustomError("set_driver_hook", kerror.ToError(err).Error()))
		}

		return otto.UndefinedValue()
	}
}
