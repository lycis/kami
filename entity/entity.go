package entity

import (
	log "github.com/Sirupsen/logrus"
	"github.com/lycis/kami/kerror"
	"github.com/lycis/kami/privilege"
	"github.com/robertkrimen/otto"
	"sync"
)

// Entity represents any in-game object that exists
type Entity struct {
	properties map[string]interface{}
	mutex      sync.Mutex

	script script_context_api
}

type script_context_api interface {
	Call(name string, this interface{}, args ...interface{}) (otto.Value, error)
	GetFunction(name string) (otto.Value, error)
	GrantPrivilege(lvl privilege.Level)
	PrivilegeLevel() privilege.Level
}

func NewEntity(ctx script_context_api) (*Entity, error) {
	e := &Entity{
		properties: make(map[string]interface{}),
		script:     ctx,
	}

	_, err := e.script.Call("$create", e)
	return e, err
}

func (e *Entity) Heartbeat() {
	f, err := e.script.GetFunction("$tick")
	if err != nil {
		log.WithError(kerror.ToError(err)).WithFields(log.Fields{"path": e.GetProp("$path"), "uuid": e.GetProp("$uuid")}).Error("Executing tick function failed.")
		return
	}

	if f.IsDefined() {
		_, err := e.script.Call("$tick", e)
		if err != nil {
			log.WithError(kerror.ToError(err)).WithFields(log.Fields{"path": e.GetProp("$path"), "uuid": e.GetProp("$uuid")}).Error("Executing tick function failed.")
			return
		}
	}
}

func (e *Entity) Call(funName string, args ...interface{}) (otto.Value, error) {
	return e.script.Call(funName, e, args...)
}

func (e Entity) Context() script_context_api {
	return e.script
}
