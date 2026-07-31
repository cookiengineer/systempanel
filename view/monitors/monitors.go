package monitors

import (
	"fmt"

	"github.com/cookiengineer/systempanel/bindings/gtk4"
	"github.com/cookiengineer/systempanel/detect"
	"github.com/cookiengineer/systempanel/model"
	"github.com/cookiengineer/systempanel/view"
)

var Descriptor = view.ViewDescriptor{
	Name:     "monitors",
	Title:    "Monitors",
	IconName: "video-display-symbolic",
	DetectFn: func() bool {
		return detect.HasProgram("xrandr")
	},
	Factory: func() view.View { return NewMonitorsView() },
}

type monitorItem struct {
	monitor      model.Monitor
	resCombo     *gtk4.ComboBoxText
	relationCombo *gtk4.ComboBoxText
	targetCombo  *gtk4.ComboBoxText
}

type MonitorsView struct {
	box         *gtk4.Box
	model       *model.MonitorModel
	listBox     *gtk4.ListBox
	rows        []*gtk4.ListBoxRow
	items       []monitorItem
}

func NewMonitorsView() *MonitorsView {
	mv := &MonitorsView{
		box:   gtk4.BoxNew(gtk4.OrientationVertical, 12),
		model: &model.MonitorModel{},
	}
	mv.box.SetMarginStart(24)
	mv.box.SetMarginEnd(24)
	mv.box.SetMarginTop(24)
	mv.box.SetMarginBottom(24)

	header := gtk4.LabelNew("Monitors")
	header.AddCSSClass("header-label")
	headerWidget := header.Widget
	mv.box.Append(&headerWidget)

	desc := gtk4.LabelNew("Configure resolution and arrangement of connected displays.")
	desc.SetWrap(true)
	desc.SetMarginBottom(12)
	dv := desc.Widget
	mv.box.Append(&dv)

	mv.listBox = gtk4.ListBoxNew()
	mv.listBox.SetSelectionMode(gtk4.SelectionNone)

	scrollW := gtk4.ScrolledWindowNew()
	scrollW.SetPolicy(gtk4.PolicyNever, gtk4.PolicyAutomatic)
	scrollW.SetHExpand(true)
	scrollW.SetVExpand(true)

	listWidget := mv.listBox.Widget
	scrollW.SetChild(&listWidget)
	sw := scrollW.Widget
	mv.box.Append(&sw)

	btnBox := gtk4.BoxNew(gtk4.OrientationHorizontal, 8)
	btnBox.SetMarginTop(4)

	refreshBtn := gtk4.ButtonNewWithLabel("Refresh")
	refreshBtn.OnClicked(func() { mv.refresh() })
	rbw := refreshBtn.Widget
	btnBox.Append(&rbw)

	spacer := gtk4.LabelNew("")
	spacer.SetHExpand(true)
	spWidget := spacer.Widget
	btnBox.Append(&spWidget)

	applyBtn := gtk4.ButtonNewWithLabel("Apply")
	applyBtn.AddCSSClass("suggested-action")
	applyBtn.OnClicked(func() { mv.apply() })
	abWidget := applyBtn.Widget
	btnBox.Append(&abWidget)

	btnBoxWidget := btnBox.Widget
	mv.box.Append(&btnBoxWidget)

	mv.refresh()

	return mv
}

func (mv *MonitorsView) refresh() {
	for _, r := range mv.rows {
		mv.listBox.Remove(r)
	}
	mv.rows = mv.rows[:0]
	mv.items = mv.items[:0]

	monitors, _ := mv.model.ListMonitors()

	connected := []model.Monitor{}
	for _, m := range monitors {
		if m.Connected {
			connected = append(connected, m)
		}
	}
	if len(connected) == 0 {
		label := gtk4.LabelNew("No monitors detected via xrandr")
		label.SetSensitive(false)
		label.SetHAlign(gtk4.AlignCenter)
		lw := label.Widget
		row := gtk4.ListBoxRowNew()
		row.SetChild(&lw)
		mv.listBox.Append(row)
		mv.rows = append(mv.rows, row)
		return
	}

	resHdr := mv.createSectionLabel("─ Resolution ─")
	mv.listBox.Append(resHdr)
	mv.rows = append(mv.rows, resHdr)

	for _, m := range connected {
		row := mv.createMonitorRow(m, connected)
		mv.listBox.Append(row)
		mv.rows = append(mv.rows, row)
	}

	arrHdr := mv.createSectionLabel("─ Arrangement ─")
	mv.listBox.Append(arrHdr)
	mv.rows = append(mv.rows, arrHdr)

	if len(connected) < 2 {
		msg := mv.createMessageRow("Connect a second monitor to configure arrangement")
		mv.listBox.Append(msg)
		mv.rows = append(mv.rows, msg)
		return
	}

	for _, m := range connected {
		row := mv.createArrangementRow(m, connected)
		mv.listBox.Append(row)
		mv.rows = append(mv.rows, row)
	}
}

