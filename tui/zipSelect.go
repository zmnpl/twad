package tui

import (
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/zmnpl/twad/helper"
)

type zipSelect struct {
	layout     *tview.Flex
	selectTree *tview.TreeView
}

func newZipImportUI() *zipSelect {
	var zui zipSelect
	zui.initZipSelect()

	zui.layout = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(zui.selectTree, 0, 1, true)
	return &zui
}

func (z *zipSelect) initZipSelect() {
	//rootDir := helper.Home() // TODO: Start from / but preselect /home/user
	rootDir := "/"
	if _, err := os.Stat(rootDir); err != nil {
		if os.IsNotExist(err) {
			// TODO
		}
	}

	var rootNode *tview.TreeNode
	z.selectTree, rootNode = newTree(rootDir)
	z.selectTree.SetTitle(dict.zipSelectTitle)
	add := makeFileTreeAddFunc(helper.FilterExtensions, ".zip.tar.gz.rar", true, true)
	add(rootNode, rootDir)

	z.selectTree.SetSelectedFunc(func(node *tview.TreeNode) {
		reference := node.GetReference()

		if reference == nil {
			return
		}
		children := node.GetChildren()
		if len(children) == 0 {
			// Load and show files in this directory.
			selPath := reference.(string)

			// check if path can at leas be read
			// otherwise return
			f, err := os.OpenFile(selPath, os.O_RDONLY, 0666)
			if err != nil && os.IsPermission(err) {
				return
			}
			defer f.Close()

			fi, err := os.Stat(selPath)
			switch {
			case err != nil:
				return // TODO: any form of info to user?
			case fi.IsDir():
				add(node, selPath)
			default:
				runZipImport(selPath, "", 0, 1, z.selectTree)
			}
		} else {
			node.SetExpanded(!node.IsExpanded())
		}
	})

	z.selectTree.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		k := event.Key()
		if k == tcell.KeyRune {
			switch event.Rune() {
			case 'q':
				app.Stop()
				return nil
			}
		}
		return event
	})
}

func (z *zipSelect) reset() {
	app.SetFocus(z.selectTree)
}
