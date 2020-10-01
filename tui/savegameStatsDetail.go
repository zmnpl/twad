package tui

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"github.com/zmnpl/twad/games/savesStats"
	st "github.com/zmnpl/twad/games/savesStats"
)

//  stats
func makeLevelStatsTable(s st.Savegame, lastFocus *tview.List) *tview.Table {
	stats := tview.NewTable().
		SetFixed(1, 1).
		SetSelectable(true, false).
		SetBorders(tableBorders).SetSeparator(':')
	stats.SetBorderPadding(0, 0, 1, 1)

	row := 0
	stats.SetCell(row, 0, tview.NewTableCell("                    ").SetTextColor(tview.Styles.SecondaryTextColor))
	stats.SetCell(row, 1, tview.NewTableCell("                    ").SetAlign(tview.AlignLeft))
	row++

	populate := func(lvl st.MapStats) {
		stats.SetCell(row, 0, tview.NewTableCell(strings.ToUpper(lvl.LevelName)).SetTextColor(tview.Styles.ContrastBackgroundColor))
		row++
		stats.SetCell(row+0, 0, tview.NewTableCell("Time").SetTextColor(tview.Styles.SecondaryTextColor))
		stats.SetCell(row+1, 0, tview.NewTableCell("# Kills").SetTextColor(tview.Styles.SecondaryTextColor))
		stats.SetCell(row+2, 0, tview.NewTableCell("# Secrets").SetTextColor(tview.Styles.SecondaryTextColor))
		stats.SetCell(row+3, 0, tview.NewTableCell("# Items").SetTextColor(tview.Styles.SecondaryTextColor))
		stats.SetCell(row+0, 1, tview.NewTableCell(fmt.Sprintf("%vm %vs", uint32(lvl.LevelTime/60), lvl.LevelTime%60)).SetAlign(tview.AlignLeft))
		stats.SetCell(row+1, 1, tview.NewTableCell(fmt.Sprintf("%v/%v", lvl.KillCount, lvl.TotalKills)).SetAlign(tview.AlignLeft))
		stats.SetCell(row+2, 1, tview.NewTableCell(fmt.Sprintf("%v/%v", lvl.SecretCount, lvl.TotalSecrets)).SetAlign(tview.AlignLeft))
		stats.SetCell(row+3, 1, tview.NewTableCell(fmt.Sprintf("%v/%v", lvl.ItemCount, lvl.TotalItems)).SetAlign(tview.AlignLeft))
		row += 4
	}

	// add totals
	populate(savesStats.SummarizeStats(s.Levels))
	stats.Select(1, 0)
	row += 1

	// add all levels
	for _, lvl := range s.Levels {
		populate(lvl)
	}

	stats.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		k := event.Key()
		if k == tcell.KeyTAB {
			app.SetFocus(lastFocus)
		}

		if k == tcell.KeyRune {
			switch event.Rune() {
			// quit app from here as well
			case 'q':
				app.Stop()
				return nil
			}
		}
		return event
	})

	return stats
}
