package tui

import (
	"fmt"
	"strings"

	"github.com/rivo/tview"
	"github.com/zmnpl/twad/games"
)

//  stats
func makeStatsTable() *tview.Table {
	stats := tview.NewTable().
		SetFixed(1, 1).
		SetSelectable(false, false).
		SetBorders(tableBorders).SetSeparator(':')

	return stats
}

func populateStats(g *games.Game) {
	statsTable.Clear()
	row := 0
	pts := float64(g.Playtime) / 1000 / 60
	saves := g.SaveCount()

	statsTable.SetCell(row, 0, tview.NewTableCell("# Savegames").SetTextColor(tview.Styles.SecondaryTextColor))
	statsTable.SetCell(row, 1, tview.NewTableCell(fmt.Sprintf("%v", saves)).SetAlign(tview.AlignLeft))
	row++
	statsTable.SetCell(row, 0, tview.NewTableCell("Playtime").SetTextColor(tview.Styles.SecondaryTextColor))
	statsTable.SetCell(row, 1, tview.NewTableCell(fmt.Sprintf("%.2f min", pts)).SetAlign(tview.AlignLeft))
	row++
	statsTable.SetCell(row, 0, tview.NewTableCell("Last Played").SetTextColor(tview.Styles.SecondaryTextColor))
	statsTable.SetCell(row, 1, tview.NewTableCell(fmt.Sprint(g.LastPlayed)).SetAlign(tview.AlignLeft))
	row++
	statsTable.SetCell(row, 0, tview.NewTableCell(""))
	//	statsTable.SetCell(row, 1, tview.NewTableCell(""))
	row++

	for k, v := range g.Stats {
		statsTable.SetCell(row, 0, tview.NewTableCell(strings.Title("# "+k)).SetTextColor(tview.Styles.SecondaryTextColor))
		statsTable.SetCell(row, 1, tview.NewTableCell(fmt.Sprintf("%v", v)).SetAlign(tview.AlignLeft))
		row++
	}

	statsTable.SetCell(row, 0, tview.NewTableCell("                    ").SetTextColor(tview.Styles.SecondaryTextColor))
	statsTable.SetCell(row, 1, tview.NewTableCell("                    ").SetAlign(tview.AlignLeft))
}
