// The dfun package specifies exposed driver functions. They do directly
// hook into driver functionality and thus allow you to call driver
// code from within a script.
package dfun

import (
	"github.com/lycis/kami/entity"
	log "github.com/Sirupsen/logrus"
)

type DriverInterface interface {
	SpawnEntity(rpath string) (*entity.Entity, error)
	Logger() *log.Logger
}

func NewProvider(driver DriverInterface) dfunProvider {
	return dfunProvider{
		driver: driver,
	}
}

type dfunProvider struct {
	driver DriverInterface
}

func (p dfunProvider) Spawn(rpath string) *entity.Entity {
	entity, err := p.driver.SpawnEntity(rpath)
	if err != nil {
		p.driver.Logger().Errorf("spawn failed: %s", err)
		panic(err)
		return nil
	}

	return entity
}