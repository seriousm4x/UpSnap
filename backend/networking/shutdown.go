package networking

import (
	"errors"
	"os/exec"
	"time"

	"github.com/pocketbase/pocketbase/models"
)

func ShutdownDevice(device *models.Record) error {
	shutdown_cmd := device.GetString("shutdown_cmd")
	if shutdown_cmd != "" {
		cmd := exec.Command(shutdown_cmd)
		if err := cmd.Run(); err != nil {
			return err
		}
	}

	// check state every second for 2 min
	start := time.Now()
	for {
		time.Sleep(1 * time.Second)
		isOnline := PingDevice(device)
		if !isOnline {
			return nil
		}
		if time.Since(start) >= 2*time.Minute {
			break
		}
	}
	return errors.New("device not offline after 2 min")
}
