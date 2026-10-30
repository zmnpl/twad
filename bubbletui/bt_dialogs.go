package bubbletui

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"slices"
	"sort"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/zmnpl/twad/base"
	"github.com/zmnpl/twad/games"
	"github.com/zmnpl/twad/helper"
	"github.com/zmnpl/twad/ports"
	"github.com/zmnpl/twad/tui/shared"
)

// confirm asks the user to confirm an action; either choice closes the dialog
func (m *btModel) confirm(title string, onOk func() tea.Cmd, atGame bool) {
	f := &btForm{
		title: title,
		buttons: []btButton{
			{shared.Dict.ConfirmText, func(m *btModel) tea.Cmd {
				m.closeDialog()
				return onOk()
			}},
			{shared.Dict.AbortText, func(m *btModel) tea.Cmd {
				m.closeDialog()
				return nil
			}},
		},
	}
	m.showDialog(f.init(1), 64, atGame)
}

// showError shows errors prominently on the screen
// if a function is supplied, the user gets the choice to proceed
// without fixing what might have cause this on his/her own risk
func (m *btModel) showError(errTitle, errText string, yolo func() tea.Cmd) {
	closeErr := func(m *btModel) tea.Cmd {
		m.closeDialog()
		return nil
	}

	f := &btForm{
		title: shared.Dict.ErrTitleStart + " - " + errTitle,
		text:  m.st.primary.Render(errText),
		warn:  true,
	}
	if yolo != nil {
		f.buttons = []btButton{
			{shared.Dict.ErrYolo, func(m *btModel) tea.Cmd {
				m.closeDialog()
				return yolo()
			}},
			{shared.Dict.ErrNotYet, closeErr},
		}
	} else {
		f.buttons = []btButton{{shared.Dict.ErrAbort, closeErr}}
	}
	m.showDialog(f.init(0), 0, false)
}

// warp
//#######################################################################

func (m *btModel) showWarp(g *games.Game) tea.Cmd {
	warpTo := m.newInput(shared.Dict.WarpText, "")
	warpTo.fieldWidth = 6
	// check every rune; pasted text arrives at once
	warpTo.accept = func(text string) bool {
		for _, r := range text {
			if !shared.WarpStringAcceptance(text, r) {
				return false
			}
		}
		return true
	}

	// map select
	maps := g.ModMaps()
	displayMaps := make([]string, 0, len(maps)+1)
	displayMaps = append(displayMaps, "")
	for k := range maps {
		displayMaps = append(displayMaps, k)
	}
	sort.Strings(displayMaps)
	mapSelect := m.newDropDown("maps", displayMaps, 0)

	skill := m.newDropDown(shared.Dict.SkillText, shared.SkillLevels, 2)

	// to record a demo, specify a name
	demoName := m.newInput(shared.Dict.DemoText, "")
	demoName.fieldWidth = 22
	demoName.label = func(st *btStyles) string {
		if demoName.value() != "" && g.DemoExists(demoName.value()) {
			return st.warn.Render(shared.Dict.DemoTextOverwrite)
		}
		return st.primary.Render(shared.Dict.DemoText)
	}

	f := &btForm{
		title:  g.Name,
		fields: []*btField{warpTo, mapSelect, skill, demoName},
		buttons: []btButton{{shared.Dict.WarpOkButton, func(m *btModel) tea.Cmd {
			episode, level := shared.SplitWarpString(warpTo.value())
			difficulty := skill.current
			demo := demoName.value()
			m.reset()

			return m.launch(g, func() {
				if mapSelect.current > 0 {
					mapName := maps[mapSelect.value()]
					if len(demo) > 0 {
						g.GoToMapRecord(mapName, difficulty, demo)
						return
					}
					g.GoToMap(mapName, difficulty)
					return
				}

				// supplying a demoname automatically starts recording
				if len(demo) > 0 {
					g.WarpRecord(episode, level, difficulty, demo)
					return
				}
				g.Warp(episode, level, difficulty)
			})
		}}},
	}
	m.showDialog(f.init(0), 62, true)
	return nil
}

// zip import
//#######################################################################

