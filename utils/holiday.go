package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/mnizarzr/dot-test/config"
)

// HolidayAPIResponse represents the response from Holiday API
type HolidayAPIResponse struct {
	Status   int `json:"status"`
	Requests struct {
		Used      int    `json:"used"`
		Available int    `json:"available"`
		Resets    string `json:"resets"`
	} `json:"requests"`
	Holidays []struct {
		Name     string `json:"name"`
		Date     string `json:"date"`
		Observed string `json:"observed,omitempty"`
		Public   bool   `json:"public"`
	} `json:"holidays"`
	Error   string `json:"error,omitempty"`
	Warning string `json:"warning,omitempty"`
}

func CheckHoliday(date time.Time) bool {
	apiKey := config.Configs.HolidayApiKey
	if apiKey == "" {
		return false
	}

	dateStr := date.Format("2006-01-02")
	url := fmt.Sprintf("https://holidayapi.com/v1/holidays?key=%s&country=US&year=%d&month=%d&day=%d&public=true",
		apiKey, date.Year(), date.Month(), date.Day())

	resp, err := http.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false
	}

	var holidayResp HolidayAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&holidayResp); err != nil {
		return false
	}

	if holidayResp.Error != "" {
		return false
	}

	for _, holiday := range holidayResp.Holidays {
		if holiday.Date == dateStr && holiday.Public {
			return true
		}
		if holiday.Observed != "" && holiday.Observed == dateStr && holiday.Public {
			return true
		}
	}

	return false
}
