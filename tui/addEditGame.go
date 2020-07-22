package tui

import (
	"strings"

	"github.com/rivo/tview"
	"github.com/zmnpl/twad/cfg"
	"github.com/zmnpl/twad/games"
)

const (
	addGame  = "Add New Game"
	editGame = "Edit Game"

	aeName       = "Name"
	aeSourcePort = "Source Port"
	aeIWAD       = "IWAD"

	aeEnvironment       = "Environment Variables"
	aeEnvironmentDetail = `Provide environment variables here; To turn VSync off entirely for example:
"vblank_mode=1"`
	aeOtherParams       = "Others"
	aeOtherParamsDetail = "Other parameters you want to pass to your ZDoom port for this game"
)

func splitParams(params string) []string {
	result := strings.Split(params, " ")
	for i := range result {
		result[i] = strings.TrimSpace(result[i])
	}
	return result
}

func indexOfItemIn(item string, list []string) (int, bool) {
	for i, val := range list {
		if val == item {
			return i, true
		}
	}
	return -1, false
}

func makeAddEditGame(g *games.Game) *tview.Flex {
	gWasNil := false
	title := editGame
	if g == nil {
		foo := games.NewGame("", "", "")
		g = &foo
		title = addGame
		gWasNil = true
	}

	inputName := tview.NewInputField().SetText(g.Name)

	inputSourcePort := tview.NewDropDown().SetFieldWidth(20).SetOptions([]string{"NA"}, nil)
	if cfg.GetInstance().SourcePorts != nil && len(cfg.GetInstance().SourcePorts) > 0 {
		inputSourcePort.SetOptions(cfg.GetInstance().SourcePorts, nil)
		if i, isIn := indexOfItemIn(g.SourcePort, cfg.GetInstance().SourcePorts); isIn {
			inputSourcePort.SetCurrentOption(i)
		} else {
			inputSourcePort.SetCurrentOption(0)
		}
	}

	inputIwad := tview.NewDropDown().SetFieldWidth(20).SetOptions([]string{"NA"}, nil)
	if cfg.GetInstance().IWADs != nil && len(cfg.GetInstance().IWADs) > 0 {
		inputIwad.SetOptions(cfg.GetInstance().IWADs, nil)
		if i, isIn := indexOfItemIn(g.Iwad, cfg.GetInstance().IWADs); isIn {
			inputIwad.SetCurrentOption(i)
		} else {
			inputIwad.SetCurrentOption(0)
		}
	}

	inputEnvVars := tview.NewInputField().SetText(g.EnvironmentString())

	inputCustomParams := tview.NewInputField().SetText(g.ParamsString())

	okButton := tview.NewButton(optionsOkButtonLabel).SetBackgroundColorActivated(tview.Styles.PrimaryTextColor).SetLabelColorActivated(tview.Styles.ContrastBackgroundColor)
	okButton.SetSelectedFunc(func() {
		g.Name = strings.TrimSpace(inputName.GetText())
		_, g.SourcePort = inputSourcePort.GetCurrentOption()
		_, g.Iwad = inputIwad.GetCurrentOption()
		g.Environment = splitParams(inputEnvVars.GetText())
		g.Parameters = splitParams(inputCustomParams.GetText())

		if gWasNil {
			games.AddGame(*g)
		}

		games.Persist()
		games.InformChangeListeners()
		appModeNormal()
	})

	// build form
	addEditGameForm := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewTextView().SetText(aeName).SetTextColor(tview.Styles.SecondaryTextColor), 1, 0, false).
		AddItem(inputName, 1, 0, true).
		AddItem(nil, 1, 0, false).
		AddItem(tview.NewTextView().SetText(aeSourcePort).SetTextColor(tview.Styles.SecondaryTextColor), 1, 0, false).
		AddItem(inputSourcePort, 1, 0, false).
		AddItem(nil, 1, 0, false).
		AddItem(tview.NewTextView().SetText(aeIWAD).SetTextColor(tview.Styles.SecondaryTextColor), 1, 0, false).
		AddItem(inputIwad, 1, 0, false).
		AddItem(nil, 1, 0, false).
		AddItem(tview.NewTextView().SetText(aeEnvironment).SetTextColor(tview.Styles.SecondaryTextColor), 1, 0, false).
		AddItem(tview.NewTextView().SetText(aeEnvironmentDetail), 2, 0, false).
		AddItem(inputEnvVars, 1, 0, false).
		AddItem(nil, 1, 0, false).
		AddItem(tview.NewTextView().SetText(aeOtherParams).SetTextColor(tview.Styles.SecondaryTextColor), 1, 0, false).
		AddItem(tview.NewTextView().SetText(aeOtherParamsDetail), 1, 0, false).
		AddItem(inputCustomParams, 1, 0, false).
		AddItem(nil, 1, 0, false).
		AddItem(nil, 1, 0, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).AddItem(okButton, 20, 0, false), 1, 0, false)
	addEditGameForm.SetBorder(true)
	addEditGameForm.SetTitle(title)
	addEditGameForm.SetBorderPadding(1, 1, 1, 1)

	// navigation path
	inputName.SetInputCapture(tabNavigate(okButton, inputSourcePort))
	inputSourcePort.SetInputCapture(tabNavigate(inputName, inputIwad))
	inputIwad.SetInputCapture(tabNavigate(inputSourcePort, inputEnvVars))
	inputEnvVars.SetInputCapture(tabNavigate(inputIwad, inputCustomParams))
	inputCustomParams.SetInputCapture(tabNavigate(inputEnvVars, okButton))
	okButton.SetInputCapture(tabNavigate(inputCustomParams, inputName))

	return addEditGameForm
}
