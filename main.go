package main

import (
	_ "embed"
	"os"

	"github.com/diamondburned/gotk4-adwaita/pkg/adw"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

//go:embed window.ui
var uiXML string

func main() {
	println("running")

	app := adw.NewApplication("dev.edwardh.rubus", gio.ApplicationFlags(gio.ApplicationDefaultFlags))

	app.ConnectActivate(func() {
		builder := gtk.NewBuilderFromString(uiXML)

		windowObject := builder.GetObject("window")
		if windowObject == nil {
			panic("no window")
		}

		win := windowObject.Cast().(*adw.ApplicationWindow)
		win.SetApplication(&app.Application)
		win.Present()
	})

	if code := app.Run(os.Args); code > 0 {
		os.Exit(code)
	}
}
