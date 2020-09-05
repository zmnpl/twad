package tui

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/rivo/tview"
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
func makeFileTreeAddFunc(fileFilter func(files []os.FileInfo) []os.FileInfo, hideUnixHidden bool) func(target *tview.TreeNode, path string) {
	return func(target *tview.TreeNode, path string) {
		files, err := ioutil.ReadDir(path)
		if fileFilter != nil {
			files = fileFilter(files)
		}

		sort.Slice(files, func(i, j int) bool {
			return strings.ToLower(files[i].Name()) < strings.ToLower(files[j].Name())
		})

		if err != nil {
			panic(err)
		}

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
