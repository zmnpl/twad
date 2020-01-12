package tui

import (
	"github.com/rivo/tview"
	"github.com/zmnpl/twad/games"
)

const (
	inputLabelName       = "Name"
	inputLabelSourcePort = "Source Port"
	inputLabelIWad       = "IWad"
	buttonAddLabel       = "Add"
	newGameFormTitle     = "Add new game"
)

func makeNewGameForm() *tview.Form {
	newForm := tview.NewForm().
		AddInputField(inputLabelName, "", 20, nil, nil).
		AddDropDown(inputLabelSourcePort, config.SourcePorts, 0, nil).
		AddDropDown(inputLabelIWad, config.IWADs, 0, nil).
		AddButton(buttonAddLabel, func() {
			name := newForm.GetFormItemByLabel(inputLabelName).(*tview.InputField).GetText()
			_, sourceport := newForm.GetFormItemByLabel(inputLabelSourcePort).(*tview.DropDown).GetCurrentOption()
			_, wad := newForm.GetFormItemByLabel(inputLabelIWad).(*tview.DropDown).GetCurrentOption()

			games.AddGame(games.NewGame(name, sourceport, wad))
		})

	newForm.SetBorder(true).SetTitle(newGameFormTitle).SetTitleAlign(tview.AlignCenter)

	return newForm
}
