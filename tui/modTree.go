package tui

import (
	"os"
	"strings"

	"github.com/zmnpl/twad/helper"

	"github.com/rivo/tview"
	"github.com/zmnpl/twad/games"
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
	modFolderTree.SetTitle(modTreeTitle)
	add := makeFileTreeAddFunc(helper.FilterExtensions, config.ModExtensions, true, false)
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
				app.SetFocus(gamesTable)
			}
		} else {
			node.SetExpanded(!node.IsExpanded())
		}
	})

	return modFolderTree
}
