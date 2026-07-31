package lsblk

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

type BlockDevice struct {
	Name        string        `json:"name"`
	KName       string        `json:"kname"`
	Size        *int64        `json:"size"`
	Type        string        `json:"type"`
	Mountpoint  *string       `json:"mountpoint"`
	Label       *string       `json:"label"`
	FSType      *string       `json:"fstype"`
	UUID        *string       `json:"uuid"`
	Model       *string       `json:"model"`
	Serial      *string       `json:"serial"`
	Transport   *string       `json:"tran"`
	Rotational  *bool         `json:"rota"`
	Removable   *bool         `json:"rm"`
	MountPoints []string      `json:"mountpoints,omitempty"`
	Children    []BlockDevice `json:"children,omitempty"`
}

type LsblkOutput struct {
	BlockDevices []BlockDevice `json:"blockdevices"`
}

func ParseJSON(data []byte) (*LsblkOutput, error) {
	var raw struct {
		BlockDevices []json.RawMessage `json:"blockdevices"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("failed to parse lsblk JSON: %w", err)
	}

	var out LsblkOutput
	for _, devData := range raw.BlockDevices {
		dev, err := parseDevice(devData)
		if err != nil {
			return nil, err
		}
		out.BlockDevices = append(out.BlockDevices, *dev)
	}
	return &out, nil
}

type rawDevice struct {
	Name        string            `json:"name"`
	KName       string            `json:"kname"`
	Size        json.RawMessage   `json:"size"`
	Type        string            `json:"type"`
	Mountpoint  *string           `json:"mountpoint"`
	Label       *string           `json:"label"`
	FSType      *string           `json:"fstype"`
	UUID        *string           `json:"uuid"`
	Model       *string           `json:"model"`
	Serial      *string           `json:"serial"`
	Transport   *string           `json:"tran"`
	Rotational  *bool             `json:"rota"`
	Removable   *bool             `json:"rm"`
	MountPoints []string          `json:"mountpoints,omitempty"`
	Children    []json.RawMessage `json:"children,omitempty"`
}

func parseDevice(data json.RawMessage) (*BlockDevice, error) {
	var raw rawDevice
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	size := parseFlexInt(raw.Size)

	dev := &BlockDevice{
		Name:        raw.Name,
		KName:       raw.KName,
		Size:        size,
		Type:        raw.Type,
		Mountpoint:  raw.Mountpoint,
		Label:       raw.Label,
		FSType:      raw.FSType,
		UUID:        raw.UUID,
		Model:       raw.Model,
		Serial:      raw.Serial,
		Transport:   raw.Transport,
		Rotational:  raw.Rotational,
		Removable:   raw.Removable,
		MountPoints: raw.MountPoints,
	}

	for _, childData := range raw.Children {
		child, err := parseDevice(childData)
		if err != nil {
			return nil, err
		}
		dev.Children = append(dev.Children, *child)
	}

	return dev, nil
}

func parseFlexInt(raw json.RawMessage) *int64 {
	if len(raw) == 0 || string(raw) == "null" {
		return nil
	}

	var n int64
	if err := json.Unmarshal(raw, &n); err == nil {
		return &n
	}

	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		if _, err := fmt.Sscanf(s, "%d", &n); err == nil {
			return &n
		}
	}

	return nil
}

func Run() (*LsblkOutput, error) {
	cols := "NAME,KNAME,SIZE,TYPE,MOUNTPOINT,LABEL,FSTYPE,UUID,MODEL,SERIAL,TRAN,ROTA,RM,MOUNTPOINTS"
	data, err := exec.Command("lsblk", "--json", "--bytes", "--output", cols).Output()
	if err != nil {
		return nil, fmt.Errorf("lsblk failed: %w", err)
	}
	return ParseJSON(data)
}

func (bd *BlockDevice) DevicePath() string {
	if bd.KName != "" {
		return "/dev/" + bd.KName
	}
	if bd.Name != "" {
		return "/dev/" + bd.Name
	}
	return ""
}

func (bd *BlockDevice) BytesSize() int64 {
	if bd.Size == nil {
		return 0
	}
	return *bd.Size
}

func (bd *BlockDevice) HumanSize() string {
	if bd.Size == nil {
		return "Unknown"
	}
	bytes := *bd.Size

	switch {
	case bytes >= 1<<40:
		return fmt.Sprintf("%.2f TB", float64(bytes)/(1<<40))
	case bytes >= 1<<30:
		return fmt.Sprintf("%.2f GB", float64(bytes)/(1<<30))
	case bytes >= 1<<20:
		return fmt.Sprintf("%.2f MB", float64(bytes)/(1<<20))
	case bytes >= 1<<10:
		return fmt.Sprintf("%.2f KB", float64(bytes)/(1<<10))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}

func (bd *BlockDevice) IsRotational() bool {
	return bd.Rotational != nil && *bd.Rotational
}

func (bd *BlockDevice) IsRemovable() bool {
	return bd.Removable != nil && *bd.Removable
}

func (bd *BlockDevice) TransportType() string {
	if bd.Transport != nil && *bd.Transport != "" {
		return *bd.Transport
	}
	return ""
}

func (bd *BlockDevice) ModelName() string {
	if bd.Model != nil && *bd.Model != "" {
		return *bd.Model
	}
	return ""
}

func (bd *BlockDevice) SerialNumber() string {
	if bd.Serial != nil && *bd.Serial != "" {
		return *bd.Serial
	}
	return ""
}

func (bd *BlockDevice) FSTypeName() string {
	if bd.FSType != nil && *bd.FSType != "" {
		return *bd.FSType
	}
	return ""
}

func (bd *BlockDevice) MountPointPath() string {
	if len(bd.MountPoints) > 0 {
		return bd.MountPoints[0]
	}
	if bd.Mountpoint != nil && *bd.Mountpoint != "" {
		return *bd.Mountpoint
	}
	return ""
}

func (bd *BlockDevice) LabelName() string {
	if bd.Label != nil && *bd.Label != "" {
		return *bd.Label
	}
	return ""
}

func (bd *BlockDevice) UUIDString() string {
	if bd.UUID != nil && *bd.UUID != "" {
		return *bd.UUID
	}
	return ""
}
