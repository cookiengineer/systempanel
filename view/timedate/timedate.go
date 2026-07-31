package timedate

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/cookiengineer/systempanel/bindings/gtk4"
	"github.com/cookiengineer/systempanel/detect"
	"github.com/cookiengineer/systempanel/model"
	"github.com/cookiengineer/systempanel/view"
	"github.com/cookiengineer/systempanel/widget"
)

var Descriptor = view.ViewDescriptor{
	Name:     "timedate",
	Title:    "Time & Date",
	IconName: "preferences-system-time-symbolic",
	DetectFn: func() bool { return detect.HasProgram("timedatectl") },
	Factory:  func() view.View { return NewTimeDateView() },
}

type TimeDateView struct {
	box     *gtk4.Box
	model   *model.TimeDateModel
	listBox *gtk4.ListBox
	rows    []*gtk4.ListBoxRow

	timeLabel   *gtk4.Label
	ntpCheck    *gtk4.CheckButton
	tzCombo     *gtk4.ComboBoxText
	yearSpin    *gtk4.SpinButton
	monthSpin   *gtk4.SpinButton
	daySpin     *gtk4.SpinButton
	hourSpin    *gtk4.SpinButton
	minSpin     *gtk4.SpinButton
	secSpin     *gtk4.SpinButton

	zones        []string
	settingCombo  bool
	parentWin    *gtk4.Window
}

func NewTimeDateView() *TimeDateView {
	tv := &TimeDateView{
		box:   gtk4.BoxNew(gtk4.OrientationVertical, 12),
		model: &model.TimeDateModel{},
	}
	tv.box.SetMarginStart(24)
	tv.box.SetMarginEnd(24)
	tv.box.SetMarginTop(24)
	tv.box.SetMarginBottom(24)

	header := gtk4.LabelNew("Time & Date")
	header.AddCSSClass("header-label")
	tv.box.Append(&header.Widget)

	desc := gtk4.LabelNew("View system time, NTP synchronization status, and change the timezone.")
	desc.SetWrap(true)
	desc.SetMarginBottom(12)
	tv.box.Append(&desc.Widget)

	tv.listBox = gtk4.ListBoxNew()
	tv.listBox.SetSelectionMode(gtk4.SelectionNone)

	scrollW := gtk4.ScrolledWindowNew()
	scrollW.SetPolicy(gtk4.PolicyNever, gtk4.PolicyAutomatic)
	scrollW.SetHExpand(true)
	scrollW.SetVExpand(true)
	scrollW.SetChild(&tv.listBox.Widget)
	tv.box.Append(&scrollW.Widget)

	tv.buildRows()

	btnBox := gtk4.BoxNew(gtk4.OrientationHorizontal, 8)
	btnBox.SetMarginTop(4)

	refreshBtn := gtk4.ButtonNewWithLabel("Refresh")
	refreshBtn.OnClicked(func() { tv.syncNTP() })
	btnBox.Append(&refreshBtn.Widget)

	spacer := gtk4.LabelNew("")
	spacer.SetHExpand(true)
	btnBox.Append(&spacer.Widget)

	saveBtn := gtk4.ButtonNewWithLabel("Save")
	saveBtn.OnClicked(func() {
		datetime := fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d",
			int(tv.yearSpin.GetValue()),
			int(tv.monthSpin.GetValue()),
			int(tv.daySpin.GetValue()),
			int(tv.hourSpin.GetValue()),
			int(tv.minSpin.GetValue()),
			int(tv.secSpin.GetValue()),
		)
		zone := tv.tzCombo.GetActiveID()

		widget.PromptForSudo(tv.parentWin, "Setting system time requires root privileges.", func(password string) {
			if password == "" {
				return
			}
			go func() {
				ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
				defer cancel()

				c := exec.CommandContext(ctx, "sudo", "-S", "-k", "timedatectl", "set-time", datetime)
				c.Stdin = strings.NewReader(password + "\n")
				c.Run()

				if zone != "" && zone != tv.model.CurrentTimezone() {
					c := exec.CommandContext(ctx, "sudo", "-S", "-k", "timedatectl", "set-timezone", zone)
					c.Stdin = strings.NewReader(password + "\n")
					c.Run()
				}

				gtk4.IdleAdd(func() {
					tv.loadAsync()
				})
			}()
		})
	})
	btnBox.Append(&saveBtn.Widget)

	tv.box.Append(&btnBox.Widget)

	go tv.loadAsync()

	return tv
}

func (tv *TimeDateView) buildRows() {
	tv.addTimeRow()
	tv.addNTPRow()
	tv.addTimezoneRow()
	tv.addManualRow()
}

