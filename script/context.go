package script

import (
	"gitlab.com/lycis/kami/entity"
	"gitlab.com/lycis/kami/kerror"
	"gitlab.com/lycis/kami/privilege"
	"github.com/robertkrimen/otto"
	"path/filepath"
	"sync"
)

// Everything that creates a script context should implement this interface
// and refer to itself so we can trace what create a script. This is especially
// important for taking over creator information like the inherited privilege level
type ContextCreator interface {
	// GetPrivilegeLevel provides the script privilege level of the creator
	GetScriptPrivilegeLevel() privilege.Level

	// GetScriptReferenceEntity provides a pointer to the entity that
	// is associated with the code that created the script context.
	//
	// This might be nil if no entity is associated with the creator (e.g.
	// specific driver code might not have an entity as the init script)
	GetScriptReferenceEntity() *entity.Entity
}

// When binding a variable to the script context that implements this
// interface it will receive the virtual machine instance used for parsing
// the script on execution.
type BindValueWithInstance interface {
	SetScriptInstance(ctx *ScriptContext)
}

// ScriptContext provides (duh!) a context for running a script.
// This is required for executing any control scripts and will take
// care of providing values and embedded functions (efuns)
type ScriptContext struct {
	libDir         string
	bindings       map[string]interface{}
	cache          *ScriptCache
	vm             *otto.Otto
	driver         DriverAPI
	privilegeLevel privilege.Level
	creator        ContextCreator
	runMtx         sync.Mutex
}

func ContextForScript(driver DriverAPI, script, libDir string, cache *ScriptCache) (ScriptContext, error) {
	ctx := NewContext(driver, libDir, cache)
	if err := ctx.RunScript(script); err != nil {
		return ScriptContext{}, err
	}

	return ctx, nil
}

// NewContext generates a new ScriptContext for running a script.
func NewContext(driver DriverAPI, libDir string, cache *ScriptCache) ScriptContext {
	return ScriptContext{
		libDir:         libDir,
		bindings:       make(map[string]interface{}),
		cache:          cache,
		driver:         driver,
		privilegeLevel: privilege.PrivilegeNone,
	}
}

// Bind allows you to expose internal values and objects to the
// executed script.
func (ctx *ScriptContext) Bind(vname string, value interface{}) {
	ctx.bindings[vname] = value
}

// RunScrupt executes the script given from the path relative
// to the drivers library directory.
func (ctx *ScriptContext) RunScript(rpath string) error {
	defer ctx.runMtx.Unlock()
	ctx.runMtx.Lock()
	content, err := ctx.LoadScript(rpath)
	if err != nil {
		return kerror.ToError(err)
	}

	if ctx.vm == nil {
		ctx.vm = otto.New()
	}

	exposeStaticFunctions(ctx)
	bindValues(ctx)

	compiledScript, err := ctx.vm.Compile(rpath, content)
	if err != nil {
		return kerror.ToError(err)
	}

	_, err = ctx.vm.Run(compiledScript)
	if err != nil {
		return kerror.ToError(err)
	}

	return nil
}

func bindValues(ctx *ScriptContext) {
	for vname, v := range ctx.bindings {
		bwi := v.(BindValueWithInstance)
		if _, ok := v.(BindValueWithInstance); ok {
			bwi.SetScriptInstance(ctx)
		}

		ctx.vm.Set(vname, v)
	}
}

func (ctx ScriptContext) LoadScript(path string) (string, error) {
	return ctx.cache.loadScript(filepath.Join(ctx.libDir, path))
}

func (ctx ScriptContext) Call(name string, this interface{}, args ...interface{}) (otto.Value, error) {
	defer ctx.runMtx.Unlock()
	ctx.runMtx.Lock()
	return ctx.vm.Call(name, this, args...)
}

func (ctx ScriptContext) GetFunction(name string) (otto.Value, error) {
	f, err := ctx.vm.Get(name)
	if err != nil {
		return otto.UndefinedValue(), kerror.ToError(err)
	}
	if !f.IsFunction() {
		return otto.UndefinedValue(), nil
	}

	return f, nil
}

func (ctx ScriptContext) RaiseError(name, message string) {
	panic(ctx.vm.MakeCustomError(name, message))
}

func (ctx ScriptContext) Vm() *otto.Otto {
	return ctx.vm
}

func (ctx ScriptContext) Driver() DriverAPI {
	return ctx.driver
}

// GrantPrivilege sets the privilege level of this script context
// to the defined one. This can be used to allow scripts access
// to protected functions or restrict their access.
//
// By default a context is created with PrivilegeBasis level.
func (ctx *ScriptContext) GrantPrivilege(lvl privilege.Level) {
	ctx.privilegeLevel = lvl
}

func (ctx ScriptContext) PrivilegeLevel() privilege.Level {
	return ctx.privilegeLevel
}

func (ctx ScriptContext) Creator() ContextCreator {
	return ctx.creator
}

func (ctx *ScriptContext) SetCreator(c ContextCreator) {
	ctx.creator = c
}
