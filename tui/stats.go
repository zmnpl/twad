package tui

import (
	"fmt"
	"strings"

	"github.com/rivo/tview"
	"github.com/zmnpl/twad/games"
)

//  stats
func makeStatsTable(g *games.Game) *tview.Table {
	stats := tview.NewTable().
		SetFixed(1, 1).
		SetSelectable(false, false).
		SetBorders(tableBorders).SetSeparator(':')
	stats.SetBorderPadding(0, 0, 1, 1)

	if g == nil {
		return stats
	}

	row := 0
	pts := float64(g.Playtime) / 1000 / 60

	stats.SetCell(row, 0, tview.NewTableCell("# Savegames").SetTextColor(tview.Styles.SecondaryTextColor))
	stats.SetCell(row, 1, tview.NewTableCell(fmt.Sprintf("%v", g.SaveCount())).SetAlign(tview.AlignLeft))
	row++
	stats.SetCell(row, 0, tview.NewTableCell("# Demos").SetTextColor(tview.Styles.SecondaryTextColor))
	stats.SetCell(row, 1, tview.NewTableCell(fmt.Sprintf("%v", g.DemoCount())).SetAlign(tview.AlignLeft))
	row++
	stats.SetCell(row, 0, tview.NewTableCell("Playtime").SetTextColor(tview.Styles.SecondaryTextColor))
	stats.SetCell(row, 1, tview.NewTableCell(fmt.Sprintf("%.2f min", pts)).SetAlign(tview.AlignLeft))
	row++
	stats.SetCell(row, 0, tview.NewTableCell("Last Played").SetTextColor(tview.Styles.SecondaryTextColor))
	stats.SetCell(row, 1, tview.NewTableCell(fmt.Sprint(g.LastPlayed)).SetAlign(tview.AlignLeft))
	row++
	stats.SetCell(row, 0, tview.NewTableCell(""))
	//	stats.SetCell(row, 1, tview.NewTableCell(""))
	row++

	for k, v := range g.Stats {
		stats.SetCell(row, 0, tview.NewTableCell(strings.Title("# "+k)).SetTextColor(tview.Styles.SecondaryTextColor))
		stats.SetCell(row, 1, tview.NewTableCell(fmt.Sprintf("%v", v)).SetAlign(tview.AlignLeft))
		row++
	}

	stats.SetCell(row, 0, tview.NewTableCell("                    ").SetTextColor(tview.Styles.SecondaryTextColor))
	stats.SetCell(row, 1, tview.NewTableCell("                    ").SetAlign(tview.AlignLeft))

	return stats
}
