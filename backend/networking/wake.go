package networking

import (
	"errors"
	"time"

	"github.com/pocketbase/pocketbase/models"
)

func WakeDevice(device *models.Record) error {
	err := SendMagicPacket(device)
	if err != nil {
		return err
	}

	// check state every second for 2 min
	start := time.Now()
	for {
		time.Sleep(1 * time.Second)
		isOnline := PingDevice(device)
		if isOnline {
			return nil
		}
		if time.Since(start) >= 2*time.Minute {
			break
		}
	}
	return errors.New("device not online after 1 min")
}
