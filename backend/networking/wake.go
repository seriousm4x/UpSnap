package networking

import (
	"time"

	"github.com/pocketbase/pocketbase/models"
	"github.com/seriousm4x/upsnap/backend/logger"
)

func WakeDevice(device *models.Record) bool {
	err := SendMagicPacket(device.GetString("mac"), device.GetString("netmask"))
	if err != nil {
		logger.Error.Println(err)
		return false
	}

	// we wait 1 minute for the device to come up
	// after that, we check the state
	time.Sleep(1 * time.Minute)
	isOnline := PingDevice(device)
	if isOnline {
		return true
	} else {
		return false
	}
}
