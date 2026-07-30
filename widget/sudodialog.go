package widget

import (
	"context"
	"os/exec"
	"strings"
	"time"

	"github.com/cookiengineer/systempanel/bindings/gtk4"
)

type SudoDialog struct {
	win      *gtk4.Window
	entry    *gtk4.Entry
	result   chan string
	parent   *gtk4.Window
}

func NewSudoDialog(parent *gtk4.Window) *SudoDialog {
	sd := &SudoDialog{
		result: make(chan string, 1),
		parent: parent,
	}
	sd.build()
	return sd
}

func (sd *SudoDialog) build() {
	sd.win = gtk4.WindowNew()
	sd.win.SetTitle("Authentication Required")
	sd.win.SetDefaultSize(400, 150)

	vbox := gtk4.BoxNew(gtk4.OrientationVertical, 12)
	vbox.SetMarginStart(24)
	vbox.SetMarginEnd(24)
	vbox.SetMarginTop(24)
	vbox.SetMarginBottom(24)

	msg := gtk4.LabelNew("Authentication is required to save connection settings.")
	msg.SetWrap(true)
	msgWidget := msg.Widget
	vbox.Append(&msgWidget)

	sd.entry = gtk4.EntryNew()
	sd.entry.SetVisibility(false)
	sd.entry.SetPlaceholder("Password")
	entryWidget := sd.entry.Widget
	vbox.Append(&entryWidget)

	btnBox := gtk4.BoxNew(gtk4.OrientationHorizontal, 8)
	btnBox.SetHAlign(gtk4.AlignEnd)

	cancelBtn := gtk4.ButtonNewWithLabel("Cancel")
	cancelBtn.OnClicked(func() { sd.result <- ""; sd.win.Close() })
	cbWidget := cancelBtn.Widget
	btnBox.Append(&cbWidget)

	okBtn := gtk4.ButtonNewWithLabel("Authenticate")
	okBtn.OnClicked(func() { sd.result <- sd.entry.GetText(); sd.win.Close() })
	okBtn.AddCSSClass("suggested-action")
	okWidget := okBtn.Widget
	btnBox.Append(&okWidget)

	btnBoxWidget := btnBox.Widget
	vbox.Append(&btnBoxWidget)

	vboxWidget := vbox.Widget
	sd.win.SetChild(&vboxWidget)

	ctrl := gtk4.EventControllerKeyNew()
	ctrl.OnKeyPressed(func(keyval, _, _ uint) bool {
		if keyval == 0xff0d {
			sd.result <- sd.entry.GetText()
			sd.win.Close()
			return true
		}
		return false
	})
	ctrl.AddToWidget(&sd.win.Widget)
}

func (sd *SudoDialog) Prompt() string {
	sd.win.Present()
	select {
	case pw := <-sd.result:
		return pw
	case <-time.After(5 * time.Minute):
		sd.win.Close()
		return ""
	}
}

func (sd *SudoDialog) RunCommand(cmd string, args ...string) error {
	password := sd.Prompt()
	if password == "" {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	c := exec.CommandContext(ctx, "sudo", append([]string{"-S", cmd}, args...)...)
	c.Stdin = strings.NewReader(password + "\n")
	return c.Run()
}
