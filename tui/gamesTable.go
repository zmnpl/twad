package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/zmnpl/twad/games"
)

const gamesTableHeaderRows = 1

// center table with mods
func (u *ui) makeGamesTable() *tview.Table {
	u.gamesTable = tview.NewTable().
		SetFixed(1, 1).
		SetSelectable(true, false).
		SetBorders(tableBorders).SetSeparator('|')
	u.gamesTable.SetBorderPadding(0, 0, 1, 2)

	u.gamesTable.SetSelectionChangedFunc(func(r int, c int) {
		if g, ok := gameAt(r); ok {
			u.selectedGameChanged(g)
		}
	})

	u.gamesTable.SetSelectedFunc(func(r int, c int) {
		if g, ok := gameAt(r); ok {
			g.Run()
		}
	})

	u.gamesTable.SetInputCapture(u.withQuit(u.handleKeyBindings))

	return u.gamesTable
}

// gamesTableColumn describes a column of the games table and how to show a game in it
type gamesTableColumn struct {
	header string
	value  func(g *games.Game) string
	color  tcell.Color
}

// evaluated when used, because header texts and colors get set up at runtime
func gamesTableColumns() []gamesTableColumn {
	return []gamesTableColumn{
		{dict.gameTableHeaderRating, (*games.Game).RatingString, tview.Styles.TitleColor},
		{dict.gameTableHeaderName, func(g *games.Game) string { return g.Name }, tview.Styles.SecondaryTextColor},
		{dict.gameTableHeaderSourcePort, func(g *games.Game) string { return g.Port }, tview.Styles.PrimaryTextColor},
		{dict.gameTableHeaderIwad, func(g *games.Game) string { return g.Iwad }, tview.Styles.PrimaryTextColor},
	}
}

// gameAt returns the game shown in the given row of the games table
func gameAt(r int) (*games.Game, bool) {
	allGames := games.Games()
	if r < gamesTableHeaderRows || r-gamesTableHeaderRows >= len(allGames) {
		return nil, false
	}
	return &allGames[r-gamesTableHeaderRows], true
}

// selectedGame returns the game selected in the games table
func (u *ui) selectedGame() (*games.Game, bool) {
	r, _ := u.gamesTable.GetSelection()
	return gameAt(r)
}

// selectedRowY returns the screen row right below the selected game within contentPages
func (u *ui) selectedRowY() int {
	r, _ := u.gamesTable.GetSelection()
	rowOffset, _ := u.gamesTable.GetOffset() // account for, when table is scrolled beyond visible screen
	_, _, _, commandPreviewHeight := u.commandPreview.GetRect()
	return commandPreviewHeight + (r + 1) - rowOffset
}

func (u *ui) populateGamesTable() {
	u.gamesTable.Clear()

	for c, col := range gamesTableColumns() {
		u.gamesTable.SetCell(0, c, tview.NewTableCell(col.header).SetTextColor(tview.Styles.SecondaryTextColor))
	}

	allGames := games.Games()
	for i := range allGames {
		u.setGameRow(i+gamesTableHeaderRows, &allGames[i])
	}
}

// setGameRow (re)draws the cells of a game in the given row
func (u *ui) setGameRow(r int, g *games.Game) {
	for c, col := range gamesTableColumns() {
		u.gamesTable.SetCell(r, c, tview.NewTableCell(col.value(g)).SetTextColor(col.color))
	}
}
