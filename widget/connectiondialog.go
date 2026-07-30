package widget

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/cookiengineer/systempanel/bindings/gtk4"
	"github.com/cookiengineer/systempanel/parsers/nmconnection"
)

type ConnectionDialog struct {
	win         *gtk4.Window
	stack       *gtk4.Stack
	connection  *nmconnection.Connection
	parentWin   *gtk4.Window
	ssid        string
	connType    string
	isNew       bool
	connPath    string

	autoconnectChk  *gtk4.CheckButton
	deviceCombo     *gtk4.ComboBoxText

	ssidEntry       *gtk4.Entry
	modeCombo       *gtk4.ComboBoxText
	bandCombo       *gtk4.ComboBoxText
	bssidEntry      *gtk4.Entry
	clonedMacCombo  *gtk4.ComboBoxText
	mtuSpin         *gtk4.SpinButton

	secCombo        *gtk4.ComboBoxText
	passwordEntry   *gtk4.Entry

	ip4Method       *gtk4.ComboBoxText
	ip4Address      *gtk4.Entry
	ip4Netmask      *gtk4.Entry
	ip4Gateway      *gtk4.Entry
	ip4DNS          *gtk4.Entry
	ip4NeverDefault *gtk4.CheckButton

	ip6Method       *gtk4.ComboBoxText
	ip6Address      *gtk4.Entry
	ip6Prefix       *gtk4.Entry
	ip6Gateway      *gtk4.Entry
	ip6DNS          *gtk4.Entry
	ip6NeverDefault *gtk4.CheckButton
}

func NewConnectionDialog(parent *gtk4.Window, ssid string, typeHints ...string) *ConnectionDialog {
	d := &ConnectionDialog{
		parentWin: parent,
		ssid:      ssid,
	}
	if len(typeHints) > 0 {
		d.connType = typeHints[0]
	}
	d.initConnection()
	d.build()
	return d
}

func (d *ConnectionDialog) initConnection() {
	profileDir := "/etc/NetworkManager/system-connections"
	entries, _ := os.ReadDir(profileDir)
	for _, e := range entries {
		if !strings.HasSuffix(e.Name(), ".nmconnection") {
			continue
		}
		path := filepath.Join(profileDir, e.Name())
		conn, err := nmconnection.Parse(path)
		if err != nil {
			continue
		}
		if conn.Wifi.SSID == d.ssid || conn.Connection.ID == d.ssid {
			d.connection = conn
			d.connPath = path
			d.isNew = false
			d.connType = conn.Connection.Type
			return
		}
	}
	if d.connType == "" {
		d.connType = "wifi"
	}
	if d.connType == "802-3-ethernet" || d.connType == "ethernet" {
		d.connection = nmconnection.NewEthernetConnection(d.ssid)
		d.connType = "802-3-ethernet"
	} else {
		d.connection = nmconnection.NewWiFiConnection(d.ssid)
	}
	d.isNew = true
}

func (d *ConnectionDialog) build() {
	d.win = gtk4.WindowNew()
	d.win.SetTitle("Connection Settings - " + d.ssid)
	d.win.SetDefaultSize(550, 480)
	d.win.SetModal(true)
	d.win.SetTransientFor(d.parentWin)

	vbox := gtk4.BoxNew(gtk4.OrientationVertical, 0)

	topBox := d.buildTopBar()
	tbWidget := topBox.Widget
	vbox.Append(&tbWidget)

	switcher := gtk4.StackSwitcherNew()
	switcher.SetMarginStart(12)
	switcher.SetMarginTop(4)
	swWidget := switcher.Widget
	vbox.Append(&swWidget)

	d.stack = gtk4.StackNew()
	d.stack.SetTransitionType(gtk4.StackTransitionCrossfade)
	d.stack.SetTransitionDuration(120)
	d.stack.SetVHomogeneous(true)
	d.stack.SetHExpand(true)
	d.stack.SetVExpand(true)

	if d.connType == "wifi" || d.connType == "802-11-wireless" {
		d.stack.AddTitled(&d.buildWifiTab().Widget, "wifi", "Wi-Fi")
		d.stack.AddTitled(&d.buildSecurityTab().Widget, "wifi-security", "Wi-Fi Security")
		d.stack.SetVisibleChildName("wifi")
	} else {
		d.stack.SetVisibleChildName("ipv4")
	}
	d.stack.AddTitled(&d.buildIPv4Tab().Widget, "ipv4", "IPv4 Settings")
	d.stack.AddTitled(&d.buildIPv6Tab().Widget, "ipv6", "IPv6 Settings")

	switcher.SetStack(d.stack)

	stackWidget := d.stack.Widget
	vbox.Append(&stackWidget)

	btmBox := d.buildBottomBar()
	bmWidget := btmBox.Widget
	vbox.Append(&bmWidget)

	vboxWidget := vbox.Widget
	d.win.SetChild(&vboxWidget)

	d.populateFromConnection()
}

