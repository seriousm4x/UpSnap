package networking

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/pocketbase/pocketbase/core"
	"github.com/seriousm4x/upsnap/logging"
)

func WakeDevice(app core.App, device *core.Record, trigger string) error {
	log := logging.Logger(app)
	log.Info("Wake triggered", "device", device.GetString("name"), "trigger", trigger)

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

    deviceIP, err := ResolveToIPAddr(device.GetString("ip"))
    if err != nil {
      return err
    }
    if deviceIP == "" {
      deviceIP = device.GetString("ip")
    }
		// Validate MAC address before replacing placeholders to prevent command injection
		deviceMAC := device.GetString("mac")
		if _, err := net.ParseMAC(deviceMAC); err != nil {
			return fmt.Errorf("invalid device MAC address: %q", deviceMAC)
		}
		wake_cmd = strings.ReplaceAll(wake_cmd, "{{ DEVICE_IP }}", deviceIP)
		wake_cmd = strings.ReplaceAll(wake_cmd, "{{ DEVICE_MAC }}", deviceMAC)

		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		cmd := exec.CommandContext(ctx, shell, shell_arg, wake_cmd)
		SetProcessAttributes(cmd)

		var stderr bytes.Buffer
		cmd.Stderr = &stderr

		if err := cmd.Start(); err != nil {
			log.Error("Failed to start wake command", "error", err)
			return err
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
					if cmd.Process != nil {
						if err := KillProcess(app, cmd.Process); err != nil {
							log.Error("Failed to kill wake command", "error", err)
						}
					}
					return fmt.Errorf("%s not online after %d seconds", device.GetString("name"), wakeTimeout)
				}
				isOnline, err := PingDevice(device)
				if err != nil {
					log.Error("Failed to ping device after wake", "error", err)
					return err
				}
				if isOnline {
					if cmd.Process != nil {
						if err := KillProcess(app, cmd.Process); err != nil {
							// Process might have already finished
						}
					}
					return nil
				}
			case err := <-done:
				if err != nil {
					if cmd.Process != nil {
						if err := KillProcess(app, cmd.Process); err != nil {
							log.Error("Failed to kill wake command", "error", err)
						}
					}
					return fmt.Errorf("%s", stderr.String())
				}
				// Command finished successfully, but we continue the loop until device is online or timeout
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
				log.Error("Failed to ping device after wake", "error", err)
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
