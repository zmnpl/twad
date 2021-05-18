package tui

import (
	"strings"

	"github.com/gdamore/tcell/v2"
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
	aeOwnCfg     = "Use Own Source Port CFG"
	aeSharedCfg  = "Use Shared Port CFG"
	aeLink       = "Mod URL"

	aeEnvironment       = "Environment Variables *"
	aeEnvironmentDetail = `* Provide environment variables here; To turn VSync off entirely for example:
"vblank_mode=1"`
	aeOtherParams       = "Others **"
	aeOtherParamsDetail = "** Other parameters you want to pass to your ZDoom port for this game"

	aeOkButton = "Ok"
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
		foo := games.NewGame("", "", "", "")
		g = &foo
		title = addGame
		gWasNil = true
	}

	ae := tview.NewForm()

	inputName := tview.NewInputField().SetText(g.Name).SetLabel(aeName).SetLabelColor(tview.Styles.SecondaryTextColor)
	ae.AddFormItem(inputName)

	inputSourcePort := tview.NewDropDown().SetOptions([]string{"NA"}, nil).SetLabel(aeSourcePort).SetLabelColor(tview.Styles.SecondaryTextColor)
	ae.AddFormItem(inputSourcePort)
	if len(cfg.Instance().SourcePorts) > 0 {
		inputSourcePort.SetOptions(cfg.Instance().SourcePorts, nil)
		if i, isIn := indexOfItemIn(g.SourcePort, cfg.Instance().SourcePorts); isIn {
			inputSourcePort.SetCurrentOption(i)
		} else {
			inputSourcePort.SetCurrentOption(0)
		}
	}
	// get shared configs for selected source port
	sharedCfgs := []string{}
	inputSourcePort.SetDoneFunc(func(key tcell.Key) {
		_, selectedPort := inputSourcePort.GetCurrentOption()
		sharedCfgs = cfg.GetSharedGameConfigs(selectedPort)
	})

	inputIwad := tview.NewDropDown().SetOptions([]string{"NA"}, nil).SetLabel(aeIWAD).SetLabelColor(tview.Styles.SecondaryTextColor)
	ae.AddFormItem(inputIwad)
	if len(cfg.Instance().IWADs) > 0 {
		inputIwad.SetOptions(cfg.Instance().IWADs, nil)
		if i, isIn := indexOfItemIn(g.Iwad, cfg.Instance().IWADs); isIn {
			inputIwad.SetCurrentOption(i)
		} else {
			inputIwad.SetCurrentOption(0)
		}
	}

	inputOwnCfg := tview.NewCheckbox().SetChecked(g.PersonalPortCfg).SetLabel(aeOwnCfg).SetLabelColor(tview.Styles.SecondaryTextColor)
	ae.AddFormItem(inputOwnCfg)

	inputSharedCfg := tview.NewInputField().SetText(g.SharedConfig).SetLabel(aeSharedCfg).SetLabelColor(tview.Styles.SecondaryTextColor)
	inputSharedCfg.SetAutocompleteFunc(
		func(currentText string) (entries []string) {
			if len(currentText) == 0 {
				return
			}
			for _, word := range sharedCfgs {
				if strings.HasPrefix(strings.ToLower(word), strings.ToLower(currentText)) {
					entries = append(entries, word)
				}
			}
			return
		})
	ae.AddFormItem(inputSharedCfg)

	inputURL := tview.NewInputField().SetText(g.Link).SetLabel(aeLink).SetLabelColor(tview.Styles.SecondaryTextColor)
	ae.AddFormItem(inputURL)

	inputEnvVars := tview.NewInputField().SetText(g.EnvironmentString()).SetLabel(aeEnvironment).SetLabelColor(tview.Styles.SecondaryTextColor)
	ae.AddFormItem(inputEnvVars)

	inputCustomParams := tview.NewInputField().SetText(g.ParamsString()).SetLabel(aeOtherParams).SetLabelColor(tview.Styles.SecondaryTextColor)
	ae.AddFormItem(inputCustomParams)

	ae.AddButton(aeOkButton, func() {
		g.Name = strings.TrimSpace(inputName.GetText())
		_, g.SourcePort = inputSourcePort.GetCurrentOption()
		_, g.Iwad = inputIwad.GetCurrentOption()
		g.PersonalPortCfg = inputOwnCfg.IsChecked()
		g.SharedConfig = inputSharedCfg.GetText()
		g.Environment = splitParams(inputEnvVars.GetText())
		g.CustomParameters = splitParams(inputCustomParams.GetText())
		g.Link = inputURL.GetText()

		if gWasNil {
			games.AddGame(*g)
		}

		games.Persist()
		games.InformChangeListeners()
		appModeNormal()
	})

	// build form
	addEditGameForm := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(ae, 0, 1, true).
		AddItem(tview.NewTextView().SetText(aeEnvironmentDetail), 2, 0, false).
		AddItem(nil, 1, 0, false).
		AddItem(tview.NewTextView().SetText(aeOtherParamsDetail), 1, 0, false)
	addEditGameForm.SetBorder(true)
	addEditGameForm.SetTitle(title)
	addEditGameForm.SetBackgroundColor(tview.Styles.PrimitiveBackgroundColor)
	addEditGameForm.SetBorderPadding(1, 1, 1, 1)

	return addEditGameForm
}
