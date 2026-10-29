package tui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/zmnpl/twad/games"
	"github.com/zmnpl/twad/helper"
)

// keyBinding describes a key and what it does
// bindings with an action get handled by the games table, the others are informational
type keyBinding struct {
	label  string // key as shown to the user; empty hides the binding from help
	desc   string
	footer bool // show in the always visible keymap

	key tcell.Key // special key
	ch  rune      // or rune

	run  func()            // action
	game func(*games.Game) // action on the selected game; skipped when there is none
}

type keyGroup struct {
	title    string
	bindings []keyBinding
}

func (u *ui) buildKeyGroups() {
	u.keyGroups = []keyGroup{
		{"Basic Navigation", []keyBinding{
			{label: "Arrows (or hjkl)", desc: "Navigate"},
			{label: "TAB", desc: "Switch Focus", key: tcell.KeyTab, game: u.focusMods},
			{label: "ENTER", desc: "Confirm"},
			{label: "F1", desc: "Help/Keymap", footer: true, key: tcell.KeyF1, run: u.showHelp},
			{label: "F2", desc: "IDGames Browser", key: tcell.KeyF2, run: u.showIdgamesBrowser},
			{label: "F3", desc: "Import Archive", key: tcell.KeyF3, run: u.showZipSelect},
		}},
		{"General Functions", []keyBinding{
			{label: "ESC", desc: "Reset UI", footer: true},
			{label: "q", desc: "Quit", footer: true},
			{label: "c", desc: "Credits/License", ch: 'c', run: u.toggleLicense},
			{label: "o", desc: "Options", footer: true, ch: 'o', run: u.showOptions},
			{label: "s", desc: "Sort Games Alphabetically", ch: 's', run: u.sortGames},
		}},
		{"Launch Selected Game", []keyBinding{
			{label: "ENTER", desc: "Run Game", footer: true},
			{label: "F9", desc: "Run Last Savegame", footer: true, key: tcell.KeyF9, game: quickload},
			{label: "w", desc: "Warp (+Record)", footer: true, ch: 'w', game: u.warpGame},
		}},
		{"Selected Game Details", []keyBinding{
			{label: "d", desc: "Demos", footer: true, ch: 'd', game: u.showDemos},
			{label: "z", desc: "Savegames", footer: true, ch: 'z', game: u.showSavegames},
			{label: "u", desc: "Visit Url", ch: 'u', game: visitUrl},
		}},
		{"Add/Edit Game", []keyBinding{
			{label: "e", desc: "Edit Game", ch: 'e', game: u.editGame},
			{label: "m", desc: "Add Mod To Game", footer: true, ch: 'm', game: u.showModTree},
			{label: "n", desc: "New Game", footer: true, ch: 'n', run: u.newGame},
			{label: "DEL", desc: "Remove Game", footer: true, key: tcell.KeyDelete, game: u.removeGame},
			{label: "+/-", desc: "Rate Game", ch: '+', game: u.rateGame(1)},
			{ch: '-', game: u.rateGame(-1)},
		}},
	}
}

func (b keyBinding) String() string {
	return fmt.Sprintf(colorTagMoreContrast+"%-6v"+colorTagPrimaryText+"%v", b.label, b.desc)
}

func (b keyBinding) matches(event *tcell.EventKey) bool {
	if b.ch != 0 {
		return event.Key() == tcell.KeyRune && event.Rune() == b.ch
	}
	return b.key != 0 && event.Key() == b.key
}

// handleKeyBindings runs the action bound to the given key
func (u *ui) handleKeyBindings(event *tcell.EventKey) *tcell.EventKey {
	for _, group := range u.keyGroups {
		for _, b := range group.bindings {
			if !b.matches(event) {
				continue
			}
			if b.run != nil {
				b.run()
				return nil
			}
			if b.game != nil {
				if g, ok := u.selectedGame(); ok {
					b.game(g)
				}
				return nil
			}
		}
	}
	return event
}

// withQuit makes q quit the app, before handing other keys to the given input capture
func (u *ui) withQuit(capture func(event *tcell.EventKey) *tcell.EventKey) func(event *tcell.EventKey) *tcell.EventKey {
	return func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune && event.Rune() == 'q' {
			u.app.Stop()
			return nil
		}
		if capture == nil {
			return event
		}
		return capture(event)
	}
}

