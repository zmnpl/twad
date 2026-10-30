package bubbletui

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/zmnpl/goidgames"
	"github.com/zmnpl/twad/base"
)

type (
	btIdgamesMsg struct {
		results []goidgames.Idgame
		err     error
	}
	btIdgameDetailMsg struct{ game goidgames.Idgame }
)

// btIdgames browses the idgames archive and downloads from it
type btIdgames struct {
	search   textinput.Model
	inSearch bool

	results []goidgames.Idgame
	details map[int]goidgames.Idgame // fetched details by id
	loading bool
	loaded  bool
	err     string

	cursor        int
	offset        int
	height        int
	detailsOffset int

	downloading bool
}

func (m *btModel) newIdgames() *btIdgames {
	ti := textinput.New()
	ti.Prompt = ""
	ti.Placeholder = "title, file name or author"
	ti.Cursor.SetMode(cursor.CursorStatic)
	return &btIdgames{search: ti, details: make(map[int]goidgames.Idgame)}
}

func (m *btModel) showIdgames() tea.Cmd {
	if m.idgames == nil {
		m.idgames = m.newIdgames()
	}
	m.content = m.idgames

	b := m.idgames
	if b.loaded || b.loading {
		return nil
	}
	b.loading = true
	return func() tea.Msg {
		results, err := goidgames.LatestFiles(50, 0)
		return btIdgamesMsg{results, err}
	}
}

func (b *btIdgames) typing() bool { return b.inSearch }

// handle takes the results of api calls
func (b *btIdgames) handle(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case btIdgamesMsg:
		b.loading, b.loaded = false, true
		b.results, b.cursor, b.offset, b.detailsOffset = msg.results, 0, 0, 0
		b.err = ""
		if msg.err != nil {
			b.err = msg.err.Error()
		}
		return b.fetchDetails()
	case btIdgameDetailMsg:
		b.details[msg.game.Id] = msg.game
	}
	return nil
}

// fetchDetails gets the full data of the selected file; search results lack the text file
func (b *btIdgames) fetchDetails() tea.Cmd {
	if b.cursor >= len(b.results) {
		return nil
	}
	id := b.results[b.cursor].Id
	if _, ok := b.details[id]; ok {
		return nil
	}
	return func() tea.Msg {
		g, err := goidgames.Get(id, "")
		if err != nil {
			return nil
		}
		return btIdgameDetailMsg{g}
	}
}

func (b *btIdgames) runSearch() tea.Cmd {
	query := strings.TrimSpace(b.search.Value())
	if len(query) < 3 {
		b.err = "Query must at least be 3 characters."
		return nil
	}
	b.loading, b.err = true, ""
	return func() tea.Msg {
		found, err := goidgames.SearchMultipleTypes(query,
			[]string{goidgames.SEARCH_TYPE_TITLE, goidgames.SEARCH_TYPE_FILENAME, goidgames.SEARCH_TYPE_AUTHOR},
			goidgames.SEARCH_SORT_RATING, goidgames.SEARCH_SORT_DESC)

		// the same file can be found by several types
		seen := make(map[int]bool)
		results := make([]goidgames.Idgame, 0, len(found))
		for _, g := range found {
			if !seen[g.Id] {
				seen[g.Id] = true
				results = append(results, g)
			}
		}
		return btIdgamesMsg{results, err}
	}
}

func (b *btIdgames) setSearch(on bool) {
	b.inSearch = on
	if on {
		b.search.Focus()
	} else {
		b.search.Blur()
	}
}

func (b *btIdgames) update(m *btModel, key tea.KeyMsg) tea.Cmd {
	k := key.String()
	if k == "tab" {
		b.setSearch(!b.inSearch)
		return nil
	}

	if b.inSearch {
		if k == "enter" {
			b.setSearch(false)
			return b.runSearch()
		}
		var cmd tea.Cmd
		b.search, cmd = b.search.Update(key)
		return cmd
	}

	switch k {
	case "/":
		b.setSearch(true)
		return nil
	case "pgdown", "ctrl+d":
		b.detailsOffset += 10
		return nil
	case "pgup", "ctrl+u":
		b.detailsOffset = max(0, b.detailsOffset-10)
		return nil
	case "enter":
		if b.downloading || b.cursor >= len(b.results) {
			return nil
		}
		g := b.results[b.cursor]
		m.confirm(fmt.Sprintf("Download %v?", g.Title), func() tea.Cmd { return m.download(g) }, false)
		return nil
	}

	before := b.cursor
	b.cursor, _ = btMove(k, b.cursor, len(b.results), b.height/2)
	if b.cursor != before {
		b.detailsOffset = 0
		return b.fetchDetails()
	}
	return nil
}

