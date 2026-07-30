package services

import (
	"unsafe"

	"github.com/cookiengineer/systempanel/bindings/gtk4"
	"github.com/cookiengineer/systempanel/detect"
	"github.com/cookiengineer/systempanel/model"
	"github.com/cookiengineer/systempanel/view"
	"github.com/cookiengineer/systempanel/widget"
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

type serviceItem struct {
	service model.ServiceUnit
	row     *gtk4.ListBoxRow
}

type ServicesView struct {
	box         *gtk4.Box
	model       *model.ServicesModel
	listBox     *gtk4.ListBox
	rowPtrToIdx map[unsafe.Pointer]int
	items       []serviceItem
	rows        []*gtk4.ListBoxRow
	parentWin   *gtk4.Window
}

func NewServicesView() *ServicesView {
	sv := &ServicesView{
		box:         gtk4.BoxNew(gtk4.OrientationVertical, 12),
		model:       &model.ServicesModel{},
		rowPtrToIdx: make(map[unsafe.Pointer]int),
	}
	sv.box.SetMarginStart(24)
	sv.box.SetMarginEnd(24)
	sv.box.SetMarginTop(24)
	sv.box.SetMarginBottom(24)

	header := gtk4.LabelNew("Services")
	header.AddCSSClass("header-label")
	sv.box.Append(&header.Widget)

	desc := gtk4.LabelNew("View and manage systemd service units and their status.")
	desc.SetWrap(true)
	desc.SetMarginBottom(12)
	dv := desc.Widget
	sv.box.Append(&dv)

	sv.listBox = gtk4.ListBoxNew()
	sv.listBox.SetSelectionMode(gtk4.SelectionSingle)
	sv.listBox.OnRowActivated(sv.onRowSelected)

	scrollW := gtk4.ScrolledWindowNew()
	scrollW.SetPolicy(gtk4.PolicyNever, gtk4.PolicyAutomatic)
	scrollW.SetHExpand(true)
	scrollW.SetVExpand(true)
	scrollW.SetChild(&sv.listBox.Widget)
	sv.box.Append(&scrollW.Widget)

	btnBox := gtk4.BoxNew(gtk4.OrientationHorizontal, 8)
	btnBox.SetMarginTop(4)

	refreshBtn := gtk4.ButtonNewWithLabel("Refresh")
	refreshBtn.OnClicked(func() { sv.refresh() })
	btnBox.Append(&refreshBtn.Widget)

	spacer := gtk4.LabelNew("")
	spacer.SetHExpand(true)
	btnBox.Append(&spacer.Widget)

	btnBoxWidget := btnBox.Widget
	sv.box.Append(&btnBoxWidget)

	sv.refresh()

	return sv
}

func (sv *ServicesView) SetParentWindow(parent *gtk4.Window) { sv.parentWin = parent }

func (sv *ServicesView) refresh() {
	for _, r := range sv.rows {
		sv.listBox.Remove(r)
	}
	sv.rows = sv.rows[:0]
	sv.items = sv.items[:0]
	sv.rowPtrToIdx = make(map[unsafe.Pointer]int)

	go func() {
		userUnits, _ := sv.model.ListUnits(true)
		systemUnits, _ := sv.model.ListUnits(false)

		gtk4.IdleAdd(func() {
			hasUser := false
			for _, u := range userUnits {
				if u.ActiveState != "inactive" || u.LoadState != "not-found" {
					hasUser = true
					break
				}
			}
			if hasUser {
				hdr := sv.createSectionLabel("─ User Services ─")
				sv.listBox.Append(hdr)
				sv.rows = append(sv.rows, hdr)
				for i := range userUnits {
					row := sv.createServiceRow(&userUnits[i])
					sv.listBox.Append(row)
					sv.rows = append(sv.rows, row)
					sv.items = append(sv.items, serviceItem{service: userUnits[i], row: row})
					sv.rowPtrToIdx[row.Widget.Ptr()] = len(sv.items) - 1
				}
			}

			hdr := sv.createSectionLabel("─ System Services ─")
			sv.listBox.Append(hdr)
			sv.rows = append(sv.rows, hdr)
			for i := range systemUnits {
				row := sv.createServiceRow(&systemUnits[i])
				sv.listBox.Append(row)
				sv.rows = append(sv.rows, row)
				sv.items = append(sv.items, serviceItem{service: systemUnits[i], row: row})
				sv.rowPtrToIdx[row.Widget.Ptr()] = len(sv.items) - 1
			}
		})
	}()
}

func (sv *ServicesView) createSectionLabel(text string) *gtk4.ListBoxRow {
	row := gtk4.ListBoxRowNew()
	row.SetSensitive(false)
	label := gtk4.LabelNew(text)
	label.SetHAlign(gtk4.AlignCenter)
	label.SetMarginTop(8)
	label.SetMarginBottom(4)
	label.SetSensitive(false)
	row.SetChild(&label.Widget)
	return row
}

func (sv *ServicesView) createServiceRow(s *model.ServiceUnit) *gtk4.ListBoxRow {
	row := gtk4.ListBoxRowNew()
	row.AddCSSClass("device-row")

	hbox := gtk4.BoxNew(gtk4.OrientationHorizontal, 8)

	statusIcon, statusText := serviceState(s)
	icon := gtk4.ImageNewFromIconName(statusIcon)
	icon.SetPixelSize(20)
	hbox.Append(&icon.Widget)

	nameLabel := gtk4.LabelNew(s.Name)
	nameLabel.SetHExpand(true)
	nameLabel.SetHAlign(gtk4.AlignStart)
	hbox.Append(&nameLabel.Widget)

	stateLabel := gtk4.LabelNew(statusText)
	stateLabel.SetSensitive(false)
	hbox.Append(&stateLabel.Widget)

	unitName := s.Name
	isUser := s.IsUser
	gearBtn := gtk4.ButtonNew()
	gearBtn.SetIconName("emblem-system-symbolic")
	gearBtn.OnClicked(func() {
		uf, _ := sv.model.LoadUnitFile(unitName, isUser)
		dlg := widget.NewServiceDialog(sv.parentWin, uf, unitName, isUser)
		dlg.Present()
	})
	hbox.Append(&gearBtn.Widget)

	hboxWidget := hbox.Widget
	row.SetChild(&hboxWidget)

	return row
}

func (sv *ServicesView) onRowSelected(row *gtk4.ListBoxRow) {
	// selection tracks which service to toggle
}

func (sv *ServicesView) onAction() {
	selected := sv.listBox.GetSelectedRow()
	if selected == nil {
		return
	}
	idx, ok := sv.rowPtrToIdx[selected.Widget.Ptr()]
	if !ok || idx >= len(sv.items) {
		return
	}
	s := sv.items[idx].service
	if s.ActiveState == "active" {
		sv.model.Stop(s.Name)
	} else {
		sv.model.Start(s.Name)
	}
	sv.refresh()
}

func serviceState(s *model.ServiceUnit) (string, string) {
	switch s.ActiveState {
	case "active":
		return "media-playback-start-symbolic", "Running"
	case "failed":
		return "dialog-error-symbolic", "Failed"
	case "inactive":
		if s.SubState == "dead" {
			return "media-playback-stop-symbolic", "Stopped"
		}
		return "media-playback-pause-symbolic", "Inactive"
	default:
		return "dialog-question-symbolic", s.ActiveState
	}
}

func (sv *ServicesView) Widget() *gtk4.Widget { return &sv.box.Widget }
func (sv *ServicesView) Name() string          { return "services" }
func (sv *ServicesView) Title() string         { return "Services" }
func (sv *ServicesView) IconName() string      { return "applications-system-symbolic" }
func (sv *ServicesView) OnShow()               { sv.refresh() }
func (sv *ServicesView) OnHide()               {}
func (sv *ServicesView) Destroy()              {}
