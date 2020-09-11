package tui

import (
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"github.com/zmnpl/twad/cfg"
)

const (
	optsErrPathDoesntExist = "doesn't exist"
	optsErrPathNoIWads     = "doesn't contain IWADs"
	optsLooksGood          = "looks good"

	optsHeader                   = "Options"
	optsOkButtonLabel            = "Save"
	optsPathLabel                = "WAD Dir"
	optsDontDOOMWADDIR           = "Do NOT set DOOMWADDIR for current session (use your shell's default)"
	optsWriteBasePathToEngineCFG = "Write the path into DOOM engines *.ini files"
	optsDontWarn                 = "Do NOT warn before deletion"
	optsSourcePortLabel          = "Source Ports"
	optsIwadsLabel               = "IWADs"
	optsNextTimeFirstStart       = "Show path selection dialog on next start"
	optsHideHeader               = "UI - Hide big DOOM logo"
	optsGamesListRelativeWitdh   = "UI - Game list relative width (1-100%)"
)

func pathHasIwad(path string) (bool, error) {
	files, err := ioutil.ReadDir(path)

	if err != nil {
		return false, err
	}

	for _, file := range files {
		for _, iwad := range cfg.KnownIwads {
			if strings.ToLower(file.Name()) == iwad {
				return true, nil
			}
		}
	}

	return false, nil
}

func makeOptions() *tview.Flex {
	o := tview.NewForm()

	path := tview.NewInputField().SetLabel(optsPathLabel).SetLabelColor(tview.Styles.SecondaryTextColor).SetText(cfg.Instance().WadDir)
	o.AddFormItem(path)
	pathDoneCheck := func() {
		// does this path exist?
		if _, err := os.Stat(path.GetText()); os.IsNotExist(err) {
			path.SetLabel(optsPathLabel + warnColor + " " + optsErrPathDoesntExist)
			return
		}

		// check if selected path contains any iwads
		if hasIwad, err := pathHasIwad(path.GetText()); !hasIwad {
			if err != nil {
				path.SetLabel(optsPathLabel + warnColor + " (" + err.Error() + ")")
			}
			path.SetLabel(optsPathLabel + warnColor + " " + optsErrPathNoIWads)
			return
		}

		path.SetLabel(optsPathLabel + goodColor + " " + optsLooksGood)
	}
	// initial check of configured path
	pathDoneCheck()
	// check after entry
	path.SetDoneFunc(func(key tcell.Key) {
		pathDoneCheck()
	})

	firstStart := tview.NewCheckbox().SetLabel(optsNextTimeFirstStart).SetLabelColor(tview.Styles.SecondaryTextColor).SetChecked(!cfg.Instance().Configured)
	o.AddFormItem(firstStart)

	sourcePorts := tview.NewInputField().SetLabel(optsSourcePortLabel).SetLabelColor(tview.Styles.SecondaryTextColor).SetText(strings.Join(cfg.Instance().SourcePorts, ","))
	o.AddFormItem(sourcePorts)

	iwads := tview.NewInputField().SetLabel(optsIwadsLabel).SetLabelColor(tview.Styles.SecondaryTextColor).SetText(strings.Join(cfg.Instance().IWADs, ","))
	o.AddFormItem(iwads)

	dontWarn := tview.NewCheckbox().SetLabel(optsDontWarn).SetLabelColor(tview.Styles.SecondaryTextColor).SetChecked(cfg.Instance().DeleteWithoutWarning)
	o.AddFormItem(dontWarn)

	printHeader := tview.NewCheckbox().SetLabel(optsHideHeader).SetLabelColor(tview.Styles.SecondaryTextColor).SetChecked(cfg.Instance().HideHeader)
	o.AddFormItem(printHeader)

	gameListRelWidth := tview.NewInputField().SetLabel(optsGamesListRelativeWitdh).SetLabelColor(tview.Styles.SecondaryTextColor).SetAcceptanceFunc(func(text string, char rune) bool {
		if text == "-" {
			return false
		}
		i, err := strconv.Atoi(text)
		return err == nil && i > 0 && i <= 100
	})
	gameListRelWidth.SetText(strconv.Itoa(cfg.Instance().GameListRelativeWidth))
	o.AddFormItem(gameListRelWidth)

	o.AddButton(optsOkButtonLabel, func() {
		c := cfg.Instance()

		c.WadDir = path.GetText()

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

		c.HideHeader = printHeader.IsChecked()
		c.DeleteWithoutWarning = dontWarn.IsChecked()
		c.GameListRelativeWidth, _ = strconv.Atoi(gameListRelWidth.GetText())
		c.Configured = !firstStart.IsChecked()

		cfg.Persist()
		cfg.EnableBasePath()
		appModeNormal()
	})

	// layout
	settingsPage := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(o, 90, 0, true).
		AddItem(tview.NewBox().SetBorder(false), 0, 1, false)

	return settingsPage
}
