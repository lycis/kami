package local

import (
	"fmt"
	"github.com/lycis/kami/driver/dfun"
	"github.com/lycis/kami/entity"
	"github.com/robertkrimen/otto"
	"sync"
	"time"
)

func (driver *LocalDriver) RunWorld() {
	driver.running = true

	if hook_func, hook_set := driver.hooks[dfun.H_WHEN_WORLD_RUN]; hook_set {
		hook_func.Call(otto.UndefinedValue())
	}

	for driver.running {
		// call heartbeat
		if time.Now().Sub(driver.lastHeartbeat) > time.Second*2 {
			driver.heartbeat()
		}

		time.Sleep(time.Millisecond * 10)
	}
}

func (driver *LocalDriver) heartbeat() {
	var hbWg sync.WaitGroup
	for path, instances := range driver.entityInstances {
		driver.Log.WithField("path", path).Debug("Calling heartbeat for instance shard.")
		hbWg.Add(len(instances))
		go driver.doHeartbeatForShard(instances, &hbWg)
		//driver.Log.Debug("Waiting for heartbeat to be executed.")
		hbWg.Wait()
		//driver.Log.Debug("Heartbeat processed.")
	}

	driver.lastHeartbeat = time.Now()
}

func (driver *LocalDriver) doHeartbeatForShard(instances []*entity.Entity, hbWg *sync.WaitGroup) {
	for _, e := range instances {
		//driver.Log.WithField("uuid", e.GetProp("$uuid")).Debug("Calling Heartbeat")
		go func() {
			defer hbWg.Done()
			if err := e.Heartbeat(); err != nil {
				hberror := err.(entity.FunctionInvocationError)
				if hb_err_func, ok := driver.hooks[dfun.H_HB_ON_ERROR]; ok {
					ovEntity, err := hberror.Entity.Context().Vm().ToValue(hberror.Entity)
					if err != nil {
						driver.Log.WithField("entity", fmt.Sprintf("%s#%s", hberror.Entity.GetProp(entity.P_SYS_PATH), hberror.Entity.GetProp(entity.P_SYS_PATH))).WithError(hberror.Err).Error("Calling heratbeat error hook failed")
					}
					hb_err_func.Call(ovEntity, hberror.Error())
				}
			}
		}()
	}
}
