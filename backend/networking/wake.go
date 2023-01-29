package networking

import (
	"fmt"
	"time"

	"github.com/pocketbase/pocketbase/models"
	"github.com/seriousm4x/upsnap/logger"
)

func WakeDevice(device *models.Record) error {
	logger.Info.Println("wake triggered for", device.GetString("name"))
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
	return fmt.Errorf(device.GetString("name"), "not online after 2 min")
}
