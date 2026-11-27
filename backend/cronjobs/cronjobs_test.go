package cronjobs

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"ganigeorgiev.com/pocketbase/pocketbase/core"
	"ganigeorgiev.com/pocketbase/pocketbase/tests"

	// register the app migrations so test apps get the real schema
	_ "github.com/seriousm4x/upsnap/migrations"
)

// pingCounter is a fake ping that records how many pings ran in total and
// the highest number that ran at the same time.
type pingCounter struct {
	active atomic.Int32
	peak   atomic.Int32
	calls  atomic.Int32
}

func (c *pingCounter) ping(d *core.Record) (bool, error) {
	n := c.active.Add(1)
	for {
		p := c.peak.Load()
		if n <= p || c.peak.CompareAndSwap(p, n) {
			break
		}
	}
	time.Sleep(200 * time.Millisecond)
	c.active.Add(-1)
	c.calls.Add(1)
	return true, nil
}

// waitForCalls waits until want pings have finished, then a little longer so
// that any extra ping started by an overlapping tick is counted as well.
func (c *pingCounter) waitForCalls(t *testing.T, want int32) {
	t.Helper()
	deadline := time.Now().Add(10 * time.Second)
	for c.calls.Load() < want || c.active.Load() > 0 {
		if time.Now().After(deadline) {
			t.Fatalf("Timed out waiting for %d pings, got %d", want, c.calls.Load())
		}
		time.Sleep(10 * time.Millisecond)
	}
	time.Sleep(300 * time.Millisecond)
	for c.active.Load() > 0 {
		time.Sleep(10 * time.Millisecond)
	}
}

// newPingTestApp sets up a test app with the given number of online devices
// and a stubbed ping. It returns the fake ping and the ping job registered by
// SetPingJobs, which the test runs directly instead of waiting for the cron.
func newPingTestApp(t *testing.T, devices int, pingCmd string) (*pingCounter, func()) {
	t.Helper()
	app, err := tests.NewTestApp(t.TempDir())
	if err != nil {
		t.Fatalf("Got unexpected error: %v", err)
	}
	t.Cleanup(app.Cleanup)

	settingsCollection, err := app.FindCollectionByNameOrId("settings_private")
	if err != nil {
		t.Fatalf("Got unexpected error: %v", err)
	}
	settings := core.NewRecord(settingsCollection)
	// the test runs the job itself, so the schedule never fires
	settings.Set("interval", "0 0 1 1 *")
	settings.Set("lazy_ping", false)
	settings.Set("track_ip_interval", "")
	if err := app.Save(settings); err != nil {
		t.Fatalf("Got unexpected error: %v", err)
	}

	collection, err := app.FindCollectionByNameOrId("devices")
	if err != nil {
		t.Fatalf("Got unexpected error: %v", err)
	}
	for i := range devices {
		device := core.NewRecord(collection)
		device.Set("name", fmt.Sprintf("device-%d", i))
		device.Set("ip", fmt.Sprintf("192.0.2.%d", i+1))
		device.Set("netmask", "255.255.255.0")
		device.Set("mac", fmt.Sprintf("02:00:00:00:00:%02x", i+1))
		device.Set("ping_cmd", pingCmd)
		// the fake ping reports online, so no status save runs after the test
		device.Set("status", "online")
		if err := app.Save(device); err != nil {
			t.Fatalf("Got unexpected error: %v", err)
		}
	}

	counter := &pingCounter{}
	orig := pingDevice
	pingDevice = counter.ping
	t.Cleanup(func() { pingDevice = orig })

	SetPingJobs(app)
	entries := CronPing.Entries()
	if len(entries) != 1 {
		t.Fatalf("Expected 1 ping job, got %d", len(entries))
	}
	return counter, entries[0].Job.Run
}

func TestPingJobLimitsConcurrentPingCommands(t *testing.T) {
	counter, tick := newPingTestApp(t, 13, "ssh relay ping -c 1 192.0.2.1")
	tick()
	counter.waitForCalls(t, 13)

	if got, limit := counter.peak.Load(), int32(cap(pingCmdSlots)); got > limit {
		t.Errorf("Expected at most %d ping commands at once, got %d", limit, got)
	}
	if got := counter.calls.Load(); got != 13 {
		t.Errorf("Expected 13 pings, got %d", got)
	}
}

func TestPingJobSkipsDevicesStillBeingPinged(t *testing.T) {
	counter, tick := newPingTestApp(t, 13, "ssh relay ping -c 1 192.0.2.1")
	// a second tick fires while the first round is still running
	tick()
	tick()
	counter.waitForCalls(t, 13)

	if got := counter.calls.Load(); got != 13 {
		t.Errorf("Expected 13 pings across both ticks, got %d", got)
	}
}

func TestPingJobDoesNotLimitBuiltInPing(t *testing.T) {
	counter, tick := newPingTestApp(t, 13, "")
	tick()
	counter.waitForCalls(t, 13)

	if got := counter.peak.Load(); got != 13 {
		t.Errorf("Expected all 13 built-in pings at once, got %d", got)
	}
}
