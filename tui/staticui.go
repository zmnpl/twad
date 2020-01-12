package tui

import (
	"fmt"

	"github.com/rivo/tview"
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

// button bar showing keys
func makeButtonBar() *tview.Flex {
	btnHome := tview.NewButton("(ESC) Reset UI")
	btnRun := tview.NewButton("(Enter) Run Game")
	btnInsert := tview.NewButton("(i) Add Game")
	btnAddMod := tview.NewButton("(a) Add Mods To Game")
	btnRemoveMod := tview.NewButton("(r) Remove Last Mod From Game")
	btnDelete := tview.NewButton("(Delete) Remove Game")
	btnLicenseAndCredits := tview.NewButton("(c) Credits/License")
	btnQuit := tview.NewButton("(q) Quit")
	buttonBar := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(btnHome, 0, 1, false).
		AddItem(btnRun, 0, 1, false).
		AddItem(btnInsert, 0, 1, false).
		AddItem(btnAddMod, 0, 1, false).
		AddItem(btnRemoveMod, 0, 1, false).
		AddItem(btnDelete, 0, 1, false).
		AddItem(btnLicenseAndCredits, 0, 1, false).
		AddItem(btnQuit, 0, 1, false)

	return buttonBar
}

// help for navigation
func makeHelpPane() *tview.Flex {
	home := tview.NewTextView().SetDynamicColors(true).SetText(" [orange](ESC)[white]   - Reset UI")
	run := tview.NewTextView().SetDynamicColors(true).SetText(" [orange](Enter)[white] - Run Game")
	insert := tview.NewTextView().SetDynamicColors(true).SetText(" [orange](i)[white]     - Add Game")
	add := tview.NewTextView().SetDynamicColors(true).SetText(" [orange](a)[white]     - Add Mod To Game")
	remove := tview.NewTextView().SetDynamicColors(true).SetText(" [orange](r)[white]     - Remove Last Mod From Game")
	delet := tview.NewTextView().SetDynamicColors(true).SetText(" [orange](Del)[white]   - Remove Game")
	license := tview.NewTextView().SetDynamicColors(true).SetText(" [orange](c)[white]     - Credits/License")
	quit := tview.NewTextView().SetDynamicColors(true).SetText(" [orange](q)[white]     - Quit")

	spacer := tview.NewTextView().SetDynamicColors(true).SetText("")

	helpArea := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(spacer, 1, 0, false).
			AddItem(home, 1, 0, false).
			AddItem(run, 1, 0, false).
			AddItem(insert, 1, 0, false).
			AddItem(add, 1, 0, false).
			AddItem(spacer, 1, 0, false),
			0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(spacer, 1, 0, false).
			AddItem(remove, 1, 0, false).
			AddItem(delet, 1, 0, false).
			AddItem(license, 1, 0, false).
			AddItem(quit, 1, 0, false).
			AddItem(spacer, 1, 0, false),
			0, 1, false)
	helpArea.SetBorder(true)
	helpArea.SetTitle("Help")

	helpPage := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false).
		AddItem(helpArea, 8, 0, true)

	return helpPage
}
