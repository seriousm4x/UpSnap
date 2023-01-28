package networking

import (
	"errors"
	"fmt"
	"os/exec"
	"runtime"
	"time"

	"github.com/pocketbase/pocketbase/models"
)

func ShutdownDevice(device *models.Record) error {
	shutdown_cmd := device.GetString("shutdown_cmd")
	if shutdown_cmd == "" {
		return fmt.Errorf("%s: no shutdown_cmd definded", device.GetString("name"))
	}

	var shell string
	var shell_arg string
	if runtime.GOOS == "windows" {
		shell = "cmd"
		shell_arg = "/C"
	} else {
		shell = "/bin/sh"
		shell_arg = "-c"
	}

	cmd := exec.Command(shell, shell_arg, shutdown_cmd)
	if err := cmd.Run(); err != nil {
		return err
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
