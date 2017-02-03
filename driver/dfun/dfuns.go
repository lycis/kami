// The dfun package specifies exposed driver functions. They do directly
// hook into driver functionality and thus allow you to call driver
// code from within a script.
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

func (p DfunProvider) Spawn(rpath string) *entity.Entity {
	entity, err := p.driver.SpawnEntity(rpath)
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