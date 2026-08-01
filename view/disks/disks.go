package disks

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/cookiengineer/systempanel/bindings/gtk4"
	"github.com/cookiengineer/systempanel/detect"
	"github.com/cookiengineer/systempanel/model"
	"github.com/cookiengineer/systempanel/parsers/smartctl"
	"github.com/cookiengineer/systempanel/view"
	"github.com/cookiengineer/systempanel/widget"
)

var Descriptor = view.ViewDescriptor{
	Name:     "disks",
	Title:    "Disks",
	IconName: "drive-harddisk-symbolic",
	DetectFn: func() bool {
		return detect.HasProgram("lsblk") && detect.HasProgram("udisksctl")
	},
	Factory: func() view.View { return NewDiskView() },
}

type DiskView struct {
	box         *gtk4.Box
	model       *model.DiskModel
	listBox     *gtk4.ListBox
	spinner     *gtk4.Spinner
	rows        []*gtk4.ListBoxRow
	disks       []model.DiskInfo
	parentWin   *gtk4.Window
	hasSmartctl bool
}

func NewDiskView() *DiskView {
	dv := &DiskView{
		box:         gtk4.BoxNew(gtk4.OrientationVertical, 12),
		model:       &model.DiskModel{},
		hasSmartctl: detect.HasProgram("smartctl"),
	}
	dv.box.SetMarginStart(24)
	dv.box.SetMarginEnd(24)
	dv.box.SetMarginTop(24)
	dv.box.SetMarginBottom(24)

	header := gtk4.LabelNew("Disks")
	header.AddCSSClass("header-label")
	dv.box.Append(&header.Widget)

	desc := gtk4.LabelNew("View connected storage devices, mount partitions, unlock encrypted volumes, and monitor drive health with S.M.A.R.T. data.")
	desc.SetWrap(true)
	desc.SetMarginBottom(12)
	descWidget := desc.Widget
	dv.box.Append(&descWidget)

	dv.listBox = gtk4.ListBoxNew()
	dv.listBox.SetSelectionMode(gtk4.SelectionNone)

	scrollW := gtk4.ScrolledWindowNew()
	scrollW.SetPolicy(gtk4.PolicyNever, gtk4.PolicyAutomatic)
	scrollW.SetHExpand(true)
	scrollW.SetVExpand(true)
	scrollW.SetChild(&dv.listBox.Widget)
	dv.box.Append(&scrollW.Widget)

	btnBox := gtk4.BoxNew(gtk4.OrientationHorizontal, 8)
	btnBox.SetMarginTop(4)

	refreshBtn := gtk4.ButtonNewWithLabel("Refresh")
	refreshBtn.SetTooltip("Refresh device list and S.M.A.R.T. health data")
	refreshBtn.OnClicked(func() { dv.refresh() })
	btnBox.Append(&refreshBtn.Widget)

	dv.spinner = gtk4.SpinnerNew()
	dv.spinner.SetSizeRequest(20, 20)
	btnBox.Append(&dv.spinner.Widget)

	dv.box.Append(&btnBox.Widget)

	dv.refresh()

	return dv
}

func (dv *DiskView) SetParentWindow(win *gtk4.Window) {
	dv.parentWin = win
}

func (dv *DiskView) refresh() {
	disks, err := dv.model.ListDisks()
	if err != nil || len(disks) == 0 {
		dv.disks = dv.disks[:0]
		dv.clearRows()
		label := gtk4.LabelNew("No storage devices detected")
		label.SetSensitive(false)
		label.SetHAlign(gtk4.AlignCenter)
		row := gtk4.ListBoxRowNew()
		row.SetChild(&label.Widget)
		dv.listBox.Append(row)
		dv.rows = append(dv.rows, row)
		return
	}

	dv.disks = disks
	dv.rebuildRows()

	if !dv.hasSmartctl || dv.parentWin == nil {
		return
	}

	widget.PromptForSudo(dv.parentWin,
		"S.M.A.R.T. health data requires root privileges to query.",
		func(password string) {
			if password == "" {
				return
			}
			gtk4.IdleAdd(func() { dv.spinner.Start() })

			go func() {
				for i := range dv.disks {
					health, err := dv.model.FetchSmartHealth(dv.disks[i].DevicePath, password)
					if err != nil {
						health = &smartctl.SmartHealth{Available: false}
					}
					dv.disks[i].Health = health
				}

				gtk4.IdleAdd(func() {
					dv.spinner.Stop()
					dv.rebuildRows()
				})
			}()
		},
	)
}

