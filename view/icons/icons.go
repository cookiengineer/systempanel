package icons

import (
	"unsafe"

	"github.com/cookiengineer/systempanel/bindings/gtk4"
	"github.com/cookiengineer/systempanel/detect"
	"github.com/cookiengineer/systempanel/model"
	"github.com/cookiengineer/systempanel/view"
)

var Descriptor = view.ViewDescriptor{
	Name:     "icons",
	Title:    "Icons",
	IconName: "application-x-theme-symbolic",
	DetectFn: func() bool { return detect.HasIcons() },
	Factory:  func() view.View { return NewIconsView() },
}

type iconItem struct {
	icon model.Icon
	row  *gtk4.ListBoxRow
}

type IconsView struct {
	box         *gtk4.Box
	model       *model.IconModel
	listBox     *gtk4.ListBox
	rowPtrToIdx map[unsafe.Pointer]int
	items       []iconItem
	rows        []*gtk4.ListBoxRow
}

func NewIconsView() *IconsView {
	iv := &IconsView{
		box:         gtk4.BoxNew(gtk4.OrientationVertical, 12),
		model:       &model.IconModel{},
		rowPtrToIdx: make(map[unsafe.Pointer]int),
	}
	iv.box.SetMarginStart(24)
	iv.box.SetMarginEnd(24)
	iv.box.SetMarginTop(24)
	iv.box.SetMarginBottom(24)

	header := gtk4.LabelNew("Icon Themes")
	header.AddCSSClass("header-label")
	iv.box.Append(&header.Widget)

	desc := gtk4.LabelNew("Select an icon theme to apply it immediately.")
	desc.SetWrap(true)
	desc.SetMarginBottom(12)
	iv.box.Append(&desc.Widget)

	iv.listBox = gtk4.ListBoxNew()
	iv.listBox.SetSelectionMode(gtk4.SelectionSingle)
	iv.listBox.OnRowActivated(iv.onRowSelected)

	scrollW := gtk4.ScrolledWindowNew()
	scrollW.SetPolicy(gtk4.PolicyNever, gtk4.PolicyAutomatic)
	scrollW.SetHExpand(true)
	scrollW.SetVExpand(true)
	scrollW.SetChild(&iv.listBox.Widget)
	iv.box.Append(&scrollW.Widget)

	btnBox := gtk4.BoxNew(gtk4.OrientationHorizontal, 8)
	btnBox.SetMarginTop(4)

	refreshBtn := gtk4.ButtonNewWithLabel("Refresh")
	refreshBtn.OnClicked(func() { iv.refresh() })
	btnBox.Append(&refreshBtn.Widget)

	iv.box.Append(&btnBox.Widget)

	iv.refresh()

	return iv
}

func (iv *IconsView) refresh() {
	for _, r := range iv.rows {
		iv.listBox.Remove(r)
	}
	iv.rows = iv.rows[:0]
	iv.items = iv.items[:0]
	iv.rowPtrToIdx = make(map[unsafe.Pointer]int)

	icons, _ := iv.model.ListIcons()
	current := iv.model.CurrentIcon()

	for _, ic := range icons {
		row := iv.createIconRow(ic, current)
		iv.listBox.Append(row)
		iv.rows = append(iv.rows, row)
		iv.items = append(iv.items, iconItem{icon: ic, row: row})
		iv.rowPtrToIdx[row.Widget.Ptr()] = len(iv.items) - 1
	}
}

func (iv *IconsView) createIconRow(ic model.Icon, current string) *gtk4.ListBoxRow {
	row := gtk4.ListBoxRowNew()
	row.AddCSSClass("device-row")

	hbox := gtk4.BoxNew(gtk4.OrientationHorizontal, 8)

	iconName := "application-x-theme-symbolic"
	if ic.Name == current {
		iconName = "object-select-symbolic"
	}
	icon := gtk4.ImageNewFromIconName(iconName)
	icon.SetPixelSize(20)
	hbox.Append(&icon.Widget)

	nameLabel := gtk4.LabelNew(ic.Name)
	nameLabel.SetHExpand(true)
	nameLabel.SetHAlign(gtk4.AlignStart)
	hbox.Append(&nameLabel.Widget)

	hboxWidget := hbox.Widget
	row.SetChild(&hboxWidget)

	return row
}

func (iv *IconsView) onRowSelected(row *gtk4.ListBoxRow) {
	idx, ok := iv.rowPtrToIdx[row.Widget.Ptr()]
	if !ok || idx >= len(iv.items) {
		return
	}
	ic := iv.items[idx].icon
	iv.model.ApplyIcon(ic.Name)
}

func (iv *IconsView) Widget() *gtk4.Widget { return &iv.box.Widget }
func (iv *IconsView) Name() string          { return "icons" }
func (iv *IconsView) Title() string         { return "Icons" }
func (iv *IconsView) IconName() string      { return "application-x-theme-symbolic" }
func (iv *IconsView) OnShow()               { iv.refresh() }
func (iv *IconsView) OnHide()               {}
func (iv *IconsView) Destroy()              {}
