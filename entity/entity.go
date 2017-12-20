package entity

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"gitlab.com/lycis/kami/kerror"
	"gitlab.com/lycis/kami/privilege"
	"github.com/robertkrimen/otto"
	"sync"
)

type FunctionInvocationError struct {
	Err    error
	Entity *Entity
}

func (f FunctionInvocationError) Error() string {
	return fmt.Sprintf(
		"function invocation failed (error: %s entity: %s)",
		f.Err.Error(),
		fmt.Sprintf("%s#%s", f.Entity.GetProp(P_SYS_PATH), f.Entity.GetProp(P_SYS_ID)))
}

// Entity represents any in-game object that exists
type Entity struct {
	properties  map[string]interface{}
	propMutex   sync.RWMutex
	scriptMutex sync.Mutex

	script script_context_api
}

type script_context_api interface {
	Call(name string, this interface{}, args ...interface{}) (otto.Value, error)
	GetFunction(name string) (otto.Value, error)
	GrantPrivilege(lvl privilege.Level)
	PrivilegeLevel() privilege.Level
	Vm() *otto.Otto
}

func NewEntity() *Entity {
	e := &Entity{
		properties: make(map[string]interface{}),
	}

	e.properties[P_SYS_ACTIVE] = true

	return e
}

func (e *Entity) Create(script script_context_api) error {
	defer e.scriptMutex.Unlock()
	e.scriptMutex.Lock()

	e.script = script
	_, err := e.script.Call("$create", e)
	return err
}

// Heartbeat will take care that the $tick function is called whenever a heartbeat of the
// driver occurs
func (e *Entity) Heartbeat() error {
	if !e.IsActive() {
		return nil
	}

	defer e.scriptMutex.Unlock()
	e.scriptMutex.Lock()

	f, err := e.script.GetFunction("$tick")
	if err != nil {
		log.WithError(kerror.ToError(err)).WithFields(log.Fields{"path": e.GetProp("$path"), "uuid": e.GetProp("$uuid")}).Error("Executing tick function failed.")
		return FunctionInvocationError{err, e}
	}

	if f.IsDefined() {
		this, err := e.script.Vm().ToValue(e)
		if err != nil {
			log.WithError(kerror.ToError(err)).WithFields(log.Fields{"path": e.GetProp("$path"), "uuid": e.GetProp("$uuid")}).Error("Executing tick function failed.")
			return FunctionInvocationError{err, e}
		}

		_, err = f.Call(this)
		if err != nil {
			log.WithError(kerror.ToError(err)).WithFields(log.Fields{"path": e.GetProp("$path"), "uuid": e.GetProp("$uuid")}).Error("Executing tick function failed.")
			return FunctionInvocationError{err, e}
		}
	}

	return nil
}

func (e *Entity) OnShutdown(reason string) {
	if !e.IsActive() {
		return
	}

	defer e.scriptMutex.Unlock()
	e.scriptMutex.Lock()

	f, err := e.script.GetFunction("$onShutdown")
	if err != nil {
		log.WithError(kerror.ToError(err)).WithFields(log.Fields{"path": e.GetProp("$path"), "uuid": e.GetProp("$uuid")}).Error("Executing $onShutdown function failed.")
		panic(FunctionInvocationError{err, e})
	}

	if f.IsDefined() {
		_, err := e.script.Call("$onShutdown", e, reason)
		if err != nil {
			log.WithError(kerror.ToError(err)).WithFields(log.Fields{"path": e.GetProp("$path"), "uuid": e.GetProp("$uuid")}).Error("Executing tick function failed.")
			panic(FunctionInvocationError{err, e})
		}
	}
}

// HasFunction returns true if a function of the given name is defined
// or will return false otherwise. This can be used to check if an entity
// has a function that you wish to call
func (e Entity) HasFunction(funName string) bool {
	defer e.scriptMutex.Unlock()
	e.scriptMutex.Lock()

	f, err := e.script.GetFunction(funName)
	if err != nil {
		return false
	}

	if f.IsDefined() {
		return true
	}

	return false
}

func (e *Entity) Call(funName string, args ...interface{}) (otto.Value, error) {
	if !e.IsActive() {
		return otto.UndefinedValue(), fmt.Errorf("function called on an inactive entity")
	}

	defer e.scriptMutex.Unlock()
	e.scriptMutex.Lock()
	return e.script.Call(funName, e, args...)
}

func (e Entity) Context() script_context_api {
	return e.script
}

func (e Entity) GetScriptPrivilegeLevel() privilege.Level {
	if e.script == nil {
		return privilege.PrivilegeNone
	}

	return e.script.PrivilegeLevel()
}

func (e *Entity) GetScriptReferenceEntity() *Entity {
	return e
}

// IsActive indicates if the entity is still active or not. Inactive entities
// are considered to be already destroyed or not-yet set up correctly and
// must not be accessed in any way.
func (e Entity) IsActive() bool {
	b, ok := e.GetProp(P_SYS_ACTIVE).(bool)
	if !ok {
		return false
	}

	return b
}