func (d *ConnectionDialog) buildTopBar() *gtk4.Box {
	hbox := gtk4.BoxNew(gtk4.OrientationHorizontal, 8)
	hbox.SetMarginStart(24)
	hbox.SetMarginEnd(24)
	hbox.SetMarginTop(12)
	hbox.SetMarginBottom(8)

	d.autoconnectChk = gtk4.CheckButtonNewWithLabel("Connect automatically")
	acWidget := d.autoconnectChk.Widget
	hbox.Append(&acWidget)

	devLabel := gtk4.LabelNew("Device:")
	dlWidget := devLabel.Widget
	hbox.Append(&dlWidget)
	d.deviceCombo = gtk4.ComboBoxTextNew()
	d.deviceCombo.SetHExpand(true)
	dcWidget := d.deviceCombo.Widget
	hbox.Append(&dcWidget)

	d.populateDevices()

	return hbox
}

func (d *ConnectionDialog) buildBottomBar() *gtk4.Box {
	hbox := gtk4.BoxNew(gtk4.OrientationHorizontal, 8)
	hbox.SetMarginStart(24)
	hbox.SetMarginEnd(24)
	hbox.SetMarginTop(12)
	hbox.SetMarginBottom(12)

	spacer := gtk4.LabelNew("")
	spacer.SetHExpand(true)
	spWidget := spacer.Widget
	hbox.Append(&spWidget)

	if !d.isNew {
		removeBtn := gtk4.ButtonNewWithLabel("Remove")
		removeBtn.AddCSSClass("destructive-action")
		removeBtn.OnClicked(d.onRemove)
		rbWidget := removeBtn.Widget
		hbox.Append(&rbWidget)
	}

	saveBtn := gtk4.ButtonNewWithLabel("Save")
	saveBtn.AddCSSClass("suggested-action")
	saveBtn.OnClicked(d.onSave)
	sbWidget := saveBtn.Widget
	hbox.Append(&sbWidget)

	return hbox
}

func (d *ConnectionDialog) buildWifiTab() *gtk4.Box {
	box := gtk4.BoxNew(gtk4.OrientationVertical, 8)
	box.SetMarginStart(24)
	box.SetMarginEnd(24)
	box.SetMarginTop(12)

	d.addField(box, "SSID:", &d.ssidEntry)

	d.modeCombo = gtk4.ComboBoxTextNew()
	d.modeCombo.Append("infrastructure", "Client")
	d.modeCombo.Append("ap", "Hotspot")
	d.modeCombo.Append("adhoc", "Ad-hoc")
	d.modeCombo.SetActive(0)
	d.addComboRow(box, "Mode:", d.modeCombo)

	d.bandCombo = gtk4.ComboBoxTextNew()
	d.bandCombo.Append("", "Automatic")
	d.bandCombo.Append("a", "A (5 GHz)")
	d.bandCombo.Append("bg", "B/G (2.4 GHz)")
	d.bandCombo.SetActive(0)
	d.addComboRow(box, "Band:", d.bandCombo)

	d.addField(box, "BSSID:", &d.bssidEntry)

	d.clonedMacCombo = gtk4.ComboBoxTextNewWithEntry()
	d.clonedMacCombo.Append("preserve", "Preserve")
	d.clonedMacCombo.Append("permanent", "Permanent")
	d.clonedMacCombo.Append("random", "Random")
	d.clonedMacCombo.Append("stable", "Stable")
	d.clonedMacCombo.SetActive(0)
	d.addComboRow(box, "Cloned MAC:", d.clonedMacCombo)

	d.mtuSpin = gtk4.SpinButtonNew(0, 65536, 1)
	d.mtuSpin.SetValue(0)
	d.addSpinRow(box, "MTU:", d.mtuSpin, "bytes")

	return box
}

