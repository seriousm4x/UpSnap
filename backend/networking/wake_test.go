//go:build !windows

package networking

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
)

func newWakeTestDevice(t *testing.T, command, pingCommand string, timeout int) (core.App, *core.Record) {
	t.Helper()
	app, err := tests.NewTestApp(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(app.Cleanup)
	device := core.NewRecord(&core.Collection{})
	device.Set("name", "test-device")
	device.Set("ip", "127.0.0.1")
	device.Set("mac", "02:00:00:00:00:01")
	device.Set("wake_cmd", command)
	device.Set("ping_cmd", pingCommand)
	device.Set("wake_timeout", timeout)
	return app, device
}

func TestCustomWakeWaitsForCommandWhenDeviceIsAlreadyOnline(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("UPSNAP_WAKE_TEST_DIR", dir)
	app, device := newWakeTestDevice(t,
		`sleep 2; printf done > "$UPSNAP_WAKE_TEST_DIR/completed"`, "exit 0", 10)

	if err := WakeDevice(app, device, "test"); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(dir, "completed"))
	if err != nil || string(data) != "done" {
		t.Fatalf("wake command was interrupted despite successful ping: data=%q err=%v", data, err)
	}
}

func TestCustomWakeReturnsFailureWhenDeviceIsAlreadyOnline(t *testing.T) {
	app, device := newWakeTestDevice(t, "sleep 2; echo wake-failed >&2; exit 7", "exit 0", 10)
	err := WakeDevice(app, device, "test")
	if err == nil || !strings.Contains(err.Error(), "wake-failed") {
		t.Fatalf("expected wake command failure despite successful ping, got %v", err)
	}
}

func TestCustomWakeStillWaitsForReachabilityAfterCommandFinishes(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("UPSNAP_WAKE_TEST_DIR", dir)
	app, device := newWakeTestDevice(t, "exit 0",
		`if [ -f "$UPSNAP_WAKE_TEST_DIR/pinged" ]; then exit 0; fi; touch "$UPSNAP_WAKE_TEST_DIR/pinged"; exit 1`, 10)

	if err := WakeDevice(app, device, "test"); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(dir, "pinged")); err != nil {
		t.Fatalf("wake returned without checking reachability: %v", err)
	}
}

func TestCustomWakeTimesOutWhenCommandSucceedsButDeviceStaysOffline(t *testing.T) {
	app, device := newWakeTestDevice(t, "exit 0", "exit 1", 2)
	err := WakeDevice(app, device, "test")
	if err == nil || !strings.Contains(err.Error(), "not online after 2 seconds") {
		t.Fatalf("expected reachability timeout, got %v", err)
	}
}

func TestCustomWakeTimeoutStillStopsRunningCommand(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("UPSNAP_WAKE_TEST_DIR", dir)
	app, device := newWakeTestDevice(t,
		`sleep 3; touch "$UPSNAP_WAKE_TEST_DIR/completed"`, "exit 0", 1)
	err := WakeDevice(app, device, "test")
	if err == nil || !strings.Contains(err.Error(), "not online after 1 seconds") {
		t.Fatalf("expected command timeout even though ping succeeds, got %v", err)
	}
	// A surviving shell/child would write this file after its sleep finishes.
	time.Sleep(3 * time.Second)
	if _, err := os.Stat(filepath.Join(dir, "completed")); !os.IsNotExist(err) {
		t.Fatalf("timed-out wake command continued running: %v", err)
	}
}
