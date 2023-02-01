package networking

import (
	"bytes"
	"fmt"
	"os/exec"
	"runtime"
	"time"

	"github.com/pocketbase/pocketbase/models"
	"github.com/seriousm4x/upsnap/logger"
)

func ShutdownDevice(device *models.Record) error {
	logger.Info.Println("shutdown triggered for", device.GetString("name"))
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
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	if stderr.Len() > 0 {
		return fmt.Errorf("%s", stderr.String())
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
	return fmt.Errorf(device.GetString("name"), "not offline after 2 min")
}