func (d *ConnectionDialog) buildSecurityTab() *gtk4.Box {
	box := gtk4.BoxNew(gtk4.OrientationVertical, 8)
	box.SetMarginStart(24)
	box.SetMarginEnd(24)
	box.SetMarginTop(12)

	d.secCombo = gtk4.ComboBoxTextNew()
	d.secCombo.Append("none", "None")
	d.secCombo.Append("wpa-psk", "WPA/WPA2/WPA3 Personal")
	d.secCombo.Append("wpa-eap", "WPA/WPA2 Enterprise")
	d.secCombo.Append("ieee8021x", "Dynamic WEP (802.1X)")
	d.secCombo.Append("wep", "WEP")
	d.secCombo.Append("sae", "WPA3 Personal (SAE)")
	d.secCombo.SetActive(1)
	d.addComboRow(box, "Security:", d.secCombo)

	d.passwordEntry = gtk4.EntryNew()
	d.passwordEntry.SetVisibility(false)
	d.passwordEntry.SetPlaceholder("Wi-Fi password")
	d.addField(box, "Password:", &d.passwordEntry)

	return box
}

func (d *ConnectionDialog) buildIPv4Tab() *gtk4.Box {
	box := gtk4.BoxNew(gtk4.OrientationVertical, 8)
	box.SetMarginStart(24)
	box.SetMarginEnd(24)
	box.SetMarginTop(12)

	d.ip4Method = gtk4.ComboBoxTextNew()
	d.ip4Method.Append("auto", "Automatic (DHCP)")
	d.ip4Method.Append("manual", "Manual")
	d.ip4Method.Append("link-local", "Link-Local Only")
	d.ip4Method.Append("shared", "Shared to other computers")
	d.ip4Method.Append("disabled", "Disabled")
	d.ip4Method.SetActive(0)
	d.addComboRow(box, "Method:", d.ip4Method)

	d.addField(box, "Address:", &d.ip4Address)
	d.addField(box, "Netmask:", &d.ip4Netmask)
	d.addField(box, "Gateway:", &d.ip4Gateway)
	d.addField(box, "DNS servers:", &d.ip4DNS)

	d.ip4NeverDefault = gtk4.CheckButtonNewWithLabel("Use this connection only for resources on its network")
	ndRow := gtk4.BoxNew(gtk4.OrientationHorizontal, 8)
	ndSpacer := gtk4.LabelNew("")
	ndSpacer.SetSizeRequest(130, -1)
	nsw := ndSpacer.Widget
	ndRow.Append(&nsw)
	ndWidget := d.ip4NeverDefault.Widget
	ndWidget.SetMarginTop(8)
	ndRow.Append(&ndWidget)
	ndRowWidget := ndRow.Widget
	box.Append(&ndRowWidget)

	return box
}

func (d *ConnectionDialog) buildIPv6Tab() *gtk4.Box {
	box := gtk4.BoxNew(gtk4.OrientationVertical, 8)
	box.SetMarginStart(24)
	box.SetMarginEnd(24)
	box.SetMarginTop(12)

	d.ip6Method = gtk4.ComboBoxTextNew()
	d.ip6Method.Append("auto", "Automatic")
	d.ip6Method.Append("manual", "Manual")
	d.ip6Method.Append("link-local", "Link-Local Only")
	d.ip6Method.Append("shared", "Shared to other computers")
	d.ip6Method.Append("ignore", "Ignore")
	d.ip6Method.Append("disabled", "Disabled")
	d.ip6Method.SetActive(0)
	d.addComboRow(box, "Method:", d.ip6Method)

	d.addField(box, "Address:", &d.ip6Address)
	d.addField(box, "Prefix:", &d.ip6Prefix)
	d.addField(box, "Gateway:", &d.ip6Gateway)
	d.addField(box, "DNS servers:", &d.ip6DNS)

	d.ip6NeverDefault = gtk4.CheckButtonNewWithLabel("Use this connection only for resources on its network")
	nd6Row := gtk4.BoxNew(gtk4.OrientationHorizontal, 8)
	nd6Spacer := gtk4.LabelNew("")
	nd6Spacer.SetSizeRequest(130, -1)
	n6sw := nd6Spacer.Widget
	nd6Row.Append(&n6sw)
	nd6Widget := d.ip6NeverDefault.Widget
	nd6Widget.SetMarginTop(8)
	nd6Row.Append(&nd6Widget)
	nd6RowWidget := nd6Row.Widget
	box.Append(&nd6RowWidget)

	return box
}

func (d *ConnectionDialog) addField(box *gtk4.Box, labelText string, entryPtr **gtk4.Entry) {
	row := gtk4.BoxNew(gtk4.OrientationHorizontal, 8)
	row.SetMarginTop(4)
	row.SetMarginBottom(4)

	lbl := gtk4.LabelNew(labelText)
	lbl.SetSizeRequest(130, -1)
	lbl.SetHAlign(gtk4.AlignEnd)
	lbl.SetMarginEnd(4)
	lw := lbl.Widget
	row.Append(&lw)

	entry := gtk4.EntryNew()
	entry.SetHExpand(true)
	ew := entry.Widget
	row.Append(&ew)
	*entryPtr = entry

	rw := row.Widget
	box.Append(&rw)
}

