package bluetooth

import (
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
	box     *gtk4.Box
	model   *model.BluetoothModel
	listBox *gtk4.ListBox
	spinner *gtk4.Spinner
	devices []btDeviceItem
}

type btDeviceItem struct {
	device model.BluetoothDevice
	row    *gtk4.ListBoxRow
}

func NewBluetoothView() *BluetoothView {
	bv := &BluetoothView{
		box:     gtk4.BoxNew(gtk4.OrientationVertical, 12),
		model:   &model.BluetoothModel{},
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
	bv.listBox.OnRowActivated(func(row *gtk4.ListBoxRow) {
		for _, item := range bv.devices {
			if item.row == row {
				if item.device.Connected {
					bv.model.Disconnect(item.device.MAC)
				} else {
					bv.model.Connect(item.device.MAC)
				}
				return
			}
		}
	})

	scrollW := gtk4.ScrolledWindowNew()
	scrollW.SetPolicy(gtk4.PolicyNever, gtk4.PolicyAutomatic)
	scrollW.SetVExpand(true)

	scrollWidget := scrollW.Widget
	listWidget := bv.listBox.Widget
	scrollW.SetChild(&listWidget)
	bv.box.Append(&scrollWidget)

	btnBox := gtk4.BoxNew(gtk4.OrientationHorizontal, 8)

	scanBtn := gtk4.ButtonNewWithLabel("Scan")
	scanBtn.OnClicked(func() { bv.scan() })
	sbw := scanBtn.Widget
	btnBox.Append(&sbw)

	bv.spinner = gtk4.SpinnerNew()
	spinnerWidget := bv.spinner.Widget
	btnBox.Append(&spinnerWidget)

	rescanBtn := gtk4.ButtonNewWithLabel("Refresh")
	rescanBtn.OnClicked(func() { bv.refresh() })
	rbw := rescanBtn.Widget
	btnBox.Append(&rbw)

	btnBoxWidget := btnBox.Widget
	bv.box.Append(&btnBoxWidget)

	bv.refresh()

	return bv
}

func (bv *BluetoothView) scan() {
	bv.spinner.Start()
	go func() {
		bv.model.Scan()
		bv.spinner.Stop()
	}()
}

func (bv *BluetoothView) refresh() {
	devices, _ := bv.model.ListDevices()

	for _, item := range bv.devices {
		bv.listBox.Remove(item.row)
	}
	bv.devices = bv.devices[:0]

	for _, d := range devices {
		row := bv.createDeviceRow(d)
		bv.listBox.Append(row)
		bv.devices = append(bv.devices, btDeviceItem{device: d, row: row})
	}
}

func (bv *BluetoothView) createDeviceRow(d model.BluetoothDevice) *gtk4.ListBoxRow {
	row := gtk4.ListBoxRowNew()
	row.AddCSSClass("device-row")

	hbox := gtk4.BoxNew(gtk4.OrientationHorizontal, 8)

	iconName := "bluetooth-active-symbolic"
	if d.Connected {
		iconName = "bluetooth-connected-symbolic"
	}
	icon := gtk4.ImageNewFromIconName(iconName)
	icon.SetPixelSize(20)
	iconWidget := icon.Widget
	hbox.Append(&iconWidget)

	nameLabel := gtk4.LabelNew(d.Name)
	nameLabel.SetHExpand(true)
	nameLabel.SetHAlign(gtk4.AlignStart)
	nlWidget := nameLabel.Widget
	hbox.Append(&nlWidget)

	status := "Available"
	if d.Paired {
		status = "Paired"
	}
	if d.Connected {
		status = "Connected"
	}
	statusLabel := gtk4.LabelNew(status)
	statusLabel.SetSensitive(false)
	slWidget := statusLabel.Widget
	hbox.Append(&slWidget)

	hboxWidget := hbox.Widget
	row.SetChild(&hboxWidget)

	return row
}

func (bv *BluetoothView) Widget() *gtk4.Widget { return &bv.box.Widget }
func (bv *BluetoothView) Name() string          { return "bluetooth" }
func (bv *BluetoothView) Title() string         { return "Bluetooth" }
func (bv *BluetoothView) IconName() string      { return "bluetooth-active-symbolic" }
func (bv *BluetoothView) OnShow()               { bv.refresh() }
func (bv *BluetoothView) OnHide()               {}
func (bv *BluetoothView) Destroy()              {}
