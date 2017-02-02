package entity

import (
	"sync"
	"github.com/lycis/kami/script"
)

// Entity represents any in-game object that exists
type Entity struct {
	properties map[string]interface{}
	propMutex sync.Mutex

	script *script.Instance
}

func NewEntity(instance *script.Instance) (*Entity, error) {
	e := &Entity {
		properties: make(map[string]interface{}),
		script: instance,
	}

	_, err := e.script.Call("create", e)
	return e, err
}

func (e *Entity) SetProp(name string, value interface{}) {
	e.propMutex.Lock()
	defer e.propMutex.Unlock()

	e.properties[name] = value
}

func (e Entity) GetProp(name string) interface{} {
	e.propMutex.Lock()
	defer e.propMutex.Unlock()

	return e.properties[name]
}
