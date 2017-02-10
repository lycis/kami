package dfun

import (
	"github.com/lycis/kami/kerror"
	"github.com/lycis/kami/privilege"
	"github.com/lycis/kami/script"
	"github.com/robertkrimen/otto"
)

type dfun_get_entity_by_id struct {
	script *script.ScriptContext
}

func create_dfun_get_entity_by_id(i *script.ScriptContext) script.ExposedFunction {
	return &dfun_get_entity_by_id{
		script: i,
	}
}

func (e dfun_get_entity_by_id) RequiredPrivilegeLevel() privilege.Level {
	return privilege.PrivilegeBasic
}

func (df dfun_get_entity_by_id) Function() func(call otto.FunctionCall) otto.Value {
	return func(call otto.FunctionCall) otto.Value {
		if !call.Argument(0).IsString() {
			panic(df.script.Vm().MakeSyntaxError("get_entity_by_id(id) requires id to be string"))
		}

		id, err := call.Argument(0).ToString()
		if err != nil {
			df.script.Vm().MakeCustomError("get_entity_by_id", kerror.ToError(err).Error())
		}

		entity := df.script.Driver().GetEntityById(id)
		val, err := df.script.Vm().ToValue(entity)
		if err != nil {
			df.script.Vm().MakeCustomError("get_entity_by_id", kerror.ToError(err).Error())
		}

		return val
	}
}
