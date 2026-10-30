package bubbletui

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/zmnpl/twad/helper"
)

type btTreeNode struct {
	name     string
	path     string
	dir      bool
	err      bool // shows why a folder couldn't be expanded; not selectable
	loaded   bool
	expanded bool
	depth    int
	children []*btTreeNode
}

// btTree browses the file system, to pick files with certain extensions
type btTree struct {
	title      string
	root       *btTreeNode
	extensions string
	hideHidden bool
	onFile     func(m *btModel, path string) tea.Cmd

	cursor int
	offset int
	height int
}

func newBTTree(title, rootDir, extensions string, hideHidden bool, onFile func(m *btModel, path string) tea.Cmd) *btTree {
	t := &btTree{
		title:      title,
		root:       &btTreeNode{name: rootDir, path: rootDir, dir: true},
		extensions: extensions,
		hideHidden: hideHidden,
		onFile:     onFile,
	}
	t.load(t.root)
	return t
}

// load adds the files and directories of the nodes path as children
func (t *btTree) load(n *btTreeNode) {
	n.loaded, n.expanded = true, true
	n.children = nil

	files, err := os.ReadDir(n.path)
	if err != nil {
		n.children = []*btTreeNode{{name: err.Error(), err: true, depth: n.depth + 1}}
		return
	}

	files = helper.FilterExtensions(files, t.extensions, true)
	sort.Slice(files, func(i, j int) bool {
		return strings.ToLower(files[i].Name()) < strings.ToLower(files[j].Name())
	})

	for _, file := range files {
		if t.hideHidden && strings.HasPrefix(file.Name(), ".") {
			continue
		}
		n.children = append(n.children, &btTreeNode{
			name:  file.Name(),
			path:  filepath.Join(n.path, file.Name()),
			dir:   file.IsDir(),
			depth: n.depth + 1,
		})
	}
}

// visible returns the nodes as shown, top to bottom
func (t *btTree) visible() []*btTreeNode {
	var nodes []*btTreeNode
	var walk func(n *btTreeNode)
	walk = func(n *btTreeNode) {
		nodes = append(nodes, n)
		if n.expanded {
			for _, c := range n.children {
				walk(c)
			}
		}
	}
	walk(t.root)
	return nodes
}

func (t *btTree) update(m *btModel, key tea.KeyMsg) tea.Cmd {
	nodes := t.visible()
	if c, ok := btMove(key.String(), t.cursor, len(nodes), t.height/2); ok {
		t.cursor = c
		return nil
	}

	n := nodes[t.cursor]
	switch key.String() {
	case "left":
		n.expanded = false
	case "enter", " ", "right":
		if n.err {
			return nil
		}
		if n.loaded {
			n.expanded = !n.expanded || key.String() == "right"
			return nil
		}
		// stat, so links to directories get followed
		fi, err := os.Stat(n.path)
		switch {
		case err != nil:
			return nil
		case fi.IsDir():
			t.load(n)
		case key.String() != "right":
			return t.onFile(m, n.path)
		}
	}
	return nil
}

func (t *btTree) view(m *btModel, width, height int) string {
	nodes := t.visible()
	t.cursor = min(t.cursor, len(nodes)-1)
	t.height = height - 2
	t.offset = btScroll(t.cursor, t.offset, t.height)
	focused := m.isFocused(t)

	lines := make([]string, 0, t.height)
	for i := t.offset; i < len(nodes) && i < t.offset+t.height; i++ {
		n := nodes[i]
		marker := "  "
		if n.dir {
			marker = "▸ "
			if n.expanded {
				marker = "▾ "
			}
		}
		text := btTruncate(strings.Repeat("  ", n.depth)+marker+n.name, width-3)

		style := m.st.secondary
		switch {
		case focused && i == t.cursor:
			style = m.st.selected
		case n.err:
			style = m.st.warn
		case n == t.root:
			style = m.st.title
		case n.dir:
			style = m.st.primary
		}
		lines = append(lines, " "+style.Render(text))
	}
	return m.st.box(t.title, strings.Join(lines, "\n"), width, height, m.st.borderColor)
}
