package bubbletui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/zmnpl/twad/base"
	"github.com/zmnpl/twad/games"
	"github.com/zmnpl/twad/tui/shared"
)

// btPanel is a part of the bubbletea ui, which can take the focus
type btPanel interface {
	update(m *btModel, key tea.KeyMsg) tea.Cmd
	view(m *btModel, width, height int) string
}

// btTyping is implemented by panels which take text input; q doesn't quit while they have the focus
type btTyping interface {
	typing() bool
}

type btFocus int

const (
	btFocusGames btFocus = iota
	btFocusSide1
	btFocusSide2
)

// btDialog is a panel shown on top of the content
type btDialog struct {
	btPanel
	width  int
	atGame bool // place right below the selected game, centered otherwise
}

// btModel holds the state of the bubbletea user interface
type btModel struct {
	cfg       *base.Cfg
	theme     shared.Theme
	st        btStyles
	logo      string
	program   *tea.Program
	keyGroups []btKeyGroup

	width, height int

	cursor      int // selected game
	offset      int // scroll offset of the games table
	tableHeight int // visible game rows of the last render, for paging
	focus       btFocus

	content btPanel    // replaces the whole content area (help, options, ...)
	detail  btPanel    // replaces the detail panes next to the games table (add/edit)
	mods    *btModList // first detail pane by default
	side1   btPanel    // shown instead of the mod list
	side2   btPanel    // shown instead of the stats
	dialog  *btDialog

	// kept between openings
	zipTree *btTree
	idgames *btIdgames

	status string
}

type (
	btGameDoneMsg     struct{}
	btStatusMsg       string
	btDownloadDoneMsg struct {
		path, title string
		err         error
	}
	btImportDoneMsg struct{ err error }
)

// Draw starts the user interface built with bubbletea
func Draw() error {
	m := &btModel{cfg: base.Config()}
	m.applyTheme()
	m.keyGroups = btKeyGroups()
	m.selectedGameChanged()

	// settings - only when first start of app
	if !m.cfg.Configured {
		m.content = m.newHello()
		m.cfg.Configured = true
	}

	m.program = tea.NewProgram(m, tea.WithAltScreen())
	_, err := m.program.Run()
	return err
}

func (m *btModel) Init() tea.Cmd {
	return tea.SetWindowTitle("twad")
}

func (m *btModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	case tea.KeyMsg:
		return m, m.handleKey(msg)
	case btGameDoneMsg:
		m.status = ""
		m.selectedGameChanged()
	case btStatusMsg:
		m.status = string(msg)
	case btDownloadDoneMsg:
		m.status = ""
		if m.idgames != nil {
			m.idgames.downloading = false
		}
		if msg.err != nil {
			m.showError("Download Failed", msg.err.Error(), nil)
			break
		}
		// suggest import of downloaded archive
		m.showZipImport(msg.path, msg.title)
	case btImportDoneMsg:
		m.status = ""
		if msg.err != nil {
			m.showError(shared.Dict.ZipImportFailed, msg.err.Error(), nil)
		}
	case btIdgamesMsg, btIdgameDetailMsg:
		if m.idgames != nil {
			return m, m.idgames.handle(msg)
		}
	}
	return m, nil
}

func (m *btModel) handleKey(key tea.KeyMsg) tea.Cmd {
	switch key.String() {
	case "ctrl+c":
		return tea.Quit
	case "esc":
		// switch back to normal mode
		m.reset()
		return nil
	}

	p := m.focused()
	if key.String() == "q" {
		if t, ok := p.(btTyping); !ok || !t.typing() {
			return tea.Quit
		}
	}
	if p == nil {
		return m.updateGames(key)
	}
	return p.update(m, key)
}

// focused returns the panel which gets the keys; nil means the games table
func (m *btModel) focused() btPanel {
	switch {
	case m.dialog != nil:
		return m.dialog.btPanel
	case m.content != nil:
		return m.content
	case m.detail != nil:
		return m.detail
	case m.focus == btFocusSide1:
		return m.side1Panel()
	case m.focus == btFocusSide2 && m.side2 != nil:
		return m.side2
	}
	return nil
}

func (m *btModel) isFocused(p btPanel) bool {
	return m.focused() == p
}

func (m *btModel) side1Panel() btPanel {
	if m.side1 != nil {
		return m.side1
	}
	return m.mods
}