func (dv *DiskView) clearRows() {
	for _, r := range dv.rows {
		dv.listBox.Remove(r)
	}
	dv.rows = dv.rows[:0]
}

func (dv *DiskView) rebuildRows() {
	dv.clearRows()
	for i := range dv.disks {
		row := dv.createDiskRow(&dv.disks[i])
		dv.listBox.Append(row)
		dv.rows = append(dv.rows, row)
	}
}

func (dv *DiskView) rescanOnly() {
	oldHealth := make(map[string]*smartctl.SmartHealth)
	for _, d := range dv.disks {
		if d.Health != nil && d.Health.Available {
			oldHealth[d.DevicePath] = d.Health
		}
	}

	disks, err := dv.model.ListDisks()
	if err != nil || len(disks) == 0 {
		dv.disks = dv.disks[:0]
		dv.clearRows()
		label := gtk4.LabelNew("No storage devices detected")
		label.SetSensitive(false)
		label.SetHAlign(gtk4.AlignCenter)
		row := gtk4.ListBoxRowNew()
		row.SetChild(&label.Widget)
		dv.listBox.Append(row)
		dv.rows = append(dv.rows, row)
		return
	}

	for i := range disks {
		if h, ok := oldHealth[disks[i].DevicePath]; ok {
			disks[i].Health = h
		}
	}

	dv.disks = disks
	dv.rebuildRows()
}

func (dv *DiskView) createDiskRow(disk *model.DiskInfo) *gtk4.ListBoxRow {
	row := gtk4.ListBoxRowNew()
	row.AddCSSClass("device-row")

	vbox := gtk4.BoxNew(gtk4.OrientationVertical, 6)
	vbox.SetMarginStart(8)
	vbox.SetMarginEnd(8)
	vbox.SetMarginTop(6)
	vbox.SetMarginBottom(6)

	topBox := gtk4.BoxNew(gtk4.OrientationHorizontal, 10)

	icon := dv.diskIcon(disk)
	icon.SetPixelSize(28)
	topBox.Append(&icon.Widget)

	infoBox := gtk4.BoxNew(gtk4.OrientationVertical, 2)

	nameLine := gtk4.BoxNew(gtk4.OrientationHorizontal, 6)
	nameText := disk.Name
	if disk.Model != "" {
		nameText += " - " + disk.Model
	}

	nameLabel := gtk4.LabelNew(nameText)
	nameLabel.SetHAlign(gtk4.AlignStart)
	nameLabel.SetMarkup(fmt.Sprintf("<b>%s</b>", strings.ReplaceAll(nameText, "&", "&amp;")))
	nameLine.Append(&nameLabel.Widget)
	infoBox.Append(&nameLine.Widget)

	detailsLine := gtk4.BoxNew(gtk4.OrientationHorizontal, 8)

	typeLabel := gtk4.LabelNew(dv.diskTypeLabel(disk))
	typeLabel.SetSensitive(false)
	detailsLine.Append(&typeLabel.Widget)

	sizeLabel := gtk4.LabelNew(formatBytes(disk.SizeBytes))
	sizeLabel.SetSensitive(false)
	detailsLine.Append(&sizeLabel.Widget)

	if disk.Serial != "" {
		serialLabel := gtk4.LabelNew("S/N: " + disk.Serial)
		serialLabel.SetSensitive(false)
		detailsLine.Append(&serialLabel.Widget)
	}

	infoBox.Append(&detailsLine.Widget)

	infoBoxWidget := infoBox.Widget
	infoBoxWidget.SetHExpand(true)
	topBox.Append(&infoBoxWidget)

	statusBox := dv.createHealthStatusBox(disk)
	if statusBox != nil {
		topBox.Append(statusBox)
	}

	vbox.Append(&topBox.Widget)

	if disk.Health != nil && disk.Health.Available {
		healthBox := dv.createHealthDetailsBox(disk.Health)
		if healthBox != nil {
			vbox.Append(healthBox)
		}
	}

	if len(disk.Partitions) > 0 {
		partSection := dv.createPartitionsSection(disk)
		vbox.Append(partSection)
	}

	vboxWidget := vbox.Widget
	row.SetChild(&vboxWidget)

	return row
}

