package script

import (
	"fmt"
	"os"
	"io/ioutil"
	"github.com/robertkrimen/otto"
)

type ScriptContext struct {
	libDir string
	bindings map[string]interface{}
}

func NewContext(libDir string) ScriptContext {
	return ScriptContext{
		libDir: libDir,
		bindings: make(map[string]interface{}),
	}
}

func (ctx *ScriptContext) Bind(vname string, value interface{}) {
	ctx.bindings[vname] = value
}

func (ctx *ScriptContext) RunScript(relativePath string) error {
	content, err := ctx.loadScript(relativePath)
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

	return nil
}

func (ctx ScriptContext) loadScript(rpath string) (string, error) {
	absPath := fmt.Sprintf("%s%s", ctx.libDir, rpath)

	if _, err := os.Stat(absPath); err != nil {
		return "", fmt.Errorf("Failed loading script '%s': %s", absPath, err)
	}

	content, err := ioutil.ReadFile(absPath)
	return string(content), err
}