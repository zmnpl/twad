package tui

import (
	"fmt"

	"github.com/rivo/tview"
)

const (
	template = colorTagMoreContrast + "%-6v" + colorTagPrimaryText + "%v"
)

var (
	// general ui
	keyNavigate     = fmt.Sprintf(template, "Arrows (or hjkl)", "Navigate")
	keyFormNav      = fmt.Sprintf(template, "TAB", "Traverse form items")
	keyConfirm      = fmt.Sprintf(template, "ENTER", "Confirm")
	keyInfoNavigate = []string{keyNavigate, keyFormNav, keyConfirm}

	// ui behaviour
	keyResetUI  = fmt.Sprintf(template, "ESC", "Reset UI")
	keyQuit     = fmt.Sprintf(template, "q", "Quit")
	keyCredits  = fmt.Sprintf(template, "c", "Credits/License")
	keyOptions  = fmt.Sprintf(template, "o", "Options")
	keyInfoMain = []string{keyResetUI, keyQuit, keyCredits, keyOptions}

	// game launching
	keyRunGame        = fmt.Sprintf(template, "ENTER", "Run Game")
	keyQuickload      = fmt.Sprintf(template, "F9", "Run Last Savegame")
	keyWarp           = fmt.Sprintf(template, "w", "Warp (+Record)")
	keyInfoGameLaunch = []string{keyRunGame, keyQuickload, keyWarp}

	// open detail views
	keyDemos           = fmt.Sprintf(template, "d", "Demos")
	keySavegameDetails = fmt.Sprintf(template, "z", "Savegame Details")
	keyVisitUrl        = fmt.Sprintf(template, "u", "Visit Url")
	keyInfoGameDetails = []string{keyDemos, keySavegameDetails, keyVisitUrl}

	// game crud ops
	keyEditGame      = fmt.Sprintf(template, "e", "Edit Game")
	keyAddMod        = fmt.Sprintf(template, "m", "Add Mod To Game")
	keyNewGame       = fmt.Sprintf(template, "n", "New Game")
	keyRemoveGame    = fmt.Sprintf(template, "DEL", "Remove Game")
	keyInfoGameTable = []string{keyEditGame, keyAddMod, keyNewGame, keyRemoveGame}

	// others
	keyImportArchive = fmt.Sprintf(template, "i", "Import Archive")
	keySortAlph      = fmt.Sprintf(template, "s", "Sort Games Alphabetically")
	keyRate          = fmt.Sprintf(template, "+/-", "Rate Game")
	keyInfoOthers    = []string{keyImportArchive, keySortAlph, keyRate}
)

func makeKeyMap() (helpPane *tview.Grid, height int) {
	keyInfoAll := keyInfoMain
	keyInfoAll = append(keyInfoAll, keyInfoGameLaunch...)
	keyInfoAll = append(keyInfoAll, keyInfoGameDetails...)
	keyInfoAll = append(keyInfoAll, keyInfoGameTable...)
	keyInfoAll = append(keyInfoAll, keyInfoOthers...)

	// could be easier / more static, but like this the layout can maybe be made more dynamic in the future
	rows := 4
	rowDimens := make([]int, rows)
	for i := range rowDimens {
		rowDimens[i] = 1
	}

	cols := 4
	colDimens := make([]int, cols)
	for i := range colDimens {
		colDimens[i] = 0
	}

	helpPane = tview.NewGrid()
	helpPane.SetBackgroundColor(tview.Styles.PrimitiveBackgroundColor).SetBorder(true)
	helpPane.SetRows(rowDimens...).SetColumns(colDimens...)

	keyItem := 0
ADDITEMS:
	for c := range colDimens {
		for r := range rowDimens {
			if keyItem >= len(keyInfoAll) {
				break ADDITEMS
			}
			helpPane.AddItem(tview.NewTextView().SetDynamicColors(true).SetText(keyInfoAll[keyItem]), r, c, 1, 1, 0, 0, false)
			keyItem++
		}
	}
	height = rows + 2
	return
}

func makeHelp() *tview.TextView {
	explanation := tview.NewTextView().SetRegions(true).SetWrap(true).SetWordWrap(true).SetDynamicColors(true)
	explanation.SetBorder(true)
	explanation.Clear()

	header := "%s%s\n"

	fmt.Fprintf(explanation, header, colorTagContrast, "Basic Navigation")
	for _, keyText := range keyInfoNavigate {
		fmt.Fprintf(explanation, "%s\n", keyText)
	}
	fmt.Fprint(explanation, "\n")

	fmt.Fprintf(explanation, header, colorTagContrast, "UI")
	for _, keyText := range keyInfoMain {
		fmt.Fprintf(explanation, "%s\n", keyText)
	}
	fmt.Fprint(explanation, "\n")

	fmt.Fprintf(explanation, header, colorTagContrast, "Launch Selected Game")
	for _, keyText := range keyInfoGameLaunch {
		fmt.Fprintf(explanation, "%s\n", keyText)
	}
	fmt.Fprint(explanation, "\n")

	fmt.Fprintf(explanation, header, colorTagContrast, "Selected Game Details")
	for _, keyText := range keyInfoGameDetails {
		fmt.Fprintf(explanation, "%s\n", keyText)
	}
	fmt.Fprint(explanation, "\n")

	fmt.Fprintf(explanation, header, colorTagContrast, "Add/Edit Game")
	for _, keyText := range keyInfoGameTable {
		fmt.Fprintf(explanation, "%s\n", keyText)
	}
	fmt.Fprint(explanation, "\n")

	fmt.Fprintf(explanation, header, colorTagContrast, "Others")
	for _, keyText := range keyInfoOthers {
		fmt.Fprintf(explanation, "%s\n", keyText)
	}
	fmt.Fprint(explanation, "\n")

	return explanation
}

func showHelp() {
	help := makeHelp()
	appModeNormal()
	detailSidePagesSub2.AddPage(pageHelpKeymap, help, true, true)
	detailSidePagesSub2.SwitchToPage(pageHelpKeymap)
	app.SetFocus(help)
}
