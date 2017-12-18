package local

import (
	"fmt"
	"gitlab.com/lycis/kami/driver/dfun"
	"gitlab.com/lycis/kami/entity"
	"github.com/robertkrimen/otto"
	"sync"
	"time"
)

// RunWorld is the actual world loop for the local driver. It will pause about
// 10 msec after every run.
// TODO Improve... maybe use events and channels to run world loop only when
// required and get rid of the pausing
func (driver *Driver) RunWorld() {
	driver.running = true

	if hookFunc, hookSet := driver.hooks[dfun.H_WHEN_WORLD_RUN]; hookSet {
		hookFunc.Call(otto.UndefinedValue())
	}

	for driver.running {
		// call heartbeat
		if time.Now().Sub(driver.lastHeartbeat) > time.Second*2 {
			driver.heartbeat()
		}

		time.Sleep(time.Millisecond * 10)
	}
}

func (driver *Driver) heartbeat() {
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

func (driver *Driver) doHeartbeatForShard(instances []*entity.Entity, hbWg *sync.WaitGroup) {
	for _, e := range instances {
		//driver.Log.WithField("uuid", e.GetProp("$uuid")).Debug("Calling Heartbeat")
		go func() {
			defer hbWg.Done()
			if err := e.Heartbeat(); err != nil {
				hberror := err.(entity.FunctionInvocationError)
				if hbErrFunc, ok := driver.hooks[dfun.H_HB_ON_ERROR]; ok {
					ovEntity, err := hberror.Entity.Context().Vm().ToValue(hberror.Entity)
					if err != nil {
						driver.Log.WithField("entity", fmt.Sprintf("%s#%s", hberror.Entity.GetProp(entity.P_SYS_PATH), hberror.Entity.GetProp(entity.P_SYS_PATH))).WithError(hberror.Err).Error("Calling heratbeat error hook failed")
					}
					hbErrFunc.Call(ovEntity, hberror.Error())
				}
			}
		}()
	}
}
