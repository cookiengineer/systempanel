package widget

import (
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"

	"github.com/cookiengineer/systempanel/bindings/gtk4"
	"github.com/cookiengineer/systempanel/parsers/systemdunit"
)

type tagPill struct {
	box    *gtk4.Box
	label  string
}

type tagListWidget struct {
	combo     *gtk4.ComboBoxText
	flow      *gtk4.Box
	pills     []tagPill
	options   []string
	onChanged func()
}

func newTagListWidget(options []string) *tagListWidget {
	tl := &tagListWidget{
		options: options,
	}
	tl.combo = gtk4.ComboBoxTextNewWithEntry()
	tl.combo.SetHExpand(true)
	for _, o := range options {
		tl.combo.Append(o, o)
	}
	tl.combo.OnChanged(func() {
		id := tl.combo.GetActiveID()
		if id != "" {
			tl.addPill(id)
			tl.combo.SetActive(-1)
		}
	})

	tl.flow = gtk4.BoxNew(gtk4.OrientationHorizontal, 4)
	tl.flow.SetHExpand(true)
	tl.flow.SetMarginTop(2)

	return tl
}

func (tl *tagListWidget) rows() []*gtk4.Widget {
	return []*gtk4.Widget{&tl.combo.Widget, &tl.flow.Widget}
}

func (tl *tagListWidget) addPill(tag string) {
	for _, p := range tl.pills {
		if p.label == tag {
			return
		}
	}

	pillBox := gtk4.BoxNew(gtk4.OrientationHorizontal, 2)
	pillBox.AddCSSClass("tag-pill")
	pillBox.SetMarginStart(2)

	lbl := gtk4.LabelNew(tag)
	pillBox.Append(&lbl.Widget)

	rmBtn := gtk4.ButtonNew()
	rmBtn.SetIconName("window-close-symbolic")
	rmBtn.AddCSSClass("flat")
	pill := tagPill{box: pillBox, label: tag}
	rmBtn.OnClicked(func() {
		tl.removePill(pill)
	})
	rmBtnW := rmBtn.Widget
	pillBox.Append(&rmBtnW)

	tl.pills = append(tl.pills, pill)
	pbWidget := pillBox.Widget
	tl.flow.Append(&pbWidget)

	if tl.onChanged != nil {
		tl.onChanged()
	}
}

func (tl *tagListWidget) removePill(pill tagPill) {
	for i, p := range tl.pills {
		if p.label == pill.label {
			tl.flow.Remove(&pill.box.Widget)
			tl.pills = append(tl.pills[:i], tl.pills[i+1:]...)
			if tl.onChanged != nil {
				tl.onChanged()
			}
			return
		}
	}
}

func (tl *tagListWidget) setValue(value string) {
	for _, p := range tl.pills {
		tl.flow.Remove(&p.box.Widget)
	}
	tl.pills = nil

	if value == "" {
		return
	}
	for _, tag := range strings.Fields(value) {
		tag = strings.TrimSpace(tag)
		if tag != "" {
			tl.addPillInner(tag)
		}
	}
}

func (tl *tagListWidget) addPillInner(tag string) {
	pillBox := gtk4.BoxNew(gtk4.OrientationHorizontal, 2)
	pillBox.AddCSSClass("tag-pill")
	pillBox.SetMarginStart(2)

	lbl := gtk4.LabelNew(tag)
	pillBox.Append(&lbl.Widget)

	rmBtn := gtk4.ButtonNew()
	rmBtn.SetIconName("window-close-symbolic")
	rmBtn.AddCSSClass("flat")
	pill := tagPill{box: pillBox, label: tag}
	rmBtn.OnClicked(func() {
		tl.removePill(pill)
	})
	rmBtnW := rmBtn.Widget
	pillBox.Append(&rmBtnW)

	tl.pills = append(tl.pills, pill)
	pbWidget := pillBox.Widget
	tl.flow.Append(&pbWidget)
}

func (tl *tagListWidget) getValue() string {
	var tags []string
	for _, p := range tl.pills {
		tags = append(tags, p.label)
	}
	return strings.Join(tags, " ")
}

func (tl *tagListWidget) SetOnChanged(fn func()) {
	tl.onChanged = fn
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

	entries      map[string]*gtk4.Entry
	textViews    map[string]*gtk4.TextView
	checks       map[string]*gtk4.CheckButton
	combos       map[string]*gtk4.ComboBoxText
	spins        map[string]*gtk4.SpinButton
	tagLists     map[string]*tagListWidget
	targets      []string
}

