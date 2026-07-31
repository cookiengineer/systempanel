package powerprofile

import (
	"unsafe"

	"github.com/cookiengineer/systempanel/bindings/gtk4"
	"github.com/cookiengineer/systempanel/detect"
	"github.com/cookiengineer/systempanel/model"
	"github.com/cookiengineer/systempanel/view"
)

var Descriptor = view.ViewDescriptor{
	Name:     "powerprofile",
	Title:    "Power Profile",
	IconName: "power-profile-performance-symbolic",
	DetectFn: func() bool { return detect.HasProgram("powerprofilesctl") },
	Factory:  func() view.View { return NewPowerProfileView() },
}

type profileItem struct {
	profile model.PowerProfile
	row     *gtk4.ListBoxRow
}

type PowerProfileView struct {
	box              *gtk4.Box
	model            *model.PowerProfileModel
	listBox          *gtk4.ListBox
	rowPtrToIdx      map[unsafe.Pointer]int
	items            []profileItem
	rows             []*gtk4.ListBoxRow
	refreshInProgress bool
}

func NewPowerProfileView() *PowerProfileView {
	pv := &PowerProfileView{
		box:         gtk4.BoxNew(gtk4.OrientationVertical, 12),
		model:       &model.PowerProfileModel{},
		rowPtrToIdx: make(map[unsafe.Pointer]int),
	}
	pv.box.SetMarginStart(24)
	pv.box.SetMarginEnd(24)
	pv.box.SetMarginTop(24)
	pv.box.SetMarginBottom(24)

	header := gtk4.LabelNew("Power Profiles")
	header.AddCSSClass("header-label")
	pv.box.Append(&header.Widget)

	desc := gtk4.LabelNew("Select a power profile to apply it immediately.")
	desc.SetWrap(true)
	desc.SetMarginBottom(12)
	pv.box.Append(&desc.Widget)

	pv.listBox = gtk4.ListBoxNew()
	pv.listBox.SetSelectionMode(gtk4.SelectionSingle)
	pv.listBox.OnRowActivated(pv.onRowSelected)

	scrollW := gtk4.ScrolledWindowNew()
	scrollW.SetPolicy(gtk4.PolicyNever, gtk4.PolicyAutomatic)
	scrollW.SetHExpand(true)
	scrollW.SetVExpand(true)
	scrollW.SetChild(&pv.listBox.Widget)
	pv.box.Append(&scrollW.Widget)

	btnBox := gtk4.BoxNew(gtk4.OrientationHorizontal, 8)
	btnBox.SetMarginTop(4)

	refreshBtn := gtk4.ButtonNewWithLabel("Refresh")
	refreshBtn.OnClicked(func() { go pv.refreshAsync() })
	btnBox.Append(&refreshBtn.Widget)

	pv.box.Append(&btnBox.Widget)

	go pv.refreshAsync()

	return pv
}

func (pv *PowerProfileView) refreshAsync() {
	profiles, err := pv.model.ListProfiles()
	if err != nil {
		return
	}
	gtk4.IdleAdd(func() {
		pv.populate(profiles)
	})
}

func (pv *PowerProfileView) populate(profiles []model.PowerProfile) {
	for _, r := range pv.rows {
		pv.listBox.Remove(r)
	}
	pv.rows = pv.rows[:0]
	pv.items = pv.items[:0]
	pv.rowPtrToIdx = make(map[unsafe.Pointer]int)

	for _, p := range profiles {
		row := pv.createProfileRow(p)
		pv.listBox.Append(row)
		pv.rows = append(pv.rows, row)
		pv.items = append(pv.items, profileItem{profile: p, row: row})
		pv.rowPtrToIdx[row.Widget.Ptr()] = len(pv.items) - 1
		if p.Active {
			pv.refreshInProgress = true
			pv.listBox.SelectRow(row)
			pv.refreshInProgress = false
		}
	}
}

func (pv *PowerProfileView) createProfileRow(p model.PowerProfile) *gtk4.ListBoxRow {
	row := gtk4.ListBoxRowNew()
	row.AddCSSClass("device-row")

	hbox := gtk4.BoxNew(gtk4.OrientationHorizontal, 8)

	nameLabel := gtk4.LabelNew(p.Name)
	nameLabel.SetHExpand(true)
	nameLabel.SetHAlign(gtk4.AlignStart)
	hbox.Append(&nameLabel.Widget)

	infoLabel := gtk4.LabelNew(p.Driver)
	if p.Degraded != "" && p.Degraded != "no" {
		infoLabel.SetText(p.Driver + " (degraded)")
	}
	infoLabel.SetSensitive(false)
	hbox.Append(&infoLabel.Widget)

	hboxWidget := hbox.Widget
	row.SetChild(&hboxWidget)

	return row
}

func (pv *PowerProfileView) onRowSelected(row *gtk4.ListBoxRow) {
	if pv.refreshInProgress {
		return
	}
	idx, ok := pv.rowPtrToIdx[row.Widget.Ptr()]
	if !ok || idx >= len(pv.items) {
		return
	}
	if pv.items[idx].profile.Active {
		return
	}
	pv.model.SetProfile(pv.items[idx].profile.Name)
	for i := range pv.items {
		pv.items[i].profile.Active = (i == idx)
	}
}

func (pv *PowerProfileView) Widget() *gtk4.Widget { return &pv.box.Widget }
func (pv *PowerProfileView) Name() string          { return "powerprofile" }
func (pv *PowerProfileView) Title() string         { return "Power Profile" }
func (pv *PowerProfileView) IconName() string      { return "power-profile-performance-symbolic" }
func (pv *PowerProfileView) OnShow()               { go pv.refreshAsync() }
func (pv *PowerProfileView) OnHide()               {}
func (pv *PowerProfileView) Destroy()              {}