func (dv *DiskView) diskIcon(disk *model.DiskInfo) *gtk4.Image {
	if disk.IsUSB || disk.IsRemovable {
		return gtk4.ImageNewFromIconName("drive-removable-media-symbolic")
	}
	if disk.IsRotational {
		return gtk4.ImageNewFromIconName("drive-harddisk-symbolic")
	}
	return gtk4.ImageNewFromIconName("drive-harddisk-solidstate-symbolic")
}

func (dv *DiskView) diskTypeLabel(disk *model.DiskInfo) string {
	if disk.IsUSB {
		return "USB Drive"
	}
	if disk.Transport == "nvme" {
		return "NVMe SSD"
	}
	if disk.IsRotational {
		if disk.Transport != "" {
			return strings.ToUpper(disk.Transport) + " HDD"
		}
		return "HDD"
	}
	if disk.Transport != "" {
		return strings.ToUpper(disk.Transport) + " SSD"
	}
	return "SSD"
}

func (dv *DiskView) createHealthStatusBox(disk *model.DiskInfo) *gtk4.Widget {
	hbox := gtk4.BoxNew(gtk4.OrientationHorizontal, 6)

	if disk.Health == nil || !disk.Health.Available {
		statusLabel := gtk4.LabelNew("S.M.A.R.T.\nunavailable")
		statusLabel.SetSensitive(false)
		statusLabel.SetSizeRequest(100, -1)
		statusLabel.SetHAlign(gtk4.AlignEnd)
		hbox.Append(&statusLabel.Widget)
		w := hbox.Widget
		return &w
	}

	vbox := gtk4.BoxNew(gtk4.OrientationVertical, 2)
	vbox.SetHAlign(gtk4.AlignEnd)

	statusClass := disk.Health.OverallStatusClass()
	var statusText string
	switch {
	case disk.Health.IsCritical():
		statusText = "FAILING"
	case disk.Health.HasWarnings():
		statusText = "Warning"
	default:
		statusText = "Healthy"
	}

	statusLabel := gtk4.LabelNew(statusText)
	statusLabel.SetHAlign(gtk4.AlignEnd)
	if statusClass != "" {
		statusLabel.AddCSSClass(statusClass)
	}
	vbox.Append(&statusLabel.Widget)

	tempStr := fmt.Sprintf("%d\xc2\xb0C", disk.Health.Temperature)
	tempLabel := gtk4.LabelNew(tempStr)
	tempLabel.SetHAlign(gtk4.AlignEnd)
	tempLabel.SetSensitive(false)
	tempClass := disk.Health.TempStatusClass()
	if tempClass != "" {
		tempLabel.AddCSSClass(tempClass)
		tempLabel.SetSensitive(true)
	}
	vbox.Append(&tempLabel.Widget)

	vboxWidget := vbox.Widget
	hbox.Append(&vboxWidget)
	w := hbox.Widget
	return &w
}

