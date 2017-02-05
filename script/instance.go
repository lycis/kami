package script

import (
	"github.com/robertkrimen/otto"
	"sync"
)

type Instance struct {
	Vm    *otto.Otto
	Cache *ScriptCache
	mutex sync.Mutex
}