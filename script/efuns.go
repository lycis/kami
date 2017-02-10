package script

import (
	"github.com/robertkrimen/otto"
)

type ExposedFunction interface {
	Function() func(call otto.FunctionCall) otto.Value
	RequiredPrivilegeLevel() PrivilegeLevel
}

type EFunCreator func(ctx *ScriptContext) ExposedFunction

var eFuns map[string]EFunCreator

func ExposeFunction(name string, f EFunCreator) {
	if eFuns == nil {
		eFuns = make(map[string]EFunCreator)
	}

	eFuns[name] = f
}

func exposeStaticFunctions(ctx *ScriptContext) {
	for name, f := range eFuns {
		efun := f(ctx)

		// only expose function that the privilege level of the
		// context grants access to
		if ctx.privilegeLevel >= efun.RequiredPrivilegeLevel() {
			ctx.Vm().Set(name, efun.Function())
		}
	}
}
