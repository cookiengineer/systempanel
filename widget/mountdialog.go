package widget

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/cookiengineer/systempanel/bindings/gtk4"
)

	type MountDialog struct {
	win      *gtk4.Window
	parent   *gtk4.Window

	mountEntry  *gtk4.Entry
	mapperEntry *gtk4.Entry
	passEntry   *gtk4.Entry
	passVisBtn  *gtk4.Button
	mapperRow   *gtk4.Box
	passRow     *gtk4.Box
	errorLbl    *gtk4.Label
	mountBtn    *gtk4.Button
	spinner     *gtk4.Spinner

	devicePath  string
	mapperPath  string
	mapperName  string
	isEncrypted bool
	isUnlocked  bool
	onDone      func()
	running     bool
}

func ShowMountDialog(parent *gtk4.Window, devicePath string, partName string, isEncrypted bool, isUnlocked bool, mapperName string, onDone func()) {
	md := &MountDialog{
		parent:      parent,
		devicePath:  devicePath,
		isEncrypted: isEncrypted,
		isUnlocked:  isUnlocked,
		onDone:      onDone,
	}

	if isEncrypted {
		if isUnlocked && mapperName != "" {
			md.mapperName = mapperName
		} else {
			md.mapperName = "systempanel_" + partName
		}
		md.mapperPath = "/dev/mapper/" + md.mapperName
	}

	md.build()
	md.win.Present()
}

func (md *MountDialog) build() {
	md.win = gtk4.WindowNew()
	md.win.SetTitle("Mount Partition")
	md.win.SetDefaultSize(480, -1)
	md.win.SetModal(true)
	md.win.SetTransientFor(md.parent)

	vbox := gtk4.BoxNew(gtk4.OrientationVertical, 12)
	vbox.SetMarginStart(24)
	vbox.SetMarginEnd(24)
	vbox.SetMarginTop(24)
	vbox.SetMarginBottom(24)

	title := gtk4.LabelNew("Mount Partition")
	title.SetMarkup("<b>Mount Partition</b>")
	title.SetHAlign(gtk4.AlignStart)
	tw := title.Widget
	vbox.Append(&tw)

	mountRow := gtk4.BoxNew(gtk4.OrientationHorizontal, 8)
	mountRow.SetHExpand(true)

	mountLabel := gtk4.LabelNew("Mount Path:")
	mountLabel.SetHAlign(gtk4.AlignStart)
	mountLabel.SetSizeRequest(100, -1)
	mountRow.Append(&mountLabel.Widget)

	md.mountEntry = gtk4.EntryNew()
	md.mountEntry.SetPlaceholder("/mnt/mydisk")
	md.mountEntry.SetHExpand(true)
	mountRow.Append(&md.mountEntry.Widget)
	vbox.Append(&mountRow.Widget)

	md.mapperRow = gtk4.BoxNew(gtk4.OrientationHorizontal, 8)
	md.mapperRow.SetHExpand(true)

	mapperLabel := gtk4.LabelNew("Mapper Path:")
	mapperLabel.SetHAlign(gtk4.AlignStart)
	mapperLabel.SetSizeRequest(100, -1)
	md.mapperRow.Append(&mapperLabel.Widget)

	md.mapperEntry = gtk4.EntryNew()
	md.mapperEntry.SetHExpand(true)
	if md.mapperPath != "" {
		md.mapperEntry.SetText(md.mapperPath)
	}
	md.mapperRow.Append(&md.mapperEntry.Widget)
	vbox.Append(&md.mapperRow.Widget)

	if !md.isEncrypted {
		md.mapperRow.Hide()
	}

	md.passRow = gtk4.BoxNew(gtk4.OrientationHorizontal, 8)
	md.passRow.SetHExpand(true)

	passLabel := gtk4.LabelNew("Password:")
	passLabel.SetHAlign(gtk4.AlignStart)
	passLabel.SetSizeRequest(100, -1)
	md.passRow.Append(&passLabel.Widget)

	md.passEntry = gtk4.EntryNew()
	md.passEntry.SetVisibility(false)
	md.passEntry.SetPlaceholder("LUKS passphrase")
	md.passEntry.SetHExpand(true)
	md.passRow.Append(&md.passEntry.Widget)

	md.passVisBtn = gtk4.ButtonNew()
	md.passVisBtn.SetIconName("view-conceal-symbolic")
	md.passVisBtn.OnClicked(md.togglePassVisibility)
	md.passRow.Append(&md.passVisBtn.Widget)
	vbox.Append(&md.passRow.Widget)

	if !md.isEncrypted || md.isUnlocked {
		md.passRow.Hide()
	}

	md.errorLbl = gtk4.LabelNew("")
	md.errorLbl.AddCSSClass("journal-err")
	md.errorLbl.SetHAlign(gtk4.AlignStart)
	md.errorLbl.SetMarginStart(108)
	elWidget := md.errorLbl.Widget
	vbox.Append(&elWidget)

	btnBox := gtk4.BoxNew(gtk4.OrientationHorizontal, 8)
	btnBox.SetMarginTop(8)

	cancelBtn := gtk4.ButtonNewWithLabel("Cancel")
	cancelBtn.OnClicked(func() { md.close() })
	cbWidget := cancelBtn.Widget
	btnBox.Append(&cbWidget)

	md.spinner = gtk4.SpinnerNew()
	md.spinner.SetSizeRequest(16, 16)
	btnBox.Append(&md.spinner.Widget)

	spacer := gtk4.LabelNew("")
	spacer.SetHExpand(true)
	btnBox.Append(&spacer.Widget)

	md.mountBtn = gtk4.ButtonNewWithLabel("Mount")
	md.mountBtn.AddCSSClass("suggested-action")
	md.mountBtn.OnClicked(func() { md.onMount() })
	btnBox.Append(&md.mountBtn.Widget)

	btnBoxWidget := btnBox.Widget
	vbox.Append(&btnBoxWidget)

	vboxWidget := vbox.Widget
	md.win.SetChild(&vboxWidget)

	ctrl := gtk4.EventControllerKeyNew()
	ctrl.OnKeyPressed(func(keyval, _, _ uint) bool {
		if keyval == 0xff0d && !md.running {
			md.onMount()
			return true
		}
		return false
	})
	ctrl.AddToWidget(&md.win.Widget)
}

