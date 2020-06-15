package tui

import (
	"strings"

	"github.com/rivo/tview"
	"github.com/zmnpl/twad/cfg"
)

const (
	optionsOkButtonLabel      = "Ok"
	optionsHeader             = "Options"
	optionsPathLabel          = "Base Path"
	optionsWarnBeforeLabel    = "Warn before deletion"
	optionsSourcePortLabel    = "Source Ports"
	optionsIwadsLabel         = "IWADs"
	optionsNextTimeFirstStart = "Show path selection on next start"
	optionsSaveDirsLabel      = "Use separate save game directories"
	optionsPrintHeaderLabel   = "Show Header"
	optionsMaxLabelLength     = 35
)

func makeOptions() *tview.Flex {
	leftColSize := optionsMaxLabelLength + 1
	rigthColSize := 1

	pathLabel := tview.NewTextView().SetText(optionsPathLabel).SetTextColor(tview.Styles.SecondaryTextColor)
	path := tview.NewInputField()
	path.SetText(cfg.GetInstance().ModBasePath)
	pathRow := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(pathLabel, leftColSize, 0, false).
		AddItem(path, 0, rigthColSize, true)

	sourcePortsLabel := tview.NewTextView().SetText(optionsSourcePortLabel).SetTextColor(tview.Styles.SecondaryTextColor)
	sourcePorts := tview.NewInputField().SetText(strings.Join(cfg.GetInstance().SourcePorts, ","))
	sourcePortsRow := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(sourcePortsLabel, leftColSize, 0, false).
		AddItem(sourcePorts, 0, rigthColSize, false)

	iwadsLabel := tview.NewTextView().SetText(optionsIwadsLabel).SetTextColor(tview.Styles.SecondaryTextColor)
	iwads := tview.NewInputField().SetText(strings.Join(cfg.GetInstance().IWADs, ","))
	iwadsRow := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(iwadsLabel, leftColSize, 0, false).
		AddItem(iwads, 0, rigthColSize, false)

	printHeaderLabel := tview.NewTextView().SetText(optionsPrintHeaderLabel).SetTextColor(tview.Styles.SecondaryTextColor)
	printHeader := tview.NewCheckbox().SetChecked(cfg.GetInstance().PrintHeader)
	printHeaderRow := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(printHeaderLabel, leftColSize, 0, false).
		AddItem(printHeader, 0, rigthColSize, false)

	warnBeforeDeleteLabel := tview.NewTextView().SetText(optionsWarnBeforeLabel).SetTextColor(tview.Styles.SecondaryTextColor)
	warn := tview.NewCheckbox().SetChecked(cfg.GetInstance().WarnBeforeDelete)
	warnRow := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(warnBeforeDeleteLabel, leftColSize, 0, false).
		AddItem(warn, 0, rigthColSize, false)

	saveGameDirsLabel := tview.NewTextView().SetText(optionsSaveDirsLabel).SetTextColor(tview.Styles.SecondaryTextColor)
	saveDirs := tview.NewCheckbox().SetChecked(cfg.GetInstance().SaveDirs)
	saveDirsRow := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(saveGameDirsLabel, leftColSize, 0, false).
		AddItem(saveDirs, 0, rigthColSize, false)

	firstStartLabel := tview.NewTextView().SetText(optionsNextTimeFirstStart).SetTextColor(tview.Styles.SecondaryTextColor)
	firstStart := tview.NewCheckbox().SetChecked(!cfg.GetInstance().Configured)
	firstStartRow := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(firstStartLabel, leftColSize, 0, false).
		AddItem(firstStart, 0, rigthColSize, false)

	okButton := tview.NewButton(optionsOkButtonLabel)
	okButton.SetBackgroundColorActivated(tview.Styles.PrimaryTextColor)
	okButton.SetLabelColorActivated(tview.Styles.ContrastBackgroundColor)
	okButtonRow := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(okButton, len(optionsOkButtonLabel)+2, 0, false).
		AddItem(nil, 0, rigthColSize, false)
	okButton.SetSelectedFunc(func() {
		config := cfg.GetInstance()
		config.ModBasePath = path.GetText()
		cfg.AddPathToCfgs()

		sps := strings.Split(sourcePorts.GetText(), ",")
		for i := range sps {
			sps[i] = strings.TrimSpace(sps[i])
		}
		config.SourcePorts = sps

		iwds := strings.Split(iwads.GetText(), ",")
		for i := range iwds {
			iwds[i] = strings.TrimSpace(iwds[i])
		}
		config.IWADs = iwds

		config.PrintHeader = printHeader.IsChecked()
		config.WarnBeforeDelete = warn.IsChecked()
		config.SaveDirs = saveDirs.IsChecked()
		config.Configured = !firstStart.IsChecked()
		cfg.Persist()
		appModeNormal()
	})

	// navigation path
	path.SetInputCapture(tabNavigate(okButton, sourcePorts))
	sourcePorts.SetInputCapture(tabNavigate(path, iwads))
	iwads.SetInputCapture(tabNavigate(sourcePorts, printHeader))
	printHeader.SetInputCapture(tabNavigate(iwads, warn))
	warn.SetInputCapture(tabNavigate(iwads, saveDirs))
	saveDirs.SetInputCapture(tabNavigate(warn, firstStart))
	firstStart.SetInputCapture(tabNavigate(saveDirs, okButton))
	okButton.SetInputCapture(tabNavigate(firstStart, path))

	options := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(pathRow, 1, 0, true).
		AddItem(nil, 1, 0, false).
		AddItem(sourcePortsRow, 1, 0, false).
		AddItem(nil, 1, 0, false).
		AddItem(iwadsRow, 1, 0, false).
		AddItem(nil, 1, 0, false).
		AddItem(printHeaderRow, 1, 0, false).
		AddItem(nil, 1, 0, false).
		AddItem(warnRow, 1, 0, false).
		AddItem(nil, 1, 0, false).
		AddItem(saveDirsRow, 1, 0, false).
		AddItem(nil, 1, 0, false).
		AddItem(firstStartRow, 1, 0, false).
		AddItem(nil, 1, 0, false).
		AddItem(okButtonRow, 1, 0, false)
	options.SetBorder(true)
	//options.SetBorderColor(accentColor)
	//options.SetTitleColor(accentColor)
	options.SetTitle(optionsHeader)
	options.SetBorderPadding(1, 1, 1, 1)

	settingsPage := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(options, 90, 0, true).
		AddItem(tview.NewBox().SetBorder(false), 0, 1, false)

	return settingsPage
}