func (tv *TimeDateView) addTimeRow() {
	row := gtk4.ListBoxRowNew()
	row.SetSensitive(false)
	row.AddCSSClass("device-row")

	hbox := gtk4.BoxNew(gtk4.OrientationHorizontal, 8)
	hbox.SetMarginStart(8)
	hbox.SetMarginEnd(8)
	hbox.SetMarginTop(6)
	hbox.SetMarginBottom(6)

	lbl := gtk4.LabelNew("Current Time")
	lbl.SetSizeRequest(160, -1)
	lbl.SetHAlign(gtk4.AlignEnd)
	lbl.SetMarginEnd(8)
	hbox.Append(&lbl.Widget)

	tv.timeLabel = gtk4.LabelNew("Loading...")
	tv.timeLabel.SetHExpand(true)
	tv.timeLabel.SetHAlign(gtk4.AlignStart)
	hbox.Append(&tv.timeLabel.Widget)

	hboxWidget := hbox.Widget
	row.SetChild(&hboxWidget)
	tv.listBox.Append(row)
	tv.rows = append(tv.rows, row)
}

func (tv *TimeDateView) addNTPRow() {
	row := gtk4.ListBoxRowNew()
	row.AddCSSClass("device-row")

	hbox := gtk4.BoxNew(gtk4.OrientationHorizontal, 8)
	hbox.SetMarginStart(8)
	hbox.SetMarginEnd(8)
	hbox.SetMarginTop(6)
	hbox.SetMarginBottom(6)

	lbl := gtk4.LabelNew("NTP Sync")
	lbl.SetSizeRequest(160, -1)
	lbl.SetHAlign(gtk4.AlignEnd)
	lbl.SetMarginEnd(8)
	hbox.Append(&lbl.Widget)

	tv.ntpCheck = gtk4.CheckButtonNewWithLabel("Enable automatic NTP synchronization")
	tv.ntpCheck.OnToggled(func() {
		if tv.settingCombo {
			return
		}
		if tv.ntpCheck.GetActive() {
			tv.model.EnableNTP()
		} else {
			tv.model.DisableNTP()
		}
	})
	hbox.Append(&tv.ntpCheck.Widget)

	hboxWidget := hbox.Widget
	row.SetChild(&hboxWidget)
	tv.listBox.Append(row)
	tv.rows = append(tv.rows, row)
}

func (tv *TimeDateView) addTimezoneRow() {
	row := gtk4.ListBoxRowNew()
	row.AddCSSClass("device-row")

	hbox := gtk4.BoxNew(gtk4.OrientationHorizontal, 8)
	hbox.SetMarginStart(8)
	hbox.SetMarginEnd(8)
	hbox.SetMarginTop(6)
	hbox.SetMarginBottom(6)

	lbl := gtk4.LabelNew("Timezone")
	lbl.SetSizeRequest(160, -1)
	lbl.SetHAlign(gtk4.AlignEnd)
	lbl.SetMarginEnd(8)
	hbox.Append(&lbl.Widget)

	tv.tzCombo = gtk4.ComboBoxTextNew()
	tv.tzCombo.SetSizeRequest(300, -1)
	tv.tzCombo.OnChanged(func() {})
	hbox.Append(&tv.tzCombo.Widget)

	hboxWidget := hbox.Widget
	row.SetChild(&hboxWidget)
	tv.listBox.Append(row)
	tv.rows = append(tv.rows, row)
}

