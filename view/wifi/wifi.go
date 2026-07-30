package wifi

import (
	"fmt"
	"sort"
	"strings"
	"unsafe"

	"github.com/cookiengineer/systempanel/bindings/gtk4"
	"github.com/cookiengineer/systempanel/detect"
	"github.com/cookiengineer/systempanel/model"
	"github.com/cookiengineer/systempanel/view"
	"github.com/cookiengineer/systempanel/widget"
)

var Descriptor = view.ViewDescriptor{
	Name:     "wifi",
	Title:    "Wi-Fi",
	IconName: "network-wireless-symbolic",
	DetectFn: func() bool {
		return detect.HasProgram("nmcli") && detect.HasWiFiHardware()
	},
	Factory: func() view.View { return NewWiFiView() },
}

type WiFiView struct {
	box         *gtk4.Box
	model       *model.WiFiModel
	listBox     *gtk4.ListBox
	connectBtn  *gtk4.Button
	networks    []networkListItem
	rowPtrToIdx map[unsafe.Pointer]int
	sectionRows []*gtk4.ListBoxRow
	parentWin   *gtk4.Window
	selectedSSID string
	selectedCfg  bool
}

type networkListItem struct {
	network    model.WiFiNetwork
	row        *gtk4.ListBoxRow
	configured bool
}

func NewWiFiView() *WiFiView {
	wv := &WiFiView{
		box:         gtk4.BoxNew(gtk4.OrientationVertical, 12),
		model:       &model.WiFiModel{},
		rowPtrToIdx: make(map[unsafe.Pointer]int),
	}
	wv.box.SetMarginStart(24)
	wv.box.SetMarginEnd(24)
	wv.box.SetMarginTop(24)
	wv.box.SetMarginBottom(24)

	header := gtk4.LabelNew("Wi-Fi Networks")
	header.AddCSSClass("header-label")
	headerWidget := header.Widget
	wv.box.Append(&headerWidget)

	wv.listBox = gtk4.ListBoxNew()
	wv.listBox.SetSelectionMode(gtk4.SelectionSingle)
	wv.listBox.OnRowActivated(wv.onRowSelected)

	scrollW := gtk4.ScrolledWindowNew()
	scrollW.SetPolicy(gtk4.PolicyNever, gtk4.PolicyAutomatic)
	scrollW.SetHExpand(true)
	scrollW.SetVExpand(true)

	scrollWidget := scrollW.Widget
	listWidget := wv.listBox.Widget
	scrollW.SetChild(&listWidget)
	wv.box.Append(&scrollWidget)

	btnBox := gtk4.BoxNew(gtk4.OrientationHorizontal, 8)
	btnBox.SetMarginTop(4)

	refreshBtn := gtk4.ButtonNewWithLabel("Refresh")
	refreshBtn.OnClicked(func() { wv.refresh() })
	rbw := refreshBtn.Widget
	btnBox.Append(&rbw)

	spacer := gtk4.LabelNew("")
	spacer.SetHExpand(true)
	spWidget := spacer.Widget
	btnBox.Append(&spWidget)

	wv.connectBtn = gtk4.ButtonNewWithLabel("Connect")
	wv.connectBtn.SetSensitive(false)
	wv.connectBtn.OnClicked(wv.onConnectClicked)
	cbWidget := wv.connectBtn.Widget
	btnBox.Append(&cbWidget)

	btnBoxWidget := btnBox.Widget
	wv.box.Append(&btnBoxWidget)

	wv.refresh()

	return wv
}

func (wv *WiFiView) SetParentWindow(parent *gtk4.Window) { wv.parentWin = parent }

func (wv *WiFiView) refresh() {
	wv.selectedSSID = ""
	wv.selectedCfg = false

	networks, _ := wv.model.Scan()
	profiles, _ := wv.model.ConfiguredSSIDs()

	for _, item := range wv.networks {
		wv.listBox.Remove(item.row)
	}
	for _, row := range wv.sectionRows {
		wv.listBox.Remove(row)
	}
	wv.networks = wv.networks[:0]
	wv.sectionRows = wv.sectionRows[:0]
	wv.rowPtrToIdx = make(map[unsafe.Pointer]int)

	for _, n := range networks {
		_, configured := profiles[n.SSID]
		wv.networks = append(wv.networks, networkListItem{network: n, configured: configured})
	}

	sort.Slice(wv.networks, func(i, j int) bool {
		if wv.networks[i].configured != wv.networks[j].configured {
			return wv.networks[i].configured
		}
		return strings.ToLower(wv.networks[i].network.SSID) < strings.ToLower(wv.networks[j].network.SSID)
	})

	insertedSep := false
	hasConfigured := false
	for _, item := range wv.networks {
		if item.configured {
			hasConfigured = true
			break
		}
	}
	if hasConfigured {
		headerRow := wv.createSectionLabel("─ Saved Networks ─")
		wv.listBox.Append(headerRow)
		wv.sectionRows = append(wv.sectionRows, headerRow)
	}
	for i := range wv.networks {
		if !wv.networks[i].configured && !insertedSep {
			sepRow := wv.createSectionLabel("─ Available Networks ─")
			wv.listBox.Append(sepRow)
			wv.sectionRows = append(wv.sectionRows, sepRow)
			insertedSep = true
		}
		row := wv.createNetworkRow(&wv.networks[i])
		wv.listBox.Append(row)
		wv.networks[i].row = row
		wv.rowPtrToIdx[row.Widget.Ptr()] = i
	}

	wv.connectBtn.SetSensitive(false)
}

