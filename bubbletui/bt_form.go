package bubbletui

import (
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

type btFieldKind int

const (
	btInput btFieldKind = iota
	btCheckbox
	btDropDown
)

// btField is one item of a form
type btField struct {
	kind  btFieldKind
	label func(st *btStyles) string

	// validation, shown next to the label; runs initially and whenever the field is left
	check  func(value string) (msg string, ok bool)
	note   string
	noteOk bool

	done    func(f *btField) // when the field is left
	changed func()           // value changed

	// input
	input      textinput.Model
	accept     func(text string) bool // refuses edits, which make the text invalid
	suggest    func(text string) []string
	fieldWidth int // 0 takes the available width

	// checkbox
	checked bool

	// drop down
	options []string
	current int
}

func btPlain(label string) func(st *btStyles) string {
	return func(st *btStyles) string { return st.primary.Render(label) }
}

func (m *btModel) newInput(label, value string) *btField {
	ti := textinput.New()
	ti.Prompt = ""
	ti.SetValue(value)
	ti.Cursor.SetMode(cursor.CursorStatic)
	ti.TextStyle = m.st.field
	ti.PlaceholderStyle = m.st.field
	ti.CompletionStyle = m.st.field.Foreground(btColorOf(m.st.title))
	ti.Cursor.TextStyle = m.st.field
	ti.Cursor.Style = m.st.buttonFocus
	return &btField{kind: btInput, label: btPlain(label), input: ti}
}

func (m *btModel) newCheckbox(label string, checked bool) *btField {
	return &btField{kind: btCheckbox, label: btPlain(label), checked: checked}
}

func (m *btModel) newDropDown(label string, options []string, current int) *btField {
	return &btField{kind: btDropDown, label: btPlain(label), options: options, current: current}
}

func btColorOf(s lipgloss.Style) lipgloss.TerminalColor {
	return s.GetForeground()
}

// value returns the text of an input or the selected option of a drop down
func (f *btField) value() string {
	switch f.kind {
	case btInput:
		return f.input.Value()
	case btDropDown:
		if f.current < len(f.options) {
			return f.options[f.current]
		}
	}
	return ""
}

func (f *btField) runCheck() {
	if f.check != nil {
		f.note, f.noteOk = f.check(f.value())
	}
}

func (f *btField) refreshSuggestions() {
	if f.suggest != nil {
		f.input.ShowSuggestions = true
		f.input.SetSuggestions(f.suggest(f.input.Value()))
	}
}

func (f *btField) blur() {
	f.input.Blur()
	if f.done != nil {
		f.done(f)
	}
	f.runCheck()
}

func (f *btField) update(key tea.KeyMsg) tea.Cmd {
	k := key.String()
	switch f.kind {
	case btCheckbox:
		if k == " " || k == "x" {
			f.checked = !f.checked
			f.onChange()
		}
	case btDropDown:
		switch k {
		case "left", "h":
			f.current = (f.current - 1 + len(f.options)) % max(1, len(f.options))
			f.onChange()
		case "right", "l", " ":
			f.current = (f.current + 1) % max(1, len(f.options))
			f.onChange()
		}
	case btInput:
		ti := &f.input
		// right at the end of the text takes the suggestion
		if k == "right" && ti.Position() == len([]rune(ti.Value())) {
			if s := ti.CurrentSuggestion(); s != "" && s != ti.Value() {
				ti.SetValue(s)
				ti.CursorEnd()
				f.refreshSuggestions()
				return nil
			}
		}
		before := ti.Value()
		var cmd tea.Cmd
		*ti, cmd = ti.Update(key)
		if f.accept != nil && !f.accept(ti.Value()) {
			ti.SetValue(before)
		}
		if ti.Value() != before {
			f.refreshSuggestions()
			f.onChange()
		}
		return cmd
	}
	return nil
}

func (f *btField) onChange() {
	if f.changed != nil {
		f.changed()
	}
}

type btButton struct {
	label string
	run   func(m *btModel) tea.Cmd
}

// btForm is a bordered list of fields with buttons below
type btForm struct {
	title   string
	text    string // shown above the fields, already styled
	fields  []*btField
	buttons []btButton
	footer  string // shown below the buttons
	warn    bool   // styled as error
	width   int    // maximum outer width; 0 takes the available width

	focus int // index over fields, then buttons
}

// forms never quit on q
func (f *btForm) typing() bool { return true }

func (f *btForm) focusedField() *btField {
	if f.focus < len(f.fields) {
		return f.fields[f.focus]
	}
	return nil
}

func (f *btForm) setFocus(i int) {
	n := len(f.fields) + len(f.buttons)
	if n == 0 {
		return
	}
	if old := f.focusedField(); old != nil {
		old.blur()
	}
	f.focus = (i%n + n) % n
	if field := f.focusedField(); field != nil && field.kind == btInput {
		field.input.Focus()
		field.refreshSuggestions()
	}
}

// init runs all checks; called once, when the form is set up
func (f *btForm) init(focus int) *btForm {
	for _, field := range f.fields {
		field.runCheck()
	}
	f.setFocus(focus)
	return f
}

func (f *btForm) update(m *btModel, key tea.KeyMsg) tea.Cmd {
	field := f.focusedField()
	switch key.String() {
	case "tab", "down":
		f.setFocus(f.focus + 1)
		return nil
	case "shift+tab", "up":
		f.setFocus(f.focus - 1)
		return nil
	case "enter":
		if field == nil {
			return f.buttons[f.focus-len(f.fields)].run(m)
		}
		f.setFocus(f.focus + 1)
		return nil
	}

	if field == nil {
		switch key.String() {
		case "left", "h":
			f.setFocus(f.focus - 1)
		case "right", "l":
			f.setFocus(f.focus + 1)
		}
		return nil
	}
	return field.update(key)
}

func (f *btForm) view(m *btModel, width, height int) string {
	st := &m.st
	if f.width > 0 {
		width = min(width, f.width)
	}
	inner := width - 4 // border and padding

	wrap := func(s string) []string {
		return strings.Split(lipgloss.NewStyle().Width(max(1, inner)).Render(s), "\n")
	}

	var lines []string
	if f.text != "" {
		lines = append(lines, wrap(f.text)...)
		lines = append(lines, "")
	}

	// labels including validation notes
	labels := make([]string, len(f.fields))
	labelWidth := 0
	for i, field := range f.fields {
		labels[i] = field.label(st)
		if field.note != "" {
			noteStyle := st.warn
			if field.noteOk {
				noteStyle = st.good
			}
			labels[i] += " " + noteStyle.Render(field.note)
		}
		labelWidth = max(labelWidth, ansi.StringWidth(labels[i]))
	}
	labelWidth = min(labelWidth, inner*2/3)

	for i, field := range f.fields {
		marker := "  "
		if i == f.focus {
			marker = st.secondary.Render("› ")
		}
		label := btTruncate(labels[i], labelWidth)
		label += strings.Repeat(" ", labelWidth-ansi.StringWidth(label))

		fieldWidth := max(3, inner-labelWidth-3)
		if field.fieldWidth > 0 {
			fieldWidth = min(fieldWidth, field.fieldWidth)
		}
		lines = append(lines, marker+label+" "+f.fieldView(st, field, i == f.focus, fieldWidth))
	}

	if len(f.buttons) > 0 {
		buttons := make([]string, len(f.buttons))
		for i, b := range f.buttons {
			style := st.button
			if f.warn {
				style = st.buttonWarn
			}
			if i == f.focus-len(f.fields) {
				style = st.buttonFocus
			}
			buttons[i] = style.Render(" " + b.label + " ")
		}
		if len(f.fields) > 0 {
			lines = append(lines, "")
		}
		lines = append(lines, "  "+strings.Join(buttons, "  "))
	}

	if f.footer != "" {
		lines = append(lines, "")
		lines = append(lines, wrap(st.primary.Render(f.footer))...)
	}

	for i := range lines {
		lines[i] = " " + lines[i]
	}
	if height <= 0 {
		height = len(lines) + 2
	}

	borderColor := st.borderColor
	if f.warn {
		borderColor = st.warnColor
	}
	return st.box(f.title, strings.Join(lines, "\n"), width, height, borderColor)
}

func (f *btForm) fieldView(st *btStyles, field *btField, focused bool, width int) string {
	switch field.kind {
	case btCheckbox:
		mark := " "
		if field.checked {
			mark = "X"
		}
		return st.field.Render(" " + mark + " ")

	case btDropDown:
		text := " " + field.value()
		if focused {
			text = "◀" + field.value() + "▶"
		}
		text = btTruncate(text, width)
		return st.field.Render(text + strings.Repeat(" ", width-ansi.StringWidth(text)))

	default:
		// leading space, like the other field kinds
		field.input.Width = width - 2
		v := st.field.Render(" ") + field.input.View()
		if pad := width - ansi.StringWidth(v); pad > 0 {
			v += st.field.Render(strings.Repeat(" ", pad))
		}
		return v
	}
}
