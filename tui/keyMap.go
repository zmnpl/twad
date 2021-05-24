package tui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	template = colorTagMoreContrast + "%-6v" + colorTagPrimaryText + "%v"
)

var (
	// general ui
	keyNavigate     = fmt.Sprintf(template, "Arrows (or hjkl)", "Navigate")
	keyFormNav      = fmt.Sprintf(template, "TAB", "Switch Focus")
	keyConfirm      = fmt.Sprintf(template, "ENTER", "Confirm")
	keyHelp         = fmt.Sprintf(template, "F1", "Help/Keymap")
	keyIdgames      = fmt.Sprintf(template, "F2", "IDGames Browser")
	keyInfoNavigate = []string{keyNavigate, keyFormNav, keyConfirm, keyIdgames}

	// general
	keyResetUI       = fmt.Sprintf(template, "ESC", "Reset UI")
	keyQuit          = fmt.Sprintf(template, "q", "Quit")
	keyCredits       = fmt.Sprintf(template, "c", "Credits/License")
	keyOptions       = fmt.Sprintf(template, "o", "Options")
	keyImportArchive = fmt.Sprintf(template, "i", "Import Archive")
	keySortAlph      = fmt.Sprintf(template, "s", "Sort Games Alphabetically")
	keyInfoMain      = []string{keyResetUI, keyQuit, keyCredits, keyOptions, keyImportArchive, keySortAlph}

	// game launching
	keyRunGame        = fmt.Sprintf(template, "ENTER", "Run Game")
	keyQuickload      = fmt.Sprintf(template, "F9", "Run Last Savegame")
	keyWarp           = fmt.Sprintf(template, "w", "Warp (+Record)")
	keyInfoGameLaunch = []string{keyRunGame, keyQuickload, keyWarp}

	// open detail views
	keyDemos           = fmt.Sprintf(template, "d", "Demos")
	keySavegameDetails = fmt.Sprintf(template, "z", "Savegames")
	keyVisitUrl        = fmt.Sprintf(template, "u", "Visit Url")
	keyInfoGameDetails = []string{keyDemos, keySavegameDetails, keyVisitUrl}

	// game crud ops
	keyEditGame      = fmt.Sprintf(template, "e", "Edit Game")
	keyAddMod        = fmt.Sprintf(template, "m", "Add Mod To Game")
	keyNewGame       = fmt.Sprintf(template, "n", "New Game")
	keyRemoveGame    = fmt.Sprintf(template, "DEL", "Remove Game")
	keyRate          = fmt.Sprintf(template, "+/-", "Rate Game")
	keyInfoGameTable = []string{keyEditGame, keyAddMod, keyNewGame, keyRemoveGame, keyRate}

	// pick "most important ones" for always visible footer keymap
	keyInfoQuickMapSelectino = []string{
		keyResetUI, keyQuit, keyOptions,
		keyRunGame, keyQuickload, keyWarp,
		keyDemos, keySavegameDetails,
		keyNewGame, keyAddMod, keyRemoveGame,
		keyHelp,
	}
)

func makeKeyMap() (helpPane *tview.Grid, height int) {
	// could be easier / more static, but like this the layout can maybe be made more dynamic in the future
	rows := 3
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
			if keyItem >= len(keyInfoQuickMapSelectino) {
				break ADDITEMS
			}
			helpPane.AddItem(tview.NewTextView().SetDynamicColors(true).SetText(keyInfoQuickMapSelectino[keyItem]), r, c, 1, 1, 0, 0, false)
			keyItem++
		}
	}
	height = rows + 2
	return
}

func makeHelp() *tview.TextView {
	explanation := tview.NewTextView().SetRegions(true).SetWrap(true).SetWordWrap(true).SetDynamicColors(true)
	explanation.SetBorder(true)
	explanation.SetTitle("Help")
	explanation.Clear()

	explanation.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		k := event.Key()
		// switch back to nowmal mode
		if k == tcell.KeyF1 {
			appModeNormal()
			return nil
		}

		if k == tcell.KeyRune {
			switch event.Rune() {

			// get out
			case 'q':
				app.Stop()
				return nil
			}
		}

		return event
	})

	header := "%s%s\n"

	fmt.Fprintf(explanation, header, colorTagContrast, "Basic Navigation")
	for _, keyText := range keyInfoNavigate {
		fmt.Fprintf(explanation, "%s\n", keyText)
	}
	fmt.Fprint(explanation, "\n")

	fmt.Fprintf(explanation, header, colorTagContrast, "General Functions")
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

	return explanation
}

func showHelp() {
	help := makeHelp()
	appModeNormal()
	detailSidePagesSub2.AddPage(pageHelpKeymap, help, true, true)
	detailSidePagesSub2.SwitchToPage(pageHelpKeymap)
	app.SetFocus(help)
}
