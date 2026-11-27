//go:build !windows

package cronjobs

import (
	"fmt"
	"testing"
	"time"

	"ganigeorgiev.com/pocketbase/pocketbase/core"
	"ganigeorgiev.com/pocketbase/pocketbase/tests"
)

func TestPingJobHungCommandsDoNotStarveOtherDevices(t *testing.T) {
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
	settings.Set("interval", "0 0 1 1 *")
	settings.Set("lazy_ping", false)
	if err := app.Save(settings); err != nil {
		t.Fatalf("Got unexpected error: %v", err)
	}

	collection, err := app.FindCollectionByNameOrId("devices")
	if err != nil {
		t.Fatalf("Got unexpected error: %v", err)
	}
	newDevice := func(i int, pingCmd string, status string) *core.Record {
		device := core.NewRecord(collection)
		device.Set("name", fmt.Sprintf("device-%d", i))
		device.Set("ip", fmt.Sprintf("192.0.2.%d", i+1))
		device.Set("netmask", "255.255.255.0")
		device.Set("mac", fmt.Sprintf("02:00:00:00:00:%02x", i+1))
		device.Set("ping_cmd", pingCmd)
		device.Set("status", status)
		if err := app.Save(device); err != nil {
			t.Fatalf("Got unexpected error: %v", err)
		}
		return device
	}
	// more hung commands than there are slots for ping commands
	hung := []*core.Record{}
	for i := range cap(pingCmdSlots) + 1 {
		hung = append(hung, newDevice(i, "sleep 60", "online"))
	}

	SetPingJobs(app)
	start := time.Now()
	tick := CronPing.Entries()[0].Job.Run
	tick()
	// wait until hung commands hold every slot, then add a healthy device
	// and let the next tick reach it
	for len(pingCmdSlots) < cap(pingCmdSlots) {
		if time.Since(start) > 5*time.Second {
			t.Fatalf("Hung ping commands did not start")
		}
		time.Sleep(10 * time.Millisecond)
	}
	healthy := newDevice(len(hung), "exit 0", "offline")
	tick()

	waitForStatus := func(d *core.Record, want string, within time.Duration) {
		t.Helper()
		for {
			fresh, err := app.FindRecordById("devices", d.Id)
			if err != nil {
				t.Fatalf("Got unexpected error: %v", err)
			}
			if fresh.GetString("status") == want {
				return
			}
			if time.Since(start) > within {
				t.Fatalf("%s is not %s after %s", d.GetString("name"), want, within)
			}
			time.Sleep(50 * time.Millisecond)
		}
	}
	// the healthy device waits at most for one round of hung commands to time out
	waitForStatus(healthy, "online", 15*time.Second)
	t.Logf("healthy device online after %s", time.Since(start).Round(100*time.Millisecond))
	for _, d := range hung {
		waitForStatus(d, "offline", 30*time.Second)
	}
	t.Logf("all hung devices offline after %s", time.Since(start).Round(100*time.Millisecond))
}
