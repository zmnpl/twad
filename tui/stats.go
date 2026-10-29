package tui

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/rivo/tview"
	"github.com/zmnpl/twad/games"
)

// stats
func makeStatsTable(g *games.Game) *tview.Table {
	stats := newKVTable()

	if g == nil {
		return stats.Table
	}

	// generic stuff
	stats.section("Overall")
	stats.kv("Added/Edited", g.AddEdit)
	stats.kv("URL", g.Link)
	stats.kv("# Savegames", g.SaveCount())
	stats.kv("# Demos", g.DemoCount())
	stats.kv("Playtime", fmt.Sprintf("%.2f min", float64(g.Playtime)/1000/60))
	stats.kv("Last Played", g.LastPlayed)

	// stats from savegames
	stats.section("Latest Savegame")
	stats.kv("# Maps", len(g.Stats))
	stats.progress(g.StatsTotal)

	// what the game printed into console
	stats.section("Pickups")
	for _, k := range slices.Sorted(maps.Keys(g.ConsoleStats)) {
		stats.kv(strings.Title("# "+k), g.ConsoleStats[k])
	}

	stats.spacer()

	return stats.Table
}
