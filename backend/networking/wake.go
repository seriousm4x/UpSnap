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

	// we wait 1 minute for the device to come up
	// after that, we check the state
	time.Sleep(1 * time.Minute)
	isOnline := PingDevice(device)
	if isOnline {
		return nil
	}
	return errors.New("device not online after 1 min")
}