func (u *ui) makeKeyMap() (helpPane *tview.Grid, height int) {
	var footerBindings []keyBinding
	for _, group := range u.keyGroups {
		for _, b := range group.bindings {
			if b.footer {
				footerBindings = append(footerBindings, b)
			}
		}
	}

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
			if keyItem >= len(footerBindings) {
				break ADDITEMS
			}
			helpPane.AddItem(tview.NewTextView().SetDynamicColors(true).SetText(footerBindings[keyItem].String()), r, c, 1, 1, 0, 0, false)
			keyItem++
		}
	}
	height = rows + 2
	return
}

func (u *ui) makeHelp() *tview.TextView {
	explanation := tview.NewTextView().SetRegions(true).SetWrap(true).SetWordWrap(true).SetDynamicColors(true)
	explanation.SetBorder(true)
	explanation.SetTitle("Help")
	explanation.Clear()

	explanation.SetInputCapture(u.withQuit(func(event *tcell.EventKey) *tcell.EventKey {
		// switch back to nowmal mode
		if event.Key() == tcell.KeyF1 {
			u.appModeNormal()
			return nil
		}
		return event
	}))

	for _, group := range u.keyGroups {
		fmt.Fprintf(explanation, "%s%s\n", colorTagContrast, group.title)
		for _, b := range group.bindings {
			if b.label != "" {
				fmt.Fprintf(explanation, "%s\n", b)
			}
		}
		fmt.Fprint(explanation, "\n")
	}

	return explanation
}

// key actions
//#######################################################################

func (u *ui) showHelp() {
	u.appModeNormal()
	u.contentPages.switchTo(pageHelpKeymap, u.makeHelp())
}

func (u *ui) showIdgamesBrowser() {
	u.contentPages.open(pageIdgamesBrowser, u.idgamesBrowser.GetRootLayout())
}

func (u *ui) showZipSelect() {
	u.contentPages.open(pageZipSelect, u.zipSelector.selectTree)
}

// c again toggles back
func (u *ui) toggleLicense() {
	if front, _ := u.detailSidePagesSub2.GetFrontPage(); front == pageLicense {
		u.appModeNormal()
		return
	}
	u.detailSidePagesSub2.switchTo(pageLicense, u.makeLicense())
}

func (u *ui) showOptions() {
	u.contentPages.switchTo(pageOptions, u.makeOptions())
}

func (u *ui) sortGames() {
	games.SortAlph()
	u.populateGamesTable()
}

func (u *ui) newGame() {
	u.detailPages.switchTo(pageAddEdit, u.makeAddEditGame(nil))
}

func (u *ui) editGame(g *games.Game) {
	u.detailPages.switchTo(pageAddEdit, u.makeAddEditGame(g))
}

// "quickload" tries to load the latest savegame
func quickload(g *games.Game) {
	g.Quickload()
}

func (u *ui) warpGame(g *games.Game) {
	u.contentPages.showAt(pageWarp, u.makeWarpRecord(*g), dialogXOffset, u.selectedRowY(), warpDialogWidth, warpDialogHeight)
}

func (u *ui) showDemos(g *games.Game) {
	demoList, err := u.makeDemoList(g)
	if err != nil {
		return
	}
	u.detailSidePagesSub1.show(pageDemos, demoList)
}

func (u *ui) showSavegames(g *games.Game) {
	savegameList, err := u.makeSavegameList(g)
	if err != nil {
		return
	}
	u.detailSidePagesSub1.show(pageSaves, savegameList)
}

func visitUrl(g *games.Game) {
	helper.Openbrowser(g.Link)
}

func (u *ui) showModTree(g *games.Game) {
	modTree, err := u.makeModTree(g)
	if err != nil {
		u.showError(dict.modTreeErrWadDir, err.Error(), nil, nil)
		return
	}
	u.detailSidePagesSub2.switchTo(pageModSelector, modTree)
}

// switch focus to mod list; only if game has mods
func (u *ui) focusMods(g *games.Game) {
	if len(g.Mods) > 0 {
		u.app.SetFocus(u.detailSidePagesSub1)
	}
}

func (u *ui) rateGame(increment int) func(*games.Game) {
	return func(g *games.Game) {
		g.Rate(increment)
		r, _ := u.gamesTable.GetSelection()
		u.setGameRow(r, g)
		games.Persist()
	}
}

func (u *ui) removeGame(g *games.Game) {
	r, _ := u.gamesTable.GetSelection()
	remove := func() {
		if r == u.gamesTable.GetRowCount()-1 {
			u.gamesTable.Select(r-1, 0)
		}
		games.RemoveGameAt(r - gamesTableHeaderRows)
	}

	if u.config.DeleteWithoutWarning {
		remove()
		return
	}
	confirm(u.contentPages, g.Name, dialogXOffset, u.selectedRowY(), u.gamesTable, remove)
}
