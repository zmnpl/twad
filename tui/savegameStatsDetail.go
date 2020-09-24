package tui

import (
	"fmt"
	"strings"

	"github.com/rivo/tview"
	st "github.com/zmnpl/twad/games/savesStats"
)

//  stats
func makeLevelStatsTable(s st.Savegame) *tview.Table {
	stats := tview.NewTable().
		SetFixed(1, 1).
		SetSelectable(false, false).
		SetBorders(tableBorders).SetSeparator(':')
	stats.SetBorderPadding(0, 0, 1, 1)

	row := 0

	stats.SetCell(row, 0, tview.NewTableCell("                    ").SetTextColor(tview.Styles.SecondaryTextColor))
	stats.SetCell(row, 1, tview.NewTableCell("                    ").SetAlign(tview.AlignLeft))
	row++

	for _, level := range s.Levels {
		stats.SetCell(row, 0, tview.NewTableCell(strings.ToUpper(level.LevelName)).SetTextColor(tview.Styles.ContrastBackgroundColor))
		row++
		stats.SetCell(row+0, 0, tview.NewTableCell("# Kills").SetTextColor(tview.Styles.SecondaryTextColor))
		stats.SetCell(row+1, 0, tview.NewTableCell("# Secrets").SetTextColor(tview.Styles.SecondaryTextColor))
		stats.SetCell(row+2, 0, tview.NewTableCell("# Items").SetTextColor(tview.Styles.SecondaryTextColor))
		stats.SetCell(row+0, 1, tview.NewTableCell(fmt.Sprintf("%v/%v", level.KillCount, level.TotalKills)).SetAlign(tview.AlignLeft))
		stats.SetCell(row+1, 1, tview.NewTableCell(fmt.Sprintf("%v/%v", level.SecretCount, level.TotalSecrets)).SetAlign(tview.AlignLeft))
		stats.SetCell(row+2, 1, tview.NewTableCell(fmt.Sprintf("%v/%v", level.ItemCount, level.TotalItems)).SetAlign(tview.AlignLeft))
		row += 3
	}

	return stats
}
