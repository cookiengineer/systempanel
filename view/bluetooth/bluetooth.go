package bluetooth

import (
	"fmt"
	"sort"
	"sync"
	"unsafe"

	"github.com/cookiengineer/systempanel/bindings/gtk4"
	"github.com/cookiengineer/systempanel/detect"
	"github.com/cookiengineer/systempanel/model"
	"github.com/cookiengineer/systempanel/view"
)

var Descriptor = view.ViewDescriptor{
	Name:     "bluetooth",
	Title:    "Bluetooth",
	IconName: "bluetooth-active-symbolic",
	DetectFn: func() bool {
		return detect.HasProgram("bluetoothctl") && detect.HasBluetoothHardware()
	},
	Factory: func() view.View { return NewBluetoothView() },
}

type BluetoothView struct {
	box         *gtk4.Box
	model       *model.BluetoothModel
	listBox     *gtk4.ListBox
	connectBtn  *gtk4.Button
	spinner     *gtk4.Spinner
	devices     []btDeviceItem
	rowPtrToIdx map[unsafe.Pointer]int
	sectionRows []*gtk4.ListBoxRow
	parentWin   *gtk4.Window
	selectedMAC string
	scanning    bool
	mu          sync.Mutex
}

type btDeviceItem struct {
	device model.BluetoothDevice
	row    *gtk4.ListBoxRow
	paired bool
}

func NewBluetoothView() *BluetoothView {
	bv := &BluetoothView{
		box:         gtk4.BoxNew(gtk4.OrientationVertical, 12),
		model:       &model.BluetoothModel{},
		rowPtrToIdx: make(map[unsafe.Pointer]int),
	}
	bv.box.SetMarginStart(24)
	bv.box.SetMarginEnd(24)
	bv.box.SetMarginTop(24)
	bv.box.SetMarginBottom(24)

	header := gtk4.LabelNew("Bluetooth Devices")
	header.AddCSSClass("header-label")
	headerWidget := header.Widget
	bv.box.Append(&headerWidget)

	bv.listBox = gtk4.ListBoxNew()
	bv.listBox.SetSelectionMode(gtk4.SelectionSingle)
	bv.listBox.OnRowActivated(bv.onRowSelected)

	scrollW := gtk4.ScrolledWindowNew()
	scrollW.SetPolicy(gtk4.PolicyNever, gtk4.PolicyAutomatic)
	scrollW.SetHExpand(true)
	scrollW.SetVExpand(true)

	scrollWidget := scrollW.Widget
	listWidget := bv.listBox.Widget
	scrollW.SetChild(&listWidget)
	bv.box.Append(&scrollWidget)

	btnBox := gtk4.BoxNew(gtk4.OrientationHorizontal, 8)
	btnBox.SetMarginTop(4)

	refreshBtn := gtk4.ButtonNewWithLabel("Refresh")
	refreshBtn.OnClicked(func() {
		bv.startScan()
	})
	rbw := refreshBtn.Widget
	btnBox.Append(&rbw)

	bv.spinner = gtk4.SpinnerNew()
	spinnerWidget := bv.spinner.Widget
	btnBox.Append(&spinnerWidget)

	spacer := gtk4.LabelNew("")
	spacer.SetHExpand(true)
	spWidget := spacer.Widget
	btnBox.Append(&spWidget)

	bv.connectBtn = gtk4.ButtonNewWithLabel("Connect")
	bv.connectBtn.SetSensitive(false)
	bv.connectBtn.OnClicked(bv.onConnectClicked)
	cbWidget := bv.connectBtn.Widget
	btnBox.Append(&cbWidget)

	btnBoxWidget := btnBox.Widget
	bv.box.Append(&btnBoxWidget)

	bv.loadInitial()

	return bv
}

func (bv *BluetoothView) SetParentWindow(parent *gtk4.Window) { bv.parentWin = parent }

func (bv *BluetoothView) loadInitial() {
	go func() {
		devices, _ := bv.model.ListDevices()
		gtk4.IdleAdd(func() {
			bv.populateList(devices)
		})
	}()
}

func (bv *BluetoothView) startScan() {
	bv.mu.Lock()
	if bv.scanning {
		bv.mu.Unlock()
		return
	}
	bv.scanning = true
	bv.mu.Unlock()

	gtk4.IdleAdd(func() { bv.spinner.Start() })

	go bv.model.ScanLoop(
		func(devices []model.BluetoothDevice) {
			gtk4.IdleAdd(func() {
				bv.populateList(devices)
			})
		},
		func() {
			gtk4.IdleAdd(func() {
				bv.spinner.Stop()
				bv.mu.Lock()
				bv.scanning = false
				bv.mu.Unlock()
			})
		},
	)
}