// selectedGame returns the game selected in the games table
func (m *btModel) selectedGame() (*games.Game, bool) {
	allGames := games.Games()
	if m.cursor < 0 || m.cursor >= len(allGames) {
		return nil, false
	}
	return &allGames[m.cursor], true
}

// selectedGameChanged updates everything showing details of the selected game
func (m *btModel) selectedGameChanged() {
	m.cursor = max(0, min(m.cursor, games.GameCount()-1))
	g, ok := m.selectedGame()
	if !ok {
		g = &games.Game{}
	}
	m.mods = m.newModList(g)
}

// reset brings the ui back into its default state
func (m *btModel) reset() {
	m.content = nil
	m.detail = nil
	m.side1 = nil
	m.side2 = nil
	m.dialog = nil
	m.focus = btFocusGames
	m.status = ""
	m.selectedGameChanged()
}

func (m *btModel) showDialog(p btPanel, width int, atGame bool) {
	m.dialog = &btDialog{btPanel: p, width: width, atGame: atGame}
}

func (m *btModel) closeDialog() {
	m.dialog = nil
}

// launch runs a game without blocking the ui
func (m *btModel) launch(g *games.Game, run func()) tea.Cmd {
	m.status = fmt.Sprintf(" Playing %v ...", g.Name)
	return func() tea.Msg {
		run()
		return btGameDoneMsg{}
	}
}

// games table
//#######################################################################

// x offset for dialogs that pop up at the selected game
const btDialogXOffset = 7

// btGamesColumn describes a column of the games table and how to show a game in it
type btGamesColumn struct {
	header string
	value  func(g *games.Game) string
	style  lipgloss.Style
}

func (m *btModel) gamesColumns() []btGamesColumn {
	return []btGamesColumn{
		{shared.Dict.GameTableHeaderRating, (*games.Game).RatingString, m.st.title},
		{shared.Dict.GameTableHeaderName, func(g *games.Game) string { return g.Name }, m.st.secondary},
		{shared.Dict.GameTableHeaderSourcePort, func(g *games.Game) string { return g.Port }, m.st.primary},
		{shared.Dict.GameTableHeaderIwad, func(g *games.Game) string { return g.Iwad }, m.st.primary},
	}
}

func (m *btModel) updateGames(key tea.KeyMsg) tea.Cmd {
	if c, ok := btMove(key.String(), m.cursor, games.GameCount(), m.tableHeight); ok {
		if c != m.cursor {
			m.cursor = c
			m.selectedGameChanged()
		}
		return nil
	}

	if key.String() == "enter" {
		if g, ok := m.selectedGame(); ok {
			return m.launch(g, func() { g.Run() })
		}
		return nil
	}

	return m.handleKeyBindings(key)
}

func (m *btModel) gamesView(width, height int) string {
	const padLeft, padRight = 1, 2
	cols := m.gamesColumns()
	allGames := games.Games()

	// column widths from content
	widths := make([]int, len(cols))
	for c, col := range cols {
		widths[c] = ansi.StringWidth(col.header)
		for i := range allGames {
			widths[c] = max(widths[c], ansi.StringWidth(col.value(&allGames[i])))
		}
	}
	// shrink the name column, when the table doesn't fit
	total := len(cols) - 1
	for _, w := range widths {
		total += w
	}
	if over := total - (width - padLeft - padRight); over > 0 {
		widths[1] = max(12, widths[1]-over)
	}

	cell := func(s string, c int) string {
		s = btTruncate(s, widths[c])
		return s + strings.Repeat(" ", widths[c]-ansi.StringWidth(s))
	}
	sep := m.st.title.Render("|")

	header := make([]string, len(cols))
	for c, col := range cols {
		header[c] = m.st.secondary.Render(cell(col.header, c))
	}
	lines := []string{strings.Join(header, sep)}

	m.tableHeight = max(1, height-1)
	m.offset = btScroll(m.cursor, m.offset, m.tableHeight)
	focused := m.focused() == nil
	for i := m.offset; i < len(allGames) && i < m.offset+m.tableHeight; i++ {
		g := &allGames[i]
		row := make([]string, len(cols))
		if i == m.cursor {
			for c, col := range cols {
				row[c] = cell(col.value(g), c)
			}
			sel := m.st.selected
			if !focused {
				sel = m.st.button
			}
			lines = append(lines, sel.Render(strings.Join(row, "|")))
			continue
		}
		for c, col := range cols {
			row[c] = col.style.Render(cell(col.value(g), c))
		}
		lines = append(lines, strings.Join(row, sep))
	}

	for i, l := range lines {
		lines[i] = strings.Repeat(" ", padLeft) + l
	}
	return btFit(strings.Join(lines, "\n"), width, height)
}

