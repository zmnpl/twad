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

func printHelpPane() (*tview.Grid, int) {
	helpPane := tview.NewGrid()
	helpPane.SetBackgroundColor(tview.Styles.PrimitiveBackgroundColor)
	helpPane.SetBorder(true)

	keyInfos := make([]string, 10)
	keyInfos = append(keyInfos, "ESC   - Reset UI")
	keyInfos = append(keyInfos, "Enter - Run Game")
	keyInfos = append(keyInfos, "q - Quit")
	keyInfos = append(keyInfos, "e - Edit Game")
	keyInfos = append(keyInfos, "n - New Game")
	keyInfos = append(keyInfos, "m - Add Mod To Game")
	keyInfos = append(keyInfos, "Del - Remove Game")
	keyInfos = append(keyInfos, "s - Sort Games Alphabetically")
	keyInfos = append(keyInfos, "+/- - Rate Game")
	keyInfos = append(keyInfos, "c - Credits/License")
	keyInfos = append(keyInfos, "o - Options")

	return helpPane, 5
}

// help for navigation
func makeHelpPane() (*tview.Flex, int) {
	template := colorTagMoreContrast + "(%v)" + colorTagPrimaryText + "%v- %v"

	home := tview.NewTextView().SetDynamicColors(true).SetText(fmt.Sprintf(template, "ESC", "   ", "Reset UI"))
	run := tview.NewTextView().SetDynamicColors(true).SetText(fmt.Sprintf(template, "Enter", " ", "Run Game"))
	quit := tview.NewTextView().SetDynamicColors(true).SetText(fmt.Sprintf(template, "q", "     ", "Quit"))

	edit := tview.NewTextView().SetDynamicColors(true).SetText(fmt.Sprintf(template, "e", "     ", "Edit Game"))
	insert := tview.NewTextView().SetDynamicColors(true).SetText(fmt.Sprintf(template, "n", "     ", "New Game"))
	add := tview.NewTextView().SetDynamicColors(true).SetText(fmt.Sprintf(template, "m", "     ", "Add Mod To Game"))

	delet := tview.NewTextView().SetDynamicColors(true).SetText(fmt.Sprintf(template, "Del", "   ", "Remove Game"))
	sort := tview.NewTextView().SetDynamicColors(true).SetText(fmt.Sprintf(template, "s", "     ", "Sort Games Alphabetically"))
	rate := tview.NewTextView().SetDynamicColors(true).SetText(fmt.Sprintf(template, "+/-", "   ", "Rate Game"))

	license := tview.NewTextView().SetDynamicColors(true).SetText(fmt.Sprintf(template, "c", "     ", "Credits/License"))
	options := tview.NewTextView().SetDynamicColors(true).SetText(fmt.Sprintf(template, "o", "     ", "Options"))

	helpArea := tview.NewFlex().SetDirection(tview.FlexColumn).AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(run, 1, 0, false).
		AddItem(home, 1, 0, false).
		AddItem(quit, 1, 0, false),
		0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(insert, 1, 0, false).
			AddItem(add, 1, 0, false).
			AddItem(edit, 1, 0, false),
			0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(rate, 1, 0, false).
			AddItem(sort, 1, 0, false).
			AddItem(delet, 1, 0, false),
			0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(options, 1, 0, false).
			AddItem(license, 1, 0, false),
			0, 1, false)
	helpArea.SetBorder(true)
	helpArea.SetBackgroundColor(tview.Styles.PrimitiveBackgroundColor)

	return helpArea, 5
}
