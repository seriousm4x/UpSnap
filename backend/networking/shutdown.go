package networking

import (
	"bytes"
	"context"
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

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	cmd := exec.CommandContext(ctx, shell, shell_arg, shutdown_cmd)
	SetProcessAttributes(cmd)
	logger.Debug.Println(cmd)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		logger.Error.Println(err)
	}

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	// check state every seconds for 2 min
	start := time.Now()
	for {
		select {
		case <-time.After(1 * time.Second):
			if time.Since(start) >= 2*time.Minute {
				if err := KillProcess(cmd.Process); err != nil {
					logger.Error.Println(err)
				}
				cancel()
				return fmt.Errorf("%s not offline after 2 minutes", device.GetString("name"))
			}
			isOnline := PingDevice(device)
			if !isOnline {
				if err := KillProcess(cmd.Process); err != nil {
					logger.Error.Println(err)
				}
				cancel()
				return nil
			}
		case err := <-done:
			if err != nil {
				if err := KillProcess(cmd.Process); err != nil {
					logger.Error.Println(err)
				}
				cancel()
				return fmt.Errorf("%s", stderr.String())
			}
		}
	}
}
