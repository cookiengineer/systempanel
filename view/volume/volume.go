package volume

import (
	"github.com/cookiengineer/systempanel/bindings/gtk4"
	"github.com/cookiengineer/systempanel/detect"
	"github.com/cookiengineer/systempanel/model"
	"github.com/cookiengineer/systempanel/view"
)

var Descriptor = view.ViewDescriptor{
	Name:     "volume",
	Title:    "Volume",
	IconName: "audio-volume-high-symbolic",
	DetectFn: func() bool {
		return detect.HasProgram("pactl") && detect.HasPulseAudio()
	},
	Factory: func() view.View { return NewVolumeView() },
}

type volumeDeviceItem struct {
	device model.VolumeDevice
	muteBtn *gtk4.Button
	scale   *gtk4.Scale
}

type VolumeView struct {
	box     *gtk4.Box
	model   *model.VolumeModel
	listBox *gtk4.ListBox
	items   []volumeDeviceItem
	rows    []*gtk4.ListBoxRow
}

func NewVolumeView() *VolumeView {
	vv := &VolumeView{
		box:   gtk4.BoxNew(gtk4.OrientationVertical, 12),
		model: model.NewVolumeModel(),
	}
	vv.box.SetMarginStart(24)
	vv.box.SetMarginEnd(24)
	vv.box.SetMarginTop(24)
	vv.box.SetMarginBottom(24)

	header := gtk4.LabelNew("Audio Devices")
	header.AddCSSClass("header-label")
	headerWidget := header.Widget
	vv.box.Append(&headerWidget)

	vv.listBox = gtk4.ListBoxNew()
	vv.listBox.SetSelectionMode(gtk4.SelectionNone)

	scrollW := gtk4.ScrolledWindowNew()
	scrollW.SetPolicy(gtk4.PolicyNever, gtk4.PolicyAutomatic)
	scrollW.SetHExpand(true)
	scrollW.SetVExpand(true)

	scrollWidget := scrollW.Widget
	listWidget := vv.listBox.Widget
	scrollW.SetChild(&listWidget)
	vv.box.Append(&scrollWidget)

	btnBox := gtk4.BoxNew(gtk4.OrientationHorizontal, 8)
	btnBox.SetMarginTop(4)

	refreshBtn := gtk4.ButtonNewWithLabel("Refresh")
	refreshBtn.OnClicked(func() { vv.refresh() })
	rbw := refreshBtn.Widget
	btnBox.Append(&rbw)

	btnBoxWidget := btnBox.Widget
	vv.box.Append(&btnBoxWidget)

	vv.refresh()

	return vv
}

func (vv *VolumeView) refresh() {
	for _, r := range vv.rows {
		vv.listBox.Remove(r)
	}
	vv.rows = vv.rows[:0]
	vv.items = vv.items[:0]

	sinks, _ := vv.model.ListSinks()
	sources, _ := vv.model.ListSources()

	if len(sources) > 0 {
		hdr := vv.createSectionLabel("─ Input Devices ─")
		vv.listBox.Append(hdr)
		vv.rows = append(vv.rows, hdr)
		for _, s := range sources {
			row := vv.createDeviceRow(s)
			vv.listBox.Append(row)
			vv.rows = append(vv.rows, row)
		}
	}

	if len(sinks) > 0 {
		hdr := vv.createSectionLabel("─ Output Devices ─")
		vv.listBox.Append(hdr)
		vv.rows = append(vv.rows, hdr)
		for _, s := range sinks {
			row := vv.createDeviceRow(s)
			vv.listBox.Append(row)
			vv.rows = append(vv.rows, row)
		}
	}

	if len(sinks) == 0 && len(sources) == 0 {
		label := gtk4.LabelNew("No audio devices found")
		label.SetSensitive(false)
		label.SetHAlign(gtk4.AlignCenter)
		lw := label.Widget
		row := gtk4.ListBoxRowNew()
		row.SetChild(&lw)
		vv.listBox.Append(row)
	}
}

func (vv *VolumeView) createSectionLabel(text string) *gtk4.ListBoxRow {
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

func (vv *VolumeView) createDeviceRow(dev model.VolumeDevice) *gtk4.ListBoxRow {
	row := gtk4.ListBoxRowNew()
	row.AddCSSClass("device-row")

	hbox := gtk4.BoxNew(gtk4.OrientationHorizontal, 8)

	muteBtn := gtk4.ButtonNew()
	vv.updateMuteIcon(muteBtn, dev)

	name := dev.Name
	isInput := dev.IsInput
	item := volumeDeviceItem{device: dev, muteBtn: muteBtn}
	vv.items = append(vv.items, item)

	muteBtn.OnClicked(func() {
		item.device.Mute = !item.device.Mute
		if isInput {
			vv.model.ToggleMuteSource(name)
		} else {
			vv.model.ToggleMuteSink(name)
		}
		vv.updateMuteIcon(muteBtn, item.device)
	})
	mbWidget := muteBtn.Widget
	hbox.Append(&mbWidget)

	nameLabel := gtk4.LabelNew(dev.Description)
	nameLabel.SetHExpand(true)
	nameLabel.SetHAlign(gtk4.AlignStart)
	nlWidget := nameLabel.Widget
	hbox.Append(&nlWidget)

	scale := gtk4.ScaleNewWithRange(gtk4.OrientationHorizontal, 0, 100, 1)
	scale.SetValue(float64(dev.Volume))
	scale.SetSizeRequest(180, -1)
	initialized := false
	scale.OnValueChanged(func(val float64) {
		if !initialized {
			initialized = true
			return
		}
		if isInput {
			return
		}
		vv.model.SetSinkVolume(name, int(val))
	})
	item.scale = scale
	sw := scale.Widget
	hbox.Append(&sw)

	hboxWidget := hbox.Widget
	row.SetChild(&hboxWidget)

	return row
}

func (vv *VolumeView) updateMuteIcon(btn *gtk4.Button, dev model.VolumeDevice) {
	if dev.Mute {
		btn.SetIconName("audio-volume-muted-symbolic")
	} else {
		btn.SetIconName("audio-volume-high-symbolic")
	}
}

func (vv *VolumeView) Widget() *gtk4.Widget { return &vv.box.Widget }
func (vv *VolumeView) Name() string          { return "volume" }
func (vv *VolumeView) Title() string         { return "Volume" }
func (vv *VolumeView) IconName() string      { return "audio-volume-high-symbolic" }
func (vv *VolumeView) OnShow()               { vv.refresh() }
func (vv *VolumeView) OnHide()               {}
func (vv *VolumeView) Destroy()              {}
