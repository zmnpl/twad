package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/zmnpl/twad/games"
	"github.com/zmnpl/twad/helper"
)

const (
	gameTableHeaderRating     = "Rating"
	gameTableHeaderName       = "Name"
	gameTableHeaderSourcePort = "SourcePort"
	gameTableHeaderIwad       = "Iwad"

	deleteGameQuestion = "Delete '%v'?"
	deleteModQuestion  = "Remove '%v' from '%v'?"
)

// center table with mods
func makeGamesTable() *tview.Table {
	gamesTable = tview.NewTable().
		SetFixed(1, 1).
		SetSelectable(true, false).
		SetBorders(tableBorders).SetSeparator('|')
	gamesTable.SetBorderPadding(0, 0, 1, 2)
	return gamesTable
}

func populateGamesTable() {
	gamesTable.Clear()
	allGames := games.Games()

	fixRows, fixCols := 1, 4
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
					cell = tview.NewTableCell(gameTableHeaderRating).SetTextColor(tview.Styles.SecondaryTextColor)
				case 1:
					cell = tview.NewTableCell(gameTableHeaderName).SetTextColor(tview.Styles.SecondaryTextColor)
				case 2:
					cell = tview.NewTableCell(gameTableHeaderSourcePort).SetTextColor(tview.Styles.SecondaryTextColor)
				case 3:
					cell = tview.NewTableCell(gameTableHeaderIwad).SetTextColor(tview.Styles.SecondaryTextColor)
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

	gamesTable.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		k := event.Key()
		r, _ := gamesTable.GetSelection()
		dialogXOffset := 7

		if k == tcell.KeyRune {
			switch event.Rune() {

			// get out
			case 'q':
				app.Stop()
				return nil

			// show credits and license
			case 'c':
				// c again to toggle
				frontPage, _ := detailSidePagesSub2.GetFrontPage()
				if frontPage == pageLicense {
					appModeNormal()
					return nil
				}
				lp := makeLicense()
				detailSidePagesSub2.AddPage(pageLicense, lp, true, true)
				detailSidePagesSub2.SwitchToPage(pageLicense)
				app.SetFocus(lp)
				return nil

			// options
			case 'o':
				optionsDiag := makeOptions()
				contentPages.AddPage(pageOptions, optionsDiag, true, false)
				contentPages.SwitchToPage(pageOptions)
				app.SetFocus(optionsDiag)

			// new game
			case 'n':
				newForm := makeAddEditGame(nil)
				detailPages.AddPage(pageAddEdit, newForm, true, false)
				detailPages.SwitchToPage(pageAddEdit)
				app.SetFocus(newForm)
				return nil

			// increase game rating
			case '+':
				allGames[r-fixRows].Rate(1)
				c := tview.NewTableCell(allGames[r-fixRows].RatingString()).SetTextColor(tview.Styles.TitleColor)
				gamesTable.SetCell(r, 0, c)
				games.Persist()

			// decrease game rating
			case '-':
				allGames[r-fixRows].Rate(-1)
				c := tview.NewTableCell(allGames[r-fixRows].RatingString()).SetTextColor(tview.Styles.TitleColor)
				gamesTable.SetCell(r, 0, c)
				games.Persist()

			// warp
			case 'w':
				if r > 0 {
					_, _, _, commandPreviewHeight := commandPreview.GetRect()
					rowOffset, _ := gamesTable.GetOffset() // account for, when table is scrolled beyond visible screen
					warp := makeWarpRecord(allGames[r-fixRows], appModeNormal, dialogXOffset, commandPreviewHeight+(r+1)-rowOffset, contentPages.Box)
					contentPages.AddPage(pageWarp, warp, true, true)
					app.SetFocus(warp)
					return nil
				}

			// demos
			case 'd':
				if r > 0 {
					g := allGames[r-fixRows]
					demoList, err := makeDemoList(&g)
					if err != nil {
						return nil
					}
					detailSidePagesSub1.AddPage(pageDemos, demoList, true, true)
					app.SetFocus(demoList)
					return nil
				}

			// savegames
			case 'z':
				if r > 0 {
					g := allGames[r-fixRows]
					savegameList, err := makeSavegameList(&g)
					if err != nil {
						return nil
					}
					detailSidePagesSub1.AddPage(pageSaves, savegameList, true, true)
					app.SetFocus(savegameList)
					return nil
				}
			// add mod to game
			case 'm':
				if r > 0 {
					modTree := makeModTree(&allGames[r-fixRows])
					detailSidePagesSub2.AddPage(pageModSelector, modTree, true, false)
					detailSidePagesSub2.SwitchToPage(pageModSelector)
					app.SetFocus(modTree)
					return nil
				}

			// open game url in browser
			case 'u':
				if r > 0 {
					g := allGames[r-fixRows]
					helper.Openbrowser(g.Link)
					return nil
				}

			case 'i':
				contentPages.SwitchToPage(pageZipImport)
				app.SetFocus(zipInput.selectTree)
				return nil

			// edit selected game
			case 'e':
				if r > 0 {
					addEdit := makeAddEditGame(&allGames[r-fixRows])
					detailPages.AddPage(pageAddEdit, addEdit, true, false)
					detailPages.SwitchToPage(pageAddEdit)
					app.SetFocus(addEdit)
					return nil
				}

			// sort alphabetically
			case 's':
				games.SortAlph()
				populateGamesTable()
				return nil

			}

		}

		if k == tcell.KeyF1 {
			showHelp()
			return nil
		}

		if k == tcell.KeyF2 {
			contentPages.SwitchToPage(pageIdgamesBrowser)
			app.SetFocus(idgamesBrowser.GetRootLayout())
			//idgamesBrowser.UpdateLatest()
			return nil
		}

		// "quickload" tries to load the latest savegame
		if k == tcell.KeyF9 && r > 0 {
			allGames[r-fixRows].Quickload()
		}

		// delete selected game
		if k == tcell.KeyDelete && r > 0 {
			remove := func() {
				if r == gamesTable.GetRowCount()-1 {
					gamesTable.Select(r-1, 0)
				}
				games.RemoveGameAt(r - fixRows)
				contentPages.RemovePage(pageYouSure)
				app.SetFocus(gamesTable)
			}

			if config.DeleteWithoutWarning {
				remove()
				return nil
			}

			rowOffset, _ := gamesTable.GetOffset() // account for, when table is scrolled beyond visible screen
			_, _, _, commandPreviewHeight := commandPreview.GetRect()
			youSure := makeYouSureBox(allGames[r-fixRows].Name, remove, appModeNormal, dialogXOffset, commandPreviewHeight+(r+1)-rowOffset, contentPages.Box)
			contentPages.AddPage(pageYouSure, youSure, true, true)
			return nil
		}

		// switch focus to mod list
		if k == tcell.KeyTAB && r >= fixRows {
			// only if game has mods
			if len(allGames[r-fixRows].Mods) > 0 {
				app.SetFocus(detailSidePagesSub1)
			}
		}

		return event
	})
}
