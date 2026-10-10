package main

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/diamondburned/gotk4-adwaita/pkg/adw"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

//go:embed window.ui
var uiXML string

func main() {
	app := adw.NewApplication("com.example.blueprint", gio.ApplicationFlagsNone)

	app.ConnectActivate(func() {
		// Initialize the builder using the embedded UI XML
		builder := gtk.NewBuilderFromString(uiXML)

		obj := builder.GetObject("main_window")
		if obj == nil {
			fmt.Println("Error: could not find main_window object")
			return
		}

		// FIX 1: Safely cast the object to your UI window type
		win := obj.Cast().(*gtk.Window)
		app.AddWindow(win)

		btnObj := builder.GetObject("nav_button")
		if btnObj != nil {

		}

		win.Present()
	})

	if code := app.Run(os.Args); code > 0 {
		os.Exit(code)
	}
}
