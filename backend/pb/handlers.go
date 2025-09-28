package pb

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/router"
	"github.com/robfig/cron/v3"
	"github.com/seriousm4x/upsnap/logger"
	"github.com/seriousm4x/upsnap/networking"
)

func HandlerWake(e *core.RequestEvent) error {
	record, err := e.App.FindFirstRecordByData("devices", "id", e.Request.PathValue("id"))
	if err != nil {
		return apis.NewNotFoundError("The device does not exist.", err)
	}

	record.Set("status", "pending")
	if err := e.App.Save(record); err != nil {
		logger.Error.Println("Failed to save record:", err)
	}

	if err := asyncCall(e, func() *router.ApiError {
		if err := networking.WakeDevice(record); err != nil {
			logger.Error.Println(err)
			record.Set("status", "offline")
			if err := e.App.Save(record); err != nil {
				logger.Error.Println("Failed to save record:", err)
			}
			return apis.NewBadRequestError(err.Error(), nil)
		}

		record.Set("status", "online")
		if err := e.App.Save(record); err != nil {
			logger.Error.Println("Failed to save record:", err)
		}

		return nil
	}); err != nil {
		return err
	}

	return e.JSON(http.StatusOK, record)
}

func HandlerSleep(e *core.RequestEvent) error {
	record, err := e.App.FindFirstRecordByData("devices", "id", e.Request.PathValue("id"))
	if err != nil {
		return apis.NewNotFoundError("The device does not exist.", err)
	}

	record.Set("status", "pending")
	if err := e.App.Save(record); err != nil {
		logger.Error.Println("Failed to save record:", err)
	}

	if err := asyncCall(e, func() *router.ApiError {
		resp, err := networking.SleepDevice(record)
		if err != nil {
			logger.Error.Println(err)
			record.Set("status", "online")
			if err := e.App.Save(record); err != nil {
				logger.Error.Println("Failed to save record:", err)
			}
			return apis.NewBadRequestError(resp.Message, nil)
		}

		record.Set("status", "offline")
		if err := e.App.Save(record); err != nil {
			logger.Error.Println("Failed to save record:", err)
		}

		return nil
	}); err != nil {
		return err
	}

	return e.JSON(http.StatusOK, nil)
}

func HandlerReboot(e *core.RequestEvent) error {
	record, err := e.App.FindFirstRecordByData("devices", "id", e.Request.PathValue("id"))
	if err != nil {
		return apis.NewNotFoundError("The device does not exist.", err)
	}

	record.Set("status", "pending")
	if err := e.App.Save(record); err != nil {
		logger.Error.Println("Failed to save record:", err)
	}

	if err := asyncCall(e, func() *router.ApiError {
		if err := networking.ShutdownDevice(record); err != nil {
			logger.Error.Println(err)
			record.Set("status", "online")
			if err := e.App.Save(record); err != nil {
				logger.Error.Println("Failed to save record:", err)
			}
			return apis.NewBadRequestError(err.Error(), nil)
		}

		// if this code is reached, the device is shutting down and not responding to pings anymore.
		// this does not mean it is ready to boot again, it might still be in the process of shutting down.
		// so we wait a little to make sure the device has shut down completely and is ready to receive wake requests.
		time.Sleep(15 * time.Second)

		if err := networking.WakeDevice(record); err != nil {
			logger.Error.Println(err)
			record.Set("status", "offline")
			if err := e.App.Save(record); err != nil {
				logger.Error.Println("Failed to save record:", err)
			}
			return apis.NewBadRequestError(err.Error(), nil)
		}

		record.Set("status", "online")
		if err := e.App.Save(record); err != nil {
			logger.Error.Println("Failed to save record:", err)
		}

		return nil
	}); err != nil {
		return err
	}

	return e.JSON(http.StatusOK, record)
}

func HandlerShutdown(e *core.RequestEvent) error {
	record, err := e.App.FindFirstRecordByData("devices", "id", e.Request.PathValue("id"))
	if err != nil {
		return apis.NewNotFoundError("The device does not exist.", err)
	}

	record.Set("status", "pending")
	if err := e.App.Save(record); err != nil {
		logger.Error.Println("Failed to save record:", err)
	}

	if err := asyncCall(e, func() *router.ApiError {
		if err := networking.ShutdownDevice(record); err != nil {
			logger.Error.Println(err)
			record.Set("status", "online")
			if err := e.App.Save(record); err != nil {
				logger.Error.Println("Failed to save record:", err)
			}
			return apis.NewBadRequestError(err.Error(), nil)
		}

		record.Set("status", "offline")
		if err := e.App.Save(record); err != nil {
			logger.Error.Println("Failed to save record:", err)
		}

		return nil
	}); err != nil {
		return err
	}

	return e.JSON(http.StatusOK, record)
}

