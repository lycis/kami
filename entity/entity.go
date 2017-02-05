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

	script *script.Instance
}

func NewEntity(instance *script.Instance) (*Entity, error) {
	e := &Entity {
		properties: make(map[string]interface{}),
		script: instance,
	}

	_, err := e.script.Call("$create", e)
	return e, err
}

func (e *Entity) SetProp(name string, value interface{}) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	e.properties[name] = value
}

func (e Entity) GetProp(name string) interface{} {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	return e.properties[name]
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