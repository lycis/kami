package script

import (
	"github.com/robertkrimen/otto"
)

type ExposedFunction interface {
	Function() func(call otto.FunctionCall) otto.Value
}

type EFunCreator func(ctx *ScriptContext) ExposedFunction

var eFuns map[string]EFunCreator

func ExposeFunction(name string, f EFunCreator)  {
	if eFuns == nil {
		eFuns = make(map[string]EFunCreator)
	}

	eFuns[name] = f
}

func exposeStaticFunctions(ctx *ScriptContext) {
	for name, f := range eFuns {
		efun := f(ctx)
		ctx.Vm().Set(name, efun.Function())
	}
}