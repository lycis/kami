package entity

const (
	P_SYS_ID = "$uuid"
	P_SYS_PATH = "$path"
	P_SYS_EXCLUSIVE = "$exclusive"
)

// SetProp assigns a value to an entity property.
func (e *Entity) SetProp(name string, value interface{}) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	e.properties[name] = value
}

// GetProp returns you the value of a previously set property
// or nil (undefined) if the property is not set.
func (e Entity) GetProp(name string) interface{} {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	return e.properties[name]
}