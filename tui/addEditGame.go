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
	aeConfig     = "Config"

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
	if cfg.Instance().SourcePorts != nil && len(cfg.Instance().SourcePorts) > 0 {
		inputSourcePort.SetOptions(cfg.Instance().SourcePorts, nil)
		if i, isIn := indexOfItemIn(g.SourcePort, cfg.Instance().SourcePorts); isIn {
			inputSourcePort.SetCurrentOption(i)
		} else {
			inputSourcePort.SetCurrentOption(0)
		}
	}

	inputIwad := tview.NewDropDown().SetOptions([]string{"NA"}, nil).SetLabel(aeIWAD).SetLabelColor(tview.Styles.SecondaryTextColor)
	ae.AddFormItem(inputIwad)
	if cfg.Instance().IWADs != nil && len(cfg.Instance().IWADs) > 0 {
		inputIwad.SetOptions(cfg.Instance().IWADs, nil)
		if i, isIn := indexOfItemIn(g.Iwad, cfg.Instance().IWADs); isIn {
			inputIwad.SetCurrentOption(i)
		} else {
			inputIwad.SetCurrentOption(0)
		}
	}

	inputConfig := tview.NewDropDown().SetOptions([]string{"NA"}, nil).SetLabel(aeConfig).SetLabelColor(tview.Styles.SecondaryTextColor)
	ae.AddFormItem(inputConfig)
	if cfg.Instance().Configs != nil && len(cfg.Instance().Configs) > 0 {
		options := cfg.Instance().Configs
		options = append([]string{"Default"}, cfg.Instance().Configs...)
		inputConfig.SetOptions(options, nil)
		if i, isIn := indexOfItemIn(g.Config, options); isIn {
			inputConfig.SetCurrentOption(i)
		} else {
			inputConfig.SetCurrentOption(0)
		}
	} else {
		// TODO: Make this get the value from cfg.go, configPortConfigsDefaultLabel 
		inputConfig.SetOptions([]string{"Default"}, nil)
		inputConfig.SetCurrentOption(0)
	}

	inputEnvVars := tview.NewInputField().SetText(g.EnvironmentString()).SetLabel(aeEnvironment).SetLabelColor(tview.Styles.SecondaryTextColor)
	ae.AddFormItem(inputEnvVars)

	inputCustomParams := tview.NewInputField().SetText(g.ParamsString()).SetLabel(aeOtherParams).SetLabelColor(tview.Styles.SecondaryTextColor)
	ae.AddFormItem(inputCustomParams)

	ae.AddButton(aeOkButton, func() {
		g.Name = strings.TrimSpace(inputName.GetText())
		_, g.SourcePort = inputSourcePort.GetCurrentOption()
		_, g.Iwad = inputIwad.GetCurrentOption()
		_, g.Config = inputConfig.GetCurrentOption()
		g.Environment = splitParams(inputEnvVars.GetText())
		g.CustomParameters = splitParams(inputCustomParams.GetText())

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
