package bubbletui

import (
	"errors"
	"fmt"
	"maps"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/zmnpl/twad/base"
	"github.com/zmnpl/twad/games"
	"github.com/zmnpl/twad/ports/savesStats"
	"github.com/zmnpl/twad/tui/shared"
)

// mod list
//#######################################################################

type btModList struct {
	btList
	g    *games.Game
	edit bool // moving the selected mod around
}

func (m *btModel) newModList(g *games.Game) *btModList {
	l := &btModList{btList: btList{title: shared.Dict.OverviewMods}, g: g}
	for _, mod := range g.Mods {
		l.items = append(l.items, btListItem{main: path.Base(mod), secondary: path.Dir(mod), warn: !base.ModOk(mod)})
	}
	return l
}

func (l *btModList) update(m *btModel, key tea.KeyMsg) tea.Cmd {
	switch key.String() {
	// switch back to game table
	case "tab":
		l.editOff()
		m.focus = btFocusGames
		return nil

	case "enter":
		if l.edit {
			l.editOff()
		} else if len(l.items) > 0 {
			l.edit = true
		}
		return nil

	case "delete":
		// when in edit mode, this is only confusing
		if l.edit || len(l.items) == 0 {
			return nil
		}
		i := l.cursor
		remove := func() tea.Cmd {
			l.remove(i)
			l.g.RemoveMod(i)
			games.Persist()
			return nil
		}
		if m.cfg.DeleteWithoutWarning {
			return remove()
		}
		m.confirm(l.g.Mods[i], remove, false)
		return nil

	// add mod
	case "m":
		return m.showModTree(l.g)
	}

	last := l.cursor
	if moved, _ := l.move(key.String()); moved && l.edit {
		// switch mod positions in game and list
		l.g.SwitchMods(last, l.cursor)
		l.items[last], l.items[l.cursor] = l.items[l.cursor], l.items[last]
	}
	return nil
}

func (l *btModList) editOff() {
	if l.edit {
		l.edit = false
		games.Persist()
	}
}

func (l *btModList) view(m *btModel, width, height int) string {
	return l.btList.view(&m.st, width, height, m.isFocused(l), l.edit)
}

// demos
//#######################################################################

type btDemoList struct {
	btList
	g     *games.Game
	demos []os.DirEntry
}

func (m *btModel) newDemoList(g *games.Game) (*btDemoList, error) {
	demos, err := g.Demos()
	if err != nil {
		return nil, err
	}
	if len(demos) == 0 {
		return nil, errors.New("no demos in demo dir")
	}

	l := &btDemoList{btList: btList{title: shared.Dict.DemosHeader}, g: g, demos: demos}
	l.populate()
	return l, nil
}

func (l *btDemoList) populate() {
	l.items = l.items[:0]
	for _, demo := range l.demos {
		t, size := time.Now(), int64(0)
		if di, err := demo.Info(); err == nil {
			t, size = di.ModTime(), di.Size()
		}
		l.items = append(l.items, btListItem{
			main:      demo.Name(),
			secondary: fmt.Sprintf("%v (%.2f KiB)", t.Format("2006-01-02 15:04"), float32(size)/1024),
		})
	}
	l.cursor = max(0, min(l.cursor, len(l.items)-1))
}

func (l *btDemoList) update(m *btModel, key tea.KeyMsg) tea.Cmd {
	switch key.String() {
	case "tab":
		m.focus = btFocusGames
		return nil

	// enter plays demo
	case "enter":
		name := l.demos[l.cursor].Name()
		return m.launch(l.g, func() { l.g.PlayDemo(name) })

	case "delete":
		name := l.demos[l.cursor].Name()
		remove := func() tea.Cmd {
			demos, err := l.g.RemoveDemo(name)
			if err != nil {
				m.showError("could not remove demo", err.Error(), nil)
				return nil
			}
			l.demos = demos
			if len(demos) == 0 {
				m.reset()
			} else {
				l.populate()
			}
			games.Persist()
			return nil
		}
		if m.cfg.DeleteWithoutWarning {
			return remove()
		}
		m.confirm(name, remove, false)
		return nil

	// start zip import from here
	case "i":
		return m.showZipSelect()
	}

	l.move(key.String())
	return nil
}

func (l *btDemoList) view(m *btModel, width, height int) string {
	return l.btList.view(&m.st, width, height, m.isFocused(l), false)
}

// savegames
//#######################################################################

type btSavegameList struct {
	btList
	savegames []*savesStats.Savegame
}

func (m *btModel) newSavegameList(g *games.Game) (*btSavegameList, error) {
	savegames := g.LoadSavegames()
	if len(savegames) == 0 {
		return nil, errors.New("no savegames available")
	}

	l := &btSavegameList{btList: btList{title: shared.Dict.SavesHeader}, savegames: savegames}
	for _, s := range savegames {
		l.items = append(l.items, btListItem{
			main:      s.Meta.Title,
			secondary: fmt.Sprintf("%v (%v)", s.FI.Name(), s.FI.ModTime().Format("2006-01-02 15:04:05")),
		})
	}
	m.side2 = newBTLevelStats(*savegames[0])
	return l, nil
}

