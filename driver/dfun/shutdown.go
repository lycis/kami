package dfun

import (
	"github.com/lycis/kami/privilege"
	"github.com/lycis/kami/script"
	"github.com/robertkrimen/otto"
)

// syntax:
//   shutdown(reason: string)
//
// shutdown is used to turn off the driver and bring it to a clean end.
// The reason argument gives an explanation for the shutdown that will be
// tracked in the log and also passed to all active entities in their
// $onShutdown function.

type dfun_shutdown struct {
	script *script.ScriptContext
}

func create_shutdown(i *script.ScriptContext) script.ExposedFunction {
	return &dfun_shutdown{
		script: i,
	}
}

func (dfun_shutdown) RequiredPrivilegeLevel() privilege.Level {
	return privilege.PrivilegeRoot
}

func (df dfun_shutdown) Function() func(call otto.FunctionCall) otto.Value {
	return func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) != 1 {
			panic(df.script.Vm().MakeSyntaxError("shutdown(reason: string): requires one parameter"))
		}

		reason, err := call.Argument(0).ToString()
		if err != nil {
			panic(df.script.Vm().MakeSyntaxError("shutdown(reason: string): requires reason to be string"))
		}

		df.script.Driver().Shutdown(reason)
		return otto.UndefinedValue()
	}
}
