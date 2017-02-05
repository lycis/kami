package entity

import (
	"sync"
	"github.com/lycis/kami/script"
	log "github.com/Sirupsen/logrus"
)

// Entity represents any in-game object that exists
type Entity struct {
	properties map[string]interface{}
	mutex      sync.Mutex

	script *script.ScriptContext
}

func NewEntity(ctx *script.ScriptContext) (*Entity, error) {
	e := &Entity {
		properties: make(map[string]interface{}),
		script: ctx,
	}

	_, err := e.script.Call("$create", e)
	return e, err
}



func (e *Entity) Heartbeat() {
	f, err := e.script.GetFunction("$tick")
	if err != nil {
		log.WithError(script.ToError(err)).WithFields(log.Fields{"path": e.GetProp("$path"), "uuid": e.GetProp("$uuid")}).Error("Executing tick function failed.")
		return
	}

	if f.IsDefined() {
		_, err := e.script.Call("$tick", e)
		if err != nil {
			log.WithError(script.ToError(err)).WithFields(log.Fields{"path": e.GetProp("$path"), "uuid": e.GetProp("$uuid")}).Error("Executing tick function failed.")
			return
		}
	}
}