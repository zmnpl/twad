package bubbletui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/zmnpl/twad/games"
	"github.com/zmnpl/twad/helper"
	"github.com/zmnpl/twad/tui/shared"
)

// btKeyBinding describes a key and what it does in the bubbletea ui
// bindings with an action get handled by the games table, the others are informational
type btKeyBinding struct {
	label  string // key as shown to the user; empty hides the binding from help
	desc   string
	footer bool // show in the always visible keymap

	keys []string // as given by tea.KeyMsg.String()

	run  func(m *btModel) tea.Cmd                // action
	game func(m *btModel, g *games.Game) tea.Cmd // action on the selected game; skipped when there is none
}

type btKeyGroup struct {
	title    string
	bindings []btKeyBinding
}

func btKeyGroups() []btKeyGroup {
	return []btKeyGroup{
		{"Basic Navigation", []btKeyBinding{
			{label: "Arrows (or jk)", desc: "Navigate"},
			{label: "TAB", desc: "Switch Focus", keys: []string{"tab"}, game: (*btModel).focusMods},
			{label: "ENTER", desc: "Confirm"},
			{label: "F1", desc: "Help/Keymap", footer: true, keys: []string{"f1"}, run: (*btModel).showHelp},
			{label: "F2", desc: "IDGames Browser", keys: []string{"f2"}, run: (*btModel).showIdgames},
			{label: "F3", desc: "Import Archive", keys: []string{"f3"}, run: (*btModel).showZipSelect},
		}},
		{"General Functions", []btKeyBinding{
			{label: "ESC", desc: "Reset UI", footer: true},
			{label: "q", desc: "Quit", footer: true},
			{label: "c", desc: "Credits/License", keys: []string{"c"}, run: (*btModel).toggleLicense},
			{label: "o", desc: "Options", footer: true, keys: []string{"o"}, run: (*btModel).showOptions},
			{label: "s", desc: "Sort Games Alphabetically", keys: []string{"s"}, run: (*btModel).sortGames},
		}},
		{"Launch Selected Game", []btKeyBinding{
			{label: "ENTER", desc: "Run Game", footer: true},
			{label: "F9", desc: "Run Last Savegame", footer: true, keys: []string{"f9"}, game: (*btModel).quickload},
			{label: "w", desc: "Warp (+Record)", footer: true, keys: []string{"w"}, game: (*btModel).showWarp},
		}},
		{"Selected Game Details", []btKeyBinding{
			{label: "d", desc: "Demos", footer: true, keys: []string{"d"}, game: (*btModel).showDemos},
			{label: "z", desc: "Savegames", footer: true, keys: []string{"z"}, game: (*btModel).showSavegames},
			{label: "u", desc: "Visit Url", keys: []string{"u"}, game: (*btModel).visitUrl},
		}},
		{"Add/Edit Game", []btKeyBinding{
			{label: "e", desc: "Edit Game", keys: []string{"e"}, game: (*btModel).editGame},
			{label: "m", desc: "Add Mod To Game", footer: true, keys: []string{"m"}, game: (*btModel).showModTree},
			{label: "n", desc: "New Game", footer: true, keys: []string{"n"}, run: (*btModel).newGame},
			{label: "DEL", desc: "Remove Game", footer: true, keys: []string{"delete"}, game: (*btModel).removeGame},
			{label: "+/-", desc: "Rate Game", keys: []string{"+"}, game: btRateGame(1)},
			{keys: []string{"-"}, game: btRateGame(-1)},
		}},
	}
}

func (b btKeyBinding) matches(key string) bool {
	for _, k := range b.keys {
		if k == key {
			return true
		}
	}
	return false
}

func (m *btModel) bindingString(b btKeyBinding) string {
	return m.st.secondary.Render(fmt.Sprintf("%-6v ", b.label)) + m.st.primary.Render(b.desc)
}

// handleKeyBindings runs the action bound to the given key
func (m *btModel) handleKeyBindings(key tea.KeyMsg) tea.Cmd {
	for _, group := range m.keyGroups {
		for _, b := range group.bindings {
			if !b.matches(key.String()) {
				continue
			}
			if b.run != nil {
				return b.run(m)
			}
			if b.game != nil {
				if g, ok := m.selectedGame(); ok {
					return b.game(m, g)
				}
			}
			return nil
		}
	}
	return nil
}

