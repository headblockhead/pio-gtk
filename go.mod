module github.com/headblockhead/pio-gtk

go 1.26.7

replace github.com/headblockhead/pio => ../pio

require (
	github.com/diamondburned/gotk4-adwaita/pkg v0.0.0-20260808200908-d4aecaa0ff32
	github.com/diamondburned/gotk4/pkg v0.4.1
	github.com/headblockhead/pio v0.0.0-20261006201519-d6338329a4ad
)

require golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