func (m *btModel) showZipImport(archivePath, folderNameSuggest string) {
	name := strings.TrimSuffix(path.Base(archivePath), path.Ext(archivePath))
	if folderNameSuggest != "" {
		name = folderNameSuggest
	}

	modName := m.newInput(shared.Dict.ZipImportToLabel, name)
	modName.check = func(v string) (string, bool) {
		if !helper.IsFileNameValid(v) {
			return shared.Dict.ZipImportToBadNameLabel, false
		}
		if _, err := os.Stat(path.Join(m.cfg.WadDir, v)); !os.IsNotExist(err) {
			return shared.Dict.ZipImportToExistsLabel, false
		}
		return "", true
	}

	f := &btForm{
		title:  shared.Dict.ZipImportTitle + " " + filepath.Base(archivePath),
		text:   m.st.warn.Render(shared.Dict.ZipImportSecurityWarn),
		fields: []*btField{modName},
		buttons: []btButton{
			{shared.Dict.ZipImportCancel, func(m *btModel) tea.Cmd {
				m.closeDialog()
				return nil
			}},
			{shared.Dict.ZipImportFormOk, func(m *btModel) tea.Cmd {
				name := modName.value()
				m.closeDialog()

				// test file name again
				if !helper.IsFileNameValid(name) {
					m.showError(shared.Dict.ZipImportNameInvalid, shared.Dict.ZipImportNameInvalidReasons, nil)
					return nil
				}
				// test if provided zip exists
				if _, err := os.Stat(archivePath); os.IsNotExist(err) {
					m.showError(shared.Dict.ZipImportArchiveNotFound, err.Error(), nil)
					return nil
				}

				m.status = " Importing " + filepath.Base(archivePath) + " ..."
				return func() tea.Msg {
					return btImportDoneMsg{base.ImportArchive(archivePath, name)}
				}
			}},
		},
	}
	m.showDialog(f.init(0), len(shared.Dict.ZipImportSecurityWarn)+6, false)
}

// add / edit game
//#######################################################################

func (m *btModel) newAddEditGame(g *games.Game) *btForm {
	isNew := g == nil
	title := shared.Dict.EditGame
	if isNew {
		port, iwad := "", ""
		if len(m.cfg.Ports) > 0 {
			port = m.cfg.Ports[0]
		}
		if len(m.cfg.IWADs) > 0 {
			iwad = m.cfg.IWADs[0]
		}
		newGame := games.NewGame("", port, "", iwad)
		g = &newGame
		title = shared.Dict.AddGame
	}

	// options from config; the games current value, when nothing is configured
	optionsOf := func(configured []string, current string) ([]string, int) {
		if len(configured) == 0 {
			return []string{current}, 0
		}
		return configured, max(0, slices.Index(configured, current))
	}

	expectedExtension := ports.ConfigFileExtension(g.Port)
	sharedCfgs := base.GetSharedGameConfigs(g.Port)

	name := m.newInput(shared.Dict.AeName, g.Name)
	portOptions, portCurrent := optionsOf(m.cfg.Ports, g.Port)
	sourcePort := m.newDropDown(shared.Dict.AeSourcePort, portOptions, portCurrent)
	iwadOptions, iwadCurrent := optionsOf(m.cfg.IWADs, g.Iwad)
	iwad := m.newDropDown(shared.Dict.AeIWAD, iwadOptions, iwadCurrent)
	ownCfg := m.newCheckbox(shared.Dict.AeOwnCfg, g.PersonalPortCfg)

	// shared config label shows the expected extension; warns when an own config is used instead
	sharedCfg := m.newInput("", g.SharedConfig)
	sharedCfg.label = func(st *btStyles) string {
		label := fmt.Sprintf(shared.Dict.AeSharedCfgT, expectedExtension)
		if ownCfg.checked {
			return st.warn.Render(label)
		}
		return st.primary.Render(label)
	}
	sharedCfg.suggest = func(text string) []string {
		if text == "" {
			return nil
		}
		return sharedCfgs
	}
	sharedCfg.done = func(f *btField) {
		if v := f.value(); v != "" && !strings.HasSuffix(strings.ToLower(v), expectedExtension) {
			f.input.SetValue(v + expectedExtension)
		}
	}

	// get shared configs for selected port
	sourcePort.changed = func() {
		sharedCfgs = base.GetSharedGameConfigs(sourcePort.value())
		expectedExtension = ports.ConfigFileExtension(sourcePort.value())
	}

	url := m.newInput(shared.Dict.AeLink, g.Link)
	envVars := m.newInput(shared.Dict.AeEnvironment, g.EnvironmentString())
	customParams := m.newInput(shared.Dict.AeOtherParams, g.ParamsString())

	f := &btForm{
		title:  title,
		fields: []*btField{name, sourcePort, iwad, ownCfg, sharedCfg, url, envVars, customParams},
		buttons: []btButton{{shared.Dict.AeOkButton, func(m *btModel) tea.Cmd {
			g.Name = strings.TrimSpace(name.value())
			g.Port = sourcePort.value()
			g.Iwad = iwad.value()
			g.PersonalPortCfg = ownCfg.checked
			g.SharedConfig = sharedCfg.value()
			g.Environment = shared.SplitParams(envVars.value())
			g.CustomParameters = shared.SplitParams(customParams.value())
			g.Link = url.value()

			if isNew {
				games.AddGame(*g)
				m.cursor = games.GameCount() - 1
			}
			if err := games.Persist(); err != nil {
				m.reset()
				m.showError("Saving games failed", err.Error(), nil)
				return nil
			}
			m.reset()
			return nil
		}}},
		footer: shared.Dict.AeEnvironmentDetail + "\n\n" + shared.Dict.AeOtherParamsDetail,
	}
	return f.init(0)
}

