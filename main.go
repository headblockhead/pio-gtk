package main

import (
	"os"

	"github.com/diamondburned/gotk4-adwaita/pkg/adw"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

func main() {
	app := adw.NewApplication("dev.edwardh.pio-gtk", gio.ApplicationDefaultFlags)
	app.ConnectActivate(func() { activate(app) })

	if code := app.Run(os.Args); code > 0 {
		os.Exit(code)
	}
}

func activate(app *adw.Application) {
	header := adw.NewHeaderBar()

	label := gtk.NewLabel("Hello, World!")
	label.AddCSSClass("title-1")
	label.SetVExpand(true)
	label.SetHExpand(true)

	box := gtk.NewBox(gtk.OrientationVertical, 0)
	box.Append(header)
	box.Append(label)

	win := adw.NewApplicationWindow(&app.Application)
	win.SetTitle("Hello World")
	win.SetDefaultSize(400, 300)
	win.SetContent(box)
	win.Present()
}
