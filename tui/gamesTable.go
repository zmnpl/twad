package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/zmnpl/twad/games"
	"github.com/zmnpl/twad/helper"
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

	gamesTable.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		k := event.Key()
		r, _ := gamesTable.GetSelection()

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
				detailSidePagesSub2.switchTo(pageLicense, makeLicense())
				return nil

			// options
			case 'o':
				contentPages.switchTo(pageOptions, makeOptions())

			// new game
			case 'n':
				detailPages.switchTo(pageAddEdit, makeAddEditGame(nil))
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
					contentPages.showAt(pageWarp, makeWarpRecord(allGames[r-fixRows]), dialogXOffset, selectedRowY(), warpDialogWidth, warpDialogHeight)
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
					detailSidePagesSub1.show(pageDemos, demoList)
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
					detailSidePagesSub1.show(pageSaves, savegameList)
					return nil
				}
			// add mod to game
			case 'm':
				if r > 0 {
					detailSidePagesSub2.switchTo(pageModSelector, makeModTree(&allGames[r-fixRows]))
					return nil
				}

			// open game url in browser
			case 'u':
				if r > 0 {
					g := allGames[r-fixRows]
					helper.Openbrowser(g.Link)
					return nil
				}

			// edit selected game
			case 'e':
				if r > 0 {
					detailPages.switchTo(pageAddEdit, makeAddEditGame(&allGames[r-fixRows]))
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
			contentPages.open(pageIdgamesBrowser, idgamesBrowser.GetRootLayout())
			//idgamesBrowser.UpdateLatest()
			return nil
		}

		if k == tcell.KeyF3 {
			contentPages.open(pageZipSelect, zipSelector.selectTree)
			return nil
		}

		if k == tcell.KeyF4 {
			// i := 0
			// for {
			// 	statusline.Write([]byte("fooo " + strconv.Itoa(i)))
			// 	//statusline.SetText("fooo " + strconv.Itoa(i))
			// 	time.Sleep(1 * time.Second)
			// 	i++
			// 	if i > 9 {
			// 		break
			// 	}
			// }
			corporate := `Leverage agile frameworks to provide a robust synopsis for high level overviews. Iterative approaches to corporate strategy foster collaborative thinking to further the overall value proposition. Organically grow the holistic world view of disruptive innovation via workplace diversity and empowerment.
Bring to the table win-win survival strategies to ensure proactive domination. At the end of the day, going forward, a new normal that has evolved from generation X is on the runway heading towards a streamlined cloud solution. User generated content in real-time will have multiple touchpoints for offshoring.
Capitalize on low hanging fruit to identify a ballpark value added activity to beta test. Override the digital divide with additional clickthroughs from DevOps. Nanotechnology immersion along the information highway will close the loop on focusing solely on the bottom line.
[yellow]Press Enter, then Tab/Backtab for word selections`
			numSelections := 0
			f := func() {
				for _, word := range strings.Split(corporate, " ") {
					if word == "the" {
						word = "[red]the[white]"
					}
					if word == "to" {
						word = fmt.Sprintf(`["%d"]to[""]`, numSelections)
						numSelections++
					}
					fmt.Fprintf(statusline, "%s ", word)
					time.Sleep(200 * time.Millisecond)
				}
			}
			go f()
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
			}

			if config.DeleteWithoutWarning {
				remove()
				return nil
			}

			confirm(contentPages, allGames[r-fixRows].Name, dialogXOffset, selectedRowY(), gamesTable, remove)
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
