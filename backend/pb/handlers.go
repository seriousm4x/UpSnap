package pb

import (
	"encoding/xml"
	"errors"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/models"
	"github.com/seriousm4x/upsnap/backend/logger"
	"github.com/seriousm4x/upsnap/backend/networking"
)

func HandlerWake(c echo.Context) error {
	record, err := App.Dao().FindFirstRecordByData("devices", "id", c.PathParam("id"))
	if err != nil {
		return apis.NewNotFoundError("The device does not exist.", err)
	}
	go func(*models.Record) {
		record.Set("status", "pending")
		App.Dao().SaveRecord(record)
		if err := networking.WakeDevice(record); err != nil {
			record.Set("status", "offline")
		} else {
			record.Set("status", "online")
		}
		App.Dao().SaveRecord(record)
	}(record)
	return c.JSON(http.StatusOK, record)
}

func HandlerShutdown(c echo.Context) error {
	record, err := App.Dao().FindFirstRecordByData("devices", "id", c.PathParam("id"))
	if err != nil {
		return apis.NewNotFoundError("The device does not exist.", err)
	}
	if err := networking.ShutdownDevice(record); err != nil {
		logger.Error.Println(err)
		return apis.NewBadRequestError(err.Error(), record)
	}
	return nil
}

type Nmaprun struct {
	Host []struct {
		Address []struct {
			Addr     string `xml:"addr,attr" binding:"required"`
			Addrtype string `xml:"addrtype,attr" binding:"required"`
			Vendor   string `xml:"vendor,attr"`
		} `xml:"address"`
	} `xml:"host"`
}

func HandlerScan(c echo.Context) error {
	if runtime.GOOS == "windows" {
		// check for admin on windows by trying to open physicaldrive
		_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
		if err != nil {
			err := errors.New("network scan requires upsnap to be run as administrator")
			return apis.NewBadRequestError(err.Error(), nil)
		}
	} else {
		// check for root on anything else
		sudoUID := os.Getenv("SUDO_UID")
		if sudoUID == "" {
			err := errors.New("network scan requires upsnap to be run as root")
			return apis.NewBadRequestError(err.Error(), nil)
		}
	}

	// check if nmap installed
	if _, err := exec.LookPath("nmap"); err != nil {
		return apis.NewBadRequestError(err.Error(), nil)
	}

	// check if scan range is valid
	allSettings, err := App.Dao().FindRecordsByExpr("settings")
	if err != nil {
		return err
	}
	settings := allSettings[0]
	scanRange := settings.GetString("scan_range")
	_, ipNet, err := net.ParseCIDR(scanRange)
	if err != nil {
		return apis.NewBadRequestError(err.Error(), nil)
	}

	// run nmap
	cmd := exec.Command("nmap", "-sn", "-oX", "-", scanRange, "--host-timeout", "500ms")
	cmdOutput, err := cmd.Output()
	if err != nil {
		return err
	}

	// unmarshal xml
	nmapOutput := Nmaprun{}
	if err := xml.Unmarshal(cmdOutput, &nmapOutput); err != nil {
		return err
	}

	type Device struct {
		Name string `json:"name"`
		IP   string `json:"ip"`
		MAC  string `json:"mac"`
	}

	// extract info from struct into data
	type Response struct {
		Netmask string   `json:"netmask"`
		Devices []Device `json:"devices"`
	}
	res := Response{}
	var nm []string
	for _, octet := range ipNet.Mask {
		nm = append(nm, strconv.Itoa(int(octet)))
	}
	res.Netmask = strings.Join(nm, ".")

	for _, host := range nmapOutput.Host {
		dev := Device{}
		for _, addr := range host.Address {
			if addr.Addrtype == "ipv4" {
				dev.IP = addr.Addr
			} else if addr.Addrtype == "mac" {
				dev.MAC = addr.Addr
			}
			if addr.Vendor != "" {
				dev.Name = addr.Vendor
			}
		}
		if dev.IP == "" || dev.MAC == "" {
			continue
		}
		if dev.Name == "" {
			dev.Name = "Unknown"
		}
		res.Devices = append(res.Devices, dev)
	}

	return c.JSON(http.StatusOK, res)
}
