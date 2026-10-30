package tui

import (
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/zmnpl/twad/base"
	"github.com/zmnpl/twad/tui/shared"
)

// function to check source port entry
func sourcePortCheck(input *tview.InputField) {
	// does this path exist?
	if _, err := os.Stat(input.GetText()); os.IsNotExist(err) {
		if shared.CommandExists(strings.TrimSpace(input.GetText())) {
			input.SetLabel(labelOk(shared.Dict.OptsSourcePortLabel, shared.Dict.OptsLooksGood))
			return
		}
		input.SetLabel(labelWarn(shared.Dict.OptsSourcePortLabel, shared.Dict.OptsErrPathDoesntExist))
		return
	}

	input.SetLabel(labelOk(shared.Dict.OptsSourcePortLabel, shared.Dict.OptsLooksGood))
}

// builder function that creates a function which autocompletes input fields
func (u *ui) autocompletePathMaker(path *tview.InputField, dirsOnly bool, extensionFilter map[string]bool) func(pathText string) (entries []string) {
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

			u.app.Draw()

		}()
		return nil
	}

	return foo
}

func (u *ui) makeOptions() *tview.Flex {
	o := tview.NewForm()

	// doomwaddir
	//#######################################################################

	// added to form later
	iwads := tview.NewInputField().SetLabel(shared.Dict.OptsIwadsLabel).SetText(strings.Join(base.Config().IWADs, ","))

	// path for doomwaddir
	doomwaddirPath := tview.NewInputField().SetLabel(shared.Dict.OptsPathLabel).SetText(base.Config().WadDir)

	// add doomwaddir to form
	o.AddFormItem(doomwaddirPath)
	// add iwads input to form
	o.AddFormItem(iwads)

	// autocompletion for path
	autocompleteDoomwaddir := u.autocompletePathMaker(doomwaddirPath, true, nil)
	doomwaddirPath.SetAutocompleteFunc(autocompleteDoomwaddir)

	// check iwad path
	doomwaddirPathCheck := func() {
		// does this path exist?
		if _, err := os.Stat(doomwaddirPath.GetText()); os.IsNotExist(err) {
			doomwaddirPath.SetLabel(labelWarn(shared.Dict.OptsPathLabel, shared.Dict.OptsErrPathDoesntExist))
			return
		}

		// check if selected path contains any iwads
		if hasIwad, err := base.PathHasIwads(doomwaddirPath.GetText()); !hasIwad {
			if err != nil {
				doomwaddirPath.SetLabel(labelWarn(shared.Dict.OptsPathLabel, "("+err.Error()+")"))
				return
			}
			doomwaddirPath.SetLabel(labelWarn(shared.Dict.OptsPathLabel, shared.Dict.OptsErrPathNoIWads))
			return
		}

		//iwads.SetText(strings.Join(base.GePathIwads(doomwaddirPath.GetText()), ","))

		doomwaddirPath.SetLabel(labelOk(shared.Dict.OptsPathLabel, shared.Dict.OptsLooksGood))
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
		sourcePort := tview.NewInputField().SetLabel(shared.Dict.OptsSourcePortLabel)
		autocompleteSourcePort := u.autocompletePathMaker(sourcePort, false, spExtensionFilter)

		// only autocomplete source ports on windows
		// on linux the simple executable name is usually enough when it is in path
		if runtime.GOOS == "windows" {
			sourcePort.SetAutocompleteFunc(autocompleteSourcePort)
		}
		sourcePort.SetDoneFunc(func(key tcell.Key) {
			sourcePortCheck(sourcePort)
		})
		if i < len(base.Config().Ports) {
			sourcePort.SetText(base.Config().Ports[i])
		}
		spInputs[i] = sourcePort
		o.AddFormItem(sourcePort)
	}

	// ui options
	//#######################################################################
	termColors := tview.NewCheckbox().SetLabel(shared.Dict.OptsUseTerminalColors).SetChecked(base.Config().UseTerminalColors)
	o.AddFormItem(termColors)

	dontWarn := tview.NewCheckbox().SetLabel(shared.Dict.OptsDontWarn).SetChecked(base.Config().DeleteWithoutWarning)
	o.AddFormItem(dontWarn)

	printHeader := tview.NewCheckbox().SetLabel(shared.Dict.OptsHideHeader).SetChecked(base.Config().HideHeader)
	o.AddFormItem(printHeader)

	gameListRelWidth := tview.NewInputField().SetLabel(shared.Dict.OptsGamesListRelativeWitdh).SetAcceptanceFunc(func(text string, char rune) bool {
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
	o.AddButton(shared.Dict.OptsOkButtonLabel, func() {
		c := base.Config()

		c.WadDir = doomwaddirPath.GetText()

		sps := make([]string, 0, base.MAX_SOURCE_PORTS)
		for i := range spInputs {
			if sp := strings.TrimSpace(spInputs[i].GetText()); sp != "" {
				sps = append(sps, sp)
			}
		}
		c.Ports = sps

		c.IWADs = nil
		for _, iwad := range strings.Split(iwads.GetText(), ",") {
			if iwad = strings.TrimSpace(iwad); iwad != "" {
				c.IWADs = append(c.IWADs, iwad)
			}
		}

		c.HideHeader = printHeader.IsChecked()
		c.DeleteWithoutWarning = dontWarn.IsChecked()
		c.GameListRelativeWidth, _ = strconv.Atoi(gameListRelWidth.GetText())
		c.UseTerminalColors = termColors.IsChecked()

		base.EnableBasePath()
		u.appModeNormal()
		if err := base.Persist(); err != nil {
			u.showError("Saving options failed", err.Error(), nil, nil)
		}
	})

	// layout
	settingsPage := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(o, 90, 0, true).
		AddItem(tview.NewBox().SetBorder(false), 0, 1, false)

	return settingsPage
}
