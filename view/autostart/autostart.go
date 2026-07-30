package autostart

import (
	"os"
	"path/filepath"

	"github.com/cookiengineer/systempanel/bindings/gtk4"
	"github.com/cookiengineer/systempanel/detect"
	"github.com/cookiengineer/systempanel/model"
	"github.com/cookiengineer/systempanel/view"
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

type AutostartView struct {
	box     *gtk4.Box
	model   *model.AutostartModel
	listBox *gtk4.ListBox
	entries []autostartItem
}

type autostartItem struct {
	entry model.AutostartEntry
	row   *gtk4.ListBoxRow
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
	headerWidget := header.Widget
	av.box.Append(&headerWidget)

	desc := gtk4.LabelNew("Applications configured to start automatically at login.")
	desc.SetWrap(true)
	desc.SetMarginBottom(12)
	descWidget := desc.Widget
	av.box.Append(&descWidget)

	av.listBox = gtk4.ListBoxNew()
	av.listBox.SetSelectionMode(gtk4.SelectionSingle)

	scrollW := gtk4.ScrolledWindowNew()
	scrollW.SetPolicy(gtk4.PolicyNever, gtk4.PolicyAutomatic)
	scrollW.SetHExpand(true)
	scrollW.SetVExpand(true)

	scrollWidget := scrollW.Widget
	listWidget := av.listBox.Widget
	scrollW.SetChild(&listWidget)
	av.box.Append(&scrollWidget)

	refreshBtn := gtk4.ButtonNewWithLabel("Refresh")
	refreshBtn.OnClicked(func() { av.refresh() })
	refreshWidget := refreshBtn.Widget
	av.box.Append(&refreshWidget)

	av.refresh()

	return av
}

func (av *AutostartView) refresh() {
	entries, _ := av.model.ListEntries()

	for _, item := range av.entries {
		av.listBox.Remove(item.row)
	}
	av.entries = av.entries[:0]

	for _, e := range entries {
		row := av.createEntryRow(e)
		av.listBox.Append(row)
		av.entries = append(av.entries, autostartItem{entry: e, row: row})
	}
}

func (av *AutostartView) createEntryRow(e model.AutostartEntry) *gtk4.ListBoxRow {
	row := gtk4.ListBoxRowNew()
	row.AddCSSClass("device-row")

	hbox := gtk4.BoxNew(gtk4.OrientationHorizontal, 8)

	iconName := e.Icon
	if iconName == "" {
		iconName = "application-x-executable-symbolic"
	}
	icon := gtk4.ImageNewFromIconName(iconName)
	icon.SetPixelSize(24)
	iconWidget := icon.Widget
	hbox.Append(&iconWidget)

	infoBox := gtk4.BoxNew(gtk4.OrientationVertical, 2)

	nameLabel := gtk4.LabelNew(e.Name)
	nameLabel.SetHAlign(gtk4.AlignStart)
	nlWidget := nameLabel.Widget
	infoBox.Append(&nlWidget)

	pathLabel := gtk4.LabelNew(filepath.Base(e.Exec))
	pathLabel.SetSensitive(false)
	pathLabel.SetHAlign(gtk4.AlignStart)
	plWidget := pathLabel.Widget
	infoBox.Append(&plWidget)

	infoBoxWidget := infoBox.Widget
	infoBox.SetHExpand(true)
	hbox.Append(&infoBoxWidget)

	toggle := gtk4.SwitchNew()
	toggle.SetActive(!e.Hidden)
	entry := e
	toggle.OnActivate(func() {
		if toggle.GetActive() {
			av.model.Enable(entry.Path)
		} else {
			av.model.Disable(entry.Path)
		}
	})
	toggleWidget := toggle.Widget
	hbox.Append(&toggleWidget)

	hboxWidget := hbox.Widget
	row.SetChild(&hboxWidget)

	return row
}

func autostartDirs() []string {
	var dirs []string
	configDir := os.Getenv("XDG_CONFIG_HOME")
	if configDir == "" {
		home, err := os.UserHomeDir()
		if err == nil {
			configDir = filepath.Join(home, ".config")
		}
	}
	if configDir != "" {
		dirs = append(dirs, filepath.Join(configDir, "autostart"))
	}
	dirs = append(dirs, "/etc/xdg/autostart")
	return dirs
}

func (av *AutostartView) Widget() *gtk4.Widget { return &av.box.Widget }
func (av *AutostartView) Name() string          { return "autostart" }
func (av *AutostartView) Title() string         { return "Autostart" }
func (av *AutostartView) IconName() string      { return "system-run-symbolic" }
func (av *AutostartView) OnShow()               { av.refresh() }
func (av *AutostartView) OnHide()               {}
func (av *AutostartView) Destroy()              {}
