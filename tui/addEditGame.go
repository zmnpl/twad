package tui

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/zmnpl/twad/base"
	"github.com/zmnpl/twad/games"
	"github.com/zmnpl/twad/portspec"
)

const (
	addGame  = "Add New Game"
	editGame = "Edit Game"

	aeName       = "Name"
	aeSourcePort = "Source Port"
	aeIWAD       = "IWAD"
	aeOwnCfg     = "Use Own Config"
	aeSharedCfgT = "Use Shared Config [%v]"
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

	port, iwad := "", ""
	if len(config.SourcePorts) > 0 {
		port = config.SourcePorts[0]
	}
	if len(config.IWADs) > 0 {
		iwad = config.IWADs[0]
	}

	if g == nil {
		foo := games.NewGame("", port, "", iwad)
		g = &foo
		title = addGame
		gWasNil = true
	}
	expectedExtension := portspec.PortConfigFileExtension(g.Port)

	// create basic form items
	inputName := tview.NewInputField().SetText(g.Name).SetLabel(aeName).SetLabelColor(tview.Styles.SecondaryTextColor)
	inputOwnCfg := tview.NewCheckbox().SetChecked(g.PersonalPortCfg).SetLabel(aeOwnCfg).SetLabelColor(tview.Styles.SecondaryTextColor)
	inputSharedCfg := tview.NewInputField().SetText(g.SharedConfig).SetLabel(fmt.Sprintf(aeSharedCfgT, expectedExtension)).SetLabelColor(tview.Styles.SecondaryTextColor)
	if g.PersonalPortCfg {
		inputSharedCfg.SetLabel(warnColor + fmt.Sprintf(aeSharedCfgT, expectedExtension))
	}
	inputSourcePort := tview.NewDropDown().SetOptions([]string{"NA"}, nil).SetLabel(aeSourcePort).SetLabelColor(tview.Styles.SecondaryTextColor)
	inputIwad := tview.NewDropDown().SetOptions([]string{"NA"}, nil).SetLabel(aeIWAD).SetLabelColor(tview.Styles.SecondaryTextColor)
	inputURL := tview.NewInputField().SetText(g.Link).SetLabel(aeLink).SetLabelColor(tview.Styles.SecondaryTextColor)
	inputEnvVars := tview.NewInputField().SetText(g.EnvironmentString()).SetLabel(aeEnvironment).SetLabelColor(tview.Styles.SecondaryTextColor)
	inputCustomParams := tview.NewInputField().SetText(g.ParamsString()).SetLabel(aeOtherParams).SetLabelColor(tview.Styles.SecondaryTextColor)

	ae := tview.NewForm()

	// functionality of form items
	// port
	if len(base.Config().SourcePorts) > 0 {
		inputSourcePort.SetOptions(base.Config().SourcePorts, nil)
		if i, isIn := indexOfItemIn(g.Port, base.Config().SourcePorts); isIn {
			inputSourcePort.SetCurrentOption(i)
		} else {
			inputSourcePort.SetCurrentOption(0)
		}
	}
	// get shared configs for selected port
	sharedCfgs := []string{}
	inputSourcePort.SetDoneFunc(func(key tcell.Key) {
		_, selectedPort := inputSourcePort.GetCurrentOption()
		sharedCfgs = base.GetSharedGameConfigs(selectedPort)
		expectedExtension = portspec.PortConfigFileExtension(selectedPort)

		inputSharedCfg.SetLabel(fmt.Sprintf(aeSharedCfgT, expectedExtension))
		if inputOwnCfg.IsChecked() {
			inputSharedCfg.SetLabel(warnColor + fmt.Sprintf(aeSharedCfgT, expectedExtension))
		}
	})

	// for iwad
	if len(base.Config().IWADs) > 0 {
		inputIwad.SetOptions(base.Config().IWADs, nil)
		if i, isIn := indexOfItemIn(g.Iwad, base.Config().IWADs); isIn {
			inputIwad.SetCurrentOption(i)
		} else {
			inputIwad.SetCurrentOption(0)
		}
	}

	// own configs
	inputOwnCfg.SetDoneFunc(func(key tcell.Key) {
		if inputOwnCfg.IsChecked() {
			inputSharedCfg.SetLabel(warnColor + fmt.Sprintf(aeSharedCfgT, expectedExtension))
			return
		}
		inputSharedCfg.SetLabel(fmt.Sprintf(aeSharedCfgT, expectedExtension))
	})

	// shared configs
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
	inputSharedCfg.SetDoneFunc(func(key tcell.Key) {
		if len(inputSharedCfg.GetText()) > 0 && !strings.HasSuffix(strings.ToLower(inputSharedCfg.GetText()), expectedExtension) {
			inputSharedCfg.SetText(inputSharedCfg.GetText() + expectedExtension)
		}
	})

	// add controls in order

	ae.AddFormItem(inputName)
	ae.AddFormItem(inputSourcePort)
	ae.AddFormItem(inputIwad)
	ae.AddFormItem(inputOwnCfg)
	ae.AddFormItem(inputSharedCfg)
	ae.AddFormItem(inputURL)
	ae.AddFormItem(inputEnvVars)
	ae.AddFormItem(inputCustomParams)

	ae.AddButton(aeOkButton, func() {
		g.Name = strings.TrimSpace(inputName.GetText())
		_, g.Port = inputSourcePort.GetCurrentOption()
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