func (tv *TimeDateView) addManualRow() {
	row := gtk4.ListBoxRowNew()
	row.AddCSSClass("device-row")

	hbox := gtk4.BoxNew(gtk4.OrientationHorizontal, 8)
	hbox.SetMarginStart(8)
	hbox.SetMarginEnd(8)
	hbox.SetMarginTop(6)
	hbox.SetMarginBottom(6)

	lbl := gtk4.LabelNew("Date & Time")
	lbl.SetSizeRequest(160, -1)
	lbl.SetHAlign(gtk4.AlignEnd)
	lbl.SetMarginEnd(8)
	hbox.Append(&lbl.Widget)

	rightBox := gtk4.BoxNew(gtk4.OrientationVertical, 4)
	rightBox.SetHExpand(true)

	now := time.Now()

	tv.yearSpin = gtk4.SpinButtonNew(2000, 2100, 1)
	tv.yearSpin.SetValue(float64(now.Year()))
	tv.yearSpin.SetSizeRequest(72, -1)

	tv.monthSpin = gtk4.SpinButtonNew(1, 12, 1)
	tv.monthSpin.SetValue(float64(now.Month()))
	tv.monthSpin.SetSizeRequest(72, -1)
	tv.monthSpin.SetText(fmt.Sprintf("%02.0f", float64(now.Month())))
	tv.monthSpin.OnValueChanged(func(v float64) {
		tv.monthSpin.SetText(fmt.Sprintf("%02.0f", v))
	})

	tv.daySpin = gtk4.SpinButtonNew(1, 31, 1)
	tv.daySpin.SetValue(float64(now.Day()))
	tv.daySpin.SetSizeRequest(72, -1)
	tv.daySpin.SetText(fmt.Sprintf("%02.0f", float64(now.Day())))
	tv.daySpin.OnValueChanged(func(v float64) {
		tv.daySpin.SetText(fmt.Sprintf("%02.0f", v))
	})

	tv.hourSpin = gtk4.SpinButtonNew(0, 23, 1)
	tv.hourSpin.SetValue(float64(now.Hour()))
	tv.hourSpin.SetSizeRequest(72, -1)
	tv.hourSpin.SetText(fmt.Sprintf("%02.0f", float64(now.Hour())))
	tv.hourSpin.OnValueChanged(func(v float64) {
		tv.hourSpin.SetText(fmt.Sprintf("%02.0f", v))
	})

	tv.minSpin = gtk4.SpinButtonNew(0, 59, 1)
	tv.minSpin.SetValue(float64(now.Minute()))
	tv.minSpin.SetSizeRequest(72, -1)
	tv.minSpin.SetText(fmt.Sprintf("%02.0f", float64(now.Minute())))
	tv.minSpin.OnValueChanged(func(v float64) {
		tv.minSpin.SetText(fmt.Sprintf("%02.0f", v))
	})

	tv.secSpin = gtk4.SpinButtonNew(0, 59, 1)
	tv.secSpin.SetValue(float64(now.Second()))
	tv.secSpin.SetSizeRequest(72, -1)
	tv.secSpin.SetText(fmt.Sprintf("%02.0f", float64(now.Second())))
	tv.secSpin.OnValueChanged(func(v float64) {
		tv.secSpin.SetText(fmt.Sprintf("%02.0f", v))
	})

	grid := gtk4.GridNew()
	grid.SetColumnSpacing(4)
	grid.SetRowSpacing(4)

	grid.Attach(&tv.yearSpin.Widget, 0, 0, 1, 1)
	dl1 := gtk4.LabelNew("-")
	dl1.SetHAlign(gtk4.AlignCenter)
	grid.Attach(&dl1.Widget, 1, 0, 1, 1)
	grid.Attach(&tv.monthSpin.Widget, 2, 0, 1, 1)
	dl2 := gtk4.LabelNew("-")
	dl2.SetHAlign(gtk4.AlignCenter)
	grid.Attach(&dl2.Widget, 3, 0, 1, 1)
	grid.Attach(&tv.daySpin.Widget, 4, 0, 1, 1)

	grid.Attach(&tv.hourSpin.Widget, 0, 1, 1, 1)
	cl1 := gtk4.LabelNew(":")
	cl1.SetHAlign(gtk4.AlignCenter)
	grid.Attach(&cl1.Widget, 1, 1, 1, 1)
	grid.Attach(&tv.minSpin.Widget, 2, 1, 1, 1)
	cl2 := gtk4.LabelNew(":")
	cl2.SetHAlign(gtk4.AlignCenter)
	grid.Attach(&cl2.Widget, 3, 1, 1, 1)
	grid.Attach(&tv.secSpin.Widget, 4, 1, 1, 1)

	rightBox.Append(&grid.Widget)

	hbox.Append(&rightBox.Widget)

	hboxWidget := hbox.Widget
	row.SetChild(&hboxWidget)
	tv.listBox.Append(row)
	tv.rows = append(tv.rows, row)
}

func (tv *TimeDateView) loadAsync() {
	status, err := tv.model.GetStatus()
	zones, _ := tv.model.ListTimezones()

	gtk4.IdleAdd(func() {
		tv.settingCombo = true
		if status != nil {
			tv.timeLabel.SetText(time.Now().Format(time.RFC3339))
			tv.ntpCheck.SetActive(status.NTPEnabled)
		}
		tv.tzCombo.RemoveAll()
		for _, z := range zones {
			tv.tzCombo.Append(z, z)
		}
		if status != nil {
			for i, z := range zones {
				if z == status.Timezone {
					tv.tzCombo.SetActive(i)
					break
				}
			}
		}
		tv.settingCombo = false
		tv.zones = zones
	})
	_ = err
}

func (tv *TimeDateView) syncNTP() {
	tv.model.EnableNTP()
	go tv.loadAsync()
}

func (tv *TimeDateView) Widget() *gtk4.Widget { return &tv.box.Widget }
func (tv *TimeDateView) Name() string          { return "timedate" }
func (tv *TimeDateView) Title() string         { return "Time & Date" }
func (tv *TimeDateView) IconName() string      { return "preferences-system-time-symbolic" }
func (tv *TimeDateView) OnShow()               { go tv.loadAsync() }
func (tv *TimeDateView) OnHide()               {}
func (tv *TimeDateView) Destroy()              {}
func (tv *TimeDateView) SetParentWindow(parent *gtk4.Window) { tv.parentWin = parent }