// options
//#######################################################################

// builder function that creates a function which suggests paths for input fields
func btPathSuggestions(dirsOnly bool, extensionFilter map[string]bool) func(text string) []string {
	cache := make(map[string][]string)

	return func(text string) []string {
		if strings.TrimSpace(text) == "" {
			return nil
		}

		dir := filepath.Dir(text)
		if entries, ok := cache[dir]; ok {
			return entries
		}

		files, err := os.ReadDir(dir)
		if err != nil {
			return nil
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
			if len(extensionFilter) > 0 && !file.IsDir() && !extensionFilter[strings.ToLower(filepath.Ext(file.Name()))] {
				continue
			}
			entries = append(entries, filepath.Join(dir, file.Name()))
		}
		cache[dir] = entries
		return entries
	}
}

func btSourcePortCheck(sp string) (string, bool) {
	sp = strings.TrimSpace(sp)
	if sp == "" {
		return "", true
	}
	if _, err := os.Stat(sp); err == nil || shared.CommandExists(sp) {
		return shared.Dict.OptsLooksGood, true
	}
	return shared.Dict.OptsErrPathDoesntExist, false
}

func btWadDirCheck(dir string) (string, bool) {
	// does this path exist?
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return shared.Dict.OptsErrPathDoesntExist, false
	}
	// check if selected path contains any iwads
	if hasIwad, err := base.PathHasIwads(dir); !hasIwad {
		if err != nil {
			return "(" + err.Error() + ")", false
		}
		return shared.Dict.OptsErrPathNoIWads, false
	}
	return shared.Dict.OptsLooksGood, true
}

func (m *btModel) newOptions() *btForm {
	c := m.cfg

	wadDir := m.newInput(shared.Dict.OptsPathLabel, c.WadDir)
	wadDir.suggest = btPathSuggestions(true, nil)
	wadDir.check = btWadDirCheck

	iwads := m.newInput(shared.Dict.OptsIwadsLabel, strings.Join(c.IWADs, ","))

	fields := []*btField{wadDir, iwads}

	// source ports
	// only suggest paths on windows; on linux the simple executable name is usually enough when it is in path
	var spSuggest func(string) []string
	if runtime.GOOS == "windows" {
		spSuggest = btPathSuggestions(false, map[string]bool{".exe": true})
	}
	spInputs := make([]*btField, base.MAX_SOURCE_PORTS)
	for i := range spInputs {
		value := ""
		if i < len(c.Ports) {
			value = c.Ports[i]
		}
		sp := m.newInput(shared.Dict.OptsSourcePortLabel, value)
		sp.suggest = spSuggest
		sp.check = btSourcePortCheck
		spInputs[i] = sp
		fields = append(fields, sp)
	}

	// ui options
	termColors := m.newCheckbox(shared.Dict.OptsUseTerminalColors, c.UseTerminalColors)
	dontWarn := m.newCheckbox(shared.Dict.OptsDontWarn, c.DeleteWithoutWarning)
	hideHeader := m.newCheckbox(shared.Dict.OptsHideHeader, c.HideHeader)
	gameListRelWidth := m.newInput(shared.Dict.OptsGamesListRelativeWitdh, strconv.Itoa(c.GameListRelativeWidth))
	gameListRelWidth.fieldWidth = 4
	gameListRelWidth.accept = func(text string) bool {
		if text == "" {
			return true
		}
		i, err := strconv.Atoi(text)
		return err == nil && i > 0 && i <= 100
	}
	fields = append(fields, termColors, dontWarn, hideHeader, gameListRelWidth)

	f := &btForm{
		title:  shared.Dict.OptsHeader,
		width:  90,
		fields: fields,
		buttons: []btButton{{shared.Dict.OptsOkButtonLabel, func(m *btModel) tea.Cmd {
			c.WadDir = wadDir.value()

			c.Ports = make([]string, 0, base.MAX_SOURCE_PORTS)
			for _, sp := range spInputs {
				if v := strings.TrimSpace(sp.value()); v != "" {
					c.Ports = append(c.Ports, v)
				}
			}

			c.IWADs = nil
			for _, iwad := range strings.Split(iwads.value(), ",") {
				if iwad = strings.TrimSpace(iwad); iwad != "" {
					c.IWADs = append(c.IWADs, iwad)
				}
			}

			c.HideHeader = hideHeader.checked
			c.DeleteWithoutWarning = dontWarn.checked
			if w, err := strconv.Atoi(gameListRelWidth.value()); err == nil {
				c.GameListRelativeWidth = w
			}
			c.UseTerminalColors = termColors.checked

			base.EnableBasePath()
			m.applyTheme()
			m.reset()
			if err := base.Persist(); err != nil {
				m.showError("Saving options failed", err.Error(), nil)
			}
			return nil
		}}},
	}
	return f.init(0)
}
