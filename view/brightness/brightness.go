package brightness

import (
	"github.com/cookiengineer/systempanel/bindings/gtk4"
	"github.com/cookiengineer/systempanel/detect"
	"github.com/cookiengineer/systempanel/model"
	"github.com/cookiengineer/systempanel/view"
)

var Descriptor = view.ViewDescriptor{
	Name:     "brightness",
	Title:    "Brightness",
	IconName: "display-brightness-symbolic",
	DetectFn: func() bool {
		return detect.HasProgram("brightnessctl")
	},
	Factory: func() view.View { return NewBrightnessView() },
}

type brightnessDeviceItem struct {
	device model.BrightnessDevice
	scale  *gtk4.Scale
}

type BrightnessView struct {
	box     *gtk4.Box
	model   *model.BrightnessModel
	listBox *gtk4.ListBox
	items   []brightnessDeviceItem
	rows    []*gtk4.ListBoxRow
}

func NewBrightnessView() *BrightnessView {
	bv := &BrightnessView{
		box:   gtk4.BoxNew(gtk4.OrientationVertical, 12),
		model: model.NewBrightnessModel(),
	}
	bv.box.SetMarginStart(24)
	bv.box.SetMarginEnd(24)
	bv.box.SetMarginTop(24)
	bv.box.SetMarginBottom(24)

	header := gtk4.LabelNew("Brightness")
	header.AddCSSClass("header-label")
	headerWidget := header.Widget
	bv.box.Append(&headerWidget)

	desc := gtk4.LabelNew("Control backlight brightness for connected displays.")
	desc.SetWrap(true)
	desc.SetMarginBottom(12)
	dv := desc.Widget
	bv.box.Append(&dv)

	bv.listBox = gtk4.ListBoxNew()
	bv.listBox.SetSelectionMode(gtk4.SelectionNone)

	scrollW := gtk4.ScrolledWindowNew()
	scrollW.SetPolicy(gtk4.PolicyNever, gtk4.PolicyAutomatic)
	scrollW.SetHExpand(true)
	scrollW.SetVExpand(true)

	listWidget := bv.listBox.Widget
	scrollW.SetChild(&listWidget)
	sw := scrollW.Widget
	bv.box.Append(&sw)

	btnBox := gtk4.BoxNew(gtk4.OrientationHorizontal, 8)
	btnBox.SetMarginTop(4)

	refreshBtn := gtk4.ButtonNewWithLabel("Refresh")
	refreshBtn.OnClicked(func() { bv.refresh() })
	rbw := refreshBtn.Widget
	btnBox.Append(&rbw)

	btnBoxWidget := btnBox.Widget
	bv.box.Append(&btnBoxWidget)

	bv.refresh()

	return bv
}

func (bv *BrightnessView) refresh() {
	for _, r := range bv.rows {
		bv.listBox.Remove(r)
	}
	bv.rows = bv.rows[:0]
	bv.items = bv.items[:0]

	backlights, _ := bv.model.ListBacklights()

	if len(backlights) == 0 {
		label := gtk4.LabelNew("No backlight devices found")
		label.SetSensitive(false)
		label.SetHAlign(gtk4.AlignCenter)
		lw := label.Widget
		row := gtk4.ListBoxRowNew()
		row.SetChild(&lw)
		bv.listBox.Append(row)
		bv.rows = append(bv.rows, row)
		return
	}

	hdr := bv.createSectionLabel("─ Backlight Devices ─")
	bv.listBox.Append(hdr)
	bv.rows = append(bv.rows, hdr)

	for _, dev := range backlights {
		row := bv.createDeviceRow(dev)
		bv.listBox.Append(row)
		bv.rows = append(bv.rows, row)
	}
}

func (bv *BrightnessView) createSectionLabel(text string) *gtk4.ListBoxRow {
	row := gtk4.ListBoxRowNew()
	row.SetSensitive(false)
	label := gtk4.LabelNew(text)
	label.SetHAlign(gtk4.AlignCenter)
	label.SetMarginTop(8)
	label.SetMarginBottom(4)
	label.SetSensitive(false)
	lw := label.Widget
	row.SetChild(&lw)
	return row
}

func (bv *BrightnessView) createDeviceRow(dev model.BrightnessDevice) *gtk4.ListBoxRow {
	row := gtk4.ListBoxRowNew()
	row.AddCSSClass("device-row")

	hbox := gtk4.BoxNew(gtk4.OrientationHorizontal, 8)

	icon := gtk4.ImageNewFromIconName("display-brightness-symbolic")
	icon.SetPixelSize(20)
	iconWidget := icon.Widget
	hbox.Append(&iconWidget)

	nameLabel := gtk4.LabelNew(dev.Name)
	nameLabel.SetHExpand(true)
	nameLabel.SetHAlign(gtk4.AlignStart)
	nlWidget := nameLabel.Widget
	hbox.Append(&nlWidget)

	scale := gtk4.ScaleNewWithRange(gtk4.OrientationHorizontal, 0, 100, 1)
	scale.SetValue(float64(dev.Percentage))
	scale.SetSizeRequest(180, -1)

	item := brightnessDeviceItem{device: dev, scale: scale}
	bv.items = append(bv.items, item)

	deviceName := dev.Name
	skipInit := true
	scale.OnValueChanged(func(val float64) {
		if skipInit {
			skipInit = false
			return
		}
		bv.model.SetBrightness(deviceName, int(val))
	})
	sw := scale.Widget
	hbox.Append(&sw)

	hboxWidget := hbox.Widget
	row.SetChild(&hboxWidget)

	return row
}

func (bv *BrightnessView) Widget() *gtk4.Widget { return &bv.box.Widget }
func (bv *BrightnessView) Name() string          { return "brightness" }
func (bv *BrightnessView) Title() string         { return "Brightness" }
func (bv *BrightnessView) IconName() string      { return "display-brightness-symbolic" }
func (bv *BrightnessView) OnShow()               { bv.refresh() }
func (bv *BrightnessView) OnHide()               {}
func (bv *BrightnessView) Destroy()              {}
