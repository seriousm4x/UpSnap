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

	// we wait 1 minute for the device to come up
	// after that, we check the state
	time.Sleep(1 * time.Minute)
	isOnline := PingDevice(device)
	if isOnline {
		return errors.New("device not offline after 1 min")
	}
	return nil
}