func (d *ConnectionDialog) addComboRow(box *gtk4.Box, labelText string, combo *gtk4.ComboBoxText) {
	row := gtk4.BoxNew(gtk4.OrientationHorizontal, 8)
	row.SetMarginTop(4)
	row.SetMarginBottom(4)

	lbl := gtk4.LabelNew(labelText)
	lbl.SetSizeRequest(130, -1)
	lbl.SetHAlign(gtk4.AlignEnd)
	lbl.SetMarginEnd(4)
	lw := lbl.Widget
	row.Append(&lw)

	cw := combo.Widget
	cw.SetHExpand(true)
	row.Append(&cw)

	rw := row.Widget
	box.Append(&rw)
}

func (d *ConnectionDialog) addSpinRow(box *gtk4.Box, labelText string, spin *gtk4.SpinButton, suffix string) {
	row := gtk4.BoxNew(gtk4.OrientationHorizontal, 8)
	row.SetMarginTop(4)
	row.SetMarginBottom(4)

	lbl := gtk4.LabelNew(labelText)
	lbl.SetSizeRequest(130, -1)
	lbl.SetHAlign(gtk4.AlignEnd)
	lbl.SetMarginEnd(4)
	lw := lbl.Widget
	row.Append(&lw)

	sw := spin.Widget
	row.Append(&sw)

	if suffix != "" {
		sl := gtk4.LabelNew(suffix)
		slw := sl.Widget
		row.Append(&slw)
	}

	rw := row.Widget
	box.Append(&rw)
}

func (d *ConnectionDialog) populateDevices() {
	out, err := exec.Command("nmcli", "-t", "-f", "DEVICE,TYPE", "device", "status").Output()
	if err != nil {
		return
	}
	d.deviceCombo.Append("", "Automatic")
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		parts := strings.Split(line, ":")
		if len(parts) < 2 {
			continue
		}
		d.deviceCombo.Append(parts[0], parts[0])
	}
	d.deviceCombo.SetActive(0)
}

func (d *ConnectionDialog) populateFromConnection() {
	c := d.connection

	d.autoconnectChk.SetActive(c.Connection.Autoconnect)
	if c.Connection.InterfaceName != "" {
		d.selectComboOrAppend(d.deviceCombo, c.Connection.InterfaceName)
	}

	if d.connType == "wifi" || d.connType == "802-11-wireless" {
		d.ssidEntry.SetText(c.Wifi.SSID)
		if c.Wifi.Band != "" {
			if c.Wifi.Band == "a" { d.bandCombo.SetActive(1) } else { d.bandCombo.SetActive(2) }
		}
		d.bssidEntry.SetText(c.Wifi.BSSID)
		if c.Wifi.ClonedMAC != "" {
			d.setComboOrText(d.clonedMacCombo, c.Wifi.ClonedMAC)
		}
		d.mtuSpin.SetValue(float64(c.Wifi.MTU))

		secType := c.WifiSecurity.KeyMgmt
		if secType == "" { secType = "wpa-psk" }
		for i := 0; i < 6; i++ {
			d.secCombo.SetActive(i)
			if d.secCombo.GetActiveID() == secType { break }
		}
		d.passwordEntry.SetText(c.WifiSecurity.PSK)
	}

	for i := 0; i < 5; i++ {
		d.ip4Method.SetActive(i)
		if d.ip4Method.GetActiveID() == c.IPv4.Method { break }
	}
	d.ip4Address.SetText(c.IPv4.Address1)
	if c.IPv4.Address1 != "" {
		parts := strings.Split(c.IPv4.Address1, ",")
		if len(parts) >= 1 { d.ip4Address.SetText(parts[0]) }
		if len(parts) >= 2 { d.ip4Netmask.SetText(parts[1]) }
	}
	d.ip4Gateway.SetText(c.IPv4.Gateway)
	d.ip4DNS.SetText(c.IPv4.DNS)
	d.ip4NeverDefault.SetActive(c.IPv4.NeverDefault)

	for i := 0; i < 6; i++ {
		d.ip6Method.SetActive(i)
		if d.ip6Method.GetActiveID() == c.IPv6.Method { break }
	}
	d.ip6Address.SetText(c.IPv6.Address1)
	if c.IPv6.Address1 != "" {
		parts := strings.Split(c.IPv6.Address1, ",")
		if len(parts) >= 1 { d.ip6Address.SetText(parts[0]) }
		if len(parts) >= 2 { d.ip6Prefix.SetText(parts[1]) }
	}
	d.ip6Gateway.SetText(c.IPv6.Gateway)
	d.ip6DNS.SetText(c.IPv6.DNS)
	d.ip6NeverDefault.SetActive(c.IPv6.NeverDefault)
}

