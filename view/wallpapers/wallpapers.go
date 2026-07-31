package wallpapers

import (
	"unsafe"

	"github.com/cookiengineer/systempanel/bindings/gtk4"
	"github.com/cookiengineer/systempanel/detect"
	"github.com/cookiengineer/systempanel/model"
	"github.com/cookiengineer/systempanel/view"
)

var Descriptor = view.ViewDescriptor{
	Name:     "wallpapers",
	Title:    "Wallpapers",
	IconName: "preferences-desktop-wallpaper-symbolic",
	DetectFn: func() bool {
		return detect.HasWallpaperTool()
	},
	Factory: func() view.View { return NewWallpapersView() },
}

type wallpaperItem struct {
	wallpaper model.Wallpaper
	row       *gtk4.ListBoxRow
}

type WallpapersView struct {
	box         *gtk4.Box
	model       *model.WallpaperModel
	listBox     *gtk4.ListBox
	modeCombo   *gtk4.ComboBoxText
	saveBtn     *gtk4.Button
	items       []wallpaperItem
	rowPtrToIdx map[unsafe.Pointer]int
	rows        []*gtk4.ListBoxRow
	selected    string
}

func NewWallpapersView() *WallpapersView {
	wv := &WallpapersView{
		box:         gtk4.BoxNew(gtk4.OrientationVertical, 12),
		model:       model.NewWallpaperModel(),
		rowPtrToIdx: make(map[unsafe.Pointer]int),
	}
	wv.box.SetMarginStart(24)
	wv.box.SetMarginEnd(24)
	wv.box.SetMarginTop(24)
	wv.box.SetMarginBottom(24)

	header := gtk4.LabelNew("Wallpapers")
	header.AddCSSClass("header-label")
	hw := header.Widget
	wv.box.Append(&hw)

	desc := gtk4.LabelNew("Browse and apply wallpapers from ~/Pictures/Wallpapers.")
	desc.SetWrap(true)
	desc.SetMarginBottom(12)
	dw := desc.Widget
	wv.box.Append(&dw)

	wv.listBox = gtk4.ListBoxNew()
	wv.listBox.SetSelectionMode(gtk4.SelectionSingle)
	wv.listBox.OnRowActivated(wv.onRowSelected)

	scrollW := gtk4.ScrolledWindowNew()
	scrollW.SetPolicy(gtk4.PolicyNever, gtk4.PolicyAutomatic)
	scrollW.SetHExpand(true)
	scrollW.SetVExpand(true)
	swWidget := scrollW.Widget
	lWidget := wv.listBox.Widget
	scrollW.SetChild(&lWidget)
	wv.box.Append(&swWidget)

	btnBox := gtk4.BoxNew(gtk4.OrientationHorizontal, 8)
	btnBox.SetMarginTop(4)

	refreshBtn := gtk4.ButtonNewWithLabel("Refresh")
	refreshBtn.OnClicked(func() { wv.refresh() })
	rbw := refreshBtn.Widget
	btnBox.Append(&rbw)

	spacer := gtk4.LabelNew("")
	spacer.SetHExpand(true)
	spacerWidget := spacer.Widget
	btnBox.Append(&spacerWidget)

	modeLabel := gtk4.LabelNew("Mode:")
	mlWidget := modeLabel.Widget
	btnBox.Append(&mlWidget)

	wv.modeCombo = gtk4.ComboBoxTextNew()
	for _, mode := range wv.model.GetModes() {
		wv.modeCombo.Append(mode, mode)
	}
	wv.modeCombo.SetActive(0)
	modeWidget := wv.modeCombo.Widget
	btnBox.Append(&modeWidget)

	wv.saveBtn = gtk4.ButtonNewWithLabel("Save")
	wv.saveBtn.AddCSSClass("suggested-action")
	wv.saveBtn.SetSensitive(false)
	wv.saveBtn.OnClicked(wv.onSaveClicked)
	sbWidget := wv.saveBtn.Widget
	btnBox.Append(&sbWidget)

	bbWidget := btnBox.Widget
	wv.box.Append(&bbWidget)

	wv.refresh()

	return wv
}

func (wv *WallpapersView) SetParentWindow(parent *gtk4.Window) {}

func (wv *WallpapersView) onRowSelected(row *gtk4.ListBoxRow) {
	wv.selected = ""
	wv.saveBtn.SetSensitive(false)

	idx, ok := wv.rowPtrToIdx[row.Widget.Ptr()]
	if ok && idx < len(wv.items) {
		wv.selected = wv.items[idx].wallpaper.Path
		wv.saveBtn.SetSensitive(true)
	}
}

func (wv *WallpapersView) onSaveClicked() {
	if wv.selected == "" {
		return
	}
	mode := wv.modeCombo.GetActiveID()
	wv.model.SetWallpaper(wv.selected, mode)
}

func (wv *WallpapersView) refresh() {
	wv.selected = ""
	wv.saveBtn.SetSensitive(false)

	for _, r := range wv.rows {
		wv.listBox.Remove(r)
	}
	wv.rows = wv.rows[:0]
	wv.items = wv.items[:0]
	wv.rowPtrToIdx = make(map[unsafe.Pointer]int)

	if !wv.model.HasWallpaperDir() {
		wv.showPlaceholder("~/Pictures/Wallpapers directory does not exist.\nCreate it and place some images there.")
		return
	}

	wallpapers, _ := wv.model.ListWallpapers()
	if len(wallpapers) == 0 {
		wv.showPlaceholder("No images found in ~/Pictures/Wallpapers.\nPlace some images there to get started.")
		return
	}

	for _, wp := range wallpapers {
		row := wv.createRow(wp)
		wv.listBox.Append(row)
		wv.rows = append(wv.rows, row)
		wv.rowPtrToIdx[row.Widget.Ptr()] = len(wv.items)
		wv.items = append(wv.items, wallpaperItem{wallpaper: wp, row: row})
	}
}

func (wv *WallpapersView) showPlaceholder(msg string) {
	row := gtk4.ListBoxRowNew()
	row.SetSensitive(false)
	label := gtk4.LabelNew(msg)
	label.SetHAlign(gtk4.AlignCenter)
	label.SetMarginTop(24)
	label.SetMarginBottom(24)
	label.SetSensitive(false)
	row.SetChild(&label.Widget)
	wv.listBox.Append(row)
	wv.rows = append(wv.rows, row)
}

func (wv *WallpapersView) createRow(wp model.Wallpaper) *gtk4.ListBoxRow {
	row := gtk4.ListBoxRowNew()
	row.AddCSSClass("device-row")

	hbox := gtk4.BoxNew(gtk4.OrientationHorizontal, 12)

	img := gtk4.ImageNewFromFile(wp.Path)
	img.SetPixelSize(64)
	imgWidget := img.Widget
	hbox.Append(&imgWidget)

	nameLabel := gtk4.LabelNew(wp.Filename)
	nameLabel.SetHAlign(gtk4.AlignStart)
	nameLabel.SetHExpand(true)
	nlWidget := nameLabel.Widget
	hbox.Append(&nlWidget)

	row.SetChild(&hbox.Widget)
	return row
}

func (wv *WallpapersView) Widget() *gtk4.Widget  { return &wv.box.Widget }
func (wv *WallpapersView) Name() string           { return "wallpapers" }
func (wv *WallpapersView) Title() string          { return "Wallpapers" }
func (wv *WallpapersView) IconName() string       { return "preferences-desktop-wallpaper-symbolic" }
func (wv *WallpapersView) OnShow()                { wv.refresh() }
func (wv *WallpapersView) OnHide()                {}
func (wv *WallpapersView) Destroy()               {}
