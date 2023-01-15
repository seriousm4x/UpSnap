package pb

import (
	"encoding/xml"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
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
		isOnline := networking.WakeDevice(record)
		if isOnline {
			record.Set("status", "online")
		} else {
			record.Set("status", "offline")
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
	shutdown_cmd := record.GetString("shutdown_cmd")
	if shutdown_cmd != "" {
		cmd := exec.Command(shutdown_cmd)
		if err := cmd.Run(); err != nil {
			logger.Error.Println(err)
			return apis.NewBadRequestError(err.Error(), record)
		}
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
	scanRange := os.Getenv("UPSNAP_SCAN_RANGE")
	_, _, err := net.ParseCIDR(scanRange)
	if err != nil {
		return apis.NewBadRequestError(err.Error(), nil)
	}

	cmd := []string{}
	if runtime.GOOS != "windows" {
		cmd = []string{"sudo"}
	}
	cmd = append(cmd, "nmap", "-sn", "-oX", "-", scanRange, "--host-timeout", "500ms")

	cmdOutput, err := exec.Command(strings.Join(cmd, " ")).Output()
	if err != nil {
		return err
	}

	nmapOutput := Nmaprun{}
	if err := xml.Unmarshal(cmdOutput, &nmapOutput); err != nil {
		return err
	}

	type Device struct {
		Name string `json:"name"`
		IP   string `json:"ip"`
		MAC  string `json:"mac"`
	}
	data := []Device{}

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
		data = append(data, dev)
	}

	return c.JSON(http.StatusOK, data)
}
