package bubbletui

import (
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"
	"github.com/zmnpl/twad/base"
)

func TestBTTagsKeepsNonColorBrackets(t *testing.T) {
	m := &btModel{cfg: &base.Cfg{}}
	m.applyTheme()

	got := ansi.Strip(m.st.tags("[orange]twad[white] - [.ini] [notacolor]"))
	if want := "twad - [.ini] [notacolor]"; got != want {
		t.Errorf("tags: got %q, want %q", got, want)
	}
}

func TestBTFit(t *testing.T) {
	got := btFit("abcdef\nx\n1\n2", 3, 2)
	if want := "abc\nx  "; got != want {
		t.Errorf("fit: got %q, want %q", got, want)
	}
}

func TestBTOverlay(t *testing.T) {
	bg := strings.Join([]string{"..........", "..........", ".........."}, "\n")
	got := ansi.Strip(btOverlay(bg, "ab\ncd", 3, 1))
	want := strings.Join([]string{"..........", "...ab.....", "...cd....."}, "\n")
	if got != want {
		t.Errorf("overlay: got\n%s\nwant\n%s", got, want)
	}
}

func TestBTMove(t *testing.T) {
	cases := []struct {
		key           string
		cursor, n     int
		want          int
		wantNavigates bool
	}{
		{"down", 0, 3, 1, true},
		{"down", 2, 3, 2, true},
		{"up", 0, 3, 0, true},
		{"end", 0, 3, 2, true},
		{"pgdown", 0, 30, 10, true},
		{"x", 1, 3, 1, false},
	}
	for _, c := range cases {
		got, ok := btMove(c.key, c.cursor, c.n, 10)
		if got != c.want || ok != c.wantNavigates {
			t.Errorf("btMove(%q, %v, %v) = %v, %v; want %v, %v", c.key, c.cursor, c.n, got, ok, c.want, c.wantNavigates)
		}
	}
}