func (dv *DiskView) createHealthDetailsBox(health *smartctl.SmartHealth) *gtk4.Widget {
	vbox := gtk4.BoxNew(gtk4.OrientationVertical, 4)
	vbox.SetMarginStart(4)
	vbox.SetMarginEnd(4)

	barBox := gtk4.BoxNew(gtk4.OrientationHorizontal, 8)

	bar := gtk4.LevelBarNew()
	bar.SetHExpand(true)

	healthPct := 1.0
	if health.IsCritical() {
		healthPct = 0.15
	} else if health.HasWarnings() {
		healthPct = 0.5
	}
	bar.SetValue(healthPct)

	barWidget := bar.Widget
	barWidget.SetHExpand(true)
	barBox.Append(&barWidget)

	if health.PowerOnHours > 0 {
		hoursStr := formatHours(int64(health.PowerOnHours))
		if hoursStr != "" {
			hoursLabel := gtk4.LabelNew(hoursStr)
			hoursLabel.SetSensitive(false)
			barBox.Append(&hoursLabel.Widget)
		}
	}

	vbox.Append(&barBox.Widget)

	criticalAttrs := dv.gatherCriticalAttributes(health)
	if len(criticalAttrs) > 0 {
		attrBox := gtk4.BoxNew(gtk4.OrientationHorizontal, 8)
		attrBox.SetHomogeneous(true)

		for _, a := range criticalAttrs {
			avBox := gtk4.BoxNew(gtk4.OrientationVertical, 1)
			avBox.SetHAlign(gtk4.AlignCenter)

			nameLabel := gtk4.LabelNew(a.Name)
			nameLabel.SetSensitive(false)

			attrClass := "disk-health-good"
			if a.WhenFailed != "" || (isCriticalAttr(a.ID) && a.RawValue > 0) {
				attrClass = "disk-health-critical"
			}
			valLabel := gtk4.LabelNew(a.RawString)
			valLabel.AddCSSClass(attrClass)

			avBox.Append(&nameLabel.Widget)
			avBox.Append(&valLabel.Widget)
			attrBox.Append(&avBox.Widget)
		}

		vbox.Append(&attrBox.Widget)
	}

	vboxWidget := vbox.Widget
	return &vboxWidget
}

func (dv *DiskView) gatherCriticalAttributes(health *smartctl.SmartHealth) []smartctl.SMARTAttributeInfo {
	var result []smartctl.SMARTAttributeInfo

	for _, attr := range health.Attributes {
		if isCriticalAttr(attr.ID) && (attr.RawValue > 0 || attr.WhenFailed != "") {
			result = append(result, attr)
		}
	}

	tempFound := false
	for _, attr := range health.Attributes {
		if (attr.ID == 190 || attr.ID == 194) && !tempFound {
			if health.Temperature >= 45 {
				result = append(result, attr)
				tempFound = true
			}
		}
	}

	wearFound := false
	for _, attr := range health.Attributes {
		if attr.ID == 177 && !wearFound && attr.RawValue > 0 {
			result = append(result, attr)
			wearFound = true
		}
	}

	return result
}

func isCriticalAttr(id int) bool {
	return id == 5 || id == 197 || id == 198
}

func formatHours(hours int64) string {
	if hours <= 0 {
		return ""
	}
	days := hours / 24
	years := days / 365
	remainderDays := days % 365
	if years > 0 {
		return fmt.Sprintf("%dy %dd", years, remainderDays)
	}
	return fmt.Sprintf("%dd", days)
}

