package efun

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"gitlab.com/lycis/kami/entity"
	"gitlab.com/lycis/kami/kerror"
	"gitlab.com/lycis/kami/privilege"
	"gitlab.com/lycis/kami/script"
	"github.com/robertkrimen/otto"
)

// grant_privilege(entity, level) set the privilege level of an
// entity to the given level. Note that it is not possible to grant
// a higher level than the caller has itself.

type efun_grant_privilege struct {
	script *script.ScriptContext
}

func create_grant_privilege(i *script.ScriptContext) script.ExposedFunction {
	return &efun_grant_privilege{
		script: i,
	}
}

func (e efun_grant_privilege) RequiredPrivilegeLevel() privilege.Level {
	return privilege.PrivilegeBasic
}

func (e efun_grant_privilege) Function() func(call otto.FunctionCall) otto.Value {
	return func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) != 2 {
			panic(e.script.Vm().MakeSyntaxError("grant_privilege(entity: (string|entity), level: int) requires two parameters"))
		}

		level, err := call.Argument(1).ToInteger()
		if err != nil {
			panic(e.script.Vm().MakeSyntaxError("grant_privilege(entity, level) requires level to be int"))
		}

		if e.script.PrivilegeLevel() < privilege.Level(level) {
			panic(e.script.Vm().MakeSyntaxError("grant_privilege(entity: (string|entity), level: int) can not grant higher level than caller"))
		}

		var target *entity.Entity
		if call.Argument(0).IsString() {
			id, err := call.Argument(0).ToString()
			if err != nil {
				panic(e.script.Vm().MakeCustomError("grant_privilege", kerror.ToError(err).Error()))
			}

			target = e.script.Driver().GetEntityById(id)
			if target == nil {
				panic(e.script.Vm().MakeCustomError("grant_privilege", fmt.Sprintf("entity with id '%s' not found", id)))
			}
		} else if call.Argument(0).IsObject() {
			o := call.Argument(0).Object()

			fret, err := o.Call("GetProp", entity.P_SYS_ID)
			if err != nil {
				panic(e.script.Vm().MakeCustomError("grant_privilege", kerror.ToError(err).Error()))
			}

			id, err := fret.ToString()
			if err != nil {
				panic(e.script.Vm().MakeCustomError("grant_privilege", kerror.ToError(err).Error()))
			}

			target = e.script.Driver().GetEntityById(id)
			if target == nil {
				panic(e.script.Vm().MakeCustomError("grant_privilege", fmt.Sprintf("entity with id '%s' not found", id)))
			}
		} else {
			panic(e.script.Vm().MakeSyntaxError("grant_privilege(entity, level) requires first argument to be string or entity"))
		}

		if target.Context() == nil {
			panic(e.script.Vm().MakeCustomError("grant_privilege", "grant_privilege must not be called during entity initialisation"))
		}

		originalLevel := target.Context().PrivilegeLevel()
		e.script.Driver().Logger().WithFields(logrus.Fields{"e": e.script.Creator().GetScriptReferenceEntity(), "from-level": originalLevel, "to-level": level, "entity": fmt.Sprintf("%s#%s", target.GetProp(entity.P_SYS_PATH), target.GetProp(entity.P_SYS_ID))}).Info("Changed privilege level.")
		target.Context().GrantPrivilege(privilege.Level(level))
		return otto.TrueValue()
	}
}