func (wv *WiFiView) createSectionLabel(text string) *gtk4.ListBoxRow {
	row := gtk4.ListBoxRowNew()
	row.SetSensitive(false)
	label := gtk4.LabelNew(text)
	label.SetHAlign(gtk4.AlignCenter)
	label.SetMarginTop(8)
	label.SetMarginBottom(4)
	label.SetSensitive(false)
	lw := label.Widget
	row.SetChild(&lw)
	return row
}

func (wv *WiFiView) createNetworkRow(item *networkListItem) *gtk4.ListBoxRow {
	row := gtk4.ListBoxRowNew()
	row.AddCSSClass("device-row")

	hbox := gtk4.BoxNew(gtk4.OrientationHorizontal, 8)

	iconName := "network-wireless-symbolic"
	if item.network.Signal > 80 {
		iconName = "network-wireless-signal-excellent-symbolic"
	} else if item.network.Signal > 60 {
		iconName = "network-wireless-signal-good-symbolic"
	} else if item.network.Signal > 40 {
		iconName = "network-wireless-signal-ok-symbolic"
	} else {
		iconName = "network-wireless-signal-weak-symbolic"
	}
	icon := gtk4.ImageNewFromIconName(iconName)
	icon.SetPixelSize(20)
	iconWidget := icon.Widget
	hbox.Append(&iconWidget)

	nameLabel := gtk4.LabelNew(item.network.SSID)
	nameLabel.SetHExpand(true)
	nameLabel.SetHAlign(gtk4.AlignStart)
	nlWidget := nameLabel.Widget
	hbox.Append(&nlWidget)

	signalText := fmt.Sprintf("%d%%", item.network.Signal)
	sigLabel := gtk4.LabelNew(signalText)
	slWidget := sigLabel.Widget
	hbox.Append(&slWidget)

	secLabel := gtk4.LabelNew(item.network.Security)
	secLabel.SetSensitive(false)
	seclWidget := secLabel.Widget
	hbox.Append(&seclWidget)

	if item.configured {
		gearBtn := gtk4.ButtonNew()
		gearBtn.SetIconName("emblem-system-symbolic")
		ssid := item.network.SSID
		gearBtn.OnClicked(func() {
			dlg := widget.NewConnectionDialog(wv.parentWin, ssid)
			dlg.Present()
		})
		gbWidget := gearBtn.Widget
		hbox.Append(&gbWidget)
	}

	hboxWidget := hbox.Widget
	row.SetChild(&hboxWidget)

	return row
}

func (wv *WiFiView) onRowSelected(row *gtk4.ListBoxRow) {
	wv.selectedSSID = ""
	wv.connectBtn.SetSensitive(false)

	idx, ok := wv.rowPtrToIdx[row.Widget.Ptr()]
	if ok {
		wv.selectedSSID = wv.networks[idx].network.SSID
		wv.selectedCfg = wv.networks[idx].configured
		wv.connectBtn.SetSensitive(true)
	}
}

func (wv *WiFiView) onConnectClicked() {
	if wv.selectedSSID == "" {
		return
	}
	if wv.selectedCfg {
		wv.model.Connect(wv.selectedSSID, "")
	} else {
		dlg := widget.NewConnectionDialog(wv.parentWin, wv.selectedSSID)
		dlg.Present()
	}
}

func (wv *WiFiView) Widget() *gtk4.Widget { return &wv.box.Widget }
func (wv *WiFiView) Name() string          { return "wifi" }
func (wv *WiFiView) Title() string         { return "Wi-Fi" }
func (wv *WiFiView) IconName() string      { return "network-wireless-symbolic" }
func (wv *WiFiView) OnShow()               { wv.refresh() }
func (wv *WiFiView) OnHide()               {}
func (wv *WiFiView) Destroy()              {}
