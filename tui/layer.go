package tui

import (
	"github.com/rivo/tview"
)

// layer is a pages container with one base page, on top of which
// transient pages (dialogs, detail views, ...) get opened and closed again
type layer struct {
	*tview.Pages
	base      string
	transient map[string]bool
}

func newLayer(base string) *layer {
	return &layer{
		Pages:     tview.NewPages(),
		base:      base,
		transient: make(map[string]bool),
	}
}

// setBase sets or replaces the primitive of the base page
func (l *layer) setBase(p tview.Primitive) {
	l.AddPage(l.base, p, true, true)
}

// add puts p on top as transient page without focusing it
func (l *layer) add(name string, p tview.Primitive) {
	l.transient[name] = true
	l.AddPage(name, p, true, true)
}

// show puts p on top as transient page and focuses it
func (l *layer) show(name string, p tview.Primitive) {
	l.add(name, p)
	app.SetFocus(p)
}

// showAt shows p as dialog of the given size at x, right below row y
// if it doesn't fit there, it gets placed above
func (l *layer) showAt(name string, p tview.Primitive, x, y, width, height int) {
	l.show(name, modalAt(p, x, anchoredY(y, height, l.Box), width, height))
}

// switchTo shows p as transient page, hides all other pages and focuses p
func (l *layer) switchTo(name string, p tview.Primitive) {
	l.show(name, p)
	l.SwitchToPage(name)
}

// open switches to an already existing page and focuses the given primitive
func (l *layer) open(name string, focus tview.Primitive) {
	l.SwitchToPage(name)
	app.SetFocus(focus)
}

// close removes a transient page and hands focus to the given primitive
// if nil was given, the apps default state will be restored
func (l *layer) close(name string, focus tview.Primitive) {
	delete(l.transient, name)
	l.RemovePage(name)
	// closed page might have hidden everything else
	if front, _ := l.GetFrontPage(); front == "" {
		l.ShowPage(l.base)
	}

	if focus == nil {
		appModeNormal()
		return
	}
	app.SetFocus(focus)
}

// reset removes all transient pages and shows the base page
func (l *layer) reset() {
	for name := range l.transient {
		l.RemovePage(name)
	}
	clear(l.transient)
	l.SwitchToPage(l.base)
}
