package tui

import (
	"github.com/rivo/tview"
	"github.com/zmnpl/twad/games"
)

const gamesTableHeaderRows = 1

// center table with mods
func makeGamesTable() *tview.Table {
	gamesTable = tview.NewTable().
		SetFixed(1, 1).
		SetSelectable(true, false).
		SetBorders(tableBorders).SetSeparator('|')
	gamesTable.SetBorderPadding(0, 0, 1, 2)
	return gamesTable
}

// selectedGame returns the game selected in the games table
func selectedGame() (*games.Game, bool) {
	r, _ := gamesTable.GetSelection()
	allGames := games.Games()
	if r < gamesTableHeaderRows || r-gamesTableHeaderRows >= len(allGames) {
		return nil, false
	}
	return &allGames[r-gamesTableHeaderRows], true
}

// selectedRowY returns the screen row right below the selected game within contentPages
func selectedRowY() int {
	r, _ := gamesTable.GetSelection()
	rowOffset, _ := gamesTable.GetOffset() // account for, when table is scrolled beyond visible screen
	_, _, _, commandPreviewHeight := commandPreview.GetRect()
	return commandPreviewHeight + (r + 1) - rowOffset
}

func populateGamesTable() {
	gamesTable.Clear()
	allGames := games.Games()

	fixRows, fixCols := gamesTableHeaderRows, 4
	rows, cols := len(allGames), 0

	for r := 0; r < rows+fixRows; r++ {
		var game games.Game
		if r > 0 {
			game = allGames[r-fixRows]
		}
		for c := 0; c < cols+fixCols; c++ {
			var cell *tview.TableCell

			if r < 1 {
				switch c {
				case 0:
					cell = tview.NewTableCell(dict.gameTableHeaderRating).SetTextColor(tview.Styles.SecondaryTextColor)
				case 1:
					cell = tview.NewTableCell(dict.gameTableHeaderName).SetTextColor(tview.Styles.SecondaryTextColor)
				case 2:
					cell = tview.NewTableCell(dict.gameTableHeaderSourcePort).SetTextColor(tview.Styles.SecondaryTextColor)
				case 3:
					cell = tview.NewTableCell(dict.gameTableHeaderIwad).SetTextColor(tview.Styles.SecondaryTextColor)
				default:
					cell = tview.NewTableCell("").SetTextColor(tview.Styles.SecondaryTextColor)
				}
			} else {
				switch c {
				case 0:
					cell = tview.NewTableCell(game.RatingString()).SetTextColor(tview.Styles.TitleColor)
				case 1:
					cell = tview.NewTableCell(game.Name).SetTextColor(tview.Styles.SecondaryTextColor)
				case 2:
					cell = tview.NewTableCell(game.Port).SetTextColor(tview.Styles.PrimaryTextColor)
				case 3:
					cell = tview.NewTableCell(game.Iwad).SetTextColor(tview.Styles.PrimaryTextColor)
				default:
					cell = tview.NewTableCell("").SetTextColor(tview.Styles.PrimaryTextColor)
				}
			}
			gamesTable.SetCell(r, c, cell)
		}
	}

	gamesTable.SetSelectionChangedFunc(func(r int, c int) {
		var g *games.Game
		switch r {
		case 0:
			return
		default:
			g = &allGames[r-fixRows]
		}
		selectedGameChanged(g)
	})

	gamesTable.SetSelectedFunc(func(r int, c int) {
		switch {
		case r > 0:
			allGames[r-fixRows].Run()
		}
	})

	gamesTable.SetInputCapture(withQuit(handleKeyBindings))
}
