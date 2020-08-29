package tui

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/rivo/tview"
	"github.com/zmnpl/twad/games"
)

const (
	modTreeTitle = "Add new mod to game"
)

// tree view for selecting additional mods TODO
func makeModTree(g *games.Game) *tview.TreeView {
	rootDir := config.WadDir
	if _, err := os.Stat(rootDir); err != nil {
		if os.IsNotExist(err) {
		}
		return nil
	}

	modFolderTree, rootNode := newTree(rootDir)
	add := makeFileTreeAddFunc(filterExtensions)
	add(rootNode, rootDir)

	modFolderTree.SetSelectedFunc(func(node *tview.TreeNode) {
		reference := node.GetReference()

		if reference == nil {
			return
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
				add(node, path)
			default:
				g.AddMod(strings.TrimPrefix(path, config.WadDir+"/"))
				selectedGameChanged(g)
			}
		} else {
			node.SetExpanded(!node.IsExpanded())
		}
	})

	return modFolderTree
}

func filterExtensions(files []os.FileInfo) []os.FileInfo {
	n := 0
	for _, f := range files {
		ext := strings.ToLower(filepath.Ext(f.Name()))
		if _, found := config.ModExtensions[ext]; found || f.IsDir() {
			files[n] = f
			n++
		}
	}
	files = files[:n]
	return files
}
