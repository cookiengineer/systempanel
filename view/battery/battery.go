package battery

import (
	"fmt"

	"github.com/cookiengineer/systempanel/bindings/gtk4"
	"github.com/cookiengineer/systempanel/detect"
	"github.com/cookiengineer/systempanel/model"
	"github.com/cookiengineer/systempanel/view"
)

var Descriptor = view.ViewDescriptor{
	Name:     "battery",
	Title:    "Battery",
	IconName: "battery-good-symbolic",
	DetectFn: func() bool {
		return detect.HasProgram("upower") && detect.HasBatteryHardware()
	},
	Factory: func() view.View { return NewBatteryView() },
}

type BatteryView struct {
	box     *gtk4.Box
	model   *model.BatteryModel
	listBox *gtk4.ListBox
	rows    []*gtk4.ListBoxRow
}

func NewBatteryView() *BatteryView {
	bv := &BatteryView{
		box:   gtk4.BoxNew(gtk4.OrientationVertical, 12),
		model: &model.BatteryModel{},
	}
	bv.box.SetMarginStart(24)
	bv.box.SetMarginEnd(24)
	bv.box.SetMarginTop(24)
	bv.box.SetMarginBottom(24)

	header := gtk4.LabelNew("upower Batteries")
	header.AddCSSClass("header-label")
	bv.box.Append(&header.Widget)

	bv.listBox = gtk4.ListBoxNew()
	bv.listBox.SetSelectionMode(gtk4.SelectionNone)

	scrollW := gtk4.ScrolledWindowNew()
	scrollW.SetPolicy(gtk4.PolicyNever, gtk4.PolicyAutomatic)
	scrollW.SetHExpand(true)
	scrollW.SetVExpand(true)
	scrollW.SetChild(&bv.listBox.Widget)
	bv.box.Append(&scrollW.Widget)

	btnBox := gtk4.BoxNew(gtk4.OrientationHorizontal, 8)
	btnBox.SetMarginTop(4)

	refreshBtn := gtk4.ButtonNewWithLabel("Refresh")
	refreshBtn.OnClicked(func() { bv.refresh() })
	btnBox.Append(&refreshBtn.Widget)

	bv.box.Append(&btnBox.Widget)

	bv.refresh()

	return bv
}

func (bv *BatteryView) refresh() {
	for _, r := range bv.rows {
		bv.listBox.Remove(r)
	}
	bv.rows = bv.rows[:0]

	batteries, err := bv.model.GetBatteries()
	if err != nil || len(batteries) == 0 {
		label := gtk4.LabelNew("No battery information available")
		label.SetSensitive(false)
		label.SetHAlign(gtk4.AlignCenter)
		row := gtk4.ListBoxRowNew()
		row.SetChild(&label.Widget)
		bv.listBox.Append(row)
		bv.rows = append(bv.rows, row)
		return
	}

	for _, b := range batteries {
		row := bv.createBatteryRow(b)
		bv.listBox.Append(row)
		bv.rows = append(bv.rows, row)
	}
}

func (bv *BatteryView) createBatteryRow(b model.BatteryInfo) *gtk4.ListBoxRow {
	row := gtk4.ListBoxRowNew()
	row.AddCSSClass("device-row")

	vbox := gtk4.BoxNew(gtk4.OrientationVertical, 4)
	vbox.SetMarginStart(8)
	vbox.SetMarginEnd(8)
	vbox.SetMarginTop(4)
	vbox.SetMarginBottom(4)

	topBox := gtk4.BoxNew(gtk4.OrientationHorizontal, 8)

	iconName := "battery-good-symbolic"
	if b.Percentage < 20 && (b.State == "discharging" || b.State == "Discharging") {
		iconName = "battery-caution-symbolic"
	} else if b.State == "charging" || b.State == "Charging" || b.State == "fully-charged" || b.State == "Fully charged" {
		iconName = "battery-full-charging-symbolic"
	}
	icon := gtk4.ImageNewFromIconName(iconName)
	icon.SetPixelSize(24)
	topBox.Append(&icon.Widget)

	infoBox := gtk4.BoxNew(gtk4.OrientationVertical, 2)
	nameLabel := gtk4.LabelNew(b.Name)
	nameLabel.SetHAlign(gtk4.AlignStart)
	if b.Model != "" {
		nameLabel.SetText(b.Name + " - " + b.Model)
	}
	infoBox.Append(&nameLabel.Widget)

	if b.Vendor != "" {
		vendorLabel := gtk4.LabelNew(b.Vendor)
		vendorLabel.SetSensitive(false)
		vendorLabel.SetHAlign(gtk4.AlignStart)
		infoBox.Append(&vendorLabel.Widget)
	}
	topBox.Append(&infoBox.Widget)

	topBox.SetHExpand(true)

	stateLabel := gtk4.LabelNew(b.State)
	stateLabel.SetSensitive(false)
	topBox.Append(&stateLabel.Widget)

	vbox.Append(&topBox.Widget)

	bar := gtk4.LevelBarNew()
	bar.SetValue(b.Percentage / 100.0)
	bar.SetHExpand(true)
	vbox.Append(&bar.Widget)

	bottomBox := gtk4.BoxNew(gtk4.OrientationHorizontal, 8)

	infoText := fmt.Sprintf("%.0f%%", b.Percentage)
	percentLabel := gtk4.LabelNew(infoText)
	bottomBox.Append(&percentLabel.Widget)

	if b.TimeToEmpty > 0 {
		hours := b.TimeToEmpty / 3600
		minutes := (b.TimeToEmpty % 3600) / 60
		timeText := fmt.Sprintf("%dh %dm remaining", hours, minutes)
		timeLabel := gtk4.LabelNew(timeText)
		timeLabel.SetSensitive(false)
		bottomBox.Append(&timeLabel.Widget)
	}
	if b.TimeToFull > 0 {
		hours := b.TimeToFull / 3600
		minutes := (b.TimeToFull % 3600) / 60
		timeText := fmt.Sprintf("%dh %dm until full", hours, minutes)
		timeLabel := gtk4.LabelNew(timeText)
		timeLabel.SetSensitive(false)
		bottomBox.Append(&timeLabel.Widget)
	}

	bottomBox.SetHExpand(true)

	if b.Capacity > 0 {
		capText := fmt.Sprintf("Health: %.0f%%", b.Capacity)
		capLabel := gtk4.LabelNew(capText)
		capLabel.SetSensitive(false)
		bottomBox.Append(&capLabel.Widget)
	}

	vbox.Append(&bottomBox.Widget)

	row.SetChild(&vbox.Widget)

	return row
}

func (bv *BatteryView) Widget() *gtk4.Widget { return &bv.box.Widget }
func (bv *BatteryView) Name() string          { return "battery" }
func (bv *BatteryView) Title() string         { return "Battery" }
func (bv *BatteryView) IconName() string      { return "battery-good-symbolic" }
func (bv *BatteryView) OnShow()               { bv.refresh() }
func (bv *BatteryView) OnHide()               {}
func (bv *BatteryView) Destroy()              {}
