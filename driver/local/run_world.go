package local

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/lycis/kami/driver/dfun"
	"github.com/lycis/kami/entity"
	"time"
)

func (driver *LocalDriver) RunWorld() {
	for {
		// call heartbeat
		if time.Now().Sub(driver.lastHeartbeat) > time.Second*2 {
			driver.heartbeat()
		}

		time.Sleep(time.Millisecond * 10)
	}
}
func (driver *LocalDriver) heartbeat() {
	for path, instances := range driver.entityInstances {
		log.WithField("path", path).Debug("Calling heartbeat for instance shard.")
		go func() {
			defer func() {
				if err := recover(); err != nil {
					hberror := err.(entity.HeartbeatError)
					if hb_err_func, ok := driver.hooks[dfun.H_HB_ON_ERROR]; ok {
						ovEntity, err := hberror.Entity.Context().Vm().ToValue(hberror.Entity)
						if err != nil {
							log.WithField("entity", fmt.Sprintf("%s#%s", hberror.Entity.GetProp(entity.P_SYS_PATH), hberror.Entity.GetProp(entity.P_SYS_PATH))).WithError(hberror.Error).Error("Calling heratbeat error hook failed")
						}
						hb_err_func.Call(ovEntity, hberror.Error.Error())
					}
				}
			}()
			for _, e := range instances {
				e.Heartbeat()
			}
		}()
	}

	driver.lastHeartbeat = time.Now()
}
