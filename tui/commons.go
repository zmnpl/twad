package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/rivo/tview"
	"github.com/zmnpl/twad/ports/savesStats"
)

// a tree with default properties
func newTree(rootDir string) (*tview.TreeView, *tview.TreeNode) {
	tree := tview.NewTreeView()
	tree.SetBorder(true)

	root := tview.NewTreeNode(rootDir).SetColor(tview.Styles.TitleColor)
	tree.SetRoot(root).
		SetCurrentNode(root)

	return tree, root
}

// A helper function which adds the files and directories of the given path
// to the given target node.
// Takes a filter function to filter files, which should not be in
func makeFileTreeAddFunc(fileFilter func(files []os.DirEntry, extensions string, includeDirs bool) []os.DirEntry, extensions string, includeDirs bool, hideUnixHidden bool) func(target *tview.TreeNode, path string) {
	return func(target *tview.TreeNode, path string) {
		files, err := os.ReadDir(path)
		if err != nil {
			// show why it can't be expanded
			target.AddChild(tview.NewTreeNode(err.Error()).
				SetColor(warnColor).
				SetSelectable(false))
			return
		}

		if fileFilter != nil {
			files = fileFilter(files, extensions, includeDirs)
		}

		sort.Slice(files, func(i, j int) bool {
			return strings.ToLower(files[i].Name()) < strings.ToLower(files[j].Name())
		})

		for _, file := range files {
			// hide hiden files
			if hideUnixHidden && strings.HasPrefix(file.Name(), ".") {
				continue
			}

			node := tview.NewTreeNode(file.Name()).
				SetReference(filepath.Join(path, file.Name())).
				SetSelectable(true)
			node.SetColor(tview.Styles.SecondaryTextColor)

			if file.IsDir() {
				node.SetColor(tview.Styles.PrimaryTextColor)
			}
			target.AddChild(node)
		}
	}
}

// a list with a title line above it, as used in the detail panes
func newTitledList(title string) (*tview.Flex, *tview.List) {
	list := tview.NewList()
	list.SetSecondaryTextColor(tview.Styles.TitleColor).SetSelectedFocusOnly(true)

	frame := tview.NewFlex().SetDirection(tview.FlexRow)
	frame.SetBorderPadding(0, 0, 1, 1)
	frame.AddItem(tview.NewTextView().
		SetText(title).
		SetTextColor(tview.Styles.SecondaryTextColor), 1, 0, false)
	frame.AddItem(list, 0, 1, true)

	return frame, list
}

// kvTable is a two column table of labeled values, filled row by row
type kvTable struct {
	*tview.Table
	row int
}

func newKVTable() *kvTable {
	t := tview.NewTable().
		SetFixed(1, 1).
		SetBorders(tableBorders).SetSeparator(':')
	t.SetBorderPadding(0, 0, 1, 1)
	return &kvTable{Table: t}
}

// section adds a highlighted heading row
func (t *kvTable) section(title string) {
	t.SetCell(t.row, 0, tview.NewTableCell(title).SetTextColor(tview.Styles.ContrastBackgroundColor))
	t.row++
}

// kv adds a row with label and value
// long labels get truncated, so they can't push the values out of sight
func (t *kvTable) kv(label string, value any) {
	t.SetCell(t.row, 0, tview.NewTableCell(label).SetTextColor(tview.Styles.SecondaryTextColor).SetMaxWidth(30))
	t.SetCell(t.row, 1, tview.NewTableCell(fmt.Sprint(value)).SetAlign(tview.AlignLeft))
	t.row++
}

// spacer adds an empty row, which keeps the columns from collapsing
func (t *kvTable) spacer() {
	t.kv("                    ", "                    ")
}

// progress adds kills, secrets and items of the given stats
func (t *kvTable) progress(s savesStats.MapStats) {
	t.kv("# Kills", fmt.Sprintf("%v/%v", s.KillCount, s.TotalKills))
	t.kv("# Secrets", fmt.Sprintf("%v/%v", s.SecretCount, s.TotalSecrets))
	t.kv("# Items", fmt.Sprintf("%v/%v", s.ItemCount, s.TotalItems))
}

// labelOk appends a message, colored as good, to a label
func labelOk(label, msg string) string {
	return label + colorTagGoodColor + " " + msg
}

// labelWarn appends a message, colored as warning, to a label
func labelWarn(label, msg string) string {
	return label + colorTagWarnColor + " " + msg
}
