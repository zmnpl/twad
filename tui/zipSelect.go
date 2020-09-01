package tui

import (
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"github.com/zmnpl/twad/cfg"
	"github.com/zmnpl/twad/helper"
)

// TODO: check if doomwaddir in config is something sane to avoid critical damage
// such as not = "/"

const (
	zipSelectTitle          = "Archive to import"
	zipImportToLabel        = "Folder name"
	zipImportToExistsLabel  = " (exists already)"
	zipImportToBadNameLabel = " (cannot use that name)"
	zipImportFormTitle      = "Mod folder name"
	zipImportFormOk         = "Import"
	zipImportCancel         = "Back to selection"
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
	rootDir := helper.Home() // TODO: Start from / but preselect /home/user
	if _, err := os.Stat(rootDir); err != nil {
		if os.IsNotExist(err) {
			// TODO
		}
	}

	var rootNode *tview.TreeNode
	z.selectTree, rootNode = newTree(rootDir)
	add := makeFileTreeAddFunc(filterKnownArchives)
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

			fi, err := os.Stat(selPath)
			switch {
			case err != nil:
				// handle the error and return
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
}

func (z *zipImportUI) initZipImportForm(archivePath string) {
	z.modNameInput = tview.NewInputField().SetLabel(zipImportToLabel).SetText(path.Base(archivePath))
	if archivePath == "" {
		z.modNameInput.SetText("")
	}

	modNameDoneCheck := func() {
		suggestedName := z.modNameInput.GetText()
		if !helper.IsFileNameValid(suggestedName) {
			//z.modNameInput.SetLabel(zipImportToLabel + warnColor + zipImportToBadNameLabel)

			showError("Cannot use that name", "Possible reasons:\n- File name contains forbidden characters\n- No permission to write this file/folder", z.modNameInput, nil)

			//app.SetFocus(zipInput.selectTree)
			// TODO: deactivate ok button
			return
		}
		if _, err := os.Stat(path.Join(cfg.Instance().WadDir, suggestedName)); !os.IsNotExist(err) {
			z.modNameInput.SetLabel(zipImportToLabel + warnColor + zipImportToExistsLabel)
			return
		}
		z.modNameInput.SetLabel(zipImportToLabel)
	}

	z.modNameInput.SetDoneFunc(func(key tcell.Key) {
		modNameDoneCheck()
	})

	// TODO: do this manually instead of with form
	// otherwise the error display cannot well be focused

	z.modNameForm = tview.NewForm().
		AddFormItem(z.modNameInput).
		AddButton(zipImportFormOk, func() {
			if _, err := os.Stat(z.zipPath); os.IsNotExist(err) {
				showError("Mod archive not found", err.Error(), zipInput.selectTree, nil)
				zipInput.reset()
				return
			}
			z.modName = z.modNameInput.GetText()

			// START ACTUAL IMPORT
			cfg.ImportArchive(z.zipPath, z.modName)

			z.reset()
		}).
		AddButton(zipImportCancel, func() {
			z.reset()
			showError("Cannot use that name", "Possible reasons:\n- File name contains forbidden characters\n- No permission to write this file/folder", z.selectTree, nil)
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

func filterKnownArchives(files []os.FileInfo) []os.FileInfo {
	knownArchives := map[string]int{".zip": 1}
	n := 0
	for _, f := range files {
		ext := strings.ToLower(filepath.Ext(f.Name()))
		if _, found := knownArchives[ext]; found || f.IsDir() {
			files[n] = f
			n++
		}
	}
	files = files[:n]
	return files
}
