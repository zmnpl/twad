package tui

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/rivo/tview"
	"github.com/zmnpl/twad/games"
)

const (
	modTreeTitle = "Add new mod to game"
)

// tree view for selecting additional mods TODO
func makeModTree(g *games.Game) *tview.TreeView {
	rootDir := config.ModBasePath
	root := tview.NewTreeNode(rootDir).SetColor(tview.Styles.TitleColor)
	modFolderTree := tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root)
	modFolderTree.SetBorder(true)
	modFolderTree.SetTitle(modTreeTitle)

	// A helper function which adds the files and directories of the given path
	// to the given target node.
	add := func(target *tview.TreeNode, path string) {
		files, err := ioutil.ReadDir(path)
		files = filterExtensions(files)

		sort.Slice(files, func(i, j int) bool {
			return strings.ToLower(files[i].Name()) < strings.ToLower(files[j].Name())
		})

		if err != nil {
			panic(err)
		}
		for _, file := range files {
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

	// Add the current directory to the root node.
	add(root, rootDir)

	// If a directory was selected, open it.
	modFolderTree.SetSelectedFunc(func(node *tview.TreeNode) {
		reference := node.GetReference()

		if reference == nil {
			return // Selecting the root node does nothing.
		}
		children := node.GetChildren()
		if len(children) == 0 {
			// Load and show files in this directory.
			path := reference.(string)

			fi, err := os.Stat(path)
			switch {
			case err != nil:
				// handle the error and return
			case fi.IsDir():
				// it's a directory
				add(node, path)
			default:
				// it's not a directory
				g.AddMod(strings.TrimPrefix(path, config.ModBasePath+"/"))
			}
		} else {
			// Collapse if visible, expand if collapsed.
			node.SetExpanded(!node.IsExpanded())
		}
	})

	return modFolderTree
}

// helper functions

func filterExtensions(files []os.FileInfo) []os.FileInfo {
	tmp := files
	files = files[:0]
	for _, v := range tmp {
		ext := strings.ToLower(filepath.Ext(v.Name()))
		if _, found := config.ModExtensions[ext]; found || v.IsDir() {
			files = append(files, v)
		}
	}
	return files
}
