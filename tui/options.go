package tui

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/zmnpl/twad/base"
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
	optsSourcePortLabel          = "Source Port"
	optsIwadsLabel               = "IWADs"
	optsHideHeader               = "UI - Hide big DOOM logo"
	optsGamesListRelativeWitdh   = "UI - Game list relative width (1-100%)"
)

// function to check source port entry
var sourcePortCheck = func(input *tview.InputField) {
	// does this path exist?
	if _, err := os.Stat(input.GetText()); os.IsNotExist(err) {
		if commandExists(strings.TrimSpace(input.GetText())) {
			input.SetLabel(optsSourcePortLabel + goodColor + " " + optsLooksGood)
			return
		}
		input.SetLabel(optsSourcePortLabel + warnColor + " " + optsErrPathDoesntExist)
		return
	}

	input.SetLabel(optsSourcePortLabel + goodColor + " " + optsLooksGood)
}

// builder function that creates a function which autocompletes input fields
var autocompletePathMaker = func(path *tview.InputField, dirsOnly bool, extensionFilter map[string]bool) func(pathText string) (entries []string) {
	var mutex sync.Mutex
	prefixMap := make(map[string][]string)
	firstStart := true

	foo := func(pathText string) (entries []string) {
		// Ignore empty text.
		prefix := strings.TrimSpace(strings.ToLower(pathText))
		if prefix == "" {
			return nil
		}

		// Do we have entries for this text already?
		mutex.Lock()
		defer mutex.Unlock()
		// Prevent autocomplete to be shown when the options panel is drawn initially
		if firstStart {
			firstStart = false
			return nil
		}
		entries, ok := prefixMap[prefix]
		if ok {
			return entries
		}

		// No entries yet get entries in goroutine
		go func() {
			dir := filepath.Dir(pathText)
			files, err := os.ReadDir(dir)

			if err != nil {
				return
			}

			entries := make([]string, 0, len(files))
			for _, file := range files {
				// dont't show hidden folders
				if strings.HasPrefix(file.Name(), ".") {
					continue
				}
				// only show dirs if that is wanted
				if dirsOnly && !file.IsDir() {
					continue
				}
				// filter specific extensions
				if len(extensionFilter) > 0 {
					_, extensionOk := extensionFilter[strings.ToLower(filepath.Ext(file.Name()))]
					if !(extensionOk || file.IsDir()) {
						continue
					}
				}
				entries = append(entries, filepath.Join(dir, file.Name()))
			}

			mutex.Lock()
			prefixMap[prefix] = entries
			mutex.Unlock()

			path.Autocomplete()

			app.Draw()

		}()
		return nil
	}

	return foo
}

func makeOptions() *tview.Flex {
	o := tview.NewForm()

	// doomwaddir
	//#######################################################################

	// added to form later
	iwads := tview.NewInputField().SetLabel(optsIwadsLabel).SetLabelColor(tview.Styles.SecondaryTextColor).SetText(strings.Join(base.Config().IWADs, ","))
	// path for doomwaddir
	doomwaddirPath := tview.NewInputField().SetLabel(optsPathLabel).SetLabelColor(tview.Styles.SecondaryTextColor).SetText(base.Config().WadDir)

	// add doomwaddir to form
	o.AddFormItem(doomwaddirPath)
	// add iwads input to form
	o.AddFormItem(iwads)

	// autocompletion for path
	autocompleteDoomwaddir := autocompletePathMaker(doomwaddirPath, true, nil)
	doomwaddirPath.SetAutocompleteFunc(autocompleteDoomwaddir)

	// check iwad path
	doomwaddirPathCheck := func() {
		// does this path exist?
		if _, err := os.Stat(doomwaddirPath.GetText()); os.IsNotExist(err) {
			doomwaddirPath.SetLabel(optsPathLabel + warnColor + " " + optsErrPathDoesntExist)
			return
		}

		// check if selected path contains any iwads
		if hasIwad, err := base.PathHasIwads(doomwaddirPath.GetText()); !hasIwad {
			if err != nil {
				doomwaddirPath.SetLabel(optsPathLabel + warnColor + " (" + err.Error() + ")")
			}
			doomwaddirPath.SetLabel(optsPathLabel + warnColor + " " + optsErrPathNoIWads)
			return
		}

		availableIwads, _ := base.GePathIwads(doomwaddirPath.GetText())
		iwads.SetText(strings.Join(availableIwads, ","))

		doomwaddirPath.SetLabel(optsPathLabel + goodColor + " " + optsLooksGood)
	}
	// initial check of configured path
	doomwaddirPathCheck()
	// check after entry
	doomwaddirPath.SetDoneFunc(func(key tcell.Key) {
		doomwaddirPathCheck()
	})

	// source ports
	//#######################################################################

	// windows exe filter
	var spExtensionFilter map[string]bool
	if runtime.GOOS == "windows" {
		spExtensionFilter = make(map[string]bool)
		spExtensionFilter[".exe"] = true
	}

	// add source port input fields
	spInputs := make([]*tview.InputField, base.MAX_SOURCE_PORTS)
	for i := 0; i < base.MAX_SOURCE_PORTS; i++ {
		sourcePort := tview.NewInputField().SetLabel(optsSourcePortLabel).SetLabelColor(tview.Styles.SecondaryTextColor)
		autocompleteSourcePort := autocompletePathMaker(sourcePort, false, spExtensionFilter)

		// only autocomplete source ports on windows
		// on linux the simple executable name is usually enough when it is in path
		if runtime.GOOS == "windows" {
			sourcePort.SetAutocompleteFunc(autocompleteSourcePort)
		}
		sourcePort.SetDoneFunc(func(key tcell.Key) {
			sourcePortCheck(sourcePort)
		})
		if i < len(base.Config().SourcePorts) {
			sourcePort.SetText(base.Config().SourcePorts[i])
		}
		spInputs[i] = sourcePort
		o.AddFormItem(sourcePort)
	}

	// ui options
	//#######################################################################
	dontWarn := tview.NewCheckbox().SetLabel(optsDontWarn).SetLabelColor(tview.Styles.SecondaryTextColor).SetChecked(base.Config().DeleteWithoutWarning)
	o.AddFormItem(dontWarn)

	printHeader := tview.NewCheckbox().SetLabel(optsHideHeader).SetLabelColor(tview.Styles.SecondaryTextColor).SetChecked(base.Config().HideHeader)
	o.AddFormItem(printHeader)

	gameListRelWidth := tview.NewInputField().SetLabel(optsGamesListRelativeWitdh).SetLabelColor(tview.Styles.SecondaryTextColor).SetAcceptanceFunc(func(text string, char rune) bool {
		if text == "-" {
			return false
		}
		i, err := strconv.Atoi(text)
		return err == nil && i > 0 && i <= 100
	})
	gameListRelWidth.SetText(strconv.Itoa(base.Config().GameListRelativeWidth))
	o.AddFormItem(gameListRelWidth)

	// ok button and processing of options
	//#######################################################################
	o.AddButton(optsOkButtonLabel, func() {
		c := base.Config()

		c.WadDir = doomwaddirPath.GetText()

		sps := make([]string, base.MAX_SOURCE_PORTS)
		for i := range spInputs {
			sps[i] = strings.TrimSpace(spInputs[i].GetText())
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

		base.Persist()
		base.EnableBasePath()
		appModeNormal()
	})

	// layout
	settingsPage := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(o, 90, 0, true).
		AddItem(tview.NewBox().SetBorder(false), 0, 1, false)

	return settingsPage
}

// as util
func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}