func (d *ConnectionDialog) writeToConnection() {
	c := d.connection

	c.Connection.ID = d.ssid
	c.Connection.Autoconnect = d.autoconnectChk.GetActive()
	c.Connection.Type = d.connType
	if did := d.deviceCombo.GetActiveID(); did != "" {
		c.Connection.InterfaceName = did
	}

	if d.connType == "wifi" || d.connType == "802-11-wireless" {
		c.Wifi.SSID = d.ssidEntry.GetText()
		c.Wifi.Mode = d.modeCombo.GetActiveID()
		c.Wifi.Band = d.bandCombo.GetActiveID()
		c.Wifi.BSSID = d.bssidEntry.GetText()
		c.Wifi.ClonedMAC = d.getComboOrText(d.clonedMacCombo)
		c.Wifi.MTU = int(d.mtuSpin.GetValue())

		c.WifiSecurity.KeyMgmt = d.secCombo.GetActiveID()
		c.WifiSecurity.PSK = d.passwordEntry.GetText()
	}

	c.IPv4.Method = d.ip4Method.GetActiveID()
	addr := d.ip4Address.GetText()
	netmask := d.ip4Netmask.GetText()
	if addr != "" {
		if netmask != "" { addr += "," + netmask }
		c.IPv4.Address1 = addr
	}
	c.IPv4.Gateway = d.ip4Gateway.GetText()
	c.IPv4.DNS = d.ip4DNS.GetText()
	c.IPv4.NeverDefault = d.ip4NeverDefault.GetActive()

	c.IPv6.Method = d.ip6Method.GetActiveID()
	addr6 := d.ip6Address.GetText()
	prefix := d.ip6Prefix.GetText()
	if addr6 != "" {
		if prefix != "" { addr6 += "," + prefix }
		c.IPv6.Address1 = addr6
	}
	c.IPv6.Gateway = d.ip6Gateway.GetText()
	c.IPv6.DNS = d.ip6DNS.GetText()
	c.IPv6.NeverDefault = d.ip6NeverDefault.GetActive()
}

func (d *ConnectionDialog) onSave() {
	d.writeToConnection()
	content := d.connection.Serialize()

	profileDir := "/etc/NetworkManager/system-connections"
	filename := sanitizeFilename(d.ssid) + ".nmconnection"
	if d.connPath != "" {
		filename = filepath.Base(d.connPath)
	}
	path := filepath.Join(profileDir, filename)

	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		tmpPath := filepath.Join("/tmp", filename+".tmp")
		os.WriteFile(tmpPath, []byte(content), 0600)
		if err := RunSudoCommand("cp", tmpPath, path); err != nil {
			sd := NewSudoDialog(d.parentWin)
			sd.RunCommand("cp", tmpPath, path)
		}
		os.Remove(tmpPath)
	}
	exec.Command("nmcli", "connection", "reload").Run()
	d.win.Close()
}

func (d *ConnectionDialog) onRemove() {
	if d.connPath == "" {
		return
	}
	exec.Command("pkexec", "rm", d.connPath).Run()
	exec.Command("nmcli", "connection", "reload").Run()
	d.win.Close()
}

func (d *ConnectionDialog) Present() {
	d.win.Present()
}

func sanitizeFilename(s string) string {
	return strings.Map(func(r rune) rune {
		if r == '/' || r == 0 {
			return -1
		}
		return r
	}, s)
}

func (d *ConnectionDialog) setComboOrText(combo *gtk4.ComboBoxText, value string) {
	for i := 0; i < 10; i++ {
		combo.SetActive(i)
		if combo.GetActiveID() == value {
			return
		}
	}
	combo.SetActive(-1)
}

func (d *ConnectionDialog) getComboOrText(combo *gtk4.ComboBoxText) string {
	if id := combo.GetActiveID(); id != "" {
		return id
	}
	return combo.GetActiveText()
}

func (d *ConnectionDialog) selectComboOrAppend(combo *gtk4.ComboBoxText, id string) {
	for i := 0; i < 100; i++ {
		combo.SetActive(i)
		if combo.GetActive() == -1 {
			combo.Append(id, id)
			combo.SetActive(i)
			return
		}
		if combo.GetActiveID() == id {
			return
		}
	}
}
