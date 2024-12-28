package networking

import (
	"fmt"
	"net"
	"testing"

	"github.com/pocketbase/pocketbase/core"
)

type Device struct {
	IP                 string
	MAC                string
	Netmask            string
	Password           string
	ExpectFailPassword bool
}

func TestSendMagicPacket(t *testing.T) {
	collection := &core.Collection{}
	devices_success := []Device{
		{
			IP:       "8.8.8.8",
			MAC:      "aa:aa:aa:aa:aa:aa",
			Netmask:  "255.0.0.0",
			Password: "",
		},
		{
			IP:       "1.1.1.1",
			MAC:      "bb:bb:bb:bb:bb:bb",
			Netmask:  "255.255.0.0",
			Password: "abcd",
		},
	}

	devices_fail := []Device{
		{
			IP:       "1.1.1.1",
			MAC:      "aa",
			Netmask:  "255.255.0.0",
			Password: "abcd",
		},
		{
			IP:       "1.1.1.1",
			MAC:      "cc:cc:cc:cc:cc:cc",
			Netmask:  "255.255.0.0",
			Password: "abcdefg",
		},
	}

	for _, v := range devices_success {
		device := core.NewRecord(collection)
		device.Set("ip", v.IP)
		device.Set("mac", v.MAC)
		device.Set("netmask", v.Netmask)
		device.Set("password", v.Password)

		if err := SendMagicPacket(device); err != nil {
			t.Fatal(err)
		}
	}
	for _, v := range devices_fail {
		device := core.NewRecord(collection)
		device.Set("ip", v.IP)
		device.Set("mac", v.MAC)
		device.Set("netmask", v.Netmask)
		device.Set("password", v.Password)

		if err := SendMagicPacket(device); err == nil {
			t.Fatal("should have cause an error")
		}
	}
}

func TestWakeUDP(t *testing.T) {
	devices := []Device{
		{
			IP:                 "8.8.8.8",
			MAC:                "aa:aa:aa:aa:aa:aa",
			Netmask:            "255.0.0.0",
			Password:           "",
			ExpectFailPassword: false,
		},
		{
			IP:                 "1.1.1.1",
			MAC:                "bb:bb:bb:bb:bb:bb",
			Netmask:            "255.255.0.0",
			Password:           "abcd",
			ExpectFailPassword: false,
		},
		{
			IP:                 "10.10.10.1",
			MAC:                "cc:cc:cc:cc:cc:cc",
			Netmask:            "255.255.0.0",
			Password:           "abcd",
			ExpectFailPassword: false,
		},
		{
			IP:                 "1.1.1.1",
			MAC:                "dd:dd:dd:dd:dd:dd",
			Netmask:            "255.255.0.0",
			Password:           "abcdefg",
			ExpectFailPassword: true,
		},
	}

	for _, v := range devices {
		parsedMac, err := net.ParseMAC(v.MAC)
		if err != nil {
			t.Fatalf(`%v (parsedMac "%v")`, err, parsedMac)
		}

		var bytePassword []byte
		if len(v.Password) == 0 || len(v.Password) == 4 || len(v.Password) == 6 {
			bytePassword = []byte(v.Password)
		} else if !v.ExpectFailPassword {
			t.Fatal("password must be 0, 4 or 6 characters long")
		}

		broadcastIp, err := getBroadcastIp(v.IP, v.Netmask)
		if err != nil {
			t.Fatal(err)
		}
		targetAddr := fmt.Sprintf("%s:%d", broadcastIp, 9)

		if err := wakeUDP(targetAddr, parsedMac, bytePassword); err != nil {
			t.Fatalf(`%v (targetAddr "%v", parsedMac "%v", bytePassword "%v")`, err, targetAddr, parsedMac, bytePassword)
		}
	}
}

func TestGetBroadcastIp(t *testing.T) {
	type address struct {
		IP              string
		Netmask         string
		BroadcastResult string
	}

	addresses := []address{
		{
			IP:              "192.168.0.1",
			Netmask:         "255.255.255.0",
			BroadcastResult: "192.168.0.255",
		}, {
			IP:              "10.1.1.1",
			Netmask:         "255.255.0.0",
			BroadcastResult: "10.1.255.255",
		}, {
			IP:              "10.2.3.4",
			Netmask:         "255.0.0.0",
			BroadcastResult: "10.255.255.255",
		},
	}

	for _, v := range addresses {
		broadcast, err := getBroadcastIp(v.IP, v.Netmask)
		if err != nil {
			t.Fatal(err)
		}
		if broadcast != v.BroadcastResult {
			t.Fatalf(`getBroadcastIp("%v, %v"), want %v, got %v`, v.IP, v.Netmask, v.BroadcastResult, broadcast)
		}
	}
}
