package tui

import (
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/zmnpl/twad/base"
	"github.com/zmnpl/twad/helper"
	"github.com/zmnpl/twad/tui/shared"
)

func (u *ui) runZipImport(archivePath string, folderNameSuggest string, xOffset int, yOffset int, handFocusBackTo tview.Primitive) {
	// sets focus to the given primitive
	// if nil was given, then the apps default state will be restored
	resetFocus := func() {
		u.contentPages.close(pageZipImport, handFocusBackTo)
	}

	modNameInput := tview.NewInputField().SetLabel(shared.Dict.ZipImportToLabel).SetText(path.Base(archivePath))

	modNameInput.SetText(strings.TrimSuffix(path.Base(archivePath), path.Ext(archivePath)))
	if folderNameSuggest != "" {
		modNameInput.SetText(folderNameSuggest)
	}

	modNameDoneCheck := func() {
		suggestedName := modNameInput.GetText()
		if !helper.IsFileNameValid(suggestedName) {
			modNameInput.SetLabel(labelWarn(shared.Dict.ZipImportToLabel, shared.Dict.ZipImportToBadNameLabel))
			return
		}
		if _, err := os.Stat(path.Join(base.Config().WadDir, suggestedName)); !os.IsNotExist(err) {
			modNameInput.SetLabel(labelWarn(shared.Dict.ZipImportToLabel, shared.Dict.ZipImportToExistsLabel))
			return
		}
		modNameInput.SetLabel(shared.Dict.ZipImportToLabel)
	}

	modNameInput.SetDoneFunc(func(key tcell.Key) {
		modNameDoneCheck()
	})

	modNameForm := tview.NewForm().
		AddFormItem(modNameInput).
		AddButton(shared.Dict.ZipImportCancel, func() {
			resetFocus()
		}).
		AddButton(shared.Dict.ZipImportFormOk, func() {
			modName := modNameInput.GetText()
			resetFocus()

			// test file name again
			if !helper.IsFileNameValid(modName) {
				u.showError(shared.Dict.ZipImportNameInvalid, shared.Dict.ZipImportNameInvalidReasons, modNameInput, nil)
				return
			}

			// test if provided zip exists
			if _, err := os.Stat(archivePath); os.IsNotExist(err) {
				u.showError(shared.Dict.ZipImportArchiveNotFound, err.Error(), handFocusBackTo, nil)
				return
			}

			// START ACTUAL IMPORT
			if err := base.ImportArchive(archivePath, modName); err != nil {
				u.showError(shared.Dict.ZipImportFailed, err.Error(), handFocusBackTo, nil)
			}
		})

	modNameForm.SetFocus(0)

	layout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewTextView().SetText(" "+shared.Dict.ZipImportSecurityWarn).SetTextColor(tcell.ColorRed), 1, 0, true).
		AddItem(modNameForm, 5, 0, false)
	layout.SetBorder(true).
		SetTitle(shared.Dict.ZipImportTitle + " " + filepath.Base(archivePath))

	// dimensions
	height := 8
	width := len(shared.Dict.ZipImportSecurityWarn) + 4

	u.contentPages.showAt(pageZipImport, layout, xOffset, yOffset, width, height)
	u.app.SetFocus(modNameForm)
}
