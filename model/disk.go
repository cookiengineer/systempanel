package model

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/cookiengineer/systempanel/parsers/lsblk"
	"github.com/cookiengineer/systempanel/parsers/smartctl"
)

type DiskInfo struct {
	Name         string
	DevicePath   string
	SizeBytes    int64
	Model        string
	Serial       string
	Transport    string
	IsRotational bool
	IsRemovable  bool
	IsUSB        bool
	Partitions   []PartitionInfo
	Health       *smartctl.SmartHealth
}

type PartitionInfo struct {
	Name        string
	DevicePath  string
	SizeBytes   int64
	FSType      string
	MountPoint  string
	Label       string
	UUID        string
	IsEncrypted bool
}

type DiskModel struct {
	observers []Observer
}

func (m *DiskModel) Refresh(ctx context.Context) error { return nil }
func (m *DiskModel) Observe(fn Observer) func()        { return func() {} }

func (m *DiskModel) ListDisks() ([]DiskInfo, error) {
	out, err := lsblk.Run()
	if err != nil {
		return nil, err
	}

	var disks []DiskInfo
	for _, dev := range out.BlockDevices {
		if dev.Type != "disk" {
			continue
		}
		disk := DiskInfo{
			Name:         dev.Name,
			DevicePath:   dev.DevicePath(),
			SizeBytes:    dev.BytesSize(),
			Model:        dev.ModelName(),
			Serial:       dev.SerialNumber(),
			Transport:    dev.TransportType(),
			IsRotational: dev.IsRotational(),
			IsRemovable:  dev.IsRemovable(),
			IsUSB:        dev.TransportType() == "usb",
		}

		for _, child := range dev.Children {
			disk.Partitions = append(disk.Partitions, PartitionInfo{
				Name:        child.Name,
				DevicePath:  child.DevicePath(),
				SizeBytes:   child.BytesSize(),
				FSType:      child.FSTypeName(),
				MountPoint:  child.MountPointPath(),
				Label:       child.LabelName(),
				UUID:        child.UUIDString(),
				IsEncrypted: child.FSTypeName() == "crypto_LUKS",
			})
		}

		if len(disk.Partitions) == 0 && dev.FSTypeName() != "" {
			disk.Partitions = append(disk.Partitions, PartitionInfo{
				Name:        dev.Name,
				DevicePath:  dev.DevicePath(),
				SizeBytes:   dev.BytesSize(),
				FSType:      dev.FSTypeName(),
				MountPoint:  dev.MountPointPath(),
				Label:       dev.LabelName(),
				UUID:        dev.UUIDString(),
				IsEncrypted: dev.FSTypeName() == "crypto_LUKS",
			})
		}

		disks = append(disks, disk)
	}

	return disks, nil
}

func (m *DiskModel) FetchSmartHealth(devicePath string, sudoPassword string) (*smartctl.SmartHealth, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	args := []string{"--xall", "--json", devicePath}
	var cmd *exec.Cmd

	if sudoPassword != "" {
		cmd = exec.CommandContext(ctx, "sudo", append([]string{"-S", "-k", "smartctl"}, args...)...)
		cmd.Stdin = strings.NewReader(sudoPassword + "\n")
	} else {
		cmd = exec.CommandContext(ctx, "smartctl", args...)
	}

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("smartctl failed for %s: %w (stderr: %s)", devicePath, err, stderr.String())
	}

	sd, err := smartctl.ParseJSON(output)
	if err != nil {
		return nil, err
	}

	return sd.ToHealth(), nil
}

func (m *DiskModel) MountPartition(devicePath string) error {
	cmd := exec.Command("udisksctl", "mount", "-b", devicePath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("mount failed: %s", strings.TrimSpace(string(output)))
	}
	return nil
}

func (m *DiskModel) UnmountPartition(devicePath string) error {
	cmd := exec.Command("udisksctl", "unmount", "-b", devicePath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("unmount failed: %s", strings.TrimSpace(string(output)))
	}
	return nil
}

func (m *DiskModel) UnlockLUKS(devicePath string) error {
	cmd := exec.Command("udisksctl", "unlock", "-b", devicePath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("unlock failed: %s", strings.TrimSpace(string(output)))
	}
	return nil
}

func (m *DiskModel) LockLUKS(devicePath string) error {
	cmd := exec.Command("udisksctl", "lock", "-b", devicePath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("lock failed: %s", strings.TrimSpace(string(output)))
	}
	return nil
}
