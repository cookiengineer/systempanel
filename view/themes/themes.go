package themes

import (
	"unsafe"

	"github.com/cookiengineer/systempanel/bindings/gtk4"
	"github.com/cookiengineer/systempanel/detect"
	"github.com/cookiengineer/systempanel/model"
	"github.com/cookiengineer/systempanel/view"
)

var Descriptor = view.ViewDescriptor{
	Name:     "themes",
	Title:    "Themes",
	IconName: "preferences-desktop-theme-symbolic",
	DetectFn: func() bool { return detect.HasThemes() },
	Factory:  func() view.View { return NewThemesView() },
}

type themeItem struct {
	theme model.Theme
	row   *gtk4.ListBoxRow
}

type ThemesView struct {
	box          *gtk4.Box
	model        *model.ThemeModel
	listBox      *gtk4.ListBox
	rowPtrToIdx  map[unsafe.Pointer]int
	items        []themeItem
	rows         []*gtk4.ListBoxRow
	currentTheme string
}

func NewThemesView() *ThemesView {
	tv := &ThemesView{
		box:         gtk4.BoxNew(gtk4.OrientationVertical, 12),
		model:       &model.ThemeModel{},
		rowPtrToIdx: make(map[unsafe.Pointer]int),
	}
	tv.box.SetMarginStart(24)
	tv.box.SetMarginEnd(24)
	tv.box.SetMarginTop(24)
	tv.box.SetMarginBottom(24)

	header := gtk4.LabelNew("GTK Themes")
	header.AddCSSClass("header-label")
	tv.box.Append(&header.Widget)

	desc := gtk4.LabelNew("Select a theme to apply it immediately.")
	desc.SetWrap(true)
	desc.SetMarginBottom(12)
	tv.box.Append(&desc.Widget)

	tv.listBox = gtk4.ListBoxNew()
	tv.listBox.SetSelectionMode(gtk4.SelectionSingle)
	tv.listBox.OnRowActivated(tv.onRowSelected)

	scrollW := gtk4.ScrolledWindowNew()
	scrollW.SetPolicy(gtk4.PolicyNever, gtk4.PolicyAutomatic)
	scrollW.SetHExpand(true)
	scrollW.SetVExpand(true)
	scrollW.SetChild(&tv.listBox.Widget)
	tv.box.Append(&scrollW.Widget)

	btnBox := gtk4.BoxNew(gtk4.OrientationHorizontal, 8)
	btnBox.SetMarginTop(4)

	refreshBtn := gtk4.ButtonNewWithLabel("Refresh")
	refreshBtn.OnClicked(func() { tv.refresh() })
	btnBox.Append(&refreshBtn.Widget)

	tv.box.Append(&btnBox.Widget)

	tv.refresh()

	return tv
}

func (tv *ThemesView) refresh() {
	for _, r := range tv.rows {
		tv.listBox.Remove(r)
	}
	tv.rows = tv.rows[:0]
	tv.items = tv.items[:0]
	tv.rowPtrToIdx = make(map[unsafe.Pointer]int)

	themes, _ := tv.model.ListThemes()
	tv.currentTheme = tv.model.CurrentTheme()

	for _, t := range themes {
		row := tv.createThemeRow(t)
		tv.listBox.Append(row)
		tv.rows = append(tv.rows, row)
		tv.items = append(tv.items, themeItem{theme: t, row: row})
		tv.rowPtrToIdx[row.Widget.Ptr()] = len(tv.items) - 1
		if t.Name == tv.currentTheme {
			tv.listBox.SelectRow(row)
		}
	}
}

func (tv *ThemesView) createThemeRow(t model.Theme) *gtk4.ListBoxRow {
	row := gtk4.ListBoxRowNew()
	row.AddCSSClass("device-row")

	hbox := gtk4.BoxNew(gtk4.OrientationHorizontal, 8)

	nameLabel := gtk4.LabelNew(t.Name)
	nameLabel.SetHExpand(true)
	nameLabel.SetHAlign(gtk4.AlignStart)
	hbox.Append(&nameLabel.Widget)

	verLabel := gtk4.LabelNew("GTK4")
	if !t.IsGTK4 {
		verLabel.SetText("GTK3")
	}
	verLabel.SetSensitive(false)
	hbox.Append(&verLabel.Widget)

	hboxWidget := hbox.Widget
	row.SetChild(&hboxWidget)

	return row
}

func (tv *ThemesView) onRowSelected(row *gtk4.ListBoxRow) {
	idx, ok := tv.rowPtrToIdx[row.Widget.Ptr()]
	if !ok || idx >= len(tv.items) {
		return
	}
	theme := tv.items[idx].theme
	tv.currentTheme = theme.Name
	tv.model.ApplyTheme(theme.Name)
}

func (tv *ThemesView) Widget() *gtk4.Widget { return &tv.box.Widget }
func (tv *ThemesView) Name() string          { return "themes" }
func (tv *ThemesView) Title() string         { return "Themes" }
func (tv *ThemesView) IconName() string      { return "preferences-desktop-theme-symbolic" }
func (tv *ThemesView) OnShow()               { tv.refresh() }
func (tv *ThemesView) OnHide()               {}
func (tv *ThemesView) Destroy()              {}
