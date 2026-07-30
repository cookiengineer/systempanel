package settings

import (
	"github.com/cookiengineer/systempanel/bindings/gtk4"
	"github.com/cookiengineer/systempanel/config"
	"github.com/cookiengineer/systempanel/view"
)

var Descriptor = view.ViewDescriptor{
	Name:     "settings",
	Title:    "Settings",
	IconName: "emblem-system-symbolic",
	DetectFn: func() bool { return true },
	Factory:  func() view.View { return NewSettingsView() },
}

type SettingsView struct {
	box      *gtk4.Box
	settings *config.Settings
	toggles  map[string]*gtk4.Switch
}

func NewSettingsView() *SettingsView {
	sv := &SettingsView{
		box:     gtk4.BoxNew(gtk4.OrientationVertical, 8),
		toggles: make(map[string]*gtk4.Switch),
	}
	sv.box.SetMarginStart(24)
	sv.box.SetMarginEnd(24)
	sv.box.SetMarginTop(24)
	sv.box.SetMarginBottom(24)

	header := gtk4.LabelNew("View Visibility")
	header.AddCSSClass("header-label")
	headerWidget := header.Widget
	sv.box.Append(&headerWidget)

	desc := gtk4.LabelNew("Enable or disable individual views in the sidebar.")
	desc.SetWrap(true)
	desc.SetMarginBottom(12)
	descWidget := desc.Widget
	sv.box.Append(&descWidget)

	sv.settings, _ = config.Load()

	for _, d := range view.Registry {
		if d.Name == "settings" {
			continue
		}
		row := sv.createToggleRow(d.Name, d.Title, sv.settings.IsVisible(d.Name))
		rowWidget := row.Widget
		sv.box.Append(&rowWidget)
	}

	aboutBox := gtk4.BoxNew(gtk4.OrientationVertical, 4)
	aboutBox.SetMarginTop(24)
	aboutHeader := gtk4.LabelNew("About")
	aboutHeader.AddCSSClass("header-label")
	aboutHW := aboutHeader.Widget
	aboutBox.Append(&aboutHW)

	aboutLabel := gtk4.LabelNew("SystemPanel v0.1.0\nA GTK4 system control panel\nBuilt with Go + CGo")
	aboutLabel.SetWrap(true)
	aboutLW := aboutLabel.Widget
	aboutBox.Append(&aboutLW)

	aboutBW := aboutBox.Widget
	sv.box.Append(&aboutBW)

	return sv
}

func (sv *SettingsView) createToggleRow(name, title string, active bool) *gtk4.Box {
	row := gtk4.BoxNew(gtk4.OrientationHorizontal, 12)
	row.SetMarginStart(0)
	row.SetMarginEnd(0)
	row.SetMarginTop(4)
	row.SetMarginBottom(4)

	label := gtk4.LabelNew(title)
	label.SetHAlign(gtk4.AlignStart)
	label.SetHExpand(true)
	lw := label.Widget
	row.Append(&lw)

	sw := gtk4.SwitchNew()
	sw.SetActive(active)
	sv.toggles[name] = sw
	sw.OnActivate(func() {
		sv.settings.SetVisible(name, sw.GetActive())
	})
	swWidget := sw.Widget
	row.Append(&swWidget)

	return row
}

func (sv *SettingsView) Widget() *gtk4.Widget { return &sv.box.Widget }
func (sv *SettingsView) Name() string          { return "settings" }
func (sv *SettingsView) Title() string         { return "Settings" }
func (sv *SettingsView) IconName() string      { return "emblem-system-symbolic" }
func (sv *SettingsView) OnShow()               {}
func (sv *SettingsView) OnHide()               {}
func (sv *SettingsView) Destroy()              {}