func (md *MountDialog) togglePassVisibility() {
	visible := md.passEntry.GetVisibility()
	if visible {
		md.passEntry.SetVisibility(false)
		md.passVisBtn.SetIconName("view-conceal-symbolic")
	} else {
		md.passEntry.SetVisibility(true)
		md.passVisBtn.SetIconName("view-reveal-symbolic")
	}
}

func (md *MountDialog) close() {
	if md.running {
		return
	}
	md.win.Hide()
	md.win.Destroy()
}

func (md *MountDialog) setRunning(running bool) {
	md.running = running
	md.mountBtn.SetSensitive(!running)
	if running {
		md.spinner.Start()
	} else {
		md.spinner.Stop()
	}
}

func (md *MountDialog) onMount() {
	if md.running {
		return
	}

	mountPath := strings.TrimSpace(md.mountEntry.GetText())
	if mountPath == "" {
		md.errorLbl.SetText("Mount Path is required.")
		return
	}

	sourceDevice := md.devicePath
	mapperName := ""
	if md.isEncrypted {
		mapperPathStr := strings.TrimSpace(md.mapperEntry.GetText())
		if mapperPathStr == "" {
			md.errorLbl.SetText("Mapper Path is required.")
			return
		}
		md.mapperPath = mapperPathStr
		mapperName = filepath.Base(mapperPathStr)
		if mapperName == "" || mapperName == "." || mapperName == "/" {
			md.errorLbl.SetText("Invalid Mapper Path.")
			return
		}
		sourceDevice = mapperPathStr

		if !md.isUnlocked {
			if strings.TrimSpace(md.passEntry.GetText()) == "" {
				md.errorLbl.SetText("Password is required to unlock the encrypted partition.")
				return
			}
		}
	}

	md.errorLbl.SetText("")

	luksPassword := strings.TrimSpace(md.passEntry.GetText())
	devicePath := md.devicePath
	encrypted := md.isEncrypted
	unlocked := md.isUnlocked
	mapper := mapperName
	source := sourceDevice
	target := mountPath

	PromptForSudo(md.parent,
		"Mounting requires root privileges.",
		func(sudoPassword string) {
			if sudoPassword == "" {
				return
			}

			gtk4.IdleAdd(func() { md.setRunning(true) })

			go func() {
				err := runMount(devicePath, encrypted, unlocked, mapper, luksPassword, source, target, sudoPassword)

				gtk4.IdleAdd(func() {
					md.setRunning(false)
					if err != nil {
						md.errorLbl.SetText(err.Error())
					} else {
						md.win.Hide()
						md.win.Destroy()
						if md.onDone != nil {
							md.onDone()
						}
					}
				})
			}()
		},
	)
}

func runMount(devicePath string, encrypted bool, unlocked bool, mapperName string, luksPassword string, sourceDevice string, mountPath string, sudoPassword string) error {
	if encrypted && !unlocked {
		if err := runCryptsetup(devicePath, mapperName, luksPassword, sudoPassword); err != nil {
			return fmt.Errorf("unlock failed: %v", err)
		}
	}

	if err := runMkdir(mountPath, sudoPassword); err != nil {
		return fmt.Errorf("mkdir failed: %v", err)
	}

	if err := runMountCmd(sourceDevice, mountPath, sudoPassword); err != nil {
		return fmt.Errorf("mount failed: %v", err)
	}

	return nil
}

func runCryptsetup(devicePath string, mapperName string, luksPassword string, sudoPassword string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "sudo", "-S", "-k", "cryptsetup", "luksOpen", devicePath, mapperName)
	input := sudoPassword + "\n" + luksPassword + "\n"
	cmd.Stdin = strings.NewReader(input)

	var stderr strings.Builder
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		errStr := strings.TrimSpace(stderr.String())
		if errStr == "" {
			errStr = err.Error()
		}
		return fmt.Errorf("%s", errStr)
	}
	return nil
}

func runMkdir(mountPath string, sudoPassword string) error {
	if _, err := os.Stat(mountPath); err == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "sudo", "-S", "-k", "mkdir", "-p", mountPath)
	cmd.Stdin = strings.NewReader(sudoPassword + "\n")

	var stderr strings.Builder
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		errStr := strings.TrimSpace(stderr.String())
		if errStr == "" {
			errStr = err.Error()
		}
		return fmt.Errorf("%s", errStr)
	}
	return nil
}

func runMountCmd(sourceDevice string, mountPath string, sudoPassword string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "sudo", "-S", "-k", "mount", sourceDevice, mountPath)
	cmd.Stdin = strings.NewReader(sudoPassword + "\n")

	var stderr strings.Builder
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		errStr := strings.TrimSpace(stderr.String())
		if errStr == "" {
			errStr = err.Error()
		}
		return fmt.Errorf("%s", errStr)
	}
	return nil
}