// layout
//#######################################################################

const btCommandPreviewHeight = 4

func (m *btModel) View() string {
	if m.width <= 0 || m.height <= 0 {
		return ""
	}

	header := m.headerView()
	footer := m.footerView()
	statusline := btFit(btTruncate(m.status, m.width), m.width, 1)

	contentHeight := m.height - lipgloss.Height(header) - lipgloss.Height(footer) - 1
	// not enough space: the content gets the space of the header
	if contentHeight < 20 {
		header = ""
		contentHeight = m.height - lipgloss.Height(footer) - 1
	}
	if contentHeight < 8 {
		footer = ""
		contentHeight = m.height - 1
	}

	content := m.contentView(m.width, contentHeight)
	if m.dialog != nil {
		content = m.overlayDialog(content, contentHeight)
	}

	parts := make([]string, 0, 4)
	for _, p := range []string{header, content, footer, statusline} {
		if p != "" {
			parts = append(parts, p)
		}
	}
	return strings.Join(parts, "\n")
}

// small or big header
func (m *btModel) headerView() string {
	if m.cfg.HideHeader {
		return btFit(m.st.tags(shared.Subtitle), m.width, 1)
	}
	return btFit(m.logo, m.width, lipgloss.Height(m.logo))
}

func (m *btModel) contentView(width, height int) string {
	if m.content != nil {
		return btFit(m.content.view(m, width, height), width, height)
	}

	g, ok := m.selectedGame()
	if !ok {
		g = &games.Game{}
	}
	preview := m.st.tags(shared.PreviewText + " " + shared.StylizeCommandList(g.CommandList(), m.theme.Tag(m.theme.SecondaryText), m.theme.Tag(m.theme.Title), m.theme.Tag(m.theme.PrimaryText)))
	preview = lipgloss.NewStyle().Width(width - 2).Render(preview)
	preview = m.st.box("", preview, width, btCommandPreviewHeight, m.st.borderColor)

	height -= btCommandPreviewHeight
	tableWidth := max(1, min(width, width*m.cfg.GameListRelativeWidth/100))
	table := m.gamesView(tableWidth, height)

	detailWidth := width - tableWidth
	if detailWidth <= 0 {
		return preview + "\n" + table
	}
	var detail string
	if m.detail != nil {
		detail = btFit(m.detail.view(m, detailWidth, height), detailWidth, height)
	} else {
		detail = m.sidesView(detailWidth, height)
	}
	return preview + "\n" + lipgloss.JoinHorizontal(lipgloss.Top, table, detail)
}

// detail panes; side by side on wide screens, stacked otherwise
func (m *btModel) sidesView(width, height int) string {
	side2 := func(w, h int) string {
		if m.side2 != nil {
			return btFit(m.side2.view(m, w, h), w, h)
		}
		g, ok := m.selectedGame()
		if !ok {
			return btFit("", w, h)
		}
		return m.statsView(g, w, h)
	}

	if width >= 75 {
		w1 := width * 4 / 10
		return lipgloss.JoinHorizontal(lipgloss.Top,
			btFit(m.side1Panel().view(m, w1, height), w1, height),
			side2(width-w1, height))
	}
	h1 := height / 2
	return btFit(m.side1Panel().view(m, width, h1), width, h1) + "\n" + side2(width, height-h1)
}

// overlayDialog draws the dialog on top of the content
func (m *btModel) overlayDialog(content string, contentHeight int) string {
	d := m.dialog
	width := m.width - 4
	if d.width > 0 {
		width = min(d.width, m.width-2)
	}
	box := d.view(m, width, 0)
	boxWidth, boxHeight := lipgloss.Width(box), lipgloss.Height(box)

	x := (m.width - boxWidth) / 2
	y := (contentHeight - boxHeight) / 2
	if d.atGame && m.content == nil {
		// right below the selected game; above it, when there is no space
		x = btDialogXOffset
		row := btCommandPreviewHeight + 1 + m.cursor - m.offset
		y = row + 1
		if y+boxHeight > contentHeight {
			y = row - boxHeight
		}
		y = max(0, y)
	}
	return btOverlay(content, box, x, y)
}
