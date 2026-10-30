package tui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/zmnpl/twad/tui/shared"
)

const (
	tviewHeader = "[orange]tview"
	creditTview = `The terminal user interface is build with tview:
https://github.com/rivo/tview`
)

// header
func (u *ui) makeHeader() *tview.TextView {
	header := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true)
	header.SetBorderPadding(0, 0, 1, 1)

	// if config.UseTerminalColors {
	// 	fmt.Fprintf(header, "%s", doomLogoTermColors)
	// } else {
	// 	fmt.Fprintf(header, "%s", doomLogo)
	// }

	fmt.Fprintf(header, "%s", u.makeLogo())

	return header
}

// first start hello
func (u *ui) makeFirstStartHello() *tview.TextView {
	hello := tview.NewTextView().SetDynamicColors(true)
	hello.SetBorder(true)
	fmt.Fprint(hello, shared.HelloText)

	hello.SetInputCapture(u.withQuit(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune && event.Rune() == 'o' {
			u.showOptions()
		}
		return event
	}))

	return hello
}

// license
func (u *ui) makeLicense() *tview.TextView {
	disclaimer := tview.NewTextView().SetDynamicColors(true).SetRegions(true)
	fmt.Fprintf(disclaimer, "%s\n", shared.DoomLogoCreditHeader)
	fmt.Fprintf(disclaimer, "%s\n\n", shared.CreditDoomLogo)
	fmt.Fprintf(disclaimer, "%s\n", tviewHeader)
	fmt.Fprintf(disclaimer, "%s\n\n", creditTview)
	fmt.Fprintf(disclaimer, "%s\n", shared.LicenseHeader)
	fmt.Fprintf(disclaimer, "%s", shared.MITLicense)
	disclaimer.SetBorder(true).SetTitle("Credits / License")

	disclaimer.SetInputCapture(u.withQuit(func(event *tcell.EventKey) *tcell.EventKey {
		// switch back to nowmal mode
		if event.Key() == tcell.KeyRune && event.Rune() == 'c' {
			u.appModeNormal()
			return nil
		}
		return event
	}))

	return disclaimer
}
