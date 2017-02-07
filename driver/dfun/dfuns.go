// The dfun package specifies exposed driver functions. They do directly
// hook into driver functionality and thus allow you to call driver
// code from within a script.
//
// These functions are available on the "_driver" objects.
//
// Example:
//     var e = _driver.SpawnInstance("/npc/john.js")
package dfun

import (
	"github.com/lycis/kami/entity"
	"github.com/lycis/kami/script"
)

func NewProvider(driver script.DriverAPI) *DfunProvider {
	return &DfunProvider{
		driver: driver,
	}
}

type DfunProvider struct {
	driver script.DriverAPI
	instance *script.ScriptContext
}

// SpawnInstance creates a new, non-exclusive and non-durable entity
// for the given script. It is a wrapper for Spawn(script, false).
func (p DfunProvider) SpawnInstance(script string) *entity.Entity {
	return p.Spawn(script, false)
}

// Spawn creates a new entity from a given script. The $create() method of
// the script will be invoked on entity creation and should be used to set
// properties or execute arbitary code on entity creation.
//
// To create an exclusive entity set the "exclusive" parameter to true.
func (p DfunProvider) Spawn(script string, exclusive bool) *entity.Entity {
	var e *entity.Entity
	var err error

	if exclusive {
		e, err = p.driver.SpawnExclusive(script)
	} else {
		e, err = p.driver.SpawnEntity(script)
	}

	if err != nil {
		p.driver.Logger().Errorf("spawn failed: %s", err)
		p.instance.RaiseError("spawn failed", err.Error())
		return nil
	}

	return e
}

/*func (p DfunProvider) FindEntityById(id string) *entity.Entity {
	return p.driver.GetEnityById(id)
}*/

func (p *DfunProvider) SetScriptInstance(instance *script.ScriptContext) {
	p.instance = instance
}