package dfun

import (
	"gitlab.com/lycis/kami/entity"
	"gitlab.com/lycis/kami/kerror"
	"gitlab.com/lycis/kami/privilege"
	"gitlab.com/lycis/kami/script"
	"github.com/robertkrimen/otto"
)

// spawn creates a new entity from a given script. The $create() method of
// the script will be invoked on entity creation and should be used to set
// properties or execute arbitary code on entity creation.
//
// To create an exclusive entity set the "exclusive" parameter to true.
type dfun_spawn struct {
	script *script.ScriptContext
}

func create_dfun_spawn(i *script.ScriptContext) script.ExposedFunction {
	return &dfun_spawn{
		script: i,
	}
}

func (e dfun_spawn) RequiredPrivilegeLevel() privilege.Level {
	return privilege.PrivilegeBasic
}

func (e dfun_spawn) Function() func(call otto.FunctionCall) otto.Value {
	return e.spawn
}

func (dfun dfun_spawn) spawn(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) < 1 {
		panic(dfun.script.Vm().MakeSyntaxError("spawn(script, exclusive=flase) requires at least one parameter"))
	}

	if !call.Argument(0).IsString() {
		panic(dfun.script.Vm().MakeSyntaxError("spawn(script, exclusive=flase) requires 'script' to be string"))
	}

	path, err := call.Argument(0).ToString()
	if err != nil {
		dfun.script.Vm().MakeCustomError("spawn error", kerror.ToError(err).Error())
	}

	exclusive := false
	if call.Argument(1).IsDefined() {
		if !call.Argument(1).IsBoolean() {
			panic(dfun.script.Vm().MakeSyntaxError("spawn(script, exclusive=flase) requires 'exclusive' to be boolean"))
		}

		exclusive, err = call.Argument(1).ToBoolean()
		if err != nil {
			panic(dfun.script.Vm().MakeCustomError("spawn error", kerror.ToError(err).Error()))
		}
	}

	var e *entity.Entity

	if exclusive {
		e, err = dfun.script.Driver().SpawnExclusive(path, dfun.script.Creator())
	} else {
		e, err = dfun.script.Driver().SpawnEntity(path, dfun.script.Creator())
	}

	if err != nil {
		panic(dfun.script.Vm().MakeCustomError("spawn error", kerror.ToError(err).Error()))
	}

	value, err := dfun.script.Vm().ToValue(e)
	if err != nil {
		dfun.script.Vm().MakeCustomError("spawn error", kerror.ToError(err).Error())
	}

	return value
}
