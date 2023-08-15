package pb

import (
	"encoding/xml"
	"net"
	"net/http"
	"os/exec"
	"strconv"
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/models"
	"github.com/seriousm4x/upsnap/logger"
	"github.com/seriousm4x/upsnap/networking"
)

func HandlerWake(c echo.Context) error {
	record, err := App.Dao().FindFirstRecordByData("devices", "id", c.PathParam("id"))
	if err != nil {
		return apis.NewNotFoundError("The device does not exist.", err)
	}
	record.Set("status", "pending")
	if err := App.Dao().SaveRecord(record); err != nil {
		logger.Error.Println("Failed to save record:", err)
	}
	go func(*models.Record) {
		if err := networking.WakeDevice(record); err != nil {
			logger.Error.Println(err)
			record.Set("status", "offline")
		} else {
			record.Set("status", "online")
		}
		if err := App.Dao().SaveRecord(record); err != nil {
			logger.Error.Println("Failed to save record:", err)
		}
	}(record)
	return c.JSON(http.StatusOK, record)
}

func HandlerSleep(c echo.Context) error {
	record, err := App.Dao().FindFirstRecordByData("devices", "id", c.PathParam("id"))
	if err != nil {
		return apis.NewNotFoundError("The device does not exist.", err)
	}
	record.Set("status", "pending")
	if err := App.Dao().SaveRecord(record); err != nil {
		logger.Error.Println("Failed to save record:", err)
	}
	resp, err := networking.SleepDevice(record)
	if err != nil {
		logger.Error.Println(err)
		record.Set("status", "online")
		if err := App.Dao().SaveRecord(record); err != nil {
			logger.Error.Println("Failed to save record:", err)
		}
		return apis.NewBadRequestError(resp.Message, nil)
	}
	record.Set("status", "offline")
	if err := App.Dao().SaveRecord(record); err != nil {
		logger.Error.Println("Failed to save record:", err)
	}
	return c.JSON(http.StatusOK, nil)
}

func HandlerShutdown(c echo.Context) error {
	record, err := App.Dao().FindFirstRecordByData("devices", "id", c.PathParam("id"))
	if err != nil {
		return apis.NewNotFoundError("The device does not exist.", err)
	}
	record.Set("status", "pending")
	if err := App.Dao().SaveRecord(record); err != nil {
		logger.Error.Println("Failed to save record:", err)
	}
	go func(*models.Record) {
		if err := networking.ShutdownDevice(record); err != nil {
			logger.Error.Println(err)
			record.Set("status", "online")
		} else {
			record.Set("status", "offline")
		}
		if err := App.Dao().SaveRecord(record); err != nil {
			logger.Error.Println("Failed to save record:", err)
		}
	}(record)
	return c.JSON(http.StatusOK, record)
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
	// check if nmap installed
	if _, err := exec.LookPath("nmap"); err != nil {
		return apis.NewBadRequestError(err.Error(), nil)
	}

	// check if scan range is valid
	allPrivateSettings, err := App.Dao().FindRecordsByExpr("settings_private")
	if err != nil {
		return err
	}
	settingsPrivate := allPrivateSettings[0]
	scanRange := settingsPrivate.GetString("scan_range")
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
		Name      string `json:"name"`
		IP        string `json:"ip"`
		MAC       string `json:"mac"`
		MACVendor string `json:"mac_vendor"`
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
				dev.MACVendor = addr.Vendor
			} else {
				dev.MACVendor = "Unknown"
			}
		}

		if dev.IP == "" || dev.MAC == "" {
			continue
		}

		names, err := net.LookupAddr(dev.IP)
		if err != nil || len(names) == 0 {
			dev.Name = dev.MACVendor
		} else {
			dev.Name = strings.TrimSuffix(names[0], ".")
		}

		if dev.Name == "" && dev.MACVendor == "" {
			continue
		}

		res.Devices = append(res.Devices, dev)
	}

	return c.JSON(http.StatusOK, res)
}
