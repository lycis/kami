package driver

import (
	"time"
	log "github.com/Sirupsen/logrus"
)

func (driver *Driver) RunWorld() {
	for {
		// call heartbeat
		if time.Now().Sub(driver.lastHeartbeat) > time.Second*2 {
			driver.heartbeat()
		}

		time.Sleep(time.Millisecond*10)
	}
}
func (driver *Driver) heartbeat() {
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
