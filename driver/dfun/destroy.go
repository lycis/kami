package dfun

import (
	"gitlab.com/lycis/kami/kerror"
	"github.com/robertkrimen/otto"
	"gitlab.com/lycis/kami/privilege"
	"gitlab.com/lycis/kami/script"
)

type dfunDestroy struct {
	script *script.ScriptContext
}

func createDfunDestroy(i *script.ScriptContext) script.ExposedFunction {
	return &dfunDestroy{
		script: i,
	}
}

func (d dfunDestroy) Function() func(call otto.FunctionCall) otto.Value {
	return d.destroy
}

func (d dfunDestroy) RequiredPrivilegeLevel() privilege.Level {
	return privilege.PrivilegeBasic
}

func (d dfunDestroy) destroy(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) < 1 {
		panic(d.script.Vm().MakeSyntaxError("destroy(id: string) requires at least one parameter"))
	}

	if !call.Argument(0).IsString() {
		panic(d.script.Vm().MakeSyntaxError("spawn(id: string) requires 'id' to be string"))
	}

	id, err := call.Argument(0).ToString()
	if err != nil {
		panic(d.script.Vm().MakeCustomError("destroy error", kerror.ToError(err).Error()))
	}

	if err := d.script.Driver().RemoveEntity(id); err != nil {
		panic(d.script.Vm().MakeCustomError("destroy failed", kerror.ToError(err).Error()))
	}
	
	return otto.TrueValue()
}
