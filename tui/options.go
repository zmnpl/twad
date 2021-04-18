package tui

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/gdamore/tcell/v2"
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
	optsSourcePortLabel          = "Source Port"
	optsIwadsLabel               = "IWADs"
	optsHideHeader               = "UI - Hide big DOOM logo"
	optsGamesListRelativeWitdh   = "UI - Game list relative width (1-100%)"
)

var (
	autocompletePathMaker = func(path *tview.InputField, dirsOnly bool, extensionFilter map[string]bool) func(pathText string) (entries []string) {
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
					if dirsOnly && !file.IsDir() {
						continue
					}
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

				// Trigger an update to the input field.
				path.Autocomplete()

				// Also redraw the screen.
				app.Draw()

			}()
			return nil
		}

		return foo
	}
)

func pathHasIwad(path string) (bool, error) {
	files, err := os.ReadDir(path)

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

	// autocompletion for path
	autocompleteDoomwaddir := autocompletePathMaker(path, true, nil)
	path.SetAutocompleteFunc(autocompleteDoomwaddir)

	sourcePortDoneCheck := func(input *tview.InputField) {
		// does this path exist?
		if _, err := os.Stat(input.GetText()); os.IsNotExist(err) {
			if commandExists(strings.TrimSpace(input.GetText())) {
				input.SetLabel(optsSourcePortLabel + goodColor + " " + optsLooksGood)
				return
			}
			input.SetLabel(optsPathLabel + warnColor + " " + optsErrPathDoesntExist)
			return
		}

		input.SetLabel(optsSourcePortLabel + goodColor + " " + optsLooksGood)
	}

	iwads := tview.NewInputField().SetLabel(optsIwadsLabel).SetLabelColor(tview.Styles.SecondaryTextColor).SetText(strings.Join(cfg.Instance().IWADs, ","))
	o.AddFormItem(iwads)

	// windows exe filter
	var spExtensionFilter map[string]bool
	if runtime.GOOS == "windows" {
		spExtensionFilter = make(map[string]bool)
		spExtensionFilter[".exe"] = true
	}

	// add source port input fields
	spCount := cfg.Instance().MaxSourcePorts
	spInputs := make([]*tview.InputField, spCount, spCount)
	for i := 0; i < spCount; i++ {
		sourcePort := tview.NewInputField().SetLabel(optsSourcePortLabel + fmt.Sprintf(" %v", i)).SetLabelColor(tview.Styles.SecondaryTextColor).SetText(cfg.Instance().SourcePorts[i])
		autocompleteSourcePort := autocompletePathMaker(sourcePort, false, spExtensionFilter)
		if runtime.GOOS == "windows" {
			sourcePort.SetAutocompleteFunc(autocompleteSourcePort)
		}
		sourcePort.SetDoneFunc(func(key tcell.Key) {
			sourcePortDoneCheck(sourcePort)
		})
		spInputs[i] = sourcePort
		o.AddFormItem(sourcePort)
	}

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

		sps := make([]string, spCount, spCount)
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

// as util
func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}
