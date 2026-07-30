package app

import (
	"log"
	"unsafe"

	"github.com/cookiengineer/systempanel/bindings/gtk4"
	"github.com/cookiengineer/systempanel/config"
	"github.com/cookiengineer/systempanel/detect"
	"github.com/cookiengineer/systempanel/view"
)

type SystemPanel struct {
	app         *gtk4.Application
	win         *gtk4.Window
	header      *gtk4.HeaderBar
	sidebar     *gtk4.ListBox
	stack       *gtk4.Stack
	settings    *config.Settings
	detection   detect.DetectionResult
	views        []view.View
	sidebarRows  map[string]*gtk4.ListBoxRow
	rowToName    map[unsafe.Pointer]string
	previousView string
}

func New(app *gtk4.Application) *SystemPanel {
	detection := detect.RunAll()

	settings, err := config.Load()
	if err != nil {
		log.Printf("Failed to load settings: %v", err)
		settings = config.DefaultSettings()
	}

	return &SystemPanel{
		app:         app,
		settings:    settings,
		detection:   detection,
		views:       make([]view.View, 0),
		sidebarRows: make(map[string]*gtk4.ListBoxRow),
		rowToName:   make(map[unsafe.Pointer]string),
	}
}

func (p *SystemPanel) Build() {
	p.win = gtk4.ApplicationWindowNew(p.app)
	p.win.SetTitle("SystemPanel")
	p.win.SetDefaultSize(900, 600)

	p.header = gtk4.HeaderBarNew()
	p.header.SetShowTitleButtons(true)
	p.win.SetTitlebar(p.header)

	titleLabel := gtk4.LabelNew("SystemPanel")
	titleLabel.SetMarkup("<b>SystemPanel</b>")
	p.header.SetTitleWidget(&titleLabel.Widget)

	settingsBtn := gtk4.ButtonNew()
	settingsBtn.SetIconName("emblem-system-symbolic")
	settingsBtn.SetTooltip("Settings")
	settingsBtn.OnClicked(p.onSettingsClicked)
	settingsBtnWidget := settingsBtn.Widget
	p.header.PackEnd(&settingsBtnWidget)

	mainBox := gtk4.BoxNew(gtk4.OrientationHorizontal, 0)
	mainBoxWidget := mainBox.Widget
	p.win.SetChild(&mainBoxWidget)

	sidebarScroll := gtk4.ScrolledWindowNew()
	sidebarScroll.SetPolicy(gtk4.PolicyNever, gtk4.PolicyAutomatic)
	sidebarScroll.SetSizeRequest(220, -1)

	p.sidebar = gtk4.ListBoxNew()
	p.sidebar.SetSelectionMode(gtk4.SelectionSingle)
	p.sidebar.OnRowActivated(p.onSidebarRowActivated)
	p.sidebar.AddCSSClass("sidebar")

	sidebarContainer := gtk4.BoxNew(gtk4.OrientationVertical, 0)
	sidebarContainer.AddCSSClass("sidebar-container")

	sidebarWidget := p.sidebar.Widget
	sidebarContainer.Append(&sidebarWidget)

	sidebarScrollWidget := sidebarScroll.Widget
	sidebarContainerWidget := sidebarContainer.Widget
	sidebarScroll.SetChild(&sidebarContainerWidget)
	mainBox.Append(&sidebarScrollWidget)

	p.stack = gtk4.StackNew()
	p.stack.SetTransitionType(gtk4.StackTransitionCrossfade)
	p.stack.SetTransitionDuration(150)
	p.stack.SetVHomogeneous(true)
	stackWidget := p.stack.Widget
	stackWidget.SetHExpand(true)
	stackWidget.SetVExpand(true)
	mainBox.Append(&stackWidget)

	p.buildSidebar()

	p.win.OnCloseRequest(func() bool {
		p.app.Quit()
		return false
	})

	p.win.Present()
}

func (p *SystemPanel) buildSidebar() {
	for _, desc := range view.Registry {
		detected := desc.DetectFn()
		row := p.createSidebarRow(desc, detected)
		p.sidebar.Append(row)
		p.sidebarRows[desc.Name] = row
		p.rowToName[row.Widget.Ptr()] = desc.Name

		viewWidget := gtk4.BoxNew(gtk4.OrientationVertical, 0)
		page := p.stack.AddTitled(&viewWidget.Widget, desc.Name, desc.Title)
		page.SetIconName(desc.IconName)

		if detected {
			lazyView := desc.Factory()
			p.views = append(p.views, lazyView)
			w := lazyView.Widget()
			viewWidget.Append(w)

			if wv, ok := lazyView.(interface{ SetParentWindow(*gtk4.Window) }); ok {
				wv.SetParentWindow(p.win)
			}
		} else {
			placeholder := gtk4.LabelNew("Not available: required dependencies missing")
			placeholder.SetHAlign(gtk4.AlignCenter)
			placeholder.SetVAlign(gtk4.AlignCenter)
			placeholder.SetSensitive(false)
			placeholderWidget := placeholder.Widget
			viewWidget.Append(&placeholderWidget)
		}
	}

	if len(p.sidebarRows) > 0 {
		for _, desc := range view.Registry {
			if desc.DetectFn() {
				p.stack.SetVisibleChildName(desc.Name)
				if row, ok := p.sidebarRows[desc.Name]; ok {
					p.sidebar.SelectRow(row)
				}
				break
			}
		}
	}
}

func (p *SystemPanel) createSidebarRow(desc view.ViewDescriptor, detected bool) *gtk4.ListBoxRow {
	row := gtk4.ListBoxRowNew()
	row.AddCSSClass("sidebar-row")
	row.SetSensitive(detected)

	hbox := gtk4.BoxNew(gtk4.OrientationHorizontal, 8)
	hbox.SetMarginStart(12)
	hbox.SetMarginEnd(12)
	hbox.SetMarginTop(4)
	hbox.SetMarginBottom(4)

	icon := gtk4.ImageNewFromIconName(desc.IconName)
	icon.SetPixelSize(20)
	icon.SetSensitive(detected)
	iconWidget := icon.Widget
	hbox.Append(&iconWidget)

	label := gtk4.LabelNew(desc.Title)
	label.AddCSSClass("sidebar-row-label")
	label.SetHAlign(gtk4.AlignStart)
	label.SetHExpand(true)
	label.SetSensitive(detected)
	labelWidget := label.Widget
	hbox.Append(&labelWidget)

	statusLabel := gtk4.LabelNew("OK")
	statusLabel.AddCSSClass("sidebar-row-label")
	statusLabel.SetSensitive(detected)
	if !detected {
		statusLabel.SetText("N/A")
	}
	statusWidget := statusLabel.Widget
	hbox.Append(&statusWidget)

	hboxWidget := hbox.Widget
	row.SetChild(&hboxWidget)

	return row
}

func (p *SystemPanel) onSidebarRowActivated(row *gtk4.ListBoxRow) {
	rowPtr := row.Widget.Ptr()
	name, ok := p.rowToName[rowPtr]
	if !ok {
		return
	}
	for _, desc := range view.Registry {
		if desc.Name == name && desc.DetectFn() {
			if name != "settings" {
				p.previousView = p.stack.GetVisibleChildName()
			}
			p.stack.SetVisibleChildName(name)
			return
		}
	}
}

func (p *SystemPanel) onSettingsClicked() {
	current := p.stack.GetVisibleChildName()
	if current == "settings" {
		if p.previousView != "" {
			p.stack.SetVisibleChildName(p.previousView)
			if row, ok := p.sidebarRows[p.previousView]; ok {
				p.sidebar.SelectRow(row)
			}
		}
		return
	}
	p.previousView = current
	p.stack.SetVisibleChildName("settings")
	if row, ok := p.sidebarRows["settings"]; ok {
		p.sidebar.SelectRow(row)
	}
}

func (p *SystemPanel) Window() *gtk4.Window {
	return p.win
}

func (p *SystemPanel) Settings() *config.Settings {
	return p.settings
}

func (p *SystemPanel) Detection() detect.DetectionResult {
	return p.detection
}
