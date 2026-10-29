package tui

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/zmnpl/twad/ports/savesStats"
)

// stats
func makeLevelStatsTable(s savesStats.Savegame, lastFocus *tview.List) *tview.Table {
	stats := newKVTable()
	stats.SetSelectable(true, false)
	stats.spacer()

	populate := func(lvl savesStats.MapStats) {
		stats.section(strings.ToUpper(lvl.LevelName))
		stats.kv("Time", fmt.Sprintf("%vm %vs", uint32(lvl.LevelTime/60), lvl.LevelTime%60))
		stats.progress(lvl)
	}

	// add totals
	populate(savesStats.SummarizeStats(s.Levels))
	stats.Select(1, 0)
	stats.row++

	// add all levels
	for _, lvl := range s.ReversedLevels() {
		populate(lvl)
	}

	stats.SetInputCapture(withQuit(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTAB {
			app.SetFocus(lastFocus)
		}
		return event
	}))

	return stats.Table
}