func formatBytes(bytes int64) string {
	if bytes <= 0 {
		return "Unknown"
	}
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

func (dv *DiskView) createPartitionsSection(disk *model.DiskInfo) *gtk4.Widget {
	grid := gtk4.GridNew()
	grid.SetColumnSpacing(8)
	grid.SetRowSpacing(2)
	grid.SetMarginStart(36)
	grid.SetMarginTop(4)

	for rowIdx, part := range disk.Partitions {
		dv.attachPartitionRow(grid, rowIdx, &part)
	}

	gridWidget := grid.Widget
	return &gridWidget
}

func (dv *DiskView) attachPartitionRow(grid *gtk4.Grid, row int, part *model.PartitionInfo) {
	partText := "/dev/" + part.Name
	if part.Label != "" {
		partText += " [" + part.Label + "]"
	}
	if part.MountPoint != "" {
		partText += " \xe2\x86\x92 " + part.MountPoint
	}

	nameLabel := gtk4.LabelNew(partText)
	nameLabel.SetHAlign(gtk4.AlignStart)
	nameLabel.SetHExpand(true)
	nameLabel.SetMarkup(fmt.Sprintf("<b>%s</b>", strings.ReplaceAll(partText, "&", "&amp;")))
	nameLabel.AddCSSClass("monospace-label")
	grid.Attach(&nameLabel.Widget, 0, row, 1, 1)

	sizeBox := gtk4.BoxNew(gtk4.OrientationHorizontal, 0)
	sizeBox.SetSizeRequest(80, -1)

	sizeLabel := gtk4.LabelNew(formatBytes(part.SizeBytes))
	sizeLabel.SetHAlign(gtk4.AlignEnd)
	sizeLabel.SetSensitive(false)
	sizeLabel.AddCSSClass("monospace-label")
	sizeBox.Append(&sizeLabel.Widget)
	grid.Attach(&sizeBox.Widget, 1, row, 1, 1)

	fsBox := gtk4.BoxNew(gtk4.OrientationHorizontal, 0)
	fsBox.SetSizeRequest(64, -1)

	fsText := ""
	if part.IsEncrypted {
		fsText = "LUKS"
	} else if part.FSType != "" {
		fsText = part.FSType
	}
	if fsText != "" {
		fsLabel := gtk4.LabelNew(fsText)
		fsLabel.SetHAlign(gtk4.AlignStart)
		fsLabel.SetSensitive(!part.IsEncrypted)
		fsLabel.AddCSSClass("monospace-label")
		if part.IsEncrypted {
			fsLabel.AddCSSClass("disk-health-warning")
		}
		fsBox.Append(&fsLabel.Widget)
	}
	grid.Attach(&fsBox.Widget, 2, row, 1, 1)

	if part.MountPoint != "" {

		unmountBtn := gtk4.ButtonNewWithLabel("Unmount")
		unmountBtn.SetSizeRequest(100, -1)
		unmountBtn.SetHAlign(gtk4.AlignEnd)
		partPath := part.DevicePath
		mapperPath := part.MapperPath
		isEnc := part.IsEncrypted
		source := partPath
		if isEnc && mapperPath != "" {
			source = mapperPath
		}

		unmountBtn.OnClicked(func() {
			widget.PromptForSudo(dv.parentWin,
				"Unmounting requires root privileges.",
				func(password string) {
					if password == "" {
						return
					}
					go func() {
						ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
						defer cancel()
						c := exec.CommandContext(ctx, "sudo", "-S", "-k", "umount", source)
						c.Stdin = strings.NewReader(password + "\n")
						c.Run()
						time.Sleep(500 * time.Millisecond)
						gtk4.IdleAdd(func() { dv.rescanOnly() })
					}()
				},
			)
		})

		btnBox := gtk4.BoxNew(gtk4.OrientationHorizontal, 0)
		btnBox.SetHAlign(gtk4.AlignEnd)
		btnBox.Append(&unmountBtn.Widget)
		grid.Attach(&btnBox.Widget, 3, row, 1, 1)

	} else if part.FSType != "" || part.IsEncrypted {

		mountBtn := gtk4.ButtonNewWithLabel("Mount")
		mountBtn.SetSizeRequest(100, -1)
		mountBtn.SetHAlign(gtk4.AlignEnd)
		p := part
		mountBtn.OnClicked(func() {
			widget.ShowMountDialog(dv.parentWin, p.DevicePath, p.Name, p.IsEncrypted, p.IsUnlocked, p.MapperName, func() {
				go func() {
					time.Sleep(500 * time.Millisecond)
					gtk4.IdleAdd(func() { dv.rescanOnly() })
				}()
			})
		})

		btnBox := gtk4.BoxNew(gtk4.OrientationHorizontal, 0)
		btnBox.SetHAlign(gtk4.AlignEnd)
		btnBox.Append(&mountBtn.Widget)
		grid.Attach(&btnBox.Widget, 3, row, 1, 1)

	}
}

func (dv *DiskView) Widget() *gtk4.Widget { return &dv.box.Widget }
func (dv *DiskView) Name() string         { return "disks" }
func (dv *DiskView) Title() string        { return "Disks" }
func (dv *DiskView) IconName() string     { return "drive-harddisk-symbolic" }
func (dv *DiskView) OnShow()              { dv.refresh() }
func (dv *DiskView) OnHide()              {}
func (dv *DiskView) Destroy()             {}
