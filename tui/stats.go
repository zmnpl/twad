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

	// generic stuff
	stats.SetCell(row, 0, tview.NewTableCell("# Savegames").SetTextColor(tview.Styles.SecondaryTextColor))
	stats.SetCell(row, 1, tview.NewTableCell(fmt.Sprintf("%v", g.SaveCount())).SetAlign(tview.AlignLeft))
	row++
	stats.SetCell(row, 0, tview.NewTableCell("# Demos").SetTextColor(tview.Styles.SecondaryTextColor))
	stats.SetCell(row, 1, tview.NewTableCell(fmt.Sprintf("%v", g.DemoCount())).SetAlign(tview.AlignLeft))
	row++
	stats.SetCell(row, 0, tview.NewTableCell("Playtime").SetTextColor(tview.Styles.SecondaryTextColor))
	stats.SetCell(row, 1, tview.NewTableCell(fmt.Sprintf("%.2f min", float64(g.Playtime)/1000/60)).SetAlign(tview.AlignLeft))
	row++
	stats.SetCell(row, 0, tview.NewTableCell("Last Played").SetTextColor(tview.Styles.SecondaryTextColor))
	stats.SetCell(row, 1, tview.NewTableCell(fmt.Sprint(g.LastPlayed)).SetAlign(tview.AlignLeft))
	row++
	stats.SetCell(row, 0, tview.NewTableCell(""))
	//	stats.SetCell(row, 1, tview.NewTableCell(""))
	row++

	// stats from savegames
	stats.SetCell(row, 0, tview.NewTableCell("# Kills").SetTextColor(tview.Styles.SecondaryTextColor))
	stats.SetCell(row+1, 0, tview.NewTableCell("# Secrets").SetTextColor(tview.Styles.SecondaryTextColor))
	stats.SetCell(row+2, 0, tview.NewTableCell("# Items").SetTextColor(tview.Styles.SecondaryTextColor))
	stats.SetCell(row, 1, tview.NewTableCell(fmt.Sprintf("%v/%v", g.StatsSum.KillCount, g.StatsSum.TotalKills)).SetAlign(tview.AlignLeft))
	stats.SetCell(row+1, 1, tview.NewTableCell(fmt.Sprintf("%v/%v", g.StatsSum.SecretCount, g.StatsSum.TotalSecrets)).SetAlign(tview.AlignLeft))
	stats.SetCell(row+2, 1, tview.NewTableCell(fmt.Sprintf("%v/%v", g.StatsSum.ItemCount, g.StatsSum.TotalItems)).SetAlign(tview.AlignLeft))
	row += 4

	// what the game printed into console
	for k, v := range g.ConsoleStats {
		stats.SetCell(row, 0, tview.NewTableCell(strings.Title("# "+k)).SetTextColor(tview.Styles.SecondaryTextColor))
		stats.SetCell(row, 1, tview.NewTableCell(fmt.Sprintf("%v", v)).SetAlign(tview.AlignLeft))
		row++
	}

	stats.SetCell(row, 0, tview.NewTableCell("                    ").SetTextColor(tview.Styles.SecondaryTextColor))
	stats.SetCell(row, 1, tview.NewTableCell("                    ").SetAlign(tview.AlignLeft))

	return stats
}
