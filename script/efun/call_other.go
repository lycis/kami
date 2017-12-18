package efun

import (
	"gitlab.com/lycis/kami/kerror"
	"gitlab.com/lycis/kami/script"
	"github.com/robertkrimen/otto"
	"gitlab.com/lycis/kami/privilege"
)

type efunCallOther struct {
	script *script.ScriptContext
}

func createCallOther(i *script.ScriptContext) script.ExposedFunction {
	return &efunCallOther{
		script: i,
	}
}

func (e efunCallOther) RequiredPrivilegeLevel() privilege.Level {
	return privilege.PrivilegeBasic
}

func (e efunCallOther) Function() func(call otto.FunctionCall) otto.Value {
	return func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) < 2 {
			panic(e.script.Vm().MakeSyntaxError("call_other(<id>, <function>[, arg1...]) expects at least two arguments"))
		}

		id := call.Argument(0)
		funName := call.Argument(1)

		if !id.IsString() {
			panic(e.script.Vm().MakeSyntaxError("call_other expects first parameter to be string"))
		}

		if !funName.IsString() {
			panic(e.script.Vm().MakeSyntaxError("call_other expects second parameter to be string"))
		}

		entity := e.script.Driver().GetEntityById(id.String())
		if entity == nil {
			panic(e.script.Vm().MakeSyntaxError("entity does not exist"))
		}

		val, err := entity.Call(funName.String(), call.ArgumentList[2:])
		if err != nil {
			panic(e.script.Vm().MakeCustomError("call_other failed", kerror.ToError(err).Error()))
		}

		return val

	}
}
