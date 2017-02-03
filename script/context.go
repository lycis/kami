package script

import (
	"fmt"
	"github.com/robertkrimen/otto"
)

// When binding a variable to the script context that implements this
// interface it will receive the virtual machine instance used for parsing
// the script on execution.
type BindValueWithInstance interface {
	SetScriptInstance(instance *Instance)
}

// ScriptContext provides (duh!) a context for running a script.
// This is required for executing any control scripts and will take
// care of providing values and embedded functions (efuns)
type ScriptContext struct {
	libDir string
	bindings map[string]interface{}
	cache *ScriptCache
	instance *Instance
}

// NewContext generates a new ScriptContext for running a script.
func NewContext(libDir string, cache *ScriptCache) ScriptContext {
	return ScriptContext{
		libDir: libDir,
		bindings: make(map[string]interface{}),
		cache: cache,
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
	absPath := fmt.Sprintf("%s%s", ctx.libDir, rpath)
	content, err := ctx.cache.loadScript(absPath)
	if err != nil {
		return ToError(err)
	}

	vm := otto.New()

	ctx.instance = &Instance{
		vm: vm,
	}


	exposeStaticFunctions(vm)
	bindValues(ctx)


	compiledScript, err := vm.Compile(rpath, content)
	if err != nil {
		return ToError(err)
	}

	_, err = vm.Run(compiledScript)
	if err != nil {
		return ToError(err)
	}

	return nil
}

func bindValues(ctx *ScriptContext) {
	for vname, v := range ctx.bindings {
		bwi := v.(BindValueWithInstance)
		if _, ok := v.(BindValueWithInstance); ok {
			bwi.SetScriptInstance(ctx.instance)
		}

		ctx.instance.vm.Set(vname, v)
	}
}

// GetInstance returns the script instance that results
// from the executed script.
func (ctx ScriptContext) GetInstance() *Instance {
	return ctx.instance
}