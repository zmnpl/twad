package tui

import (
	"strconv"
	"strings"

	"github.com/rivo/tview"
	"github.com/zmnpl/twad/cfg"
)

const (
	optionsOkButtonLabel          = "Ok"
	optionsHeader                 = "Options"
	optionsPathLabel              = "Base Path"
	optionsWarnBeforeLabel        = "Warn before deletion"
	optionsSourcePortLabel        = "Source Ports"
	optionsIwadsLabel             = "IWADs"
	optionsNextTimeFirstStart     = "Show path selection on next start"
	optionsSaveDirsLabel          = "Use separate save game directories"
	optionsPrintHeaderLabel       = "Show header"
	optionsGamesListRelativeWitdh = "Game list relative width (% 1-100)"
	optionsOkButtonText           = "That's it!"
)

func makeOptions() *tview.Flex {
	o := tview.NewForm()

	path := tview.NewInputField().SetLabel(optionsPathLabel).SetLabelColor(tview.Styles.SecondaryTextColor).SetText(cfg.GetInstance().ModBasePath)
	o.AddFormItem(path)

	sourcePorts := tview.NewInputField().SetLabel(optionsSourcePortLabel).SetLabelColor(tview.Styles.SecondaryTextColor).SetText(strings.Join(cfg.GetInstance().SourcePorts, ","))
	o.AddFormItem(sourcePorts)

	iwads := tview.NewInputField().SetLabel(optionsIwadsLabel).SetLabelColor(tview.Styles.SecondaryTextColor).SetText(strings.Join(cfg.GetInstance().IWADs, ","))
	o.AddFormItem(iwads)

	printHeader := tview.NewCheckbox().SetLabel(optionsPrintHeaderLabel).SetLabelColor(tview.Styles.SecondaryTextColor).SetChecked(cfg.GetInstance().PrintHeader)
	o.AddFormItem(printHeader)

	warn := tview.NewCheckbox().SetLabel(optionsWarnBeforeLabel).SetLabelColor(tview.Styles.SecondaryTextColor).SetChecked(cfg.GetInstance().WarnBeforeDelete)
	o.AddFormItem(warn)

	saveDirs := tview.NewCheckbox().SetLabel(optionsSaveDirsLabel).SetLabelColor(tview.Styles.SecondaryTextColor).SetChecked(cfg.GetInstance().SaveDirs)
	o.AddFormItem(saveDirs)

	gameListRelWidth := tview.NewInputField().SetLabel(optionsGamesListRelativeWitdh).SetLabelColor(tview.Styles.SecondaryTextColor).SetAcceptanceFunc(func(text string, char rune) bool {
		if text == "-" {
			return false
		}
		i, err := strconv.Atoi(text)
		return err == nil && i > 0 && i <= 100
	})
	gameListRelWidth.SetText(strconv.Itoa(cfg.GetInstance().GameListRelativeWidth))
	o.AddFormItem(gameListRelWidth)

	firstStart := tview.NewCheckbox().SetLabel(optionsNextTimeFirstStart).SetLabelColor(tview.Styles.SecondaryTextColor).SetChecked(!cfg.GetInstance().Configured)
	o.AddFormItem(firstStart)

	o.AddButton("Cool", func() {
		c := cfg.GetInstance()

		c.ModBasePath = path.GetText()
		cfg.AddPathToCfgs()

		sps := strings.Split(sourcePorts.GetText(), ",")
		for i := range sps {
			sps[i] = strings.TrimSpace(sps[i])
		}
		c.SourcePorts = sps

		iwds := strings.Split(iwads.GetText(), ",")
		for i := range iwds {
			iwds[i] = strings.TrimSpace(iwds[i])
		}
		c.IWADs = iwds

		c.PrintHeader = printHeader.IsChecked()
		c.WarnBeforeDelete = warn.IsChecked()
		c.SaveDirs = saveDirs.IsChecked()
		c.GameListRelativeWidth, _ = strconv.Atoi(gameListRelWidth.GetText())
		c.Configured = !firstStart.IsChecked()

		cfg.Persist()
		appModeNormal()

	})

	// layout
	settingsPage := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(o, 90, 0, true).
		AddItem(tview.NewBox().SetBorder(false), 0, 1, false)

	return settingsPage
}
