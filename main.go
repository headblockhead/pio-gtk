package main

import (
	"os"

	"github.com/diamondburned/gotk4-adwaita/pkg/adw"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

func main() {
	app := adw.NewApplication("com.example.HelloWorld", gio.ApplicationFlagsNone)
	app.ConnectActivate(func() { activate(app) })

	if code := app.Run(os.Args); code > 0 {
		os.Exit(code)
	}
}

func activate(app *adw.Application) {
	// Header bar with the default window controls (title, close button, etc.)
	header := adw.NewHeaderBar()

	// The actual "Hello, World!" content
	label := gtk.NewLabel("Hello, World!")
	label.AddCSSClass("title-1") // Adwaita style class for large text
	label.SetVExpand(true)
	label.SetHExpand(true)

	// Stack the header bar above the content
	box := gtk.NewBox(gtk.OrientationVertical, 0)
	box.Append(header)
	box.Append(label)

	// adw.ApplicationWindow expects a *gtk.Application, which adw.Application embeds.
	win := adw.NewApplicationWindow(&app.Application)
	win.SetTitle("Hello World")
	win.SetDefaultSize(400, 300)
	win.SetContent(box)
	win.Present()
}
