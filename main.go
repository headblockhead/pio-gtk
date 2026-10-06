package main

import (
	"fmt"
	"os"

	"github.com/diamondburned/gotk4-adwaita/pkg/adw"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/headblockhead/pio"
	"github.com/headblockhead/pio/conn"
	"github.com/headblockhead/pio/gpio"
	"github.com/headblockhead/pio/sm"
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
	label.SetVExpand(true)
	label.SetHExpand(true)

	s := conn.NewSimulation()
	r := pio.NewRP2040("rp2040-0")
	s.AddNet("EBD")
	s.AddNet("EBC")
	p6, _ := r.Connection(6)
	p7, _ := r.Connection(7)
	s.Connect(p6, "EBD")
	s.Connect(p7, "EBC")
	s.AddTicker(r)

	pio0c, _ := r.PIOConfigurator(0)
	mw0 := pio0c.MemoryWriter()

	mw0.Write(0, 0xc020) //  0: irq    wait 0
	//     .wrap_target
	mw0.Write(1, 0x80a0)  //  1: pull   block
	mw0.Write(2, 0x2047)  //  2: wait   0 irq, 7
	mw0.Write(3, 0xf780)  //  3: set    pindirs, 0      side 0 [7]
	mw0.Write(4, 0xe027)  //  4: set    x, 7
	mw0.Write(5, 0x80e0)  //  5: pull   ifempty block
	mw0.Write(6, 0x6041)  //  6: out    y, 1
	mw0.Write(7, 0x006b)  //  7: jmp    !y, 11
	mw0.Write(8, 0xfa81)  //  8: set    pindirs, 1      side 1 [2]
	mw0.Write(9, 0x00cc)  //  9: jmp    pin, 12
	mw0.Write(10, 0x0000) // 10: jmp    0
	mw0.Write(11, 0xfb80) // 11: set    pindirs, 0      side 1 [3]
	mw0.Write(12, 0x10e6) // 12: jmp    !osre, 6        side 0
	mw0.Write(13, 0x0045) // 13: jmp    x--, 5
	mw0.Write(14, 0xff81) // 14: set    pindirs, 1      side 1 [7]
	mw0.Write(15, 0xc021) // 15: irq    wait 1
	//     .wrap

	sm0c, _ := pio0c.SMConfigurator(0)
	sm0c.SetBaseSetPin(6)
	sm0c.SetPinCountSet(1)
	sm0c.SetJumpPin(6)
	sm0c.SetBaseSidesetPin(7)
	sm0c.SetBitCountSideset(2)
	sm0c.SetFIFOJoin(sm.FIFOJoinTX)
	sm0c.SetOutShiftMovesRight(false)
	sm0c.SetWrapFromAddress(15)
	sm0c.SetWrapToAddress(1)

	gp6, _ := r.GPIOConfigurator(6)
	gp7, _ := r.GPIOConfigurator(7)
	gp6.SetFunction(gpio.FunctionPIO0)
	gp7.SetFunction(gpio.FunctionPIO0)

	pad6, _ := r.PadConfigurator(6)
	pad7, _ := r.PadConfigurator(7)
	pad6.SetPullUp(true)
	pad7.SetPullUp(true)

	gp6.SetOutputEnableOverride(gpio.OverrideInvert)
	gp7.SetOutputEnableOverride(gpio.OverrideInvert)

	sm0c.SetEnabled(true)

	fw := sm0c.FIFOTXWriter()
	fw.Write(0xAAAAAAAA)
	fw.Write(0xFF00FF00)

	pio0o, _ := r.PIOObserver(0)

	button := gtk.NewButtonWithLabel("Do stuff")
	button.SetHAlign(gtk.AlignCenter)
	button.SetVAlign(gtk.AlignCenter)
	button.ConnectClicked(func() {
		p6, _ := r.PadObserver(6)
		p7, _ := r.PadObserver(7)
		label.SetLabel(fmt.Sprintf("pad6: %v, pad7: %v", p6.Output(), p7.Output()))

		fmt.Println("button clicked")
		fmt.Println(fmt.Sprintf("irq: %d", pio0o.IRQs()))

		err := s.Tick()
		if err != nil {
			panic(err)
		}

		err = pio0c.IRQClear(1)
		if err != nil {
			panic(err)
		}
	})

	box := gtk.NewBox(gtk.OrientationVertical, 0)
	box.Append(header)
	box.Append(label)
	box.Append(button)

	win := adw.NewApplicationWindow(&app.Application)
	win.SetTitle("Hello World")
	win.SetDefaultSize(400, 300)
	win.SetContent(box)
	win.Present()
}
