package bubbletui

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/gdamore/tcell/v2"
	"github.com/zmnpl/twad/tui/shared"
)

// btStyles are the lipgloss styles of the bubbletea ui
// they get derived from the shared theme, so both uis look alike
type btStyles struct {
	primary      lipgloss.Style
	secondary    lipgloss.Style
	tertiary     lipgloss.Style
	title        lipgloss.Style
	contrast     lipgloss.Style
	moreContrast lipgloss.Style
	warn         lipgloss.Style
	good         lipgloss.Style

	borderColor lipgloss.TerminalColor
	warnColor   lipgloss.TerminalColor

	field        lipgloss.Style // input fields, checkboxes, drop downs
	button       lipgloss.Style
	buttonWarn   lipgloss.Style // buttons of error dialogs
	buttonFocus  lipgloss.Style
	selected     lipgloss.Style // selected row of a focused list or table
	selectedEdit lipgloss.Style // selected row of a list in edit mode
}

// btColor translates a tcell color into a lipgloss color
func btColor(c tcell.Color) lipgloss.TerminalColor {
	switch {
	case !c.Valid():
		return lipgloss.NoColor{}
	case c.IsRGB():
		return lipgloss.Color(c.CSS())
	default:
		// palette color
		return lipgloss.Color(strconv.Itoa(int(c - tcell.ColorValid)))
	}
}

// applyTheme selects the theme like the tview ui does and derives the styles from it
func (m *btModel) applyTheme() {
	t := shared.SelectTheme(m.cfg.UseTerminalColors)
	m.theme = t

	fg := func(c tcell.Color) lipgloss.Style { return lipgloss.NewStyle().Foreground(btColor(c)) }

	m.st = btStyles{
		primary:      fg(t.PrimaryText),
		secondary:    fg(t.SecondaryText),
		tertiary:     fg(t.TertiaryText),
		title:        fg(t.Title),
		contrast:     fg(t.ContrastBackground),
		moreContrast: fg(t.MoreContrastBackground),
		warn:         fg(t.Warn),
		good:         fg(t.Good),

		borderColor: btColor(t.Border),
		warnColor:   btColor(t.Warn),

		field:        fg(t.PrimaryText).Background(btColor(t.ContrastBackground)),
		button:       fg(t.PrimaryText).Background(btColor(t.ContrastBackground)),
		buttonWarn:   fg(tcell.ColorWhite).Background(btColor(t.Warn)),
		buttonFocus:  fg(t.PrimitiveBackground).Background(btColor(t.PrimaryText)),
		selected:     fg(t.PrimitiveBackground).Background(btColor(t.PrimaryText)),
		selectedEdit: fg(t.PrimitiveBackground).Background(btColor(t.TertiaryText)),
	}

	// with padding, like the tview header
	logo := strings.Split(m.st.tags(shared.Logo(t.ColorName(t.Title), t.ColorName(t.SecondaryText))), "\n")
	for i := range logo {
		logo[i] = " " + logo[i]
	}
	m.logo = strings.Join(logo, "\n")
}

// box draws a rounded border with an optional title around content of the given outer size
func (st *btStyles) box(title, content string, width, height int, borderColor lipgloss.TerminalColor) string {
	if width < 2 || height < 2 {
		return ""
	}
	inner := btFit(content, width-2, height-2)
	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Render(inner)

	if title == "" {
		return box
	}

	// put the title into the top border
	lines := strings.Split(box, "\n")
	t := btTruncate(" "+title+" ", width-4)
	border := lipgloss.NewStyle().Foreground(borderColor)
	top := border.Render("╭─") + st.title.Render(t) + border.Render(strings.Repeat("─", max(0, width-3-lipgloss.Width(t)))+"╮")
	lines[0] = top
	return strings.Join(lines, "\n")
}

// tview color tags like [orange] or [#FFA500]
var btTagPattern = regexp.MustCompile(`\[([a-zA-Z]+|#[0-9a-fA-F]{6})\]`)

// tags renders text, that contains tview color tags, with the according colors
// brackets, which don't contain a known color, are kept as they are
func (st *btStyles) tags(s string) string {
	var b strings.Builder
	style := st.primary
	last := 0
	for _, loc := range btTagPattern.FindAllStringSubmatchIndex(s, -1) {
		c := tcell.GetColor(s[loc[2]:loc[3]])
		if !c.Valid() {
			continue
		}
		b.WriteString(btRenderLines(style, s[last:loc[0]]))
		style = lipgloss.NewStyle().Foreground(btColor(c))
		last = loc[1]
	}
	b.WriteString(btRenderLines(style, s[last:]))
	return b.String()
}

// btRenderLines renders each line on its own, so lipgloss doesn't pad them to a block
func btRenderLines(style lipgloss.Style, s string) string {
	lines := strings.Split(s, "\n")
	for i, l := range lines {
		if l != "" {
			lines[i] = style.Render(l)
		}
	}
	return strings.Join(lines, "\n")
}
