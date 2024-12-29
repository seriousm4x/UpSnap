package pb

import (
	"encoding/xml"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/seriousm4x/upsnap/logger"
	"github.com/seriousm4x/upsnap/networking"
)

func HandlerWake(e *core.RequestEvent) error {
	record, err := App.FindFirstRecordByData("devices", "id", e.Request.PathValue("id"))
	if err != nil {
		return apis.NewNotFoundError("The device does not exist.", err)
	}
	record.Set("status", "pending")
	if err := App.Save(record); err != nil {
		logger.Error.Println("Failed to save record:", err)
	}
	go func(r *core.Record) {
		if err := networking.WakeDevice(r); err != nil {
			logger.Error.Println(err)
			r.Set("status", "offline")
		} else {
			r.Set("status", "online")
		}
		if err := App.Save(r); err != nil {
			logger.Error.Println("Failed to save record:", err)
		}
	}(record)
	return e.JSON(http.StatusOK, record)
}

func HandlerSleep(e *core.RequestEvent) error {
	record, err := App.FindFirstRecordByData("devices", "id", e.Request.PathValue("id"))
	if err != nil {
		return apis.NewNotFoundError("The device does not exist.", err)
	}
	record.Set("status", "pending")
	if err := App.Save(record); err != nil {
		logger.Error.Println("Failed to save record:", err)
	}
	resp, err := networking.SleepDevice(record)
	if err != nil {
		logger.Error.Println(err)
		record.Set("status", "online")
		if err := App.Save(record); err != nil {
			logger.Error.Println("Failed to save record:", err)
		}
		return apis.NewBadRequestError(resp.Message, nil)
	}
	record.Set("status", "offline")
	if err := App.Save(record); err != nil {
		logger.Error.Println("Failed to save record:", err)
	}
	return e.JSON(http.StatusOK, nil)
}

func HandlerReboot(e *core.RequestEvent) error {
	record, err := App.FindFirstRecordByData("devices", "id", e.Request.PathValue("id"))
	if err != nil {
		return apis.NewNotFoundError("The device does not exist.", err)
	}
	record.Set("status", "pending")
	if err := App.Save(record); err != nil {
		logger.Error.Println("Failed to save record:", err)
	}
	go func(r *core.Record) {
		if err := networking.ShutdownDevice(r); err != nil {
			logger.Error.Println(err)
			r.Set("status", "online")
		} else {
			time.Sleep(15 * time.Second) // some devices might not respond to ping but are still shutting down
			if err := networking.WakeDevice(r); err != nil {
				logger.Error.Println(err)
				r.Set("status", "offline")
			} else {
				r.Set("status", "online")
			}
		}
		if err := App.Save(r); err != nil {
			logger.Error.Println("Failed to save record:", err)
		}
	}(record)
	return e.JSON(http.StatusOK, record)
}

func HandlerShutdown(e *core.RequestEvent) error {
	record, err := App.FindFirstRecordByData("devices", "id", e.Request.PathValue("id"))
	if err != nil {
		return apis.NewNotFoundError("The device does not exist.", err)
	}
	record.Set("status", "pending")
	if err := App.Save(record); err != nil {
		logger.Error.Println("Failed to save record:", err)
	}
	go func(r *core.Record) {
		if err := networking.ShutdownDevice(r); err != nil {
			logger.Error.Println(err)
			r.Set("status", "online")
		} else {
			r.Set("status", "offline")
		}
		if err := App.Save(r); err != nil {
			logger.Error.Println("Failed to save record:", err)
		}
	}(record)
	return e.JSON(http.StatusOK, record)
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

func HandlerScan(e *core.RequestEvent) error {
	// check if nmap installed
	nmap, err := exec.LookPath("nmap")
	if err != nil {
		return apis.NewBadRequestError(err.Error(), nil)
	}

	// check if scan range is valid
	allPrivateSettings, err := App.FindAllRecords("settings_private")
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
	timeout := os.Getenv("UPSNAP_SCAN_TIMEOUT")
	if timeout == "" {
		timeout = "500ms"
	}
	cmd := exec.Command(nmap, "-sn", "-oX", "-", scanRange, "--host-timeout", timeout)
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

	return e.JSON(http.StatusOK, res)
}

func HandlerInitSuperuser(e *core.RequestEvent) error {
	superusersCollection, err := e.App.FindCollectionByNameOrId(core.CollectionNameSuperusers)
	if err != nil {
		return e.NotFoundError("Failed to retrieve superusers collection", err)
	}

	totalSuperusers, err := e.App.CountRecords(superusersCollection)
	if err != nil {
		return e.InternalServerError("Failed to retrieve superusers count", err)
	}

	if totalSuperusers > 0 {
		return e.BadRequestError("An initial superuser already exists", nil)
	}

	data := struct {
		Email           string `json:"email" form:"email"`
		Password        string `json:"password" form:"password"`
		PasswordConfirm string `json:"password_confirm" form:"password-confirm"`
	}{}
	err = e.BindBody(&data)
	if err != nil {
		return e.BadRequestError("Failed to read request data", err)
	}

	if data.Password != data.PasswordConfirm {
		return e.BadRequestError("Password don't match", err)
	}

	record := core.NewRecord(superusersCollection)
	record.SetEmail(data.Email)
	record.SetPassword(data.Password)
	err = App.Save(record)
	if err != nil {
		return e.BadRequestError("Failed to create initial superuser", err)
	}

	return apis.RecordAuthResponse(e, record, "", nil)
}
