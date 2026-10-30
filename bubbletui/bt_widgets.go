package bubbletui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/x/ansi"
	"github.com/zmnpl/twad/ports/savesStats"
)

// btFit cuts or pads s to exactly the given size
func btFit(s string, width, height int) string {
	if width <= 0 || height <= 0 {
		return ""
	}
	lines := strings.Split(s, "\n")
	if len(lines) > height {
		lines = lines[:height]
	}
	for len(lines) < height {
		lines = append(lines, "")
	}
	for i, l := range lines {
		l = ansi.Truncate(l, width, "")
		if w := ansi.StringWidth(l); w < width {
			l += strings.Repeat(" ", width-w)
		}
		lines[i] = l
	}
	return strings.Join(lines, "\n")
}

// btTruncate shortens a single line to the given width, marking the cut with an ellipsis
func btTruncate(s string, width int) string {
	if width <= 0 {
		return ""
	}
	return ansi.Truncate(s, width, "…")
}

// btOverlay draws fg on top of bg at x/y
func btOverlay(bg, fg string, x, y int) string {
	bgLines := strings.Split(bg, "\n")
	for i, fl := range strings.Split(fg, "\n") {
		row := y + i
		if row < 0 || row >= len(bgLines) {
			continue
		}
		b := bgLines[row]
		left := ansi.Truncate(b, x, "")
		if w := ansi.StringWidth(left); w < x {
			left += strings.Repeat(" ", x-w)
		}
		right := ansi.Cut(b, x+ansi.StringWidth(fl), ansi.StringWidth(b))
		bgLines[row] = left + "\x1b[0m" + fl + "\x1b[0m" + right
	}
	return strings.Join(bgLines, "\n")
}

// btScroll keeps a cursor within the visible window of a scrolled list
func btScroll(cursor, offset, visible int) int {
	if visible < 1 {
		return cursor
	}
	if cursor < offset {
		return cursor
	}
	if cursor >= offset+visible {
		return cursor - visible + 1
	}
	return offset
}

// btMove applies navigation keys to a cursor over n items
// returns the new cursor and if the key was a navigation key
func btMove(key string, cursor, n, page int) (int, bool) {
	switch key {
	case "up", "k":
		cursor--
	case "down", "j":
		cursor++
	case "home", "g":
		cursor = 0
	case "end", "G":
		cursor = n - 1
	case "pgup":
		cursor -= max(1, page)
	case "pgdown":
		cursor += max(1, page)
	default:
		return cursor, false
	}
	return max(0, min(cursor, n-1)), true
}

// list with two lines per item, as used in the detail panes
//#######################################################################

type btListItem struct {
	main      string
	secondary string
	warn      bool
}

type btList struct {
	title  string
	items  []btListItem
	cursor int
	offset int
	height int // visible rows of the last render, for paging
}

// move handles navigation keys and reports if the cursor moved
func (l *btList) move(key string) (moved, handled bool) {
	before := l.cursor
	l.cursor, handled = btMove(key, l.cursor, len(l.items), l.height/2)
	return l.cursor != before, handled
}

func (l *btList) remove(i int) {
	l.items = append(l.items[:i], l.items[i+1:]...)
	l.cursor = max(0, min(l.cursor, len(l.items)-1))
}

func (l *btList) view(st *btStyles, width, height int, focused, edit bool) string {
	inner := width - 2 // padding left and right
	l.height = height - 1
	visible := max(1, l.height/2)
	l.offset = btScroll(l.cursor, l.offset, visible)

	lines := []string{" " + st.secondary.Render(btTruncate(l.title, inner))}
	for i := l.offset; i < len(l.items) && i < l.offset+visible; i++ {
		it := l.items[i]
		main := btTruncate(it.main, inner)
		sec := btTruncate(it.secondary, inner)
		switch {
		case focused && i == l.cursor:
			sel := st.selected
			if edit {
				sel = st.selectedEdit
			}
			main = sel.Render(main + strings.Repeat(" ", max(0, inner-ansi.StringWidth(main))))
		case it.warn:
			main = st.warn.Render(main)
		default:
			main = st.primary.Render(main)
		}
		lines = append(lines, " "+main, " "+st.title.Render(sec))
	}
	return btFit(strings.Join(lines, "\n"), width, height)
}

// two column table of labeled values, filled row by row
//#######################################################################

type btKVRow struct {
	section bool
	label   string
	value   string
}

type btKV struct {
	rows []btKVRow
}

// section adds a highlighted heading row
func (t *btKV) section(title string) {
	t.rows = append(t.rows, btKVRow{section: true, label: title})
}

func (t *btKV) kv(label string, value any) {
	t.rows = append(t.rows, btKVRow{label: label, value: fmt.Sprint(value)})
}

func (t *btKV) spacer() {
	t.rows = append(t.rows, btKVRow{})
}

// progress adds kills, secrets and items of the given stats
func (t *btKV) progress(s savesStats.MapStats) {
	t.kv("# Kills", fmt.Sprintf("%v/%v", s.KillCount, s.TotalKills))
	t.kv("# Secrets", fmt.Sprintf("%v/%v", s.SecretCount, s.TotalSecrets))
	t.kv("# Items", fmt.Sprintf("%v/%v", s.ItemCount, s.TotalItems))
}

// lines renders the table; long labels get truncated, so they can't push the values out of sight
func (t *btKV) lines(st *btStyles, width int) []string {
	labelWidth := 0
	for _, r := range t.rows {
		if !r.section {
			labelWidth = max(labelWidth, min(30, ansi.StringWidth(r.label)))
		}
	}

	lines := make([]string, 0, len(t.rows))
	for _, r := range t.rows {
		var line string
		switch {
		case r.section:
			line = st.contrast.Render(r.label)
		case r.label == "" && r.value == "":
			line = ""
		default:
			label := btTruncate(r.label, labelWidth)
			label += strings.Repeat(" ", labelWidth-ansi.StringWidth(label))
			line = st.secondary.Render(label) + st.title.Render(" : ") + st.primary.Render(r.value)
		}
		lines = append(lines, " "+btTruncate(line, width-2))
	}
	return lines
}
