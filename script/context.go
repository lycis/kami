package script

import (
	"fmt"
	"github.com/robertkrimen/otto"
)

type ScriptContext struct {
	libDir string
	bindings map[string]interface{}
	cache *ScriptCache
	instance *Instance
}

func NewContext(libDir string, cache *ScriptCache) ScriptContext {
	return ScriptContext{
		libDir: libDir,
		bindings: make(map[string]interface{}),
		cache: cache,
	}
}

func (ctx *ScriptContext) Bind(vname string, value interface{}) {
	ctx.bindings[vname] = value
}

func (ctx *ScriptContext) RunScript(rpath string) error {
	absPath := fmt.Sprintf("%s%s", ctx.libDir, rpath)
	content, err := ctx.cache.loadScript(absPath)
	if err != nil {
		return err
	}

	vm := otto.New()

	for vname, v := range ctx.bindings {
		vm.Set(vname, v)
	}

	_, err = vm.Run(content)
	if err != nil {
		return err
	}

	ctx.instance = &Instance {
		vm: vm,
	}

	return nil
}

func (ctx ScriptContext) GetInstance() *Instance {
	return ctx.instance
}