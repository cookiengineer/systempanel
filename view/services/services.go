package services

import (
	"github.com/cookiengineer/systempanel/bindings/gtk4"
	"github.com/cookiengineer/systempanel/detect"
	"github.com/cookiengineer/systempanel/model"
	"github.com/cookiengineer/systempanel/view"
)

var Descriptor = view.ViewDescriptor{
	Name:     "services",
	Title:    "Services",
	IconName: "applications-system-symbolic",
	DetectFn: func() bool {
		return detect.HasProgram("systemctl") && detect.HasSystemd()
	},
	Factory: func() view.View { return NewServicesView() },
}

type ServicesView struct {
	box     *gtk4.Box
	model   *model.ServicesModel
	listBox *gtk4.ListBox
	units   []unitItem
	filter  string
}

type unitItem struct {
	unit model.SystemdUnit
	row  *gtk4.ListBoxRow
}

func NewServicesView() *ServicesView {
	sv := &ServicesView{
		box:   gtk4.BoxNew(gtk4.OrientationVertical, 12),
		model: &model.ServicesModel{},
	}
	sv.box.SetMarginStart(24)
	sv.box.SetMarginEnd(24)
	sv.box.SetMarginTop(24)
	sv.box.SetMarginBottom(24)

	header := gtk4.LabelNew("Systemd Services")
	header.AddCSSClass("header-label")
	headerWidget := header.Widget
	sv.box.Append(&headerWidget)

	filterBox := gtk4.BoxNew(gtk4.OrientationHorizontal, 8)
	filterBox.SetMarginBottom(8)

	filterLabel := gtk4.LabelNew("Filter:")
	flWidget := filterLabel.Widget
	filterBox.Append(&flWidget)

	filterEntry := gtk4.EntryNew()
	filterEntry.SetPlaceholder("Search units...")
	filterEntry.OnChanged(func() {
		sv.filter = filterEntry.GetText()
		sv.refreshList()
	})
	filterEntry.SetHExpand(true)
	feWidget := filterEntry.Widget
	filterBox.Append(&feWidget)

	filterBoxWidget := filterBox.Widget
	sv.box.Append(&filterBoxWidget)

	sv.listBox = gtk4.ListBoxNew()
	sv.listBox.SetSelectionMode(gtk4.SelectionSingle)
	sv.listBox.OnRowActivated(func(row *gtk4.ListBoxRow) {
		for _, item := range sv.units {
			if item.row == row {
				if item.unit.ActiveState == "active" {
					sv.model.Stop(item.unit.Name)
				} else {
					sv.model.Start(item.unit.Name)
				}
				return
			}
		}
	})

	scrollW := gtk4.ScrolledWindowNew()
	scrollW.SetPolicy(gtk4.PolicyNever, gtk4.PolicyAutomatic)
	scrollW.SetHExpand(true)
	scrollW.SetVExpand(true)

	scrollWidget := scrollW.Widget
	listWidget := sv.listBox.Widget
	scrollW.SetChild(&listWidget)
	sv.box.Append(&scrollWidget)

	refreshBtn := gtk4.ButtonNewWithLabel("Refresh")
	refreshBtn.OnClicked(func() { sv.refresh() })
	refreshWidget := refreshBtn.Widget
	sv.box.Append(&refreshWidget)

	sv.refresh()

	return sv
}

func (sv *ServicesView) refresh() {
	units, _ := sv.model.ListUnits()
	sv.units = sv.units[:0]
	for _, u := range units {
		sv.units = append(sv.units, unitItem{unit: u})
	}
	sv.refreshList()
}

func (sv *ServicesView) refreshList() {
	for _, item := range sv.units {
		if item.row != nil {
			sv.listBox.Remove(item.row)
		}
	}

	for i := range sv.units {
		u := sv.units[i].unit
		if sv.filter != "" && !contains(u.Name, sv.filter) && !contains(u.Description, sv.filter) {
			continue
		}
		row := sv.createUnitRow(u)
		sv.listBox.Append(row)
		sv.units[i].row = row
	}
}

func (sv *ServicesView) createUnitRow(u model.SystemdUnit) *gtk4.ListBoxRow {
	row := gtk4.ListBoxRowNew()
	row.AddCSSClass("device-row")

	hbox := gtk4.BoxNew(gtk4.OrientationHorizontal, 8)

	status := ""
	stateClass := ""
	switch u.ActiveState {
	case "active":
		status = "●"
		stateClass = "service-active"
	case "failed":
		status = "✗"
		stateClass = "service-failed"
	default:
		status = "○"
		stateClass = "service-inactive"
	}

	statusLabel := gtk4.LabelNew(status)
	statusLabel.AddCSSClass(stateClass)
	slWidget := statusLabel.Widget
	hbox.Append(&slWidget)

	nameLabel := gtk4.LabelNew(u.Name)
	nameLabel.SetHExpand(true)
	nameLabel.SetHAlign(gtk4.AlignStart)
	nlWidget := nameLabel.Widget
	hbox.Append(&nlWidget)

	descLabel := gtk4.LabelNew(u.Description)
	descLabel.SetSensitive(false)
	descLabel.SetHAlign(gtk4.AlignEnd)
	dlWidget := descLabel.Widget
	hbox.Append(&dlWidget)

	hboxWidget := hbox.Widget
	row.SetChild(&hboxWidget)

	return row
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchSubstring(s, substr)
}

func searchSubstring(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}

func (sv *ServicesView) Widget() *gtk4.Widget { return &sv.box.Widget }
func (sv *ServicesView) Name() string          { return "services" }
func (sv *ServicesView) Title() string         { return "Services" }
func (sv *ServicesView) IconName() string      { return "applications-system-symbolic" }
func (sv *ServicesView) OnShow()               { sv.refresh() }
func (sv *ServicesView) OnHide()               {}
func (sv *ServicesView) Destroy()              {}
