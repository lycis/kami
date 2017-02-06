package local

import (
	"time"
	log "github.com/Sirupsen/logrus"
)

func (driver *LocalDriver) RunWorld() {
	for {
		// call heartbeat
		if time.Now().Sub(driver.lastHeartbeat) > time.Second*2 {
			driver.heartbeat()
		}

		time.Sleep(time.Millisecond*10)
	}
}
func (driver *LocalDriver) heartbeat() {
	for path, instances := range driver.entityInstances {
		log.WithField("path", path).Debug("Calling heartbeat for instance shard.")
		go func() {
			for _, e := range instances {
				e.Heartbeat()
			}
		}()
	}

	driver.lastHeartbeat = time.Now()
}
