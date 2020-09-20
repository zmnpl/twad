package tui

import (
	"os"
	"path"
	"strings"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"github.com/zmnpl/twad/cfg"
	"github.com/zmnpl/twad/helper"
)

const (
	zipSelectTitle          = "Select archive"
	zipImportToLabel        = "Folder name"
	zipImportToExistsLabel  = "exists already"
	zipImportToBadNameLabel = "cannot use that name"
	zipImportFormTitle      = "Import to"
	zipImportFormOk         = "Import"
	zipImportCancel         = "Back"
)

type zipImportUI struct {
	selectTree   *tview.TreeView
	modNameInput *tview.InputField
	modNameForm  *tview.Form

	zipPath string
	modName string
}

func newZipImportUI() *zipImportUI {
	var zui zipImportUI
	zui.initZipSelect()
	zui.initZipImportForm("")
	return &zui
}

func (z *zipImportUI) initZipSelect() {
	//rootDir := helper.Home() // TODO: Start from / but preselect /home/user
	rootDir := "/"
	if _, err := os.Stat(rootDir); err != nil {
		if os.IsNotExist(err) {
			// TODO
		}
	}

	var rootNode *tview.TreeNode
	z.selectTree, rootNode = newTree(rootDir)
	z.selectTree.SetTitle(zipSelectTitle)
	add := makeFileTreeAddFunc(helper.FilterExtensions, ".zip", true, true)
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
			defer f.Close()
			if err != nil && os.IsPermission(err) {
				return
			}

			fi, err := os.Stat(selPath)
			switch {
			case err != nil:
				return // TODO: any form of info to user?
			case fi.IsDir():
				add(node, selPath)
			default:
				z.zipPath = selPath
				z.modNameInput.SetText(strings.TrimSuffix(path.Base(selPath), path.Ext(selPath)))
				app.SetFocus(z.modNameForm)
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

func (z *zipImportUI) initZipImportForm(archivePath string) {
	z.modNameInput = tview.NewInputField().SetLabel(zipImportToLabel).SetText(path.Base(archivePath))
	if archivePath == "" {
		z.modNameInput.SetText("")
	}

	modNameDoneCheck := func() {
		suggestedName := z.modNameInput.GetText()
		if !helper.IsFileNameValid(suggestedName) {
			z.modNameInput.SetLabel(zipImportToLabel + warnColor + " " + zipImportToBadNameLabel)
			return
		}
		if _, err := os.Stat(path.Join(cfg.Instance().WadDir, suggestedName)); !os.IsNotExist(err) {
			z.modNameInput.SetLabel(zipImportToLabel + warnColor + " " + zipImportToExistsLabel)
			return
		}
		z.modNameInput.SetLabel(zipImportToLabel)
	}

	z.modNameInput.SetDoneFunc(func(key tcell.Key) {
		modNameDoneCheck()
	})

	z.modNameForm = tview.NewForm().
		AddFormItem(z.modNameInput).
		AddButton(zipImportFormOk, func() {
			z.modName = z.modNameInput.GetText()

			// test file name again
			if !helper.IsFileNameValid(z.modName) {
				showError("Cannot use that name", "Possible reasons:\n- File name contains forbidden characters\n- No permission to write this file/folder", z.modNameInput, nil)
				return
			}

			// test if provided zip exists
			if _, err := os.Stat(z.zipPath); os.IsNotExist(err) {
				showError("Mod archive not found", err.Error(), zipInput.selectTree, nil)
				zipInput.reset()
				return
			}

			// START ACTUAL IMPORT
			if err := cfg.ImportArchive(z.zipPath, z.modName); err != nil {
				showError("Could not import zip", err.Error(), zipInput.selectTree, nil)
			}
			z.reset()
		}).
		AddButton(zipImportCancel, func() {
			z.reset()
		})

	z.modNameForm.
		SetBorder(true).
		SetTitle(zipImportFormTitle)
	z.modNameForm.SetFocus(0)
}

func (z *zipImportUI) reset() {
	z.modNameInput.SetText("").SetLabel(zipImportToLabel)
	z.modNameForm.SetFocus(0)
	z.modName = ""
	z.zipPath = ""
	app.SetFocus(z.selectTree)
}
