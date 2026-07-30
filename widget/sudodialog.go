package widget

import (
	"context"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/cookiengineer/systempanel/bindings/gtk4"
)

var (
	sessionPassword string
	sessionMu       sync.Mutex
)

type SudoDialog struct {
	win      *gtk4.Window
	entry    *gtk4.Entry
	errorLbl *gtk4.Label
	parent   *gtk4.Window
	message  string
	callback func(string)
}

func NewSudoDialog(parent *gtk4.Window) *SudoDialog {
	return NewSudoDialogWithMessage(parent, "Authentication is required to modify system settings.")
}

func NewSudoDialogWithMessage(parent *gtk4.Window, message string) *SudoDialog {
	sd := &SudoDialog{
		parent:  parent,
		message: message,
	}
	sd.build()
	return sd
}

func (sd *SudoDialog) build() {
	sd.win = gtk4.WindowNew()
	sd.win.SetTitle("Authentication Required")
	sd.win.SetDefaultSize(460, 240)
	sd.win.SetModal(true)
	sd.win.SetTransientFor(sd.parent)

	vbox := gtk4.BoxNew(gtk4.OrientationVertical, 12)
	vbox.SetMarginStart(24)
	vbox.SetMarginEnd(24)
	vbox.SetMarginTop(24)
	vbox.SetMarginBottom(24)

	headerBox := gtk4.BoxNew(gtk4.OrientationHorizontal, 12)

	iconWrapper := gtk4.BoxNew(gtk4.OrientationHorizontal, 0)
	iconWrapper.SetSizeRequest(90, -1)
	iconWrapper.SetHAlign(gtk4.AlignEnd)
	icon := gtk4.ImageNewFromIconName("dialog-password-symbolic")
	icon.SetPixelSize(48)
	iconWrapper.Append(&icon.Widget)
	headerBox.Append(&iconWrapper.Widget)

	msgBox := gtk4.BoxNew(gtk4.OrientationVertical, 4)
	title := gtk4.LabelNew("Authentication Required")
	title.SetMarkup("<b>Authentication Required</b>")
	title.SetHAlign(gtk4.AlignStart)
	tw := title.Widget
	msgBox.Append(&tw)

	msg := gtk4.LabelNew(sd.message)
	msg.SetWrap(true)
	msg.SetHAlign(gtk4.AlignStart)
	mw := msg.Widget
	msgBox.Append(&mw)

	msgBoxWidget := msgBox.Widget
	headerBox.Append(&msgBoxWidget)
	headerBoxWidget := headerBox.Widget
	vbox.Append(&headerBoxWidget)

	passBox := gtk4.BoxNew(gtk4.OrientationHorizontal, 8)
	passLabel := gtk4.LabelNew("Password:")
	passLabel.SetSizeRequest(90, -1)
	passLabel.SetHAlign(gtk4.AlignEnd)
	passLabel.SetMarginEnd(4)
	plWidget := passLabel.Widget
	passBox.Append(&plWidget)

	sd.entry = gtk4.EntryNew()
	sd.entry.SetVisibility(false)
	sd.entry.SetPlaceholder("Enter your password")
	sd.entry.SetHExpand(true)
	entryWidget := sd.entry.Widget
	passBox.Append(&entryWidget)
	passBoxWidget := passBox.Widget
	vbox.Append(&passBoxWidget)

	sd.errorLbl = gtk4.LabelNew("")
	sd.errorLbl.AddCSSClass("journal-err")
	sd.errorLbl.SetHAlign(gtk4.AlignStart)
	sd.errorLbl.SetMarginStart(100)
	elWidget := sd.errorLbl.Widget
	vbox.Append(&elWidget)

	btnBox := gtk4.BoxNew(gtk4.OrientationHorizontal, 8)
	btnBox.SetMarginTop(8)

	cancelBtn := gtk4.ButtonNewWithLabel("Cancel")
	cancelBtn.OnClicked(func() { sd.cancel() })
	cbWidget := cancelBtn.Widget
	btnBox.Append(&cbWidget)

	spacer := gtk4.LabelNew("")
	spacer.SetHExpand(true)
	btnBox.Append(&spacer.Widget)

	okBtn := gtk4.ButtonNewWithLabel("Authenticate")
	okBtn.AddCSSClass("suggested-action")
	okBtn.OnClicked(func() { sd.authenticate() })
	okWidget := okBtn.Widget
	btnBox.Append(&okWidget)

	btnBoxWidget := btnBox.Widget
	vbox.Append(&btnBoxWidget)

	vboxWidget := vbox.Widget
	sd.win.SetChild(&vboxWidget)

	ctrl := gtk4.EventControllerKeyNew()
	ctrl.OnKeyPressed(func(keyval, _, _ uint) bool {
		if keyval == 0xff0d {
			sd.authenticate()
			return true
		}
		return false
	})
	ctrl.AddToWidget(&sd.win.Widget)
}

func (sd *SudoDialog) cancel() {
	if sd.callback != nil {
		sd.callback("")
	}
	sd.win.Close()
}

func (sd *SudoDialog) authenticate() {
	pw := sd.entry.GetText()
	if pw == "" {
		return
	}
	if sd.checkPassword(pw) {
		sessionMu.Lock()
		sessionPassword = pw
		sessionMu.Unlock()
		if sd.callback != nil {
			sd.callback(pw)
		}
		sd.win.Close()
	} else {
		sd.entry.SetText("")
		sd.errorLbl.SetText("Incorrect password. Please try again.")
	}
}

func (sd *SudoDialog) checkPassword(password string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "sudo", "-S", "-k", "true")
	cmd.Stdin = strings.NewReader(password + "\n")
	return cmd.Run() == nil
}

func (sd *SudoDialog) Show(callback func(password string)) {
	sd.callback = callback
	sd.errorLbl.SetText("")
	sd.entry.SetText("")

	sessionMu.Lock()
	if sessionPassword != "" {
		sessionMu.Unlock()
		callback(sessionPassword)
		return
	}
	sessionMu.Unlock()

	sd.win.Present()
}

func PromptForSudo(parent *gtk4.Window, message string, callback func(password string)) {
	dialog := NewSudoDialogWithMessage(parent, message)
	dialog.Show(callback)
}

func (sd *SudoDialog) RunCommand(cmd string, args ...string) {
	sd.Show(func(password string) {
		if password == "" {
			return
		}

		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			c := exec.CommandContext(ctx, "sudo", append([]string{"-S", "-k", cmd}, args...)...)
			c.Stdin = strings.NewReader(password + "\n")
			c.Run()
		}()
	})
}

func InvalidateSudo() {
	sessionMu.Lock()
	sessionPassword = ""
	sessionMu.Unlock()
}

func RunSudoCommand(cmd string, args ...string) error {
	sessionMu.Lock()
	pw := sessionPassword
	sessionMu.Unlock()

	if pw == "" {
		return exec.Command("pkexec", append([]string{cmd}, args...)...).Run()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	c := exec.CommandContext(ctx, "sudo", append([]string{"-S", "-k", cmd}, args...)...)
	c.Stdin = strings.NewReader(pw + "\n")
	return c.Run()
}
