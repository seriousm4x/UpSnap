package networking

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/pocketbase/pocketbase/models"
	"github.com/seriousm4x/upsnap/logger"
)

type SolResponse struct {
	Message string `json:"message"`
}

func SleepDevice(device *models.Record) (SolResponse, error) {
	logger.Info.Println("Sleep triggered for", device.GetString("name"))

	var solResp SolResponse
	var url string

	if device.GetBool("sol_auth") {
		url = fmt.Sprintf("http://%s:%s@%s:%d/sleep?format=JSON",
			device.GetString("sol_user"), device.GetString("sol_password"), device.GetString("ip"), device.GetInt("sol_port"))
	} else {
		url = fmt.Sprintf("http://%s:%d/sleep?format=JSON", device.GetString("ip"), device.GetInt("sol_port"))
	}

	resp, err := http.Get(url)
	if err != nil {
		return solResp, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return solResp, err
		}

		var solResp SolResponse
		if err := json.Unmarshal(body, &solResp); err != nil {
			return solResp, err
		}

		return solResp, errors.New("status code was not 200")
	}

	return solResp, nil
}
