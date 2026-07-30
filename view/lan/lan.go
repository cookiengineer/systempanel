package lan

import (
	"os/exec"
	"strings"
	"unsafe"

	"github.com/cookiengineer/systempanel/bindings/gtk4"
	"github.com/cookiengineer/systempanel/detect"
	"github.com/cookiengineer/systempanel/view"
	"github.com/cookiengineer/systempanel/widget"
)

var Descriptor = view.ViewDescriptor{
	Name:     "lan",
	Title:    "LAN",
	IconName: "network-wired-symbolic",
	DetectFn: func() bool {
		return detect.HasProgram("nmcli")
	},
	Factory: func() view.View { return NewLANView() },
}

type LANView struct {
	box         *gtk4.Box
	list        *gtk4.ListBox
	connectBtn  *gtk4.Button
	conns       []lanConnection
	rowPtrToIdx map[unsafe.Pointer]int
	parentWin   *gtk4.Window
	selected    string
}

type lanConnection struct {
	name   string
	uuid   string
	device string
	row    *gtk4.ListBoxRow
}

func NewLANView() *LANView {
	lv := &LANView{
		box:         gtk4.BoxNew(gtk4.OrientationVertical, 12),
		rowPtrToIdx: make(map[unsafe.Pointer]int),
	}
	lv.box.SetMarginStart(24)
	lv.box.SetMarginEnd(24)
	lv.box.SetMarginTop(24)
	lv.box.SetMarginBottom(24)

	header := gtk4.LabelNew("Ethernet Connections")
	header.AddCSSClass("header-label")
	headerWidget := header.Widget
	lv.box.Append(&headerWidget)

	lv.list = gtk4.ListBoxNew()
	lv.list.SetSelectionMode(gtk4.SelectionSingle)
	lv.list.OnRowActivated(lv.onRowSelected)

	scrollW := gtk4.ScrolledWindowNew()
	scrollW.SetPolicy(gtk4.PolicyNever, gtk4.PolicyAutomatic)
	scrollW.SetHExpand(true)
	scrollW.SetVExpand(true)

	scrollWidget := scrollW.Widget
	listWidget := lv.list.Widget
	scrollW.SetChild(&listWidget)
	lv.box.Append(&scrollWidget)

	btnBox := gtk4.BoxNew(gtk4.OrientationHorizontal, 8)
	btnBox.SetMarginTop(4)

	refreshBtn := gtk4.ButtonNewWithLabel("Refresh")
	refreshBtn.OnClicked(func() { lv.refresh() })
	rbw := refreshBtn.Widget
	btnBox.Append(&rbw)

	spacer := gtk4.LabelNew("")
	spacer.SetHExpand(true)
	spWidget := spacer.Widget
	btnBox.Append(&spWidget)

	lv.connectBtn = gtk4.ButtonNewWithLabel("Connect")
	lv.connectBtn.OnClicked(lv.onConnectClicked)
	cbWidget := lv.connectBtn.Widget
	btnBox.Append(&cbWidget)

	btnBoxWidget := btnBox.Widget
	lv.box.Append(&btnBoxWidget)

	lv.refresh()

	return lv
}

func (lv *LANView) SetParentWindow(parent *gtk4.Window) { lv.parentWin = parent }

func (lv *LANView) refresh() {
	lv.selected = ""

	for _, c := range lv.conns {
		lv.list.Remove(c.row)
	}
	lv.conns = lv.conns[:0]
	lv.rowPtrToIdx = make(map[unsafe.Pointer]int)

	out, err := exec.Command("nmcli", "-t", "-f", "NAME,UUID,TYPE,DEVICE", "connection", "show").Output()
	if err != nil {
		return
	}
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		parts := strings.Split(line, ":")
		if len(parts) >= 4 && parts[2] == "802-3-ethernet" {
			conn := lanConnection{name: parts[0], uuid: parts[1], device: parts[3]}
			conn.row = lv.createRow(&conn)
			lv.list.Append(conn.row)
			lv.conns = append(lv.conns, conn)
			lv.rowPtrToIdx[conn.row.Widget.Ptr()] = len(lv.conns) - 1
		}
	}

	lv.connectBtn.SetSensitive(false)
}

func (lv *LANView) createRow(c *lanConnection) *gtk4.ListBoxRow {
	row := gtk4.ListBoxRowNew()
	row.AddCSSClass("device-row")

	hbox := gtk4.BoxNew(gtk4.OrientationHorizontal, 8)

	icon := gtk4.ImageNewFromIconName("network-wired-symbolic")
	icon.SetPixelSize(20)
	iconWidget := icon.Widget
	hbox.Append(&iconWidget)

	nameLabel := gtk4.LabelNew(c.name)
	nameLabel.SetHExpand(true)
	nameLabel.SetHAlign(gtk4.AlignStart)
	nlWidget := nameLabel.Widget
	hbox.Append(&nlWidget)

	devLabel := gtk4.LabelNew(c.device)
	devLabel.SetSensitive(false)
	dlWidget := devLabel.Widget
	hbox.Append(&dlWidget)

	gearBtn := gtk4.ButtonNew()
	gearBtn.SetIconName("emblem-system-symbolic")
	connName := c.name
	gearBtn.OnClicked(func() {
		dlg := widget.NewConnectionDialog(lv.parentWin, connName)
		dlg.Present()
	})
	gbWidget := gearBtn.Widget
	hbox.Append(&gbWidget)

	hboxWidget := hbox.Widget
	row.SetChild(&hboxWidget)

	return row
}

func (lv *LANView) onRowSelected(row *gtk4.ListBoxRow) {
	lv.selected = ""
	lv.connectBtn.SetSensitive(false)

	idx, ok := lv.rowPtrToIdx[row.Widget.Ptr()]
	if ok && idx < len(lv.conns) {
		lv.selected = lv.conns[idx].name
		lv.connectBtn.SetSensitive(true)
	}
}

func (lv *LANView) onConnectClicked() {
	if lv.selected != "" {
		exec.Command("nmcli", "connection", "up", lv.selected).Run()
		return
	}
	if len(lv.conns) == 0 {
		dlg := widget.NewConnectionDialog(lv.parentWin, "Ethernet Connection", "802-3-ethernet")
		dlg.Present()
	}
}

func (lv *LANView) Widget() *gtk4.Widget { return &lv.box.Widget }
func (lv *LANView) Name() string          { return "lan" }
func (lv *LANView) Title() string         { return "LAN" }
func (lv *LANView) IconName() string      { return "network-wired-symbolic" }
func (lv *LANView) OnShow()               { lv.refresh() }
func (lv *LANView) OnHide()               {}
func (lv *LANView) Destroy()              {}
