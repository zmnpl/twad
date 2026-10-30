package tui

import (
	"os"
	"path/filepath"

	"github.com/zmnpl/twad/helper"
	"github.com/zmnpl/twad/tui/shared"

	"github.com/rivo/tview"
	"github.com/zmnpl/twad/games"
)

// tree view for selecting additional mods
func (u *ui) makeModTree(g *games.Game) (*tview.TreeView, error) {
	rootDir := u.config.WadDir
	if _, err := os.Stat(rootDir); err != nil {
		return nil, err
	}

	modFolderTree, rootNode := newTree(rootDir)
	modFolderTree.SetTitle(shared.Dict.ModTreeTitle)
	add := makeFileTreeAddFunc(helper.FilterExtensions, u.config.ModExtensions, true, false)
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
				mod, err := filepath.Rel(u.config.WadDir, path)
				if err != nil {
					return
				}
				g.AddMod(filepath.ToSlash(mod))
				u.detailSidePagesSub2.close(pageModSelector, u.gamesTable)
				u.selectedGameChanged(g)
			}
		} else {
			node.SetExpanded(!node.IsExpanded())
		}
	})

	return modFolderTree, nil
}
