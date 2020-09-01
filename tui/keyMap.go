package tui

import (
	"fmt"

	"github.com/rivo/tview"
)

const (
	keyResetUIText    = "Reset UI"
	keyRunGameText    = "Run Game"
	keyQuickload      = "Run Last Savegame"
	keyWarp           = "Warp (+Record)"
	keyDemos          = "Demos"
	keyQuitText       = "Quit"
	keyEditGameText   = "Edit Game"
	keyNewGameText    = "New Game"
	keyAddModText     = "Add Mod To Game"
	keyRemoveGameText = "Remove Game"
	keyImportArchive  = "Import Archive"
	keySortAlphText   = "Sort Games Alphabetically"
	keyRateText       = "Rate Game"
	keyCreditsText    = "Credits/License"
	keyOptionsText    = "Options"
)

var (
	keyInfosMain []string
)

func init() {
	template := colorTagMoreContrast + "%-6v" + colorTagPrimaryText + "%v"

	keyInfosMain = make([]string, 0, 10)
	keyInfosMain = append(keyInfosMain, fmt.Sprintf(template, "ESC", keyResetUIText))
	keyInfosMain = append(keyInfosMain, fmt.Sprintf(template, "ENTER", keyRunGameText))
	keyInfosMain = append(keyInfosMain, fmt.Sprintf(template, "F9", keyQuickload))
	keyInfosMain = append(keyInfosMain, fmt.Sprintf(template, "w", keyWarp))
	keyInfosMain = append(keyInfosMain, fmt.Sprintf(template, "d", keyDemos))
	keyInfosMain = append(keyInfosMain, fmt.Sprintf(template, "e", keyEditGameText))
	keyInfosMain = append(keyInfosMain, fmt.Sprintf(template, "n", keyNewGameText))
	keyInfosMain = append(keyInfosMain, fmt.Sprintf(template, "m", keyAddModText))
	keyInfosMain = append(keyInfosMain, fmt.Sprintf(template, "DEL", keyRemoveGameText))
	keyInfosMain = append(keyInfosMain, fmt.Sprintf(template, "i", keyImportArchive))
	keyInfosMain = append(keyInfosMain, fmt.Sprintf(template, "s", keySortAlphText))
	keyInfosMain = append(keyInfosMain, fmt.Sprintf(template, "+/-", keyRateText))
	keyInfosMain = append(keyInfosMain, fmt.Sprintf(template, "c", keyCreditsText))
	keyInfosMain = append(keyInfosMain, fmt.Sprintf(template, "o", keyOptionsText))
	keyInfosMain = append(keyInfosMain, fmt.Sprintf(template, "q", keyQuitText))
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
