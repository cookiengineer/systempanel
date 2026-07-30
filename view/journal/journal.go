package journal

import (
	"fmt"
	"strings"

	"github.com/cookiengineer/systempanel/bindings/gtk4"
	"github.com/cookiengineer/systempanel/detect"
	"github.com/cookiengineer/systempanel/model"
	"github.com/cookiengineer/systempanel/view"
)

var Descriptor = view.ViewDescriptor{
	Name:     "journal",
	Title:    "Journal",
	IconName: "text-editor-symbolic",
	DetectFn: func() bool { return detect.HasProgram("journalctl") },
	Factory:  func() view.View { return NewJournalView() },
}

type JournalView struct {
	box     *gtk4.Box
	model   *model.JournalModel
	listBox *gtk4.ListBox
	rows    []*gtk4.ListBoxRow
}

func NewJournalView() *JournalView {
	jv := &JournalView{
		box:   gtk4.BoxNew(gtk4.OrientationVertical, 12),
		model: &model.JournalModel{},
	}
	jv.box.SetMarginStart(24)
	jv.box.SetMarginEnd(24)
	jv.box.SetMarginTop(24)
	jv.box.SetMarginBottom(24)

	header := gtk4.LabelNew("systemd Journal")
	header.AddCSSClass("header-label")
	jv.box.Append(&header.Widget)

	jv.listBox = gtk4.ListBoxNew()
	jv.listBox.SetSelectionMode(gtk4.SelectionNone)

	scrollW := gtk4.ScrolledWindowNew()
	scrollW.SetPolicy(gtk4.PolicyAutomatic, gtk4.PolicyAutomatic)
	scrollW.SetHExpand(true)
	scrollW.SetVExpand(true)
	scrollW.SetChild(&jv.listBox.Widget)
	jv.box.Append(&scrollW.Widget)

	btnBox := gtk4.BoxNew(gtk4.OrientationHorizontal, 8)
	btnBox.SetMarginTop(4)

	refreshBtn := gtk4.ButtonNewWithLabel("Refresh")
	refreshBtn.OnClicked(func() { jv.refresh() })
	btnBox.Append(&refreshBtn.Widget)

	jv.box.Append(&btnBox.Widget)

	jv.refresh()

	return jv
}

func (jv *JournalView) refresh() {
	for _, r := range jv.rows {
		jv.listBox.Remove(r)
	}
	jv.rows = jv.rows[:0]

	entries, err := jv.model.Fetch(100)
	if err != nil {
		return
	}

	for _, e := range entries {
		row := jv.createEntryRow(e)
		jv.listBox.Append(row)
		jv.rows = append(jv.rows, row)
	}
}

func (jv *JournalView) createEntryRow(e model.JournalEntry) *gtk4.ListBoxRow {
	row := gtk4.ListBoxRowNew()

	vbox := gtk4.BoxNew(gtk4.OrientationVertical, 2)
	vbox.SetMarginStart(8)
	vbox.SetMarginEnd(8)
	vbox.SetMarginTop(2)
	vbox.SetMarginBottom(2)

	priorityClass := fmt.Sprintf("journal-%s", e.PriorityName)
	if priorityClass == "journal-" {
		priorityClass = "journal-info"
	}

	topBox := gtk4.BoxNew(gtk4.OrientationHorizontal, 8)

	timeStr := e.Timestamp.Format("15:04:05")
	timeLabel := gtk4.LabelNew(timeStr)
	timeLabel.SetSensitive(false)
	topBox.Append(&timeLabel.Widget)

	if e.Unit != "" {
		unitLabel := gtk4.LabelNew(e.Unit)
		unitLabel.AddCSSClass(priorityClass)
		topBox.Append(&unitLabel.Widget)
	}

	vbox.Append(&topBox.Widget)

	message := strings.TrimSpace(e.Message)
	if len(message) > 300 {
		message = message[:300] + "..."
	}
	msgLabel := gtk4.LabelNew(message)
	msgLabel.SetHAlign(gtk4.AlignStart)
	msgLabel.SetWrap(true)
	msgLabel.AddCSSClass(priorityClass)
	vbox.Append(&msgLabel.Widget)

	row.SetChild(&vbox.Widget)

	return row
}

func (jv *JournalView) Widget() *gtk4.Widget { return &jv.box.Widget }
func (jv *JournalView) Name() string          { return "journal" }
func (jv *JournalView) Title() string         { return "Journal" }
func (jv *JournalView) IconName() string      { return "text-editor-symbolic" }
func (jv *JournalView) OnShow()               { jv.refresh() }
func (jv *JournalView) OnHide()               {}
func (jv *JournalView) Destroy()              {}
