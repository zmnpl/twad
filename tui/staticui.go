package tui

import (
	"fmt"

	"github.com/rivo/tview"
)

const (
	subtitle  = "[orange]twad[white] - [orange]t[white]erminal [orange]wad[white] launcher[orange]"
	subtitle2 = "twad - terminal wad manager and launcher"

	tviewHeader = "[orange]tview"
	creditTview = `The terminal user interface is build with tview:
https://github.com/rivo/tview`

	doomLogoCreditHeader = "[orange]DOOM Logo"
	creditDoomLogo       = `DOOM and Quake are registered trademarks of id Software, Inc. The DOOM, Quake and id logos are trademarks of id Software, Inc.

The ASCII version of the DOOM logo is Copyright © 1994 by F.P. de Vries.

This logo is work from Frans P. de Vries who originally made it and nicely granted me permission to use it here

Details can be found in this little piece of video game history:
http://www.gamers.org/~fpv/doomlogo.html`

	licenseHeader = "[orange]License"
	mitLicense    = `MIT License

Copyright (c) 2020 Simon Paul

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.`
)

//  header
func makeHeader() *tview.TextView {
	header := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true)
	fmt.Fprintf(header, "%s", doomLogo)

	return header
}

//  license
func makeLicense() *tview.TextView {
	disclaimer := tview.NewTextView().SetDynamicColors(true).SetRegions(true)
	fmt.Fprintf(disclaimer, "%s\n", doomLogoCreditHeader)
	fmt.Fprintf(disclaimer, "%s\n\n", creditDoomLogo)
	fmt.Fprintf(disclaimer, "%s\n", tviewHeader)
	fmt.Fprintf(disclaimer, "%s\n\n", creditTview)
	fmt.Fprintf(disclaimer, "%s\n", licenseHeader)
	fmt.Fprintf(disclaimer, "%s", mitLicense)
	disclaimer.SetBorder(true).SetTitle("Credits / License")

	return disclaimer
}

// help for navigation
func makeHelpPane() (*tview.Flex, int) {
	home := tview.NewTextView().SetDynamicColors(true).SetText(" [orange](ESC)[white]   - Reset UI")
	run := tview.NewTextView().SetDynamicColors(true).SetText(" [orange](Enter)[white] - Run Game")
	insert := tview.NewTextView().SetDynamicColors(true).SetText(" [orange](i)[white]     - Add New Game")
	edit := tview.NewTextView().SetDynamicColors(true).SetText(" [orange](e)[white]     - Edit Game")
	add := tview.NewTextView().SetDynamicColors(true).SetText(" [orange](a)[white]     - Add Mod To Game")
	delet := tview.NewTextView().SetDynamicColors(true).SetText(" [orange](Del)[white]   - Remove Game")
	license := tview.NewTextView().SetDynamicColors(true).SetText(" [orange](c)[white]     - Credits/License")
	quit := tview.NewTextView().SetDynamicColors(true).SetText(" [orange](q)[white]     - Quit")
	options := tview.NewTextView().SetDynamicColors(true).SetText(" [orange](o)[white]     - Options")
	sort := tview.NewTextView().SetDynamicColors(true).SetText(" [orange](s)[white]     - Sort Games Alphabetically")

	helpArea := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(run, 1, 0, false).
			AddItem(home, 1, 0, false).
			AddItem(quit, 1, 0, false),
			0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(insert, 1, 0, false).
			AddItem(edit, 1, 0, false).
			AddItem(add, 1, 0, false),
			0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(sort, 1, 0, false).
			AddItem(options, 1, 0, false).
			AddItem(delet, 1, 0, false),
			0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(license, 1, 0, false),
			0, 1, false)
	helpArea.SetBorder(true)

	helpPane := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false).
		AddItem(helpArea, 5, 0, false)

	return helpPane, 5
}
