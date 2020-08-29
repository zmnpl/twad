package tui

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"github.com/zmnpl/twad/cfg"
)

const (
	zipSelectTitle = "Archive to import"
)

// tree view for selecting additional mods TODO
func makeZipSelect() *tview.TreeView {
	rootDir := cfg.Home()
	if _, err := os.Stat(rootDir); err != nil {
		if os.IsNotExist(err) {
		}
		return nil
	}

	zipSelect, rootNode := newTree(rootDir)
	add := makeFileTreeAddFunc(filterKnownArchives)
	add(rootNode, rootDir)

	zipSelect.SetSelectedFunc(func(node *tview.TreeNode) {
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
				extractArchive()
				cfg.ImportArchive(path, "")
			}
		} else {
			node.SetExpanded(!node.IsExpanded())
		}
	})

	return zipSelect
}

func extractArchive() {

}

// help for navigation
func archiveImportFolderName() *tview.Flex {
	modNameInput := tview.NewInputField().SetLabel("").SetText("")
	modNameDoneCheck := func() {
		// does this path exist?
		// TODO: check if valid folder name
		// if not, dactivate ok button

		// TODO: check if exists already
		// if yes, warn with label
	}

	modNameInput.SetDoneFunc(func(key tcell.Key) {
		modNameDoneCheck()
	})

	archiveImportForm := tview.NewForm().
		AddFormItem(modNameInput).
		AddButton("ok", func() {
		})

	archiveImportForm.
		SetBorder(true).
		SetTitle("")
	archiveImportForm.SetFocus(1)

	youSureLayout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false).
		AddItem(archiveImportForm, 10, 0, true).
		AddItem(nil, 0, 1, false)

	return youSureLayout
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
