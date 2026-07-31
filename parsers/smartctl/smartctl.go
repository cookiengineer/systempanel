package smartctl

import (
	"encoding/json"
	"fmt"
)

type SMARTData struct {
	Device struct {
		Name     string `json:"name"`
		InfoName string `json:"info_name"`
		Type     string `json:"type"`
		Protocol string `json:"protocol"`
	} `json:"device"`

	ModelFamily     string `json:"model_family"`
	ModelName       string `json:"model_name"`
	SerialNumber    string `json:"serial_number"`
	FirmwareVersion string `json:"firmware_version"`

	SMARTStatus struct {
		Passed bool `json:"passed"`
	} `json:"smart_status"`

	Temperature struct {
		Current int    `json:"current"`
		Unit    string `json:"unit"`
	} `json:"temperature"`

	PowerOnTime struct {
		Hours   int `json:"hours"`
		Minutes int `json:"minutes"`
	} `json:"power_on_time"`

	ATASmartAttributes struct {
		Table    []SMARTAttribute `json:"table"`
		Revision int              `json:"revision"`
	} `json:"ata_smart_attributes"`

	ATASmartData struct {
		SelfTest struct {
			Status struct {
				Value  int    `json:"value"`
				String string `json:"string"`
				Passed bool   `json:"passed"`
			} `json:"status"`
		} `json:"self_test"`
	} `json:"ata_smart_data"`

	SCTStatus struct {
		FormatVersion int    `json:"format_version"`
		SMARTStatus   string `json:"sct_status"`
		DeviceState   string `json:"device_state"`
		Temperature   struct {
			Current int `json:"current"`
			Min     int `json:"min"`
			Max     int `json:"max"`
		} `json:"temperature"`
	} `json:"ata_sct_status"`

	DeviceStatistics map[string]interface{} `json:"ata_device_statistics"`
}

type SMARTAttribute struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Value      int    `json:"value"`
	Worst      int    `json:"worst"`
	Thresh     int    `json:"thresh"`
	WhenFailed string `json:"when_failed"`
	Flags      struct {
		Value         int    `json:"value"`
		String        string `json:"string"`
		Prefailure    bool   `json:"prefailure"`
		UpdatedOnline bool   `json:"updated_online"`
		Performance   bool   `json:"performance"`
		ErrorRate     bool   `json:"error_rate"`
		EventCount    bool   `json:"event_count"`
		AutoKeep      bool   `json:"auto_keep"`
	} `json:"flags"`
	Raw struct {
		String string `json:"string"`
		Value  int64  `json:"value"`
	} `json:"raw"`
}

type SMARTAttributeInfo struct {
	ID         int
	Name       string
	Value      int
	Worst      int
	Threshold  int
	WhenFailed string
	Prefailure bool
	RawValue   int64
	RawString  string
}

type SmartHealth struct {
	Available       bool
	SmartPassed     bool
	Temperature     int
	PowerOnHours    int64
	WarningTemp     int
	WorstTemp       int
	ModelName       string
	SerialNumber    string
	FirmwareVersion string
	Attributes      []SMARTAttributeInfo
}

func ParseJSON(data []byte) (*SMARTData, error) {
	var sd SMARTData
	if err := json.Unmarshal(data, &sd); err != nil {
		return nil, fmt.Errorf("failed to parse smartctl JSON: %w", err)
	}
	return &sd, nil
}

func (sd *SMARTData) ToHealth() *SmartHealth {
	h := &SmartHealth{
		Available:       true,
		SmartPassed:     sd.SMARTStatus.Passed,
		Temperature:     sd.Temperature.Current,
		ModelName:       sd.ModelName,
		SerialNumber:    sd.SerialNumber,
		FirmwareVersion: sd.FirmwareVersion,
	}

	totalHours := int64(sd.PowerOnTime.Hours)
	if sd.PowerOnTime.Minutes > 0 {
		totalHours += int64(sd.PowerOnTime.Minutes) / 60
	}
	h.PowerOnHours = totalHours

	if sd.SCTStatus.Temperature.Current > 0 {
		h.Temperature = sd.SCTStatus.Temperature.Current
		if sd.SCTStatus.Temperature.Max > h.WorstTemp {
			h.WorstTemp = sd.SCTStatus.Temperature.Max
		}
	}

	h.WarningTemp = extractWarningTemperature(sd)

	for _, attr := range sd.ATASmartAttributes.Table {
		h.Attributes = append(h.Attributes, SMARTAttributeInfo{
			ID:         attr.ID,
			Name:       attr.Name,
			Value:      attr.Value,
			Worst:      attr.Worst,
			Threshold:  attr.Thresh,
			WhenFailed: attr.WhenFailed,
			Prefailure: attr.Flags.Prefailure,
			RawValue:   attr.Raw.Value,
			RawString:  attr.Raw.String,
		})
	}

	return h
}

