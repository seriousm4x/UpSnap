package networking

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/pocketbase/pocketbase/core"
	"github.com/seriousm4x/upsnap/logger"
)

func ShutdownDevice(device *core.Record) error {
	logger.Info.Println("Shutdown triggered for", device.GetString("name"))
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

	shutdown_cmd = strings.ReplaceAll(shutdown_cmd, "{{ DEVICE_IP }}", device.GetString("ip"))
	shutdown_cmd = strings.ReplaceAll(shutdown_cmd, "{{ DEVICE_MAC }}", device.GetString("mac"))

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	cmd := exec.CommandContext(ctx, shell, shell_arg, shutdown_cmd)
	SetProcessAttributes(cmd)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		logger.Error.Println(err)
	}

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	shutdownTimeout := device.GetInt("shutdown_timeout")
	if shutdownTimeout <= 0 {
		shutdownTimeout = 120
	}

	start := time.Now()

	for {
		select {
		case <-time.After(1 * time.Second):
			if time.Since(start) >= time.Duration(shutdownTimeout)*time.Second {
				if err := KillProcess(cmd.Process); err != nil {
					logger.Error.Println(err)
				}
				return fmt.Errorf("%s not offline after %d seconds", device.GetString("name"), shutdownTimeout)
			}
			isOnline, err := PingDevice(device)
			if err != nil {
				logger.Error.Println(err)
				return err
			}
			if !isOnline {
				if err := KillProcess(cmd.Process); err != nil {
					logger.Error.Println(err)
				}
				return nil
			}
		case err := <-done:
			if err != nil {
				if err := KillProcess(cmd.Process); err != nil {
					logger.Error.Println(err)
				}
				return fmt.Errorf("%s", stderr.String())
			}
		}
	}
}
