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
	log "github.com/Sirupsen/logrus"
	"github.com/lycis/kami/script"
)

type DriverInterface interface {
	SpawnEntity(rpath string) (*entity.Entity, error)
	Logger() *log.Logger
}

func NewProvider(driver DriverInterface) *DfunProvider {
	return &DfunProvider{
		driver: driver,
	}
}

type DfunProvider struct {
	driver DriverInterface
	instance *script.Instance
}

// SpawnInstance creates a new, non-exclusive and non-durable entity
// for the given script. It is a wrapper for Spawn(script, false).
func (p DfunProvider) SpawnInstance(script string) *entity.Entity {
	return p.Spawn(script, false)
}

// Spawn creates a new entity from a given script. The create() method of
// the script will be invoked on entity creation and should be used to set
// properties or execute arbitary code on entity creation.
//
// To create an exclusive entity set the "exclusive" parameter to true.
func (p DfunProvider) Spawn(script string, exclusive bool) *entity.Entity {
	entity, err := p.driver.SpawnEntity(script)
	if err != nil {
		p.driver.Logger().Errorf("spawn failed: %s", err)
		p.instance.RaiseError("spawn failed", err.Error())
		return nil
	}

	return entity
}

func (p *DfunProvider) SetScriptInstance(instance *script.Instance) {
	p.instance = instance
}