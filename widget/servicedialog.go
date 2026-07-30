package widget

import (
	"os"
	"os/exec"
	"sort"
	"strings"

	"github.com/cookiengineer/systempanel/bindings/gtk4"
	"github.com/cookiengineer/systempanel/parsers/systemdunit"
)

type capCheck struct {
	check *gtk4.CheckButton
	label string
}

type ServiceDialog struct {
	win       *gtk4.Window
	stack     *gtk4.Stack
	unitFile  *systemdunit.UnitFile
	parentWin *gtk4.Window
	unitName  string
	isUser    bool
	advanced  bool
	modeBtn   *gtk4.Button
	entries   map[string]*gtk4.Entry
	combos    map[string]*gtk4.ComboBoxText
	checks    map[string]*gtk4.CheckButton
	capChecks map[string][]capCheck
	targets   []string
}

func NewServiceDialog(parent *gtk4.Window, uf *systemdunit.UnitFile, unitName string, isUser bool) *ServiceDialog {
	sd := &ServiceDialog{
		parentWin: parent,
		unitFile:  uf,
		unitName:  unitName,
		isUser:    isUser,
		entries:   make(map[string]*gtk4.Entry),
		combos:    make(map[string]*gtk4.ComboBoxText),
		checks:    make(map[string]*gtk4.CheckButton),
		capChecks: make(map[string][]capCheck),
	}
	sd.loadTargets()
	sd.build()
	return sd
}

func (sd *ServiceDialog) loadTargets() {
	out, err := exec.Command("systemctl", "list-units", "--type=target", "--all", "--no-legend", "--no-pager").Output()
	if err != nil {
		return
	}
	for _, line := range strings.Split(string(out), "\n") {
		fields := strings.Fields(line)
		if len(fields) > 0 {
			sd.targets = append(sd.targets, fields[0])
		}
	}
}

func (sd *ServiceDialog) build() {
	sd.win = gtk4.WindowNew()
	sd.win.SetTitle("Service Settings - " + sd.unitName)
	sd.win.SetDefaultSize(650, 520)
	sd.win.SetModal(true)
	sd.win.SetTransientFor(sd.parentWin)

	vbox := gtk4.BoxNew(gtk4.OrientationVertical, 0)

	switcher := gtk4.StackSwitcherNew()
	switcher.SetMarginStart(12)
	switcher.SetMarginTop(4)
	vbox.Append(&switcher.Widget)

	sd.stack = gtk4.StackNew()
	sd.stack.SetTransitionType(gtk4.StackTransitionCrossfade)
	sd.stack.SetTransitionDuration(120)
	sd.stack.SetVHomogeneous(true)
	sd.stack.SetHExpand(true)
	sd.stack.SetVExpand(true)

	sections := []string{"Unit", "Service", "Socket", "Install"}
	firstSection := ""

	for _, secName := range sections {
		keys := sd.unitFile.SectionKeys(secName)
		knownKeys := systemdunit.KnownSectionKeys(secName)
		if len(keys) == 0 && knownKeys == nil {
			continue
		}
		tab := sd.buildSectionTab(secName)
		sd.stack.AddTitled(&tab.Widget, secName, secName)
		if firstSection == "" {
			firstSection = secName
		}
	}

	for secName, sec := range sd.unitFile.Sections {
		if isStandardSection(secName) || len(sec.Keys) == 0 {
			continue
		}
		tab := sd.buildSectionTab(secName)
		sd.stack.AddTitled(&tab.Widget, secName, secName)
		if firstSection == "" {
			firstSection = secName
		}
	}

	if firstSection != "" {
		sd.stack.SetVisibleChildName(firstSection)
	}
	switcher.SetStack(sd.stack)

	vbox.Append(&sd.stack.Widget)

	sep := gtk4.BoxNew(gtk4.OrientationHorizontal, 0)
	sep.SetSizeRequest(-1, 1)
	vbox.Append(&sep.Widget)

	btnBox := gtk4.BoxNew(gtk4.OrientationHorizontal, 8)
	btnBox.SetMarginStart(24)
	btnBox.SetMarginEnd(24)
	btnBox.SetMarginTop(12)
	btnBox.SetMarginBottom(12)

	cancelBtn := gtk4.ButtonNewWithLabel("Cancel")
	cancelBtn.OnClicked(func() { sd.win.Close() })
	btnBox.Append(&cancelBtn.Widget)

	sd.modeBtn = gtk4.ButtonNewWithLabel("Advanced")
	sd.modeBtn.OnClicked(sd.onToggleMode)
	btnBox.Append(&sd.modeBtn.Widget)

	spacer := gtk4.LabelNew("")
	spacer.SetHExpand(true)
	btnBox.Append(&spacer.Widget)

	saveBtn := gtk4.ButtonNewWithLabel("Save")
	saveBtn.AddCSSClass("suggested-action")
	saveBtn.OnClicked(sd.onSave)
	btnBox.Append(&saveBtn.Widget)

	vbox.Append(&btnBox.Widget)

	sd.win.SetChild(&vbox.Widget)
}

func (sd *ServiceDialog) getSectionKeys(sectionName string) []string {
	if sd.advanced {
		fileKeys := sd.unitFile.SectionKeys(sectionName)
		knownKeys := systemdunit.KnownSectionKeys(sectionName)
		if knownKeys == nil {
			return fileKeys
		}
		keySet := make(map[string]bool)
		var allKeys []string
		for _, k := range fileKeys {
			if !keySet[k] {
				keySet[k] = true
				allKeys = append(allKeys, k)
			}
		}
		for _, k := range knownKeys {
			if !keySet[k] {
				keySet[k] = true
				allKeys = append(allKeys, k)
			}
		}
		return systemdunit.SortKeysByKnown(sectionName, allKeys)
	}
	return sd.unitFile.SectionKeys(sectionName)
}

func (sd *ServiceDialog) onToggleMode() {
	sd.advanced = !sd.advanced
	if sd.advanced {
		sd.modeBtn.SetLabel("Simple")
	} else {
		sd.modeBtn.SetLabel("Advanced")
	}
	sd.rebuildTabs()
}

func (sd *ServiceDialog) rebuildTabs() {
	sd.entries = make(map[string]*gtk4.Entry)
	sd.combos = make(map[string]*gtk4.ComboBoxText)
	sd.checks = make(map[string]*gtk4.CheckButton)
	sd.capChecks = make(map[string][]capCheck)

	sections := []string{"Unit", "Service", "Socket", "Install"}
	for _, secName := range sections {
		child := sd.stack.GetChildByName(secName)
		if child != nil {
			sd.stack.Remove(child)
		}
		newBox := sd.buildSectionTab(secName)
		sd.stack.AddTitled(&newBox.Widget, secName, secName)
	}
}

func (sd *ServiceDialog) buildSectionTab(sectionName string) *gtk4.Box {
	box := gtk4.BoxNew(gtk4.OrientationVertical, 8)
	box.SetMarginStart(24)
	box.SetMarginEnd(24)
	box.SetMarginTop(12)

	if sd.unitFile == nil {
		return box
	}

	scrollW := gtk4.ScrolledWindowNew()
	scrollW.SetPolicy(gtk4.PolicyNever, gtk4.PolicyAutomatic)
	scrollW.SetHExpand(true)
	scrollW.SetVExpand(true)

	innerBox := gtk4.BoxNew(gtk4.OrientationVertical, 4)

	allKeys := sd.getSectionKeys(sectionName)

	for _, key := range allKeys {
		value := sd.unitFile.Get(sectionName, key)

		if key == "CapabilityBoundingSet" {
			row := sd.createCapSetRow(sectionName, key, value)
			innerBox.Append(&row.Widget)
		} else if isBooleanKey(sectionName, key) {
			row := sd.createCheckRow(sectionName, key, value)
			innerBox.Append(&row.Widget)
		} else if isRelationshipKey(sectionName, key) {
			row := sd.createComboRow(sectionName, key, value)
			innerBox.Append(&row.Widget)
		} else {
			row := sd.createEntryRow(sectionName, key, value)
			innerBox.Append(&row.Widget)
		}
	}

	scrollW.SetChild(&innerBox.Widget)
	box.Append(&scrollW.Widget)

	scrollWidget := scrollW.Widget
	scrollWidget.SetVExpand(true)

	return box
}

func (sd *ServiceDialog) createEntryRow(section, key, value string) *gtk4.Box {
	row := gtk4.BoxNew(gtk4.OrientationHorizontal, 8)
	row.SetMarginTop(2)
	row.SetMarginBottom(2)

	lbl := gtk4.LabelNew(key)
	lbl.SetSizeRequest(220, -1)
	lbl.SetHAlign(gtk4.AlignStart)
	lbl.SetMarginEnd(8)
	lbl.SetTooltip(key)
	row.Append(&lbl.Widget)

	entry := gtk4.EntryNew()
	entry.SetText(value)
	entry.SetHExpand(true)
	row.Append(&entry.Widget)

	sd.entries[section+"/"+key] = entry

	return row
}

func (sd *ServiceDialog) createComboRow(section, key, value string) *gtk4.Box {
	row := gtk4.BoxNew(gtk4.OrientationHorizontal, 8)
	row.SetMarginTop(2)
	row.SetMarginBottom(2)

	lbl := gtk4.LabelNew(key)
	lbl.SetSizeRequest(220, -1)
	lbl.SetHAlign(gtk4.AlignStart)
	lbl.SetMarginEnd(8)
	lbl.SetTooltip(key)
	row.Append(&lbl.Widget)

	combo := gtk4.ComboBoxTextNewWithEntry()
	for _, t := range sd.targets {
		combo.Append(t, t)
	}
	combo.SetHExpand(true)

	if value != "" {
		targetIdx := -1
		for i, t := range sd.targets {
			if t == value {
				targetIdx = i
				break
			}
		}
		if targetIdx >= 0 {
			combo.SetActive(targetIdx)
		} else {
			combo.Append(value, value)
			combo.SetActive(len(sd.targets))
		}
	}

	row.Append(&combo.Widget)

	sd.combos[section+"/"+key] = combo

	return row
}

func (sd *ServiceDialog) createCheckRow(section, key, value string) *gtk4.Box {
	row := gtk4.BoxNew(gtk4.OrientationHorizontal, 8)
	row.SetMarginTop(2)
	row.SetMarginBottom(2)

	lbl := gtk4.LabelNew(key)
	lbl.SetSizeRequest(220, -1)
	lbl.SetHAlign(gtk4.AlignStart)
	lbl.SetMarginEnd(8)
	lbl.SetTooltip(key)
	row.Append(&lbl.Widget)

	check := gtk4.CheckButtonNew()
	check.SetActive(value == "yes" || value == "true" || value == "1")
	row.Append(&check.Widget)

	sd.checks[section+"/"+key] = check

	return row
}

func (sd *ServiceDialog) createCapSetRow(section, key, value string) *gtk4.Box {
	row := gtk4.BoxNew(gtk4.OrientationHorizontal, 8)
	row.SetMarginTop(4)
	row.SetMarginBottom(4)

	lbl := gtk4.LabelNew(key)
	lbl.SetSizeRequest(220, -1)
	lbl.SetHAlign(gtk4.AlignStart)
	lbl.SetMarginEnd(8)
	lbl.SetVAlign(gtk4.AlignStart)
	lbl.SetTooltip(key)
	row.Append(&lbl.Widget)

	activeSet := make(map[string]bool)
	for _, cap := range strings.Fields(value) {
		activeSet[strings.TrimSpace(cap)] = true
	}

	caps := linuxCapabilities()
	sort.Strings(caps)

	rightBox := gtk4.BoxNew(gtk4.OrientationVertical, 2)
	rightBox.SetHExpand(true)

	var checks []capCheck
	perRow := 4
	for start := 0; start < len(caps); start += perRow {
		end := start + perRow
		if end > len(caps) {
			end = len(caps)
		}
		cbRow := gtk4.BoxNew(gtk4.OrientationHorizontal, 4)
		for _, cap := range caps[start:end] {
			chk := gtk4.CheckButtonNewWithLabel(cap)
			chk.SetActive(activeSet[cap])
			chk.SetMarginStart(4)
			cbRow.Append(&chk.Widget)
			checks = append(checks, capCheck{check: chk, label: cap})
		}
		rightBox.Append(&cbRow.Widget)
	}

	row.Append(&rightBox.Widget)

	sd.capChecks[section+"/"+key] = checks

	return row
}

func (sd *ServiceDialog) onSave() {
	if sd.unitFile == nil {
		return
	}

	for entryKey, entry := range sd.entries {
		section, key := entryKeyToSectionKey(entryKey)
		value := strings.TrimSpace(entry.GetText())

		sd.unitFile.Set(section, key, value)
	}

	for comboKey, combo := range sd.combos {
		section, key := entryKeyToSectionKey(comboKey)
		id := combo.GetActiveID()
		value := strings.TrimSpace(combo.GetActiveText())
		if id != "" {
			value = id
		}

		sd.unitFile.Set(section, key, value)
	}

	for checkKey, check := range sd.checks {
		section, key := entryKeyToSectionKey(checkKey)
		if check.GetActive() {
			sd.unitFile.Set(section, key, "yes")
		} else {
			sd.unitFile.Set(section, key, "")
		}
	}

	for capKey, checks := range sd.capChecks {
		section, key := entryKeyToSectionKey(capKey)
		var enabled []string
		for _, c := range checks {
			if c.check.GetActive() {
				enabled = append(enabled, c.label)
			}
		}
		if len(enabled) > 0 {
			sd.unitFile.Set(section, key, strings.Join(enabled, " "))
		} else {
			sd.unitFile.Set(section, key, "")
		}
	}

	content := sd.unitFile.Serialize()

	if sd.isUser {
		path := sd.unitFile.Path
		if path == "" {
			path = sd.unitName
		}
		os.WriteFile(path, []byte(content), 0644)
	} else {
		path := sd.unitFile.Path
		if path == "" {
			path = sd.unitName
		}
		os.WriteFile("/tmp/systempanel-service.tmp", []byte(content), 0644)
		RunSudoCommand("cp", "/tmp/systempanel-service.tmp", path)
		os.Remove("/tmp/systempanel-service.tmp")
	}

	sd.win.Close()
}

func (sd *ServiceDialog) Present() {
	sd.win.Present()
}

func entryKeyToSectionKey(entryKey string) (string, string) {
	parts := strings.SplitN(entryKey, "/", 2)
	if len(parts) != 2 {
		return "", ""
	}
	return parts[0], parts[1]
}

func isStandardSection(name string) bool {
	switch name {
	case "Unit", "Service", "Socket", "Install", "Timer", "Mount", "Swap", "Path", "Slice", "Scope":
		return true
	}
	return false
}

func isRelationshipKey(section, key string) bool {
	relKeys := map[string]bool{
		"After": true, "Before": true, "Wants": true, "Requires": true,
		"Requisite": true, "BindsTo": true, "PartOf": true, "Upholds": true,
		"Conflicts": true, "OnFailure": true, "OnSuccess": true,
	}
	installKeys := map[string]bool{
		"WantedBy": true, "RequiredBy": true, "UpheldBy": true, "Also": true,
	}
	switch section {
	case "Unit":
		return relKeys[key]
	case "Install":
		return installKeys[key]
	}
	return false
}

func isBooleanKey(section, key string) bool {
	booleanKeys := map[string]bool{
		"RemainAfterExit": true, "GuessMainPID": true, "RootDirectoryStartOnly": true,
		"NonBlocking": true, "TTYReset": true, "TTYVHangup": true, "TTYVTDisallocate": true,
		"PrivateTmp": true, "PrivateDevices": true, "PrivateNetwork": true,
		"PrivateUsers": true, "PrivateMounts": true, "PrivateIPC": true,
		"ProtectHome": true, "ProtectSystem": true, "ProtectHostname": true,
		"ProtectKernelTunables": true, "ProtectKernelModules": true,
		"ProtectKernelLogs": true, "ProtectClock": true, "ProtectControlGroups": true,
		"LockPersonality": true, "MemoryDenyWriteExecute": true,
		"RestrictRealtime": true, "RestrictSUIDSGID": true, "RemoveIPC": true,
		"NoNewPrivileges": true, "DynamicUser": true,
		"IgnoreOnIsolate": true, "StopWhenUnneeded": true,
		"RefuseManualStart": true, "RefuseManualStop": true,
		"AllowIsolate": true, "DefaultDependencies": true,
		"MountAPIVFS": true, "ProtectProc": true,
		"IPAccounting": true, "CPUAccounting": true, "MemoryAccounting": true,
		"TasksAccounting": true, "IOAccounting": true,
	}
	return booleanKeys[key]
}

func linuxCapabilities() []string {
	return []string{
		"CAP_AUDIT_CONTROL", "CAP_AUDIT_READ", "CAP_AUDIT_WRITE",
		"CAP_BLOCK_SUSPEND", "CAP_BPF", "CAP_CHECKPOINT_RESTORE",
		"CAP_CHOWN", "CAP_DAC_OVERRIDE", "CAP_DAC_READ_SEARCH",
		"CAP_FOWNER", "CAP_FSETID", "CAP_IPC_LOCK", "CAP_IPC_OWNER",
		"CAP_KILL", "CAP_LEASE", "CAP_LINUX_IMMUTABLE",
		"CAP_MAC_ADMIN", "CAP_MAC_OVERRIDE", "CAP_MKNOD",
		"CAP_NET_ADMIN", "CAP_NET_BIND_SERVICE", "CAP_NET_BROADCAST",
		"CAP_NET_RAW", "CAP_PERFMON", "CAP_SETFCAP", "CAP_SETGID",
		"CAP_SETPCAP", "CAP_SETUID", "CAP_SYS_ADMIN", "CAP_SYS_BOOT",
		"CAP_SYS_CHROOT", "CAP_SYS_MODULE", "CAP_SYS_NICE",
		"CAP_SYS_PACCT", "CAP_SYS_PTRACE", "CAP_SYS_RAWIO",
		"CAP_SYS_RESOURCE", "CAP_SYS_TIME", "CAP_SYS_TTY_CONFIG",
		"CAP_SYSLOG", "CAP_WAKE_ALARM",
	}
}
