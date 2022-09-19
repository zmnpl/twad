package tui

import (
	"os"
	"path"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/zmnpl/twad/base"
	"github.com/zmnpl/twad/helper"
)

func runZipImport(archivePath string, handFocusBackTo tview.Primitive) {
	pageZipImport := "zipimport"

	// sets focus to the given primitive
	// if nil was given, then the apps default state will be restored
	resetFocus := func() {
		contentPages.RemovePage(pageZipImport)
		if handFocusBackTo != nil {
			app.SetFocus(handFocusBackTo)
		} else {
			appModeNormal()
		}
	}

	importSecurityWarning := tview.NewTextView().SetText(dict.zipImportSecurityWarn).SetTextColor(tcell.ColorRed)

	modNameInput := tview.NewInputField().SetLabel(dict.zipImportToLabel).SetText(path.Base(archivePath))
	modNameInput.SetText(strings.TrimSuffix(path.Base(archivePath), path.Ext(archivePath)))

	modNameDoneCheck := func() {
		suggestedName := modNameInput.GetText()
		if !helper.IsFileNameValid(suggestedName) {
			modNameInput.SetLabel(dict.zipImportToLabel + warnColor + " " + dict.zipImportToBadNameLabel)
			return
		}
		if _, err := os.Stat(path.Join(base.Config().WadDir, suggestedName)); !os.IsNotExist(err) {
			modNameInput.SetLabel(dict.zipImportToLabel + warnColor + " " + dict.zipImportToExistsLabel)
			return
		}
		modNameInput.SetLabel(dict.zipImportToLabel)
	}

	modNameInput.SetDoneFunc(func(key tcell.Key) {
		modNameDoneCheck()
	})

	modNameForm := tview.NewForm().
		AddFormItem(modNameInput).
		AddButton(dict.zipImportFormOk, func() {
			modName := modNameInput.GetText()

			// test file name again
			if !helper.IsFileNameValid(modName) {
				showError("Cannot use that name", "Possible reasons:\n- File name contains forbidden characters\n- No permission to write this file/folder", modNameInput, nil)
				return
			}

			// test if provided zip exists
			if _, err := os.Stat(archivePath); os.IsNotExist(err) {
				showError("Mod archive not found", err.Error(), zipInput.selectTree, nil)
				zipInput.reset()
				return
			}

			// START ACTUAL IMPORT
			if err := base.ImportArchive(archivePath, modName); err != nil {
				showError("Could not import zip", err.Error(), zipInput.selectTree, nil)
			}
			resetFocus()
		}).
		AddButton(dict.zipImportCancel, func() {
			resetFocus()
		})

	modNameForm.
		SetBorder(true).
		SetTitle(dict.zipImportFormTitle)

	modNameForm.SetFocus(0)

	layout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(importSecurityWarning, 1, 0, true).
		AddItem(modNameForm, 7, 0, false)

	contentPages.AddPage(pageZipImport, layout, true, true)
	app.SetFocus(modNameForm)
}