func (l *btSavegameList) update(m *btModel, key tea.KeyMsg) tea.Cmd {
	if key.String() == "tab" {
		m.focus = btFocusSide2
		return nil
	}
	if moved, _ := l.move(key.String()); moved {
		m.side2 = newBTLevelStats(*l.savegames[l.cursor])
	}
	return nil
}

func (l *btSavegameList) view(m *btModel, width, height int) string {
	return l.btList.view(&m.st, width, height, m.isFocused(l), false)
}

// stats of the levels of a savegame
type btLevelStats struct {
	stats  btKV
	cursor int
	offset int
	height int
}

func newBTLevelStats(s savesStats.Savegame) *btLevelStats {
	t := &btLevelStats{}
	t.stats.spacer()

	add := func(lvl savesStats.MapStats) {
		t.stats.section(strings.ToUpper(lvl.LevelName))
		t.stats.kv("Time", fmt.Sprintf("%vm %vs", uint32(lvl.LevelTime/60), lvl.LevelTime%60))
		t.stats.progress(lvl)
	}

	// totals first, then all levels
	add(savesStats.SummarizeStats(s.Levels))
	t.cursor = 1
	t.stats.spacer()
	for _, lvl := range s.ReversedLevels() {
		add(lvl)
	}
	return t
}

func (t *btLevelStats) update(m *btModel, key tea.KeyMsg) tea.Cmd {
	if key.String() == "tab" {
		m.focus = btFocusSide1
		return nil
	}
	t.cursor, _ = btMove(key.String(), t.cursor, len(t.stats.rows), t.height/2)
	return nil
}

func (t *btLevelStats) view(m *btModel, width, height int) string {
	lines := t.stats.lines(&m.st, width)
	t.height = height
	t.offset = btScroll(t.cursor, t.offset, height)
	if m.isFocused(t) && t.cursor < len(lines) {
		lines[t.cursor] = m.st.selected.Render(btFit(ansi.Strip(lines[t.cursor]), width, 1))
	}
	return btFit(strings.Join(lines[min(t.offset, len(lines)):], "\n"), width, height)
}

// stats of the selected game
//#######################################################################

func (m *btModel) statsView(g *games.Game, width, height int) string {
	var stats btKV

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

	return btFit(strings.Join(stats.lines(&m.st, width), "\n"), width, height)
}

// mod tree
//#######################################################################

func (m *btModel) newModTree(g *games.Game) (*btTree, error) {
	rootDir := m.cfg.WadDir
	if _, err := os.Stat(rootDir); err != nil {
		return nil, err
	}

	return newBTTree(shared.Dict.ModTreeTitle, rootDir, m.cfg.ModExtensions, false, func(m *btModel, path string) tea.Cmd {
		mod, err := filepath.Rel(m.cfg.WadDir, path)
		if err != nil {
			return nil
		}
		g.AddMod(filepath.ToSlash(mod))
		m.side2 = nil
		m.focus = btFocusGames
		m.selectedGameChanged()
		return nil
	}), nil
}

// scrollable text, as used for help, license and the like
//#######################################################################

type btText struct {
	title  string
	text   string // already styled
	keys   map[string]func(m *btModel) tea.Cmd
	offset int
	lines  int
	height int
}

func (t *btText) update(m *btModel, key tea.KeyMsg) tea.Cmd {
	if f, ok := t.keys[key.String()]; ok {
		return f(m)
	}
	t.offset, _ = btMove(key.String(), t.offset, max(1, t.lines-t.height+1), t.height/2)
	return nil
}

func (t *btText) view(m *btModel, width, height int) string {
	wrapped := lipgloss.NewStyle().Width(max(1, width-4)).Render(t.text)
	lines := strings.Split(wrapped, "\n")
	t.lines, t.height = len(lines), height-2
	t.offset = max(0, min(t.offset, t.lines-t.height))

	content := make([]string, 0, len(lines))
	for _, l := range lines[t.offset:] {
		content = append(content, " "+l)
	}
	return m.st.box(t.title, strings.Join(content, "\n"), width, height, m.st.borderColor)
}

// first start hello
func (m *btModel) newHello() *btText {
	return &btText{
		text: m.st.tags(shared.HelloText),
		keys: map[string]func(m *btModel) tea.Cmd{"o": (*btModel).showOptions},
	}
}

// license; own type, so c can toggle it
type btLicense struct{ btText }

func (m *btModel) newLicense() *btLicense {
	text := fmt.Sprintf("%s\n%s\n\n%s\n%s\n\n%s\n%s",
		shared.DoomLogoCreditHeader, shared.CreditDoomLogo,
		"[orange]Bubble Tea", "The terminal user interface is built with Bubble Tea:\nhttps://github.com/charmbracelet/bubbletea",
		shared.LicenseHeader, shared.MITLicense)

	return &btLicense{btText{
		title: "Credits / License",
		text:  m.st.tags(text),
		keys: map[string]func(m *btModel) tea.Cmd{
			// switch back to normal mode
			"c": func(m *btModel) tea.Cmd { m.reset(); return nil },
		},
	}}
}
