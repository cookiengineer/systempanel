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
	Factory:  func() view.View { return NewVolumeView() },
}

type VolumeView struct {
	box     *gtk4.Box
	model   *model.VolumeModel
	sliders map[string]*gtk4.Scale
	mutes   map[string]*gtk4.Button
}

func NewVolumeView() *VolumeView {
	vv := &VolumeView{
		box:     gtk4.BoxNew(gtk4.OrientationVertical, 12),
		model:   model.NewVolumeModel(),
		sliders: make(map[string]*gtk4.Scale),
		mutes:   make(map[string]*gtk4.Button),
	}
	vv.box.SetMarginStart(24)
	vv.box.SetMarginEnd(24)
	vv.box.SetMarginTop(24)
	vv.box.SetMarginBottom(24)

	header := gtk4.LabelNew("Audio Output")
	header.AddCSSClass("header-label")
	headerWidget := header.Widget
	vv.box.Append(&headerWidget)

	vv.buildDevices()

	refreshBtn := gtk4.ButtonNewWithLabel("Refresh")
	refreshBtn.OnClicked(func() {
		vv.refresh()
	})
	refreshWidget := refreshBtn.Widget
	vv.box.Append(&refreshWidget)

	return vv
}

func (vv *VolumeView) buildDevices() {
	sinks, _ := vv.model.ListSinks()

	if len(sinks) == 0 {
		label := gtk4.LabelNew("No audio devices found")
		label.SetSensitive(false)
		labelWidget := label.Widget
		vv.box.Append(&labelWidget)
		return
	}

	for _, sink := range sinks {
		row := gtk4.BoxNew(gtk4.OrientationHorizontal, 12)
		row.AddCSSClass("device-row")

		muteBtn := gtk4.ButtonNew()
		muteBtn.SetIconName("audio-volume-muted-symbolic")
		muteBtn.OnClicked(func() {
			vv.model.ToggleMuteSink(sink.Name)
		})
		mbWidget := muteBtn.Widget
		row.Append(&mbWidget)

		nameLabel := gtk4.LabelNew(sink.Description)
		nameLabel.SetHExpand(true)
		nameLabel.SetHAlign(gtk4.AlignStart)
		nlWidget := nameLabel.Widget
		row.Append(&nlWidget)

		scale := gtk4.ScaleNew(gtk4.OrientationHorizontal)
		scale.SetValue(float64(sink.Volume))
		scale.SetRange(0, 100)
		scale.SetSizeRequest(150, -1)
		scale.OnValueChanged(func(val float64) {
			vv.model.SetSinkVolume(sink.Name, int(val))
		})
		vv.sliders[sink.Name] = scale
		scaleWidget := scale.Widget
		row.Append(&scaleWidget)

		rowWidget := row.Widget
		vv.box.Append(&rowWidget)
	}
}

func (vv *VolumeView) refresh() {
	for _, s := range vv.sliders {
		s.SetValue(50)
	}
}

func (vv *VolumeView) Widget() *gtk4.Widget { return &vv.box.Widget }
func (vv *VolumeView) Name() string          { return "volume" }
func (vv *VolumeView) Title() string         { return "Volume" }
func (vv *VolumeView) IconName() string      { return "audio-volume-high-symbolic" }
func (vv *VolumeView) OnShow()               {}
func (vv *VolumeView) OnHide()               {}
func (vv *VolumeView) Destroy()              {}
