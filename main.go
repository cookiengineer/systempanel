package main

import (
	"log"
	"os"

	"github.com/cookiengineer/systempanel/app"
	"github.com/cookiengineer/systempanel/bindings/gtk4"
	"github.com/cookiengineer/systempanel/view"

	_ "github.com/cookiengineer/systempanel/css"

	"github.com/cookiengineer/systempanel/view/settings"
	"github.com/cookiengineer/systempanel/view/volume"
	"github.com/cookiengineer/systempanel/view/wifi"
	"github.com/cookiengineer/systempanel/view/bluetooth"
	"github.com/cookiengineer/systempanel/view/batteries"
	"github.com/cookiengineer/systempanel/view/brightness"
	"github.com/cookiengineer/systempanel/view/monitors"
	"github.com/cookiengineer/systempanel/view/wallpapers"
	"github.com/cookiengineer/systempanel/view/services"
	"github.com/cookiengineer/systempanel/view/autostart"
	"github.com/cookiengineer/systempanel/view/journal"
	"github.com/cookiengineer/systempanel/view/lan"
	"github.com/cookiengineer/systempanel/view/themes"
	"github.com/cookiengineer/systempanel/view/icons"
	"github.com/cookiengineer/systempanel/view/powerprofile"
	"github.com/cookiengineer/systempanel/view/timedate"
	"github.com/cookiengineer/systempanel/view/disks"
)

func init() {
	registerViews()
}

func registerViews() {
	view.Registry = []view.ViewDescriptor{
		lan.Descriptor,
		wifi.Descriptor,
		bluetooth.Descriptor,
		volume.Descriptor,
		monitors.Descriptor,
		wallpapers.Descriptor,
		brightness.Descriptor,
		batteries.Descriptor,
		disks.Descriptor,
		services.Descriptor,
		autostart.Descriptor,
		journal.Descriptor,
		powerprofile.Descriptor,
		timedate.Descriptor,
		themes.Descriptor,
		icons.Descriptor,
		settings.Descriptor,
	}
}

func main() {
	gtkApp := gtk4.ApplicationNew("com.github.cookiengineer.systempanel")

	gtkApp.OnActivate(func() {
		panel := app.New(gtkApp)
		panel.Build()
	})

	code := gtkApp.Run(os.Args)
	if code != 0 {
		log.Printf("Application exited with code %d", code)
	}
}
