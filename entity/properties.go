package entity

import (
	"fmt"
)

const (
	P_SYS_ID        = "$uuid"
	P_SYS_PATH      = "$path"
	P_SYS_EXCLUSIVE = "$exclusive"
	P_SYS_ACTIVE    = "$active"
)

// SetProp assigns a value to an entity property.
func (e *Entity) SetProp(name string, value interface{}) {
	if !e.IsActive() {
		panic(fmt.Errorf("access to inactive entity"))
	}

	e.propMutex.Lock()
	defer e.propMutex.Unlock()

	e.properties[name] = value
}

// GetProp returns you the value of a previously set property
// or nil (undefined) if the property is not set.
func (e Entity) GetProp(name string) interface{} {
	if name != P_SYS_ACTIVE && !e.IsActive() {
		panic(fmt.Errorf("access to inactive entity"))
	}

	e.propMutex.RLock()
	defer e.propMutex.RUnlock()

	return e.properties[name]
}
