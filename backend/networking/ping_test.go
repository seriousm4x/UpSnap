package networking

import (
	"testing"

	"github.com/pocketbase/pocketbase/models"
)

func TestPingDevice(t *testing.T) {
	collection := &models.Collection{}
	device := models.NewRecord(collection)
	device.Set("ip", "8.8.8.8")

	isUp := PingDevice(device)
	if !isUp {
		t.Fatalf("IP %s is not up", device.GetString("ip"))
	}
}

func TestCheckPort(t *testing.T) {
	ip := "8.8.8.8"
	port := "443"
	isOpen := CheckPort(ip, port)
	if !isOpen {
		t.Fatalf("Port %s:%s is not open", ip, port)
	}
}
