package networking

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"time"

	"github.com/pocketbase/pocketbase/core"
	"github.com/seriousm4x/upsnap/logger"
)

func WakeDevice(device *core.Record) error {
	logger.Info.Println("Wake triggered for", device.GetString("name"))

	wakeTimeout := device.GetInt("wake_timeout")
	if wakeTimeout <= 0 {
		wakeTimeout = 120
	}

	wake_cmd := device.GetString("wake_cmd")
	if wake_cmd != "" {
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
		defer cancel()

		cmd := exec.CommandContext(ctx, shell, shell_arg, wake_cmd)
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

		start := time.Now()

		for {
			select {
			case <-time.After(1 * time.Second):
				if time.Since(start) >= time.Duration(wakeTimeout)*time.Second {
					if err := KillProcess(cmd.Process); err != nil {
						logger.Error.Println(err)
					}
					return fmt.Errorf("%s not online after %d seconds", device.GetString("name"), wakeTimeout)
				}
				isOnline, err := PingDevice(device)
				if err != nil {
					logger.Error.Println(err)
					return err
				}
				if isOnline {
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
	} else {
		err := SendMagicPacket(device)
		if err != nil {
			return err
		}

		start := time.Now()
		for {
			time.Sleep(1 * time.Second)
			isOnline, err := PingDevice(device)
			if err != nil {
				logger.Error.Println(err)
				return err
			}
			if isOnline {
				return nil
			}
			if time.Since(start) >= time.Duration(wakeTimeout)*time.Second {
				break
			}
		}
		return fmt.Errorf("%s not online after %d seconds", device.GetString("name"), wakeTimeout)
	}
}