func HandlerWakeGroup(e *core.RequestEvent) error {
	records, err := e.App.FindRecordsByFilter("devices", "groups.id = {:grpId} && status = 'offline'", "", 0, 0,
		dbx.Params{"grpId": e.Request.PathValue("id")},
	)
	if err != nil {
		return apis.NewNotFoundError("No devices in group", err)
	}

	for _, record := range records {
		go func() {
			record.Set("status", "pending")
			if err := e.App.Save(record); err != nil {
				logger.Error.Println("Failed to save record:", err)
			}

			if err := networking.WakeDevice(record); err != nil {
				logger.Error.Println(err)
				record.Set("status", "offline")
				if err := e.App.Save(record); err != nil {
					logger.Error.Println("Failed to save record:", err)
				}
				return
			}

			record.Set("status", "online")
			if err := e.App.Save(record); err != nil {
				logger.Error.Println("Failed to save record:", err)
			}
		}()
	}

	return e.JSON(http.StatusOK, records)
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
	allPrivateSettings, err := e.App.FindAllRecords("settings_private")
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
	err = e.App.Save(record)
	if err != nil {
		return e.BadRequestError("Failed to create initial superuser", err)
	}

	return apis.RecordAuthResponse(e, record, "", nil)
}

type ValidateCronRequest struct {
	Cron string `json:"cron"`
}

func HandlerValidateCron(e *core.RequestEvent) error {
	var body ValidateCronRequest

	if err := e.BindBody(&body); err != nil {
		logger.Error.Println(err)
		return e.BadRequestError("invalid request", err)
	}

	p := cron.NewParser(cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	if _, err := p.Parse(body.Cron); err != nil {
		return e.BadRequestError("failed to parse cron expression", err)
	}

	return e.JSON(200, "valid")
}

func asyncCall(e *core.RequestEvent, fn func() *router.ApiError) *router.ApiError {
	isAsync, err := strconv.ParseBool(e.Request.URL.Query().Get("async"))
	if err != nil {
		isAsync = false
	}

	if !isAsync {
		return fn()
	}

	go func() {
		_ = fn()
	}()

	return nil
}

func HandlerWebsiteManifest(e *core.RequestEvent) error {
	scheme := "http"
	if e.Request.TLS != nil || e.Request.Header.Get("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}
	baseUrl := fmt.Sprintf("%s://%s", scheme, e.Request.Host)

	manifest := map[string]any{
		"name":             "UpSnap",
		"short_name":       "UpSnap",
		"description":      "A simple wake on lan web app written with SvelteKit, Go and PocketBase.",
		"theme_color":      "#55BCD9",
		"background_color": "#55BCD9",
		"display":          "standalone",
		"scope":            baseUrl,
		"start_url":        baseUrl,
		"icons": []map[string]string{
			{
				"src":     "/icon_192.png",
				"sizes":   "192x192",
				"type":    "image/png",
				"purpose": "any",
			},
			{
				"src":     "/icon_512.png",
				"sizes":   "512x512",
				"type":    "image/png",
				"purpose": "any",
			},
			{
				"src":     "/maskable_192.png",
				"sizes":   "192x192",
				"type":    "image/png",
				"purpose": "maskable",
			},
			{
				"src":     "/maskable_512.png",
				"sizes":   "512x512",
				"type":    "image/png",
				"purpose": "maskable",
			},
			{
				"src":   "/gopher.svg",
				"sizes": "any",
				"type":  "image/svg+xml",
			},
		},
	}

	publicSettings, err := e.App.FindFirstRecordByFilter("settings_public", "")
	if err != nil {
		return writeManifest(e, manifest)
	}

	// override title if set
	if title := publicSettings.GetString("website_title"); title != "" {
		manifest["name"] = title
		manifest["short_name"] = title
	}

	// override icons if set
	if favicon := publicSettings.GetString("favicon"); favicon != "" {

		manifest["icons"] = []map[string]string{
			{
				"src":     fmt.Sprintf("%s/api/files/settings_public/%s/%s", baseUrl, publicSettings.Id, favicon),
				"sizes":   "any",
				"purpose": "any",
			},
		}
	}

	return writeManifest(e, manifest)
}

func writeManifest(e *core.RequestEvent, manifest map[string]any) error {
	e.Response.Header().Set("Content-Type", "application/manifest+json")
	e.Response.Header().Set("Cache-Control", "no-cache, no-store")
	enc := json.NewEncoder(e.Response)
	return enc.Encode(manifest)
}
