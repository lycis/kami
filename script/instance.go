package script

import (
	"github.com/robertkrimen/otto"
	"sync"
)

type Instance struct {
	vm *otto.Otto
	mutex sync.Mutex
}

func (instance Instance) Call(name string, this interface{}, args ...interface{}) (otto.Value, error) {
	instance.mutex.Lock()
	defer instance.mutex.Unlock()
	return instance.vm.Call(name, this, args...)
}

func (i Instance) GetFunction(name string) (otto.Value, error) {
	f, err := i.vm.Get(name)
	if err != nil {
		return otto.UndefinedValue(), ToError(err)
	}
	if !f.IsFunction() {
		return otto.UndefinedValue(), nil
	}

	return  f, nil
}

func (instance Instance) RaiseError(name, message string) {
	panic(instance.vm.MakeCustomError(name, message))
}