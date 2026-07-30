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
	box   *gtk4.Box
	model *model.BatteryModel
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

	header := gtk4.LabelNew("Battery")
	header.AddCSSClass("header-label")
	headerWidget := header.Widget
	bv.box.Append(&headerWidget)

	bv.buildBatteryInfo()

	refreshBtn := gtk4.ButtonNewWithLabel("Refresh")
	refreshBtn.OnClicked(func() {
		bv.refresh()
	})
	refreshWidget := refreshBtn.Widget
	bv.box.Append(&refreshWidget)

	return bv
}

func (bv *BatteryView) buildBatteryInfo() {
	infos, err := bv.model.GetBatteries()
	if err != nil || len(infos) == 0 {
		label := gtk4.LabelNew("No battery information available")
		label.SetSensitive(false)
		labelWidget := label.Widget
		bv.box.Append(&labelWidget)
		return
	}

	for _, info := range infos {
		card := gtk4.BoxNew(gtk4.OrientationVertical, 8)
		card.AddCSSClass("device-row")
		card.SetMarginBottom(12)

		nameLabel := gtk4.LabelNew(info.Model)
		if info.Model == "" {
			nameLabel.SetText(info.NativePath)
		}
		cardWidget := card.Widget
		bv.box.Append(&cardWidget)

		bar := gtk4.LevelBarNew()
		bar.SetValue(info.Percentage / 100.0)
		bar.SetHExpand(true)
		barWidget := bar.Widget
		card.Append(&barWidget)

		infoText := fmt.Sprintf("%.0f%% - %s", info.Percentage, info.State)
		if info.TimeToEmpty > 0 {
			hours := info.TimeToEmpty / 3600
			minutes := (info.TimeToEmpty % 3600) / 60
			infoText += fmt.Sprintf(" (%dh %dm remaining)", hours, minutes)
		}
		infoLabel := gtk4.LabelNew(infoText)
		ilWidget := infoLabel.Widget
		card.Append(&ilWidget)

		if info.Capacity > 0 {
			capacityText := fmt.Sprintf("Capacity: %.1f%%", info.Capacity)
			capLabel := gtk4.LabelNew(capacityText)
			capLabel.SetSensitive(false)
			clWidget := capLabel.Widget
			card.Append(&clWidget)
		}
	}
}

func (bv *BatteryView) refresh() {}

func (bv *BatteryView) Widget() *gtk4.Widget { return &bv.box.Widget }
func (bv *BatteryView) Name() string          { return "battery" }
func (bv *BatteryView) Title() string         { return "Battery" }
func (bv *BatteryView) IconName() string      { return "battery-good-symbolic" }
func (bv *BatteryView) OnShow()               {}
func (bv *BatteryView) OnHide()               {}
func (bv *BatteryView) Destroy()              {}
