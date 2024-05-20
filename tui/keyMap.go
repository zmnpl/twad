package tui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	lineTemplate string

	// general ui
	keyNavigate      string
	keyFormNav       string
	keyConfirm       string
	keyHelp          string
	keyIdgames       string
	keyImportArchive string
	keyInfoNavigate  = []string{keyNavigate, keyFormNav, keyConfirm, keyIdgames, keyImportArchive}

	// general
	keyResetUI  string
	keyQuit     string
	keyCredits  string
	keyOptions  string
	keySortAlph string
	keyInfoMain = []string{keyResetUI, keyQuit, keyCredits, keyOptions, keySortAlph}

	// game launching
	keyRunGame        string
	keyQuickload      string
	keyWarp           string
	keyInfoGameLaunch = []string{keyRunGame, keyQuickload, keyWarp}

	// open detail views
	keyDemos           string
	keySavegameDetails string
	keyVisitUrl        string
	keyInfoGameDetails = []string{keyDemos, keySavegameDetails, keyVisitUrl}

	// game crud ops
	keyEditGame      string
	keyAddMod        string
	keyNewGame       string
	keyRemoveGame    string
	keyRate          string
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

func initKeyMapEntries() {
	lineTemplate = colorTagMoreContrast + "%-6v" + colorTagPrimaryText + "%v"

	// general ui
	keyNavigate = fmt.Sprintf(lineTemplate, "Arrows (or hjkl)", "Navigate")
	keyFormNav = fmt.Sprintf(lineTemplate, "TAB", "Switch Focus")
	keyConfirm = fmt.Sprintf(lineTemplate, "ENTER", "Confirm")
	keyHelp = fmt.Sprintf(lineTemplate, "F1", "Help/Keymap")
	keyIdgames = fmt.Sprintf(lineTemplate, "F2", "IDGames Browser")
	keyImportArchive = fmt.Sprintf(lineTemplate, "F3", "Import Archive")
	keyInfoNavigate = []string{keyNavigate, keyFormNav, keyConfirm, keyIdgames, keyImportArchive}

	// general
	keyResetUI = fmt.Sprintf(lineTemplate, "ESC", "Reset UI")
	keyQuit = fmt.Sprintf(lineTemplate, "q", "Quit")
	keyCredits = fmt.Sprintf(lineTemplate, "c", "Credits/License")
	keyOptions = fmt.Sprintf(lineTemplate, "o", "Options")
	keySortAlph = fmt.Sprintf(lineTemplate, "s", "Sort Games Alphabetically")
	keyInfoMain = []string{keyResetUI, keyQuit, keyCredits, keyOptions, keySortAlph}

	// game launching
	keyRunGame = fmt.Sprintf(lineTemplate, "ENTER", "Run Game")
	keyQuickload = fmt.Sprintf(lineTemplate, "F9", "Run Last Savegame")
	keyWarp = fmt.Sprintf(lineTemplate, "w", "Warp (+Record)")
	keyInfoGameLaunch = []string{keyRunGame, keyQuickload, keyWarp}

	// open detail views
	keyDemos = fmt.Sprintf(lineTemplate, "d", "Demos")
	keySavegameDetails = fmt.Sprintf(lineTemplate, "z", "Savegames")
	keyVisitUrl = fmt.Sprintf(lineTemplate, "u", "Visit Url")
	keyInfoGameDetails = []string{keyDemos, keySavegameDetails, keyVisitUrl}

	// game crud ops
	keyEditGame = fmt.Sprintf(lineTemplate, "e", "Edit Game")
	keyAddMod = fmt.Sprintf(lineTemplate, "m", "Add Mod To Game")
	keyNewGame = fmt.Sprintf(lineTemplate, "n", "New Game")
	keyRemoveGame = fmt.Sprintf(lineTemplate, "DEL", "Remove Game")
	keyRate = fmt.Sprintf(lineTemplate, "+/-", "Rate Game")
	keyInfoGameTable = []string{keyEditGame, keyAddMod, keyNewGame, keyRemoveGame, keyRate}

	// pick "most important ones" for always visible footer keymap
	keyInfoQuickMapSelectino = []string{
		keyResetUI, keyQuit, keyOptions,
		keyRunGame, keyQuickload, keyWarp,
		keyDemos, keySavegameDetails,
		keyNewGame, keyAddMod, keyRemoveGame,
		keyHelp,
	}
}

func makeKeyMap() (helpPane *tview.Grid, height int) {
	initKeyMapEntries()

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

	contentPages.AddPage(pageHelpKeymap, help, true, true)
	contentPages.SwitchToPage(pageHelpKeymap)

	app.SetFocus(help)
}