func NewServiceDialog(parent *gtk4.Window, uf *systemdunit.UnitFile, unitName string, isUser bool) *ServiceDialog {
	sd := &ServiceDialog{
		parentWin: parent,
		unitFile:  uf,
		unitName:  unitName,
		isUser:    isUser,
		entries:   make(map[string]*gtk4.Entry),
		textViews: make(map[string]*gtk4.TextView),
		checks:    make(map[string]*gtk4.CheckButton),
		combos:    make(map[string]*gtk4.ComboBoxText),
		spins:     make(map[string]*gtk4.SpinButton),
		tagLists:  make(map[string]*tagListWidget),
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
	sort.Strings(sd.targets)
}

func (sd *ServiceDialog) build() {
	sd.win = gtk4.WindowNew()
	sd.win.SetTitle("Service Settings - " + sd.unitName)
	sd.win.SetDefaultSize(860, 680)
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

func (sd *ServiceDialog) buildSectionTab(sectionName string) *gtk4.Box {
	box := gtk4.BoxNew(gtk4.OrientationVertical, 0)
	box.SetHExpand(true)
	box.SetVExpand(true)

	scrollW := gtk4.ScrolledWindowNew()
	scrollW.SetPolicy(gtk4.PolicyNever, gtk4.PolicyAutomatic)
	scrollW.SetHExpand(true)
	scrollW.SetVExpand(true)

	innerBox := gtk4.BoxNew(gtk4.OrientationVertical, 0)
	innerBox.SetMarginStart(24)
	innerBox.SetMarginEnd(24)
	innerBox.SetMarginTop(12)
	innerBox.SetMarginBottom(12)

	groups := systemdunit.GetSectionGroups(sectionName)
	for _, group := range groups {
		if !sd.advanced && !sd.groupHasValue(sectionName, group) {
			continue
		}
		sd.buildGroupBox(innerBox, sectionName, group)
	}

	scrollW.SetChild(&innerBox.Widget)
	box.Append(&scrollW.Widget)

	return box
}

func (sd *ServiceDialog) propertyLabel(key string) *gtk4.Label {
	lbl := gtk4.LabelNew(key)
	lbl.SetSizeRequest(210, -1)
	lbl.SetHAlign(gtk4.AlignStart)
	lbl.SetXAlign(0)
	lbl.SetMarginEnd(8)
	lbl.SetTooltip(key)
	return lbl
}

func (sd *ServiceDialog) groupHasValue(sectionName string, group systemdunit.PropertyGroup) bool {
	for _, key := range group.Keys {
		value := sd.unitFile.Get(sectionName, key)
		if value != "" {
			return true
		}
	}
	return false
}

func (sd *ServiceDialog) buildGroupBox(parent *gtk4.Box, sectionName string, group systemdunit.PropertyGroup) {
	groupBox := gtk4.BoxNew(gtk4.OrientationVertical, 4)
	groupBox.SetMarginTop(8)
	groupBox.SetMarginBottom(8)

	header := gtk4.LabelNew(group.Name)
	header.AddCSSClass("group-header")
	header.SetHAlign(gtk4.AlignStart)
	header.SetMarginBottom(4)
	groupBox.Append(&header.Widget)

	hasContent := false
	for _, key := range group.Keys {
		meta, ok := systemdunit.GetPropertyMeta(sectionName, key)
		if !ok {
			continue
		}
		value := sd.unitFile.Get(sectionName, key)
		row := sd.buildPropertyRow(sectionName, key, value, meta)
		if row != nil {
			groupBox.Append(&row.Widget)
			hasContent = true
		}
	}

	if !hasContent {
		return
	}

	groupBoxW := groupBox.Widget
	parent.Append(&groupBoxW)
}

func (sd *ServiceDialog) buildPropertyRow(section, key, value string, meta systemdunit.PropertyMeta) *gtk4.Box {
	switch meta.PropType {
	case systemdunit.PropBoolean:
		return sd.createCheckRow(section, key, value)
	case systemdunit.PropEnum:
		return sd.createEnumComboRow(section, key, value, meta.EnumValues)
	case systemdunit.PropNumber:
		return sd.createSpinRow(section, key, value, meta.Min, meta.Max, meta.Step)
	case systemdunit.PropFilePath:
		return sd.createFilePathRow(section, key, value, meta.MimeTypes)
	case systemdunit.PropFolderPath:
		return sd.createFolderPathRow(section, key, value)
	case systemdunit.PropMultiLine:
		return sd.createMultiLineRow(section, key, value)
	case systemdunit.PropTagList:
		return sd.createTagListRow(section, key, value)
	case systemdunit.PropTarget:
		return sd.createTargetRow(section, key, value)
	default:
		return sd.createEntryRow(section, key, value)
	}
}

func (sd *ServiceDialog) createEntryRow(section, key, value string) *gtk4.Box {
	row := gtk4.BoxNew(gtk4.OrientationHorizontal, 8)
	row.SetMarginTop(2)
	row.SetMarginBottom(2)

	lbl := sd.propertyLabel(key)
	row.Append(&lbl.Widget)

	entry := gtk4.EntryNew()
	entry.SetText(value)
	entry.SetHExpand(true)
	row.Append(&entry.Widget)

	sd.entries[section+"/"+key] = entry
	return row
}

func (sd *ServiceDialog) createCheckRow(section, key, value string) *gtk4.Box {
	row := gtk4.BoxNew(gtk4.OrientationHorizontal, 8)
	row.SetMarginTop(2)
	row.SetMarginBottom(2)

	lbl := sd.propertyLabel(key)
	row.Append(&lbl.Widget)

	check := gtk4.CheckButtonNew()
	check.SetActive(value == "yes" || value == "true" || value == "1")
	row.Append(&check.Widget)

	sd.checks[section+"/"+key] = check
	return row
}

func (sd *ServiceDialog) createEnumComboRow(section, key, value string, enumValues []string) *gtk4.Box {
	row := gtk4.BoxNew(gtk4.OrientationHorizontal, 8)
	row.SetMarginTop(2)
	row.SetMarginBottom(2)

	lbl := sd.propertyLabel(key)
	row.Append(&lbl.Widget)

	combo := gtk4.ComboBoxTextNew()
	combo.SetHExpand(true)

	selectedIdx := 0
	for i, ev := range enumValues {
		combo.Append(ev, ev)
		if ev == value {
			selectedIdx = i
		}
	}
	combo.SetActive(selectedIdx)

	row.Append(&combo.Widget)

	sd.combos[section+"/"+key] = combo
	return row
}

func (sd *ServiceDialog) createSpinRow(section, key, value string, min, max, step float64) *gtk4.Box {
	row := gtk4.BoxNew(gtk4.OrientationHorizontal, 8)
	row.SetMarginTop(2)
	row.SetMarginBottom(2)

	lbl := sd.propertyLabel(key)
	row.Append(&lbl.Widget)

	if step == 0 {
		step = 1
	}
	spinMax := max
	if spinMax == 0 && min == 0 {
		spinMax = 4294967295
	}
	spin := gtk4.SpinButtonNew(min, spinMax, step)
	if value != "" {
		if v, err := strconv.ParseFloat(value, 64); err == nil {
			spin.SetValue(v)
		}
	}
	spin.SetHExpand(true)
	row.Append(&spin.Widget)

	sd.spins[section+"/"+key] = spin
	return row
}

func (sd *ServiceDialog) createFilePathRow(section, key, value string, mimeTypes []string) *gtk4.Box {
	row := gtk4.BoxNew(gtk4.OrientationHorizontal, 8)
	row.SetMarginTop(2)
	row.SetMarginBottom(2)

	lbl := sd.propertyLabel(key)
	row.Append(&lbl.Widget)

	entry := gtk4.EntryNew()
	entry.SetText(value)
	entry.SetHExpand(true)
	row.Append(&entry.Widget)

	browseBtn := gtk4.ButtonNewWithLabel("Browse...")
	browseBtn.AddCSSClass("flat")
	entryRef := entry
	parentWin := sd.win
	browseBtn.OnClicked(func() {
		dialog := gtk4.FileDialogNew()
		dialog.SetTitle("Select File - " + key)
		if len(mimeTypes) > 0 {
			filter := gtk4.FileFilterNew()
			for _, mime := range mimeTypes {
				filter.AddMimeType(mime)
			}
			filter.SetName("Supported Files")
			dialog.SetDefaultFilter(filter)
		}
		dialog.Open(parentWin, func(path string, err string) {
			if path != "" {
				entryRef.SetText(path)
			}
		})
	})
	btnW := browseBtn.Widget
	row.Append(&btnW)

	sd.entries[section+"/"+key] = entry
	return row
}

func (sd *ServiceDialog) createFolderPathRow(section, key, value string) *gtk4.Box {
	row := gtk4.BoxNew(gtk4.OrientationHorizontal, 8)
	row.SetMarginTop(2)
	row.SetMarginBottom(2)

	lbl := sd.propertyLabel(key)
	row.Append(&lbl.Widget)

	entry := gtk4.EntryNew()
	entry.SetText(value)
	entry.SetHExpand(true)
	row.Append(&entry.Widget)

	browseBtn := gtk4.ButtonNewWithLabel("Browse...")
	browseBtn.AddCSSClass("flat")
	entryRef := entry
	parentWin := sd.win
	browseBtn.OnClicked(func() {
		dialog := gtk4.FileDialogNew()
		dialog.SetTitle("Select Folder - " + key)
		dialog.SelectFolder(parentWin, func(path string, err string) {
			if path != "" {
				entryRef.SetText(path)
			}
		})
	})
	btnW := browseBtn.Widget
	row.Append(&btnW)

	sd.entries[section+"/"+key] = entry
	return row
}

func (sd *ServiceDialog) createMultiLineRow(section, key, value string) *gtk4.Box {
	row := gtk4.BoxNew(gtk4.OrientationHorizontal, 8)
	row.SetMarginTop(2)
	row.SetMarginBottom(2)

	lbl := sd.propertyLabel(key)
	lbl.SetVAlign(gtk4.AlignStart)
	lbl.SetMarginTop(4)
	row.Append(&lbl.Widget)

	scrollFrame := gtk4.ScrolledWindowNew()
	scrollFrame.SetPolicy(gtk4.PolicyAutomatic, gtk4.PolicyAutomatic)
	scrollFrame.SetHExpand(true)
	scrollFrame.SetSizeRequest(-1, 80)

	tv := gtk4.TextViewNew()
	tv.SetText(value)
	tv.SetHExpand(true)
	tv.SetVExpand(true)
	scrollFrame.SetChild(&tv.Widget)

	row.Append(&scrollFrame.Widget)

	sd.textViews[section+"/"+key] = tv
	return row
}

func (sd *ServiceDialog) createTagListRow(section, key, value string) *gtk4.Box {
	row := gtk4.BoxNew(gtk4.OrientationHorizontal, 8)
	row.SetMarginTop(4)
	row.SetMarginBottom(4)

	lbl := sd.propertyLabel(key)
	lbl.SetVAlign(gtk4.AlignStart)
	lbl.SetMarginTop(2)
	row.Append(&lbl.Widget)

	rightBox := gtk4.BoxNew(gtk4.OrientationVertical, 4)
	rightBox.SetHExpand(true)

	options := systemdunit.GetTagOptions(key)
	tl := newTagListWidget(options)
	tl.setValue(value)

	for _, w := range tl.rows() {
		rightBox.Append(w)
	}

	sd.tagLists[section+"/"+key] = tl

	rightBoxW := rightBox.Widget
	row.Append(&rightBoxW)
	return row
}

func (sd *ServiceDialog) createTargetRow(section, key, value string) *gtk4.Box {
	row := gtk4.BoxNew(gtk4.OrientationHorizontal, 8)
	row.SetMarginTop(2)
	row.SetMarginBottom(2)

	lbl := sd.propertyLabel(key)
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

func (sd *ServiceDialog) onToggleMode() {
	sd.advanced = !sd.advanced
	if sd.advanced {
		sd.modeBtn.SetLabel("Simple")
	} else {
		sd.modeBtn.SetLabel("Advanced")
	}

	sections := []string{"Unit", "Service", "Socket", "Install"}
	for _, secName := range sections {
		child := sd.stack.GetChildByName(secName)
		if child != nil {
			sd.stack.Remove(child)
		}
	}

	sd.entries = make(map[string]*gtk4.Entry)
	sd.textViews = make(map[string]*gtk4.TextView)
	sd.checks = make(map[string]*gtk4.CheckButton)
	sd.combos = make(map[string]*gtk4.ComboBoxText)
	sd.spins = make(map[string]*gtk4.SpinButton)
	sd.tagLists = make(map[string]*tagListWidget)

	visibleName := sd.stack.GetVisibleChildName()
	for _, secName := range sections {
		newBox := sd.buildSectionTab(secName)
		sd.stack.AddTitled(&newBox.Widget, secName, secName)
	}

	if visibleName != "" {
		sd.stack.SetVisibleChildName(visibleName)
	}
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

	for tvKey, tv := range sd.textViews {
		section, key := entryKeyToSectionKey(tvKey)
		value := strings.TrimSpace(tv.GetText())
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

	for comboKey, combo := range sd.combos {
		section, key := entryKeyToSectionKey(comboKey)
		id := combo.GetActiveID()
		value := strings.TrimSpace(combo.GetActiveText())
		if id != "" {
			value = id
		}
		sd.unitFile.Set(section, key, value)
	}

	for spinKey, spin := range sd.spins {
		section, key := entryKeyToSectionKey(spinKey)
		v := spin.GetValue()
		sd.unitFile.Set(section, key, strconv.FormatFloat(v, 'f', -1, 64))
	}

	for tagKey, tl := range sd.tagLists {
		section, key := entryKeyToSectionKey(tagKey)
		value := tl.getValue()
		sd.unitFile.Set(section, key, value)
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