func extractWarningTemperature(sd *SMARTData) int {
	for _, attr := range sd.ATASmartAttributes.Table {
		if (attr.ID == 190 || attr.ID == 194) && attr.Name == "Temperature_Celsius" {
			if attr.Raw.Value > 0 && attr.Raw.Value < 0xFFFF {
				upper := (attr.Raw.Value >> 16) & 0xFF
				if upper > 0 && upper < 200 {
					return int(upper)
				}
			}
		}
	}

	if sd.SCTStatus.Temperature.Max > 0 {
		return sd.SCTStatus.Temperature.Max + 5
	}

	return 0
}

func (h *SmartHealth) TempStatusClass() string {
	if !h.Available {
		return ""
	}
	if h.WarningTemp > 0 {
		delta := h.WarningTemp - h.Temperature
		if delta <= 0 {
			return "disk-health-critical"
		}
		if delta <= 5 {
			return "disk-health-warning"
		}
		if delta <= 15 {
			return "disk-health-ok"
		}
		return "disk-health-good"
	}
	switch {
	case h.Temperature <= 0:
		return ""
	case h.Temperature >= 60:
		return "disk-health-critical"
	case h.Temperature >= 50:
		return "disk-health-warning"
	case h.Temperature >= 40:
		return "disk-health-ok"
	default:
		return "disk-health-good"
	}
}

func (h *SmartHealth) OverallStatusClass() string {
	if !h.Available {
		return ""
	}
	if !h.SmartPassed {
		return "disk-health-critical"
	}
	for _, attr := range h.Attributes {
		if attr.WhenFailed != "" {
			return "disk-health-critical"
		}
	}
	hasWarnings := false
	for _, attr := range h.Attributes {
		if attr.ID == 5 || attr.ID == 197 || attr.ID == 198 {
			if attr.RawValue > 0 {
				hasWarnings = true
				break
			}
		}
	}
	if hasWarnings {
		return "disk-health-warning"
	}
	return "disk-health-good"
}

func (h *SmartHealth) IsCritical() bool {
	if !h.Available {
		return false
	}
	if !h.SmartPassed {
		return true
	}
	for _, attr := range h.Attributes {
		if attr.WhenFailed != "" {
			return true
		}
	}
	return false
}

func (h *SmartHealth) HasWarnings() bool {
	if !h.Available {
		return false
	}
	for _, attr := range h.Attributes {
		if attr.ID == 5 || attr.ID == 197 || attr.ID == 198 {
			if attr.RawValue > 0 {
				return true
			}
		}
	}
	if h.Temperature >= 55 && h.WarningTemp <= 0 {
		return true
	}
	if h.WarningTemp > 0 && h.Temperature >= h.WarningTemp-5 {
		return true
	}
	return false
}

var KeyAttributeIDs = map[int]string{
	1:   "Raw Read Error Rate",
	5:   "Reallocated Sectors",
	7:   "Seek Error Rate",
	9:   "Power-On Hours",
	10:  "Spin Retry Count",
	12:  "Power Cycle Count",
	177: "Wear Leveling Count",
	183: "Runtime Bad Block",
	187: "Reported Uncorrectable",
	188: "Command Timeout",
	190: "Airflow Temperature",
	194: "Temperature Celsius",
	196: "Reallocation Events",
	197: "Current Pending Sectors",
	198: "Uncorrectable Sectors",
	199: "UDMA CRC Errors",
	200: "Write Error Rate",
	202: "Data Address Mark Errors",
	230: "Media Wearout Indicator",
}

func IsKeyAttribute(id int) bool {
	_, ok := KeyAttributeIDs[id]
	return ok
}
