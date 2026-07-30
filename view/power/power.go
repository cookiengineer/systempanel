package power

import (
	"context"

	"github.com/cookiengineer/systempanel/bindings/gtk4"
	"github.com/cookiengineer/systempanel/detect"
	"github.com/cookiengineer/systempanel/model"
	"github.com/cookiengineer/systempanel/view"
)

var Descriptor = view.ViewDescriptor{
	Name:     "power",
	Title:    "Power",
	IconName: "system-shutdown-symbolic",
	DetectFn: func() bool {
		return detect.HasProgram("systemctl") && detect.HasSystemd()
	},
	Factory: func() view.View {
		return NewPowerView()
	},
}

type PowerView struct {
	box    *gtk4.Box
	model  *model.PowerModel
	showing bool
}

func NewPowerView() *PowerView {
	pv := &PowerView{
		box:   gtk4.BoxNew(gtk4.OrientationVertical, 12),
		model: &model.PowerModel{},
	}
	pv.box.SetMarginStart(24)
	pv.box.SetMarginEnd(24)
	pv.box.SetMarginTop(24)
	pv.box.SetMarginBottom(24)

	header := gtk4.LabelNew("Power Actions")
	header.AddCSSClass("header-label")
	header.SetHAlign(gtk4.AlignCenter)
	headerWidget := header.Widget
	pv.box.Append(&headerWidget)

	grid := gtk4.BoxNew(gtk4.OrientationHorizontal, 12)
	grid.SetHAlign(gtk4.AlignCenter)
	grid.SetHExpand(true)

	actions := []struct {
		name     string
		cssClass string
		icon     string
		action   func()
	}{
		{"Shutdown", "shutdown", "system-shutdown-symbolic", func() {
			pv.model.PowerOff()
			gtk4AppQuit()
		}},
		{"Reboot", "reboot", "system-reboot-symbolic", func() {
			pv.model.Reboot()
		}},
		{"Suspend", "suspend", "system-suspend-symbolic", func() {
			pv.model.Suspend()
		}},
	}

	for _, a := range actions {
		btn := pv.createButton(a.name, a.cssClass, a.icon, a.action)
		btnWidget := btn.Widget
		grid.Append(&btnWidget)
	}

	hibernateBtn := pv.createButton("Hibernate", "hibernate", "system-hibernate-symbolic", func() {
		pv.model.Hibernate()
	})
	hibernateBtn.SetSensitive(detect.HasProgram("systemctl"))
	hibernateWidget := hibernateBtn.Widget
	grid.Append(&hibernateWidget)

	gridWidget := grid.Widget
	pv.box.Append(&gridWidget)

	noteLabel := gtk4.LabelNew("Actions execute via systemctl")
	noteLabel.SetHAlign(gtk4.AlignCenter)
	noteLabel.SetSensitive(false)
	noteWidget := noteLabel.Widget
	pv.box.Append(&noteWidget)

	return pv
}

func (pv *PowerView) createButton(name, cssClass, iconName string, action func()) *gtk4.Button {
	btn := gtk4.ButtonNew()
	btn.AddCSSClass("power-button")
	btn.AddCSSClass(cssClass)
	btn.SetHExpand(true)

	box := gtk4.BoxNew(gtk4.OrientationVertical, 8)
	box.SetHAlign(gtk4.AlignCenter)
	box.SetMarginTop(4)
	box.SetMarginBottom(4)

	icon := gtk4.ImageNewFromIconName(iconName)
	icon.SetPixelSize(32)
	iconWidget := icon.Widget
	box.Append(&iconWidget)

	label := gtk4.LabelNew(name)
	label.AddCSSClass("power-button-label")
	labelWidget := label.Widget
	box.Append(&labelWidget)

	boxWidget := box.Widget
	btn.SetChild(&boxWidget)

	btn.OnClicked(action)
	return btn
}

func (pv *PowerView) Widget() *gtk4.Widget { return &pv.box.Widget }
func (pv *PowerView) Name() string         { return "power" }
func (pv *PowerView) Title() string        { return "Power" }
func (pv *PowerView) IconName() string     { return "system-shutdown-symbolic" }
func (pv *PowerView) OnShow()              { pv.showing = true }
func (pv *PowerView) OnHide()              { pv.showing = false }
func (pv *PowerView) Destroy()             {}

func gtk4AppQuit() {
	ctx := context.Background()
	_ = ctx
}