func (mv *MonitorsView) createSectionLabel(text string) *gtk4.ListBoxRow {
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

func (mv *MonitorsView) createMessageRow(text string) *gtk4.ListBoxRow {
	row := gtk4.ListBoxRowNew()
	row.SetSensitive(false)
	label := gtk4.LabelNew(text)
	label.SetSensitive(false)
	label.SetHAlign(gtk4.AlignCenter)
	label.SetMarginTop(8)
	label.SetMarginBottom(8)
	lw := label.Widget
	row.SetChild(&lw)
	return row
}

func (mv *MonitorsView) createMonitorRow(m model.Monitor, all []model.Monitor) *gtk4.ListBoxRow {
	row := gtk4.ListBoxRowNew()
	row.AddCSSClass("device-row")

	hbox := gtk4.BoxNew(gtk4.OrientationHorizontal, 8)

	iconName := "video-display-symbolic"
	if m.Primary {
		iconName = "starred-symbolic"
	}
	icon := gtk4.ImageNewFromIconName(iconName)
	icon.SetPixelSize(20)
	iconWidget := icon.Widget
	hbox.Append(&iconWidget)

	infoBox := gtk4.BoxNew(gtk4.OrientationVertical, 2)
	displayName := m.Name
	if m.Primary {
		displayName += " (primary)"
	}
	nameLabel := gtk4.LabelNew(displayName)
	nameLabel.SetHAlign(gtk4.AlignStart)
	nlWidget := nameLabel.Widget
	infoBox.Append(&nlWidget)

	resLabel := gtk4.LabelNew(m.Resolution)
	resLabel.SetSensitive(false)
	resLabel.SetHAlign(gtk4.AlignStart)
	rlWidget := resLabel.Widget
	infoBox.Append(&rlWidget)

	infoBoxWidget := infoBox.Widget
	infoBox.SetHExpand(true)
	hbox.Append(&infoBoxWidget)

	combo := gtk4.ComboBoxTextNew()
	combo.SetHExpand(true)
	selectedIdx := 0
	for i, mode := range m.Modes {
		label := fmt.Sprintf("%dx%d", mode.Width, mode.Height)
		if mode.Preferred {
			label += " (preferred)"
		}
		if mode.Refresh > 0 {
			label += fmt.Sprintf(" %.1f Hz", mode.Refresh)
		}
		modeKey := fmt.Sprintf("%dx%d", mode.Width, mode.Height)
		combo.Append(modeKey, label)
		if mode.Current {
			selectedIdx = i
		}
	}
	combo.SetActive(selectedIdx)
	cw := combo.Widget
	hbox.Append(&cw)

	mv.items = append(mv.items, monitorItem{monitor: m, resCombo: combo})

	hboxWidget := hbox.Widget
	row.SetChild(&hboxWidget)

	return row
}

func (mv *MonitorsView) createArrangementRow(m model.Monitor, all []model.Monitor) *gtk4.ListBoxRow {
	row := gtk4.ListBoxRowNew()
	row.AddCSSClass("device-row")

	hbox := gtk4.BoxNew(gtk4.OrientationHorizontal, 8)
	hbox.SetHAlign(gtk4.AlignCenter)

	nameLabel := gtk4.LabelNew(m.Name)
	nlWidget := nameLabel.Widget
	hbox.Append(&nlWidget)

	relationCombo := gtk4.ComboBoxTextNew()
	relationCombo.Append("left-of", "Left of")
	relationCombo.Append("right-of", "Right of")
	relationCombo.Append("above", "Above")
	relationCombo.Append("below", "Below")
	relationCombo.Append("same-as", "Same as")
	relationCombo.SetActive(0)
	relationCombo.SetSizeRequest(100, -1)
	rcWidget := relationCombo.Widget
	hbox.Append(&rcWidget)

	targetCombo := gtk4.ComboBoxTextNew()
	for _, other := range all {
		if other.Name != m.Name {
			targetCombo.Append(other.Name, other.Name)
		}
	}
	targetCombo.SetActive(0)
	targetCombo.SetSizeRequest(100, -1)
	tcWidget := targetCombo.Widget
	hbox.Append(&tcWidget)

	idx := -1
	for i, item := range mv.items {
		if item.monitor.Name == m.Name {
			idx = i
			break
		}
	}
	if idx >= 0 {
		mv.items[idx].relationCombo = relationCombo
		mv.items[idx].targetCombo = targetCombo
	}

	hboxWidget := hbox.Widget
	row.SetChild(&hboxWidget)

	return row
}

func (mv *MonitorsView) apply() {
	for _, item := range mv.items {
		if item.resCombo != nil {
			mode := item.resCombo.GetActiveID()
			if mode != "" {
				mv.model.SetResolution(item.monitor.Name, mode)
			}
		}
		if item.relationCombo != nil && item.targetCombo != nil {
			relation := item.relationCombo.GetActiveID()
			target := item.targetCombo.GetActiveID()
			if relation != "" && target != "" && relation != "same-as" {
				mv.model.SetPosition(item.monitor.Name, relation, target)
			} else if relation == "same-as" && target != "" {
				mv.model.SetPosition(item.monitor.Name, "same-as", target)
			}
		}
	}
}

func (mv *MonitorsView) Widget() *gtk4.Widget { return &mv.box.Widget }
func (mv *MonitorsView) Name() string          { return "monitors" }
func (mv *MonitorsView) Title() string         { return "Monitors" }
func (mv *MonitorsView) IconName() string      { return "video-display-symbolic" }
func (mv *MonitorsView) OnShow()               { mv.refresh() }
func (mv *MonitorsView) OnHide()               {}
func (mv *MonitorsView) Destroy()              {}