func (bv *BluetoothView) populateList(devices []model.BluetoothDevice) {
	bv.selectedMAC = ""

	for _, item := range bv.devices {
		bv.listBox.Remove(item.row)
	}
	for _, row := range bv.sectionRows {
		bv.listBox.Remove(row)
	}
	bv.devices = bv.devices[:0]
	bv.sectionRows = bv.sectionRows[:0]
	bv.rowPtrToIdx = make(map[unsafe.Pointer]int)

	for _, d := range devices {
		bv.devices = append(bv.devices, btDeviceItem{
			device: d,
			paired: d.Paired,
		})
	}

	sort.Slice(bv.devices, func(i, j int) bool {
		if bv.devices[i].paired != bv.devices[j].paired {
			return bv.devices[i].paired
		}
		return bv.devices[i].device.Name < bv.devices[j].device.Name
	})

	idx := 0
	hasPaired := false
	for _, item := range bv.devices {
		if item.paired {
			hasPaired = true
			break
		}
	}
	if hasPaired {
		hdr := bv.createSectionLabel("─ Paired Devices ─")
		bv.listBox.Append(hdr)
		bv.sectionRows = append(bv.sectionRows, hdr)
	}

	insertedSep := false
	for i := range bv.devices {
		if !bv.devices[i].paired && !insertedSep {
			sep := bv.createSectionLabel("─ Available Devices ─")
			bv.listBox.Append(sep)
			bv.sectionRows = append(bv.sectionRows, sep)
			insertedSep = true
		}
		row := bv.createDeviceRow(&bv.devices[i])
		bv.listBox.Append(row)
		bv.devices[i].row = row
		bv.rowPtrToIdx[row.Widget.Ptr()] = idx
		idx++
	}

	bv.connectBtn.SetSensitive(false)
}

func (bv *BluetoothView) createSectionLabel(text string) *gtk4.ListBoxRow {
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

func (bv *BluetoothView) createDeviceRow(item *btDeviceItem) *gtk4.ListBoxRow {
	row := gtk4.ListBoxRowNew()
	row.AddCSSClass("device-row")

	hbox := gtk4.BoxNew(gtk4.OrientationHorizontal, 8)

	iconName := "bluetooth-active-symbolic"
	if item.device.Connected {
		iconName = "bluetooth-connected-symbolic"
	}
	icon := gtk4.ImageNewFromIconName(iconName)
	icon.SetPixelSize(20)
	iconWidget := icon.Widget
	hbox.Append(&iconWidget)

	nameLabel := gtk4.LabelNew(item.device.Name)
	nameLabel.SetHExpand(true)
	nameLabel.SetHAlign(gtk4.AlignStart)
	nlWidget := nameLabel.Widget
	hbox.Append(&nlWidget)

	if item.device.RSSI != 0 {
		pct := model.RSSIToPercent(item.device.RSSI)
		sigLabel := gtk4.LabelNew(fmt.Sprintf("%d%%", pct))
		sigLabel.SetSensitive(false)
		slWidget := sigLabel.Widget
		hbox.Append(&slWidget)
	}

	status := "Disconnected"
	if item.device.Connected {
		status = "Connected"
	}
	statusLabel := gtk4.LabelNew(status)
	statusLabel.SetSensitive(false)
	slWidget := statusLabel.Widget
	hbox.Append(&slWidget)

	if item.paired {
		unpairBtn := gtk4.ButtonNewWithLabel("Forget")
		mac := item.device.MAC
		unpairBtn.OnClicked(func() {
			bv.model.Forget(mac)
			bv.loadInitial()
		})
		ubWidget := unpairBtn.Widget
		hbox.Append(&ubWidget)
	}

	hboxWidget := hbox.Widget
	row.SetChild(&hboxWidget)

	return row
}

func (bv *BluetoothView) onRowSelected(row *gtk4.ListBoxRow) {
	bv.selectedMAC = ""
	bv.connectBtn.SetSensitive(false)

	idx, ok := bv.rowPtrToIdx[row.Widget.Ptr()]
	if ok && idx < len(bv.devices) {
		bv.selectedMAC = bv.devices[idx].device.MAC
		bv.connectBtn.SetSensitive(true)
	}
}

func (bv *BluetoothView) onConnectClicked() {
	if bv.selectedMAC == "" {
		return
	}
	bv.model.Connect(bv.selectedMAC)
	bv.loadInitial()
}

func (bv *BluetoothView) Widget() *gtk4.Widget { return &bv.box.Widget }
func (bv *BluetoothView) Name() string          { return "bluetooth" }
func (bv *BluetoothView) Title() string         { return "Bluetooth" }
func (bv *BluetoothView) IconName() string      { return "bluetooth-active-symbolic" }
func (bv *BluetoothView) OnShow()               { bv.loadInitial() }
func (bv *BluetoothView) OnHide()               {}
func (bv *BluetoothView) Destroy()              {}