// footerView is the always visible keymap
func (m *btModel) footerView() string {
	const rows, cols = 3, 4

	var footerBindings []btKeyBinding
	for _, group := range m.keyGroups {
		for _, b := range group.bindings {
			if b.footer {
				footerBindings = append(footerBindings, b)
			}
		}
	}

	colWidth := (m.width - 2) / cols
	lines := make([]string, rows)
	for i, b := range footerBindings {
		if i >= rows*cols {
			break
		}
		r := i % rows
		lines[r] += btFit(btTruncate(m.bindingString(b), colWidth-1), colWidth, 1)
	}
	return m.st.box("", strings.Join(lines, "\n"), m.width, rows+2, m.st.borderColor)
}

func (m *btModel) newHelp() *btText {
	var b strings.Builder
	for _, group := range m.keyGroups {
		fmt.Fprintf(&b, "%s\n", m.st.title.Render(group.title))
		for _, binding := range group.bindings {
			if binding.label != "" {
				fmt.Fprintf(&b, "%s\n", m.bindingString(binding))
			}
		}
		b.WriteString("\n")
	}

	return &btText{
		title: "Help",
		text:  b.String(),
		keys: map[string]func(m *btModel) tea.Cmd{
			// switch back to normal mode
			"f1": func(m *btModel) tea.Cmd { m.reset(); return nil },
		},
	}
}

// key actions
//#######################################################################

func (m *btModel) showHelp() tea.Cmd {
	m.reset()
	m.content = m.newHelp()
	return nil
}

func (m *btModel) showZipSelect() tea.Cmd {
	if m.zipTree == nil {
		//rootDir := helper.Home() // TODO: Start from / but preselect /home/user
		m.zipTree = newBTTree(shared.Dict.ZipSelectTitle, "/", ".zip.tar.gz.rar", true, func(m *btModel, path string) tea.Cmd {
			m.showZipImport(path, "")
			return nil
		})
	}
	m.content = m.zipTree
	return nil
}

// c again toggles back
func (m *btModel) toggleLicense() tea.Cmd {
	if _, isLicense := m.side2.(*btLicense); isLicense {
		m.reset()
		return nil
	}
	m.side2 = m.newLicense()
	m.focus = btFocusSide2
	return nil
}

func (m *btModel) showOptions() tea.Cmd {
	m.content = m.newOptions()
	return nil
}

func (m *btModel) sortGames() tea.Cmd {
	games.SortAlph()
	m.selectedGameChanged()
	return nil
}

func (m *btModel) newGame() tea.Cmd {
	m.detail = m.newAddEditGame(nil)
	return nil
}

func (m *btModel) editGame(g *games.Game) tea.Cmd {
	m.detail = m.newAddEditGame(g)
	return nil
}

// "quickload" tries to load the latest savegame
func (m *btModel) quickload(g *games.Game) tea.Cmd {
	return m.launch(g, func() { g.Quickload() })
}

func (m *btModel) showDemos(g *games.Game) tea.Cmd {
	demoList, err := m.newDemoList(g)
	if err != nil {
		return nil
	}
	m.side1 = demoList
	m.focus = btFocusSide1
	return nil
}

func (m *btModel) showSavegames(g *games.Game) tea.Cmd {
	savegameList, err := m.newSavegameList(g)
	if err != nil {
		return nil
	}
	m.side1 = savegameList
	m.focus = btFocusSide1
	return nil
}

func (m *btModel) visitUrl(g *games.Game) tea.Cmd {
	helper.Openbrowser(g.Link)
	return nil
}

func (m *btModel) showModTree(g *games.Game) tea.Cmd {
	modTree, err := m.newModTree(g)
	if err != nil {
		m.showError(shared.Dict.ModTreeErrWadDir, err.Error(), nil)
		return nil
	}
	m.side1 = nil
	m.side2 = modTree
	m.focus = btFocusSide2
	return nil
}

// switch focus to mod list; only if game has mods
func (m *btModel) focusMods(g *games.Game) tea.Cmd {
	if len(g.Mods) > 0 {
		m.side1 = nil
		m.focus = btFocusSide1
	}
	return nil
}

func btRateGame(increment int) func(*btModel, *games.Game) tea.Cmd {
	return func(m *btModel, g *games.Game) tea.Cmd {
		g.Rate(increment)
		games.Persist()
		return nil
	}
}

func (m *btModel) removeGame(g *games.Game) tea.Cmd {
	i := m.cursor
	remove := func() tea.Cmd {
		games.RemoveGameAt(i)
		m.selectedGameChanged()
		return nil
	}

	if m.cfg.DeleteWithoutWarning {
		return remove()
	}
	m.confirm(g.Name, remove, true)
	return nil
}
