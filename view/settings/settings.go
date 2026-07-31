package settings

import (
	"os/exec"

	"github.com/cookiengineer/systempanel/bindings/gtk4"
	"github.com/cookiengineer/systempanel/config"
	"github.com/cookiengineer/systempanel/parsers/gtktheme"
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
	box          *gtk4.Box
	settings     *config.Settings
	toggles      map[string]*gtk4.Switch
	onVisibility func()
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

	sv.settings, _ = config.Load()

	header := gtk4.LabelNew("Settings")
	header.AddCSSClass("header-label")
	sv.box.Append(&header.Widget)

	desc := gtk4.LabelNew("Enable or disable individual views in the sidebar.")
	desc.SetWrap(true)
	desc.SetMarginBottom(12)
	sv.box.Append(&desc.Widget)

	for _, d := range view.Registry {
		if d.Name == "settings" {
			continue
		}
		row := sv.createToggleRow(d.Name, d.Title, sv.settings.IsVisible(d.Name))
		sv.box.Append(&row.Widget)
	}

	sep := gtk4.LabelNew("")
	sep.SetMarginTop(8)
	sv.box.Append(&sep.Widget)

	themeHeader := gtk4.LabelNew("Appearance")
	themeHeader.AddCSSClass("header-label")
	sv.box.Append(&themeHeader.Widget)

	darkRow := sv.createToggleRow("dark-mode", "Dark Mode", gtk4.GetDarkTheme())
	sv.box.Append(&darkRow.Widget)

	aboutBox := sv.buildAboutSection()
	sv.box.Append(&aboutBox.Widget)

	return sv
}

func (sv *SettingsView) SetVisibilityCallback(fn func()) {
	sv.onVisibility = fn
}

func (sv *SettingsView) createToggleRow(name, title string, active bool) *gtk4.Box {
	row := gtk4.BoxNew(gtk4.OrientationHorizontal, 12)
	row.SetMarginTop(4)
	row.SetMarginBottom(4)

	label := gtk4.LabelNew(title)
	label.SetHExpand(true)
	label.SetHAlign(gtk4.AlignStart)
	row.Append(&label.Widget)

	sw := gtk4.SwitchNew()
	sw.SetActive(active)
	sv.toggles[name] = sw
	sw.OnActivate(func() {
		if name == "dark-mode" {
			dark := sw.GetActive()
			gtk4.SetDarkTheme(dark)
			gtktheme.SaveDarkTheme(dark)
			return
		}
		sv.settings.SetVisible(name, sw.GetActive())
		if sv.onVisibility != nil {
			sv.onVisibility()
		}
	})
	row.Append(&sw.Widget)

	return row
}

func (sv *SettingsView) buildAboutSection() *gtk4.Box {
	box := gtk4.BoxNew(gtk4.OrientationVertical, 6)
	box.SetMarginTop(24)

	header := gtk4.LabelNew("About")
	header.AddCSSClass("header-label")
	box.Append(&header.Widget)

	info := gtk4.LabelNew("SystemPanel — adaptive GTK4 system control panel")
	info.SetWrap(true)
	box.Append(&info.Widget)

	version := gtk4.LabelNew("Version: 0.1.0")
	version.SetSensitive(false)
	box.Append(&version.Widget)

	websiteBtn := gtk4.ButtonNew()
	websiteBtn.SetLabel("Website: cookie.engineer")
	websiteBtn.OnClicked(func() {
		exec.Command("xdg-open", "https://cookie.engineer").Start()
	})
	websiteBtn.SetHAlign(gtk4.AlignStart)
	box.Append(&websiteBtn.Widget)

	repoBtn := gtk4.ButtonNew()
	repoBtn.SetLabel("Repository: github.com/cookiengineer/systempanel")
	repoBtn.OnClicked(func() {
		exec.Command("xdg-open", "https://github.com/cookiengineer/systempanel").Start()
	})
	repoBtn.SetHAlign(gtk4.AlignStart)
	box.Append(&repoBtn.Widget)

	return box
}

func (sv *SettingsView) Widget() *gtk4.Widget { return &sv.box.Widget }
func (sv *SettingsView) Name() string          { return "settings" }
func (sv *SettingsView) Title() string         { return "Settings" }
func (sv *SettingsView) IconName() string      { return "emblem-system-symbolic" }
func (sv *SettingsView) OnShow()               {}
func (sv *SettingsView) OnHide()               {}
func (sv *SettingsView) Destroy()              {}
