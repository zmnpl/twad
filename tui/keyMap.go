package tui

import (
	"fmt"

	"github.com/rivo/tview"
)

const (
	template = colorTagMoreContrast + "%-6v" + colorTagPrimaryText + "%v"
)

var (
	keyResetUI         = fmt.Sprintf(template, "ESC", "Reset UI")
	keyRunGame         = fmt.Sprintf(template, "ENTER", "Run Game")
	keyQuickload       = fmt.Sprintf(template, "F9", "Run Last Savegame")
	keyWarp            = fmt.Sprintf(template, "w", "Warp (+Record)")
	keyDemos           = fmt.Sprintf(template, "d", "Demos")
	keyQuit            = fmt.Sprintf(template, "q", "Quit")
	keyEditGame        = fmt.Sprintf(template, "e", "Edit Game")
	keyNewGame         = fmt.Sprintf(template, "n", "New Game")
	keyAddMod          = fmt.Sprintf(template, "m", "Add Mod To Game")
	keyRemoveGame      = fmt.Sprintf(template, "DEL", "Remove Game")
	keyImportArchive   = fmt.Sprintf(template, "i", "Import Archive")
	keySortAlph        = fmt.Sprintf(template, "s", "Sort Games Alphabetically")
	keyRateText        = fmt.Sprintf(template, "+/-", "Rate Game")
	keyCredits         = fmt.Sprintf(template, "c", "Credits/License")
	keyOptions         = fmt.Sprintf(template, "o", "Options")
	keySavegameDetails = fmt.Sprintf(template, "z", "Savegame Details")

	keyInfosMain []string
)

func init() {
	keyInfosMain = make([]string, 0, 10)
	keyInfosMain = append(keyInfosMain, keyResetUI)
	keyInfosMain = append(keyInfosMain, keyOptions)
	keyInfosMain = append(keyInfosMain, keyCredits)
	keyInfosMain = append(keyInfosMain, keyQuit)
	keyInfosMain = append(keyInfosMain, keyRunGame)
	keyInfosMain = append(keyInfosMain, keyQuickload)
	keyInfosMain = append(keyInfosMain, keyWarp)
	keyInfosMain = append(keyInfosMain, keyDemos)
	keyInfosMain = append(keyInfosMain, keyEditGame)
	keyInfosMain = append(keyInfosMain, keyNewGame)
	keyInfosMain = append(keyInfosMain, keyAddMod)
	keyInfosMain = append(keyInfosMain, keyRemoveGame)
	keyInfosMain = append(keyInfosMain, keySavegameDetails)
	keyInfosMain = append(keyInfosMain, keyImportArchive)
	keyInfosMain = append(keyInfosMain, keySortAlph)
	keyInfosMain = append(keyInfosMain, keyRateText)

}

func makeKeyMap() (*tview.Grid, int) {
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

	helpPane := tview.NewGrid()
	helpPane.SetBackgroundColor(tview.Styles.PrimitiveBackgroundColor).SetBorder(true)
	helpPane.SetRows(rowDimens...).SetColumns(colDimens...)

	keyItem := 0
ADDITEMS:
	for c := range colDimens {
		for r := range rowDimens {
			if keyItem >= len(keyInfosMain) {
				break ADDITEMS
			}
			helpPane.AddItem(tview.NewTextView().SetDynamicColors(true).SetText(keyInfosMain[keyItem]), r, c, 1, 1, 0, 0, false)
			keyItem++
		}
	}

	return helpPane, rows + 2
}

// func makeHelp() *tview.TextView {
// 	explanation := tview.NewTextView().SetRegions(true).SetWrap(true).SetWordWrap(true).SetDynamicColors(true)
// 	fmt.Fprintf(explanation, "%s\n\nPoint me to the highlighted directory:\n", setupPathExplain)
// 	fmt.Fprintf(explanation, "%s", setupPathExample)
// 	fmt.Fprintf(explanation, "\n\n%s", setupOkHint)

// 	return explanation
// }
