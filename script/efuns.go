package script

import (
	"github.com/robertkrimen/otto"
)

type ExposedFunction func(call otto.FunctionCall) otto.Value

type ExposedFunctionProvider interface {
	ExposedFunctions() map[string]ExposedFunction
}

var eFunProviders []ExposedFunctionProvider

func RegisterExposedFunctionProvider(p ExposedFunctionProvider) {
	eFunProviders = append(eFunProviders, p)
}

func exposeStaticFunctions(vm *otto.Otto) {
	for _, provider := range eFunProviders {
		for name, fp := range provider.ExposedFunctions() {
			vm.Set(name, fp)
		}
	}
}