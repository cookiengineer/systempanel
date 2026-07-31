package autostart

import (
	"github.com/cookiengineer/systempanel/bindings/gtk4"
	"github.com/cookiengineer/systempanel/detect"
	"github.com/cookiengineer/systempanel/model"
	"github.com/cookiengineer/systempanel/parsers/desktop"
	"github.com/cookiengineer/systempanel/view"
	"github.com/cookiengineer/systempanel/widget"
)

var Descriptor = view.ViewDescriptor{
	Name:     "autostart",
	Title:    "Autostart",
	IconName: "system-run-symbolic",
	DetectFn: func() bool {
		return detect.HasProgram("systemctl") && detect.HasAutostartDir()
	},
	Factory: func() view.View { return NewAutostartView() },
}

type autostartItem struct {
	entry model.AutostartEntry
	row   *gtk4.ListBoxRow
}

type AutostartView struct {
	box       *gtk4.Box
	model     *model.AutostartModel
	listBox   *gtk4.ListBox
	items     []autostartItem
	rows      []*gtk4.ListBoxRow
	parentWin *gtk4.Window
}

func NewAutostartView() *AutostartView {
	av := &AutostartView{
		box:   gtk4.BoxNew(gtk4.OrientationVertical, 12),
		model: &model.AutostartModel{},
	}
	av.box.SetMarginStart(24)
	av.box.SetMarginEnd(24)
	av.box.SetMarginTop(24)
	av.box.SetMarginBottom(24)

	header := gtk4.LabelNew("Autostart Applications")
	header.AddCSSClass("header-label")
	av.box.Append(&header.Widget)

	desc := gtk4.LabelNew("Manage applications that automatically start when you log in.")
	desc.SetWrap(true)
	desc.SetMarginBottom(12)
	av.box.Append(&desc.Widget)

	av.listBox = gtk4.ListBoxNew()
	av.listBox.SetSelectionMode(gtk4.SelectionNone)

	scrollW := gtk4.ScrolledWindowNew()
	scrollW.SetPolicy(gtk4.PolicyNever, gtk4.PolicyAutomatic)
	scrollW.SetHExpand(true)
	scrollW.SetVExpand(true)
	scrollW.SetChild(&av.listBox.Widget)
	av.box.Append(&scrollW.Widget)

	btnBox := gtk4.BoxNew(gtk4.OrientationHorizontal, 8)
	btnBox.SetMarginTop(4)

	refreshBtn := gtk4.ButtonNewWithLabel("Refresh")
	refreshBtn.OnClicked(func() { av.refresh() })
	btnBox.Append(&refreshBtn.Widget)

	spacer := gtk4.LabelNew("")
	spacer.SetHExpand(true)
	btnBox.Append(&spacer.Widget)

	createBtn := gtk4.ButtonNewWithLabel("Create")
	createBtn.AddCSSClass("suggested-action")
	createBtn.OnClicked(av.onCreate)
	btnBox.Append(&createBtn.Widget)

	av.box.Append(&btnBox.Widget)

	av.refresh()

	return av
}

func (av *AutostartView) SetParentWindow(parent *gtk4.Window) { av.parentWin = parent }

func (av *AutostartView) refresh() {
	for _, r := range av.rows {
		av.listBox.Remove(r)
	}
	av.rows = av.rows[:0]
	av.items = av.items[:0]

	entries, _ := av.model.ListEntries()

	for _, e := range entries {
		row := av.createRow(e)
		av.listBox.Append(row)
		av.rows = append(av.rows, row)
		av.items = append(av.items, autostartItem{entry: e, row: row})
	}
}

func (av *AutostartView) createRow(e model.AutostartEntry) *gtk4.ListBoxRow {
	row := gtk4.ListBoxRowNew()
	row.AddCSSClass("device-row")

	hbox := gtk4.BoxNew(gtk4.OrientationHorizontal, 8)

	iconName := e.Icon
	if iconName == "" {
		iconName = "application-x-executable-symbolic"
	}
	icon := gtk4.ImageNewFromIconName(iconName)
	icon.SetPixelSize(20)
	hbox.Append(&icon.Widget)

	infoBox := gtk4.BoxNew(gtk4.OrientationVertical, 2)
	nameLabel := gtk4.LabelNew(e.Name)
	nameLabel.SetHAlign(gtk4.AlignStart)
	infoBox.Append(&nameLabel.Widget)

	if e.Comment != "" {
		commentLabel := gtk4.LabelNew(e.Comment)
		commentLabel.SetSensitive(false)
		commentLabel.SetHAlign(gtk4.AlignStart)
		infoBox.Append(&commentLabel.Widget)
	}

	infoBox.SetHExpand(true)
	hbox.Append(&infoBox.Widget)

	statusLabel := gtk4.LabelNew("Enabled")
	if e.Hidden {
		statusLabel.SetText("Disabled")
		statusLabel.SetSensitive(false)
	}
	hbox.Append(&statusLabel.Widget)

	path := e.Path
	gearBtn := gtk4.ButtonNew()
	gearBtn.SetIconName("emblem-system-symbolic")
	gearBtn.OnClicked(func() {
		entry, _ := desktop.Parse(path)
		dlg := widget.NewDesktopDialog(av.parentWin, entry, path)
		dlg.Present()
	})
	hbox.Append(&gearBtn.Widget)

	row.SetChild(&hbox.Widget)

	return row
}

func (av *AutostartView) onCreate() {
	dlg := widget.NewDesktopDialog(av.parentWin, nil, "")
	dlg.Present()
}

func (av *AutostartView) Widget() *gtk4.Widget { return &av.box.Widget }
func (av *AutostartView) Name() string          { return "autostart" }
func (av *AutostartView) Title() string         { return "Autostart" }
func (av *AutostartView) IconName() string      { return "system-run-symbolic" }
func (av *AutostartView) OnShow()               { av.refresh() }
func (av *AutostartView) OnHide()               {}
func (av *AutostartView) Destroy()              {}