// download fetches the file in the background and reports the progress in the status line
func (m *btModel) download(g goidgames.Idgame) tea.Cmd {
	m.idgames.downloading = true
	p := m.program
	return func() tea.Msg {
		path, err := base.DownloadIdGame(g, base.DOWNLOAD_PATH(), func(total uint64) {
			if p != nil {
				p.Send(btStatusMsg(fmt.Sprintf(" Downloading... %v complete", total)))
			}
		})
		return btDownloadDoneMsg{path: path, title: g.Title, err: err}
	}
}

var btTextfileKey = regexp.MustCompile(`^(\S.*?):(.*)`)

func (b *btIdgames) view(m *btModel, width, height int) string {
	st := &m.st

	// search
	b.search.Width = max(1, width-5)
	b.search.TextStyle = st.primary
	searchBorder := st.borderColor
	if b.inSearch {
		searchBorder = btColorOf(st.secondary)
	}
	search := st.box("Search idgames (TAB to switch, ENTER to search)", " "+b.search.View(), width, 3, searchBorder)

	// download path
	dlPath := btFit(" "+st.moreContrast.Render("Download to:")+" "+st.primary.Render(base.DOWNLOAD_PATH()), width, 1)

	bodyHeight := max(3, height-4)
	listWidth := width * 55 / 100
	list := b.listView(m, listWidth, bodyHeight)
	details := b.detailsView(m, width-listWidth, bodyHeight)

	return search + "\n" + lipgloss.JoinHorizontal(lipgloss.Top, list, details) + "\n" + dlPath
}

func (b *btIdgames) listView(m *btModel, width, height int) string {
	st := &m.st
	inner := width - 4
	b.height = height - 3

	var lines []string
	switch {
	case b.loading:
		lines = []string{st.primary.Render("Loading ...")}
	case b.err != "":
		lines = []string{st.warn.Render(b.err)}
	case b.loaded && len(b.results) == 0:
		lines = []string{st.primary.Render("Nothing found")}
	default:
		const ratingWidth, dateWidth = 5, 10
		authorWidth := max(4, inner/4)
		titleWidth := max(4, inner-ratingWidth-dateWidth-authorWidth-3)
		cell := func(s string, w int) string {
			s = btTruncate(s, w)
			return s + strings.Repeat(" ", w-ansi.StringWidth(s))
		}
		row := func(rating, title, author, date string) string {
			return cell(rating, ratingWidth) + " " + cell(title, titleWidth) + " " + cell(author, authorWidth) + " " + cell(date, dateWidth)
		}

		lines = append(lines, st.secondary.Render(row("Rating", "Title", "Author", "Date")))
		b.offset = btScroll(b.cursor, b.offset, b.height)
		for i := b.offset; i < len(b.results) && i < b.offset+b.height; i++ {
			g := b.results[i]
			r := row(btIdgamesRating(g.Rating), g.Title, g.Author, g.Date)
			if i == b.cursor {
				style := st.button
				if !b.inSearch && m.isFocused(b) {
					style = st.selected
				}
				lines = append(lines, style.Render(r))
				continue
			}
			lines = append(lines, st.primary.Render(r))
		}
	}

	for i := range lines {
		lines[i] = " " + lines[i]
	}
	return st.box("Files", strings.Join(lines, "\n"), width, height, st.borderColor)
}

func (b *btIdgames) detailsView(m *btModel, width, height int) string {
	st := &m.st
	if b.cursor >= len(b.results) || b.loading {
		return st.box("Details", "", width, height, st.borderColor)
	}

	g := b.results[b.cursor]
	text := g.Textfile
	if d, ok := b.details[g.Id]; ok {
		g, text = d, d.Textfile
	}
	if text == "" {
		text = fmt.Sprintf("Title: %v\nAuthor: %v\nDate: %v\nSize: %v\n\n%v", g.Title, g.Author, g.Date, g.Size, g.Description)
	}

	// stylize the text file a bit
	wrapped := lipgloss.NewStyle().Width(max(1, width-4)).Render(strings.ReplaceAll(text, "\r", ""))
	lines := strings.Split(wrapped, "\n")
	for i, line := range lines {
		switch {
		case strings.HasPrefix(strings.TrimSpace(line), "====="):
			line = st.moreContrast.Render(line)
		case btTextfileKey.MatchString(line):
			parts := btTextfileKey.FindStringSubmatch(line)
			line = st.moreContrast.Render(parts[1]+":") + st.primary.Render(parts[2])
		default:
			line = st.primary.Render(line)
		}
		lines[i] = " " + line
	}

	b.detailsOffset = max(0, min(b.detailsOffset, len(lines)-(height-2)))
	return st.box("Details (PGUP/PGDN)", strings.Join(lines[b.detailsOffset:], "\n"), width, height, st.borderColor)
}

// helper to make a string from the games rating
func btIdgamesRating(rating float32) string {
	r := max(0, min(5, int(rating)))
	return strings.Repeat("*", r) + strings.Repeat("-", 5-r)
}
