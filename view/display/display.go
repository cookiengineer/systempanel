package display

import (
	"fmt"

	"github.com/cookiengineer/systempanel/bindings/gtk4"
	"github.com/cookiengineer/systempanel/detect"
	"github.com/cookiengineer/systempanel/model"
	"github.com/cookiengineer/systempanel/view"
)

var Descriptor = view.ViewDescriptor{
	Name:     "display",
	Title:    "Display",
	IconName: "video-display-symbolic",
	DetectFn: func() bool { return detect.HasProgram("brightnessctl") },
	Factory:  func() view.View { return NewDisplayView() },
}

type DisplayView struct {
	box     *gtk4.Box
	model   *model.DisplayModel
	sliders map[string]*gtk4.Scale
	labels  map[string]*gtk4.Label
}

func NewDisplayView() *DisplayView {
	dv := &DisplayView{
		box:     gtk4.BoxNew(gtk4.OrientationVertical, 12),
		model:   &model.DisplayModel{},
		sliders: make(map[string]*gtk4.Scale),
		labels:  make(map[string]*gtk4.Label),
	}
	dv.box.SetMarginStart(24)
	dv.box.SetMarginEnd(24)
	dv.box.SetMarginTop(24)
	dv.box.SetMarginBottom(24)

	header := gtk4.LabelNew("Display Brightness")
	header.AddCSSClass("header-label")
	headerWidget := header.Widget
	dv.box.Append(&headerWidget)

	dv.buildDevices()

	return dv
}

func (dv *DisplayView) buildDevices() {
	devices, err := dv.model.ListDevices()
	if err != nil || len(devices) == 0 {
		label := gtk4.LabelNew("No brightness-controllable devices found")
		label.SetSensitive(false)
		labelWidget := label.Widget
		dv.box.Append(&labelWidget)
		return
	}

	for _, dev := range devices {
		info, err := dv.model.GetDevice(dev.Name)
		if err != nil {
			continue
		}

		row := gtk4.BoxNew(gtk4.OrientationVertical, 4)
		row.SetMarginBottom(12)

		nameLabel := gtk4.LabelNew(dev.Name)
		nameLabel.SetHAlign(gtk4.AlignStart)
		nlWidget := nameLabel.Widget
		row.Append(&nlWidget)

		slider := gtk4.ScaleNewWithRange(gtk4.OrientationHorizontal, 0, float64(info.Max), 1)
		slider.SetValue(float64(info.Current))
		slider.SetHExpand(true)
		devName := dev.Name
		slider.OnValueChanged(func(val float64) {
			dv.model.SetBrightness(devName, int(val))
		})
		dv.sliders[dev.Name] = slider
		sliderWidget := slider.Widget
		row.Append(&sliderWidget)

		pctLabel := gtk4.LabelNew(fmt.Sprintf("%.0f%%", info.Percentage))
		pctLabel.SetHAlign(gtk4.AlignEnd)
		dv.labels[dev.Name] = pctLabel
		plWidget := pctLabel.Widget
		row.Append(&plWidget)

		rowWidget := row.Widget
		dv.box.Append(&rowWidget)
	}
}

func (dv *DisplayView) Widget() *gtk4.Widget { return &dv.box.Widget }
func (dv *DisplayView) Name() string          { return "display" }
func (dv *DisplayView) Title() string         { return "Display" }
func (dv *DisplayView) IconName() string      { return "video-display-symbolic" }
func (dv *DisplayView) OnShow()               {}
func (dv *DisplayView) OnHide()               {}
func (dv *DisplayView) Destroy()              {}
