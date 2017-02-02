package script

import "github.com/robertkrimen/otto"

type Instance struct {
	vm *otto.Otto
}

func (instance Instance) Call(name string, this interface{}, args ...interface{}) (otto.Value, error) {
	return instance.vm.Call(name, this, args...)
}