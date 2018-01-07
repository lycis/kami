package dfun

import (
	"gitlab.com/lycis/kami/kerror"
	"github.com/robertkrimen/otto"
	"gitlab.com/lycis/kami/privilege"
	"gitlab.com/lycis/kami/script"
)

type dfunSendUserEvent struct {
	script *script.ScriptContext
}

func createDfunSendUserEvent(i *script.ScriptContext) script.ExposedFunction {
	return &dfunSendUserEvent{
		script: i,
	}
}

func (d dfunSendUserEvent) Function() func(call otto.FunctionCall) otto.Value {
	return d.sendEvent
}

func (d dfunSendUserEvent) RequiredPrivilegeLevel() privilege.Level {
	return privilege.PrivilegeBasic
}

func (d dfunSendUserEvent) sendEvent(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) < 1 {
		panic(d.script.Vm().MakeSyntaxError("send_user_event(id: string, payload: string) requires two parameters"))
	}

	if !call.Argument(0).IsString() {
		panic(d.script.Vm().MakeSyntaxError("send_user_event(id: string, payload) requires 'id' to be string"))
	}

	id, err := call.Argument(0).ToString()
	if err != nil {
		panic(d.script.Vm().MakeCustomError("send_user_event argument 'id' conversion failed", kerror.ToError(err).Error()))
	}

	payload, err := call.Argument(1).ToString()
	if err != nil {
		panic(d.script.Vm().MakeCustomError("send_user_event argument 'payload' conversion failed", kerror.ToError(err).Error()))
	}

	if err := d.script.Driver().QueueUserEvent(id, payload); err != nil {
		panic(d.script.Vm().MakeCustomError("send_user_event failed", kerror.ToError(err).Error()))
	}
	
	return otto.TrueValue()
}
