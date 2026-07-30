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
	entries []model.JournalEntry
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

	header := gtk4.LabelNew("System Journal")
	header.AddCSSClass("header-label")
	headerWidget := header.Widget
	jv.box.Append(&headerWidget)

	jv.listBox = gtk4.ListBoxNew()
	jv.listBox.SetSelectionMode(gtk4.SelectionNone)

	scrollW := gtk4.ScrolledWindowNew()
	scrollW.SetPolicy(gtk4.PolicyAutomatic, gtk4.PolicyAutomatic)
	scrollW.SetHExpand(true)
	scrollW.SetVExpand(true)

	scrollWidget := scrollW.Widget
	listWidget := jv.listBox.Widget
	scrollW.SetChild(&listWidget)
	jv.box.Append(&scrollWidget)

	btnBox := gtk4.BoxNew(gtk4.OrientationHorizontal, 8)

	refreshBtn := gtk4.ButtonNewWithLabel("Refresh")
	refreshBtn.OnClicked(func() { jv.refresh() })
	rbw := refreshBtn.Widget
	btnBox.Append(&rbw)

	clearBtn := gtk4.ButtonNewWithLabel("Clear")
	clearBtn.OnClicked(func() {
		jv.entries = nil
		jv.refreshList()
	})
	cbw := clearBtn.Widget
	btnBox.Append(&cbw)

	btnBoxWidget := btnBox.Widget
	jv.box.Append(&btnBoxWidget)

	jv.refresh()

	return jv
}

func (jv *JournalView) refresh() {
	entries, err := jv.model.Fetch(100)
	if err != nil {
		return
	}
	jv.entries = entries
	jv.refreshList()
}

func (jv *JournalView) refreshList() {
	for {
		row := jv.listBox.GetSelectedRow()
		if row == nil {
			break
		}
		jv.listBox.Remove(row)
	}

	for _, e := range jv.entries {
		row := jv.createEntryRow(e)
		jv.listBox.Append(row)
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
	tlWidget := timeLabel.Widget
	topBox.Append(&tlWidget)

	if e.Unit != "" {
		unitLabel := gtk4.LabelNew(e.Unit)
		unitLabel.AddCSSClass(priorityClass)
		ulWidget := unitLabel.Widget
		topBox.Append(&ulWidget)
	}

	topBoxWidget := topBox.Widget
	vbox.Append(&topBoxWidget)

	message := strings.TrimSpace(e.Message)
	if len(message) > 200 {
		message = message[:200] + "..."
	}
	msgLabel := gtk4.LabelNew(message)
	msgLabel.SetHAlign(gtk4.AlignStart)
	msgLabel.SetWrap(true)
	msgLabel.AddCSSClass(priorityClass)
	mlWidget := msgLabel.Widget
	vbox.Append(&mlWidget)

	vboxWidget := vbox.Widget
	row.SetChild(&vboxWidget)

	return row
}

func (jv *JournalView) Widget() *gtk4.Widget { return &jv.box.Widget }
func (jv *JournalView) Name() string          { return "journal" }
func (jv *JournalView) Title() string         { return "Journal" }
func (jv *JournalView) IconName() string      { return "text-editor-symbolic" }
func (jv *JournalView) OnShow()               { jv.refresh() }
func (jv *JournalView) OnHide()               {}
func (jv *JournalView) Destroy()              {}
