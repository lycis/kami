package efun

import (
	"github.com/robertkrimen/otto"
	"github.com/lycis/kami/script"
)

type efunInclude struct {
	script *script.ScriptContext
}

func CreateIncludeEfun(i *script.ScriptContext) script.ExposedFunction {
	return &efunInclude{
		script: i,
	}
}

func (e efunInclude) Function() func(call otto.FunctionCall) otto.Value {
	return func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) < 1 {
			panic(e.script.Vm().MakeSyntaxError("include(<file>[, ...] requires at least one file"))
		}

		vm := e.script.Vm()

		for i := 0; i < len(call.ArgumentList); i++ {
			file := call.Argument(0)
			if !file.IsString() {
				panic(vm.MakeSyntaxError("include only takes string parameters"))
			}

			content, err := e.script.LoadScript(file.String())
			if err != nil {
				panic(vm.MakeCustomError("include", err.Error()))
			}

			script, err := vm.Compile(file.String(), content)
			if err != nil {
				panic(vm.MakeCustomError("include", err.Error()))
			}

			_, err = vm.Run(script)
			if err != nil {
				panic(vm.MakeCustomError("include", err.Error()))
			}
		}

		return otto.UndefinedValue()
	}
}
