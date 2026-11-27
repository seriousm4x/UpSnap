//go:build !windows

package networking

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"ganigeorgiev.com/pocketbase/pocketbase/core"
)

func TestPingCmdTimeoutKillsCommandAndChildren(t *testing.T) {
	orig := pingCmdTimeout
	pingCmdTimeout = 500 * time.Millisecond
	t.Cleanup(func() { pingCmdTimeout = orig })

	dir := t.TempDir()
	t.Setenv("UPSNAP_PING_TEST_DIR", dir)
	device := core.NewRecord(&core.Collection{})
	device.Set("ip", "192.0.2.1")
	// a child that would outlive the shell if only the shell were killed
	device.Set("ping_cmd", `(sleep 2; touch "$UPSNAP_PING_TEST_DIR/survived") & sleep 30`)

	start := time.Now()
	isUp, err := PingDevice(device)
	if err != nil {
		t.Fatalf("Got unexpected error: %v", err)
	}
	if isUp {
		t.Errorf("Expected a timed out ping command to report offline")
	}
	if elapsed := time.Since(start); elapsed > 2*time.Second {
		t.Errorf("Expected ping to return after the timeout, took %s", elapsed)
	}

	time.Sleep(2500 * time.Millisecond)
	if _, err := os.Stat(filepath.Join(dir, "survived")); err == nil {
		t.Errorf("Child process of the ping command kept running after the timeout")
	}
}
