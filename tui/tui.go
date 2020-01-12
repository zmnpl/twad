package tui

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/zmnpl/twad/cfg"
	"github.com/zmnpl/twad/games"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

const (
	previewBackgroundColor = tcell.ColorRoyalBlue
	accentColor            = tcell.ColorOrange

	pageStats       = "stats"
	pageNewForm     = "newform"
	pageModSelector = "modselector"
	pageSettings    = "settings"
	pageMain        = "main"
	pageHelp        = "help"
	pageLicense     = "license"

	tableBorders = false
)

var (
	config         *cfg.Cfg
	app            *tview.Application
	gamesTable     *tview.Table
	statsTable     *tview.Table
	commandPreview *tview.TextView
	actionPager    *tview.Pages
	newForm        *tview.Form
	modTree        *tview.TreeView
	licensePage    *tview.TextView

	bigMainPager *tview.Pages
)

func init() {
	config = cfg.GetInstance()
	games.RegisterChangeListener(whenGamesChanged)

	// ui style
	tview.Styles.PrimitiveBackgroundColor = tcell.ColorBlack
	tview.Styles.ContrastBackgroundColor = tcell.ColorRoyalBlue
	tview.Styles.MoreContrastBackgroundColor = tcell.ColorOrange
	tview.Styles.BorderColor = tcell.ColorRoyalBlue
	tview.Styles.TitleColor = tcell.ColorRoyalBlue
	tview.Styles.GraphicsColor = tcell.ColorRoyalBlue
	tview.Styles.PrimaryTextColor = tcell.ColorWhite
	tview.Styles.SecondaryTextColor = tcell.ColorOrange
	tview.Styles.TertiaryTextColor = tcell.ColorHotPink
	tview.Styles.InverseTextColor = tcell.ColorLemonChiffon
	tview.Styles.ContrastSecondaryTextColor = tcell.ColorPeachPuff
}

// Draw performs all necessary steps to start the ui
func Draw() {
	// init basic primitives
	app = tview.NewApplication()
	gamesTable = makeGamesTable()
	commandPreview = makeCommandPreview()
	actionPager = makeActionPager()
	selectedGameChanged(&games.Game{})
	populateGamesTable()

	// center with main content
	bigMainPager = tview.NewPages()

	// main page with games table and stats
	mainPage := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(commandPreview, 1, 0, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(gamesTable, 0, 2, true).
			AddItem(tview.NewTextView(), 2, 0, false).
			AddItem(actionPager, 0, 1, false), 0, 1, true).
		AddItem(makeButtonBar(), 1, 0, false)

	bigMainPager.AddPage(pageMain, mainPage, true, true)

	// settings - only when first start of app
	if !config.Configured {
		settingsPage := makeSettingsPage()
		bigMainPager.AddPage(pageSettings, settingsPage, true, true)
	}

	// main layout
	canvas := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(makeHeader(), 20, 0, false).
		AddItem(bigMainPager, 0, 1, true)

	// capture input
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		k := event.Key()

		if k == tcell.KeyRune {
			switch event.Rune() {
			// show credits and license
			case 'c':
				frontPage, _ := actionPager.GetFrontPage()
				if frontPage == pageLicense {
					appModeNormal()
					return nil
				}
				actionPager.SwitchToPage(pageLicense)
				app.SetFocus(licensePage)
				return nil

			case 'q':
				app.Stop()
			}

		}

		// show help at bottom of screen
		if k == tcell.KeyF1 {
			frontPage, _ := bigMainPager.GetFrontPage()
			if frontPage == pageHelp {
				appModeNormal()
				return nil
			}
			help := makeHelpPane()
			app.SetFocus(help)
			bigMainPager.AddPage(pageHelp, help, true, true)
			return nil
		}

		// switch back to nowmal mode
		if k == tcell.KeyESC {
			appModeNormal()
			return nil
		}

		return event
	})

	// run app
	if err := app.SetRoot(canvas, true).SetFocus(bigMainPager).Run(); err != nil {
		panic(err)
	}
}

// update functions
func selectedGameChanged(g *games.Game) {
	populateCommandPreview(g.String())
	populateStats(g)
}

func whenGamesChanged() {
	populateGamesTable()
}

func appModeNormal() {
	actionPager.SwitchToPage(pageStats)
	bigMainPager.SwitchToPage(pageMain)
	//bigMainPager.RemovePage(pageHelp)
	app.SetFocus(gamesTable)
}

// make ui elements

// TODO
func createActionArea() {}

// settings page
func makeSettingsPage() *tview.Flex {
	basePathPreview := tview.NewTextView()
	basePathPreview.SetBackgroundColor(previewBackgroundColor)
	fmt.Fprintf(basePathPreview, "mods path: ")
	pathSelector := makePathSelectionTree(basePathPreview)

	explanation := tview.NewTextView().SetRegions(true).SetWrap(true).SetWordWrap(true).SetDynamicColors(true)
	fmt.Fprintf(explanation, "%s\n\nExample:\n", setupPathExplain)
	fmt.Fprintf(explanation, "%s", setupPathExample)
	fmt.Fprintf(explanation, "\n\n%s", setupOkHint)

	settingsFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	settingsFlex.SetBorder(true)
	settingsFlex.SetTitle("Setup")
	settingsFlex.SetBorderColor(accentColor)
	settingsFlex.SetTitleColor(accentColor)

	settingsPage := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(settingsFlex, 78, 0, true).
		AddItem(tview.NewBox().SetBorder(false), 0, 1, false)

	settingsFlex.AddItem(explanation, 11, 0, false).
		AddItem(basePathPreview, 1, 0, false).
		AddItem(pathSelector, 0, 1, true)

	return settingsPage
}

// command preview
func makeCommandPreview() *tview.TextView {
	commandPreview = tview.NewTextView().
		SetDynamicColors(true)
	commandPreview.SetBackgroundColor(previewBackgroundColor)
	fmt.Fprintf(commandPreview, "")

	return commandPreview
}

//  stats
func makeStatsTable() *tview.Table {
	stats := tview.NewTable().
		SetFixed(1, 1).
		SetSelectable(false, false).
		SetBorders(tableBorders).SetSeparator(':')

	return stats
}

// center table with mods
func makeGamesTable() *tview.Table {
	gamesTable = tview.NewTable().
		SetFixed(1, 1).
		SetSelectable(true, false).
		SetBorders(tableBorders).SetSeparator('|')

	return gamesTable
}

// action pager, which holds stats and the "new" form
func makeActionPager() *tview.Pages {
	actionPager = tview.NewPages()
	actionPager.SetTitleAlign(tview.AlignLeft)

	statsTable = makeStatsTable()
	actionPager.AddPage(pageStats, statsTable, true, true)

	newForm = makeNewGameForm()
	actionPager.AddPage(pageNewForm, newForm, true, false)

	licensePage = makeLicense()
	actionPager.AddPage(pageLicense, licensePage, true, false)

	return actionPager
}

// tree view for selecting additional mods TODO
func makePathSelectionTree(preview *tview.TextView) *tview.TreeView {
	rootDir := "/"
	root := tview.NewTreeNode(rootDir).SetColor(tview.Styles.TitleColor)
	modFolderTree := tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root)

	// A helper function which adds the files and directories of the given path
	// to the given target node.
	add := func(target *tview.TreeNode, path string) {
		files, err := ioutil.ReadDir(path)
		sort.Slice(files, func(i, j int) bool {
			return strings.ToLower(files[i].Name()) < strings.ToLower(files[j].Name())
		})

		if err != nil {
			//panic(err)
		}
		for _, file := range files {
			if !file.IsDir() {
				continue
			}
			node := tview.NewTreeNode(file.Name()).
				SetReference(filepath.Join(path, file.Name())).
				SetSelectable(true)
			node.SetColor(tview.Styles.PrimaryTextColor)

			target.AddChild(node)
		}
	}

	// Add the current directory to the root node.
	add(root, rootDir)

	// If a directory was selected, open it.
	modFolderTree.SetSelectedFunc(func(node *tview.TreeNode) {
		reference := node.GetReference()

		if reference == nil {
			return // Selecting the root node does nothing.
		}
		children := node.GetChildren()
		if len(children) == 0 {
			// Load and show files in this directory.
			path := reference.(string)

			fi, err := os.Stat(path)
			switch {
			case err != nil:
				// handle the error and return
			case fi.IsDir():
				// it's a directory
				add(node, path)
			}
		} else {
			// Collapse if visible, expand if collapsed.
			node.SetExpanded(!node.IsExpanded())
		}
	})

	modFolderTree.SetChangedFunc(func(node *tview.TreeNode) {
		reference := node.GetReference()
		if reference == nil {
			return // Selecting the root node does nothing.
		}
		preview.Clear()
		fmt.Fprintf(preview, "mod path: %s", reference.(string))
		config.ModBasePath = reference.(string)
	})

	modFolderTree.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		k := event.Key()

		switch k {
		case tcell.KeyCtrlO:
			config.Configured = true
			cfg.AddPathToCfgs()
			err := cfg.Persist()
			if err != nil {
				// TODO - handle this
			}
			appModeNormal()
			return nil
		}

		return event
	})

	return modFolderTree
}

// tree view for selecting additional mods TODO
func makeModTree(g *games.Game) *tview.TreeView {
	rootDir := config.ModBasePath
	root := tview.NewTreeNode(rootDir).SetColor(tview.Styles.TitleColor)
	modFolderTree := tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root)
	modFolderTree.SetBorder(true)
	modFolderTree.SetTitle("Add Mod To Game")

	// A helper function which adds the files and directories of the given path
	// to the given target node.
	add := func(target *tview.TreeNode, path string) {
		files, err := ioutil.ReadDir(path)
		files = filterExtensions(files)

		sort.Slice(files, func(i, j int) bool {
			return strings.ToLower(files[i].Name()) < strings.ToLower(files[j].Name())
		})

		if err != nil {
			panic(err)
		}
		for _, file := range files {
			node := tview.NewTreeNode(file.Name()).
				SetReference(filepath.Join(path, file.Name())).
				SetSelectable(true)
			node.SetColor(tview.Styles.SecondaryTextColor)
			if file.IsDir() {
				node.SetColor(tview.Styles.PrimaryTextColor)
			}
			target.AddChild(node)
		}
	}

	// Add the current directory to the root node.
	add(root, rootDir)

	// If a directory was selected, open it.
	modFolderTree.SetSelectedFunc(func(node *tview.TreeNode) {
		reference := node.GetReference()

		if reference == nil {
			return // Selecting the root node does nothing.
		}
		children := node.GetChildren()
		if len(children) == 0 {
			// Load and show files in this directory.
			path := reference.(string)

			fi, err := os.Stat(path)
			switch {
			case err != nil:
				// handle the error and return
			case fi.IsDir():
				// it's a directory
				add(node, path)
			default:
				// it's not a directory
				g.AddMod(strings.TrimPrefix(path, config.ModBasePath+"/"))
			}
		} else {
			// Collapse if visible, expand if collapsed.
			node.SetExpanded(!node.IsExpanded())
		}
	})

	return modFolderTree
}

func makeNewGameForm() *tview.Form {
	newForm := tview.NewForm().
		AddInputField(inputLabelName, "", 20, nil, nil).
		AddDropDown(inputLabelSourcePort, config.SourcePorts, 0, nil).
		AddDropDown(inputLabelIWad, config.IWADs, 0, nil).
		AddButton("Add", func() {
			name := newForm.GetFormItemByLabel(inputLabelName).(*tview.InputField).GetText()
			_, sourceport := newForm.GetFormItemByLabel(inputLabelSourcePort).(*tview.DropDown).GetCurrentOption()
			_, wad := newForm.GetFormItemByLabel(inputLabelIWad).(*tview.DropDown).GetCurrentOption()

			games.AddGame(games.NewGame(name, sourceport, wad))
		})

	newForm.SetBorder(true).SetTitle("Add new game").SetTitleAlign(tview.AlignCenter)

	return newForm
}

// populate ui elements

func populateGamesTable() {
	gamesTable.Clear()
	allGames := games.GetInstance()

	rows, cols := len(allGames), games.MaxModCount()-1
	fixRows, fixCols := 1, 3

	for r := 0; r < rows+fixRows; r++ {
		var game games.Game
		if r > 0 {
			game = allGames[r-fixRows]
		}
		for c := 0; c < cols+fixCols; c++ {
			var cell *tview.TableCell

			if r < 1 {
				switch c {
				case 0:
					cell = tview.NewTableCell("Name").SetTextColor(tview.Styles.SecondaryTextColor)
				case 1:
					cell = tview.NewTableCell("Source Port").SetTextColor(tview.Styles.SecondaryTextColor)
				case 2:
					cell = tview.NewTableCell("Iwad").SetTextColor(tview.Styles.SecondaryTextColor)
				case 3:
					cell = tview.NewTableCell("Mods / Parameters").SetTextColor(tview.Styles.SecondaryTextColor)
				default:
					cell = tview.NewTableCell("").SetTextColor(tview.Styles.SecondaryTextColor)
				}
			} else {
				switch c {
				case 0:
					cell = tview.NewTableCell(game.Name).SetTextColor(tview.Styles.SecondaryTextColor)
				case 1:
					cell = tview.NewTableCell(game.SourcePort).SetTextColor(tview.Styles.PrimaryTextColor)
				case 2:
					cell = tview.NewTableCell(game.Iwad).SetTextColor(tview.Styles.PrimaryTextColor)
				default:
					i := c - fixCols
					if i < len(game.Mods) {
						cell = tview.NewTableCell(game.Mods[i]).SetTextColor(tview.Styles.PrimaryTextColor)
					} else {
						cell = tview.NewTableCell("").SetTextColor(tview.Styles.PrimaryTextColor)
					}
				}
			}
			gamesTable.SetCell(r, c, cell)
		}
	}

	makeModTreeMaker := func(selectedGame *games.Game) func() *tview.TreeView {
		return func() *tview.TreeView {
			return makeModTree(selectedGame)
		}
	}
	modTreeMaker := makeModTreeMaker(&games.Game{})

	//makeCellPulser := func(c *tview.TableCell) func() {
	//	return func() {
	//		r, g, b := tcell.ColorOrange.RGB()
	//		c.BackgroundColor = tcell.NewRGBColor(r, g, b)
	//		for i := 0; i <= 1000; i += 16 {
	//			time.Sleep(16 * time.Millisecond)
	//			r = r * 2
	//			g = g * 2
	//			b = b * 2
	//			c.BackgroundColor = tcell.NewRGBColor(r, g, b)
	//			app.Draw()
	//		}
	//	}
	//}
	//cellPulser := makeCellPulser(tview.NewTableCell(""))

	gamesTable.SetSelectionChangedFunc(func(r int, c int) {
		var g *games.Game
		//var cell *tview.TableCell
		switch r {
		case 0:
			g = &games.Game{}
			//cell = tview.NewTableCell("")
		default:
			g = &allGames[r-fixRows]
			//cell = gamesTable.GetCell(r, len(g.Mods)+3)
		}
		selectedGameChanged(g)
		modTreeMaker = makeModTreeMaker(g)
		//cellPulser = makeCellPulser(cell)
	})

	gamesTable.SetSelectedFunc(func(r int, c int) {
		switch {
		case r > 0:
			allGames[r-fixRows].Run()
		}
	})

	gamesTable.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		k := event.Key()
		r, _ := gamesTable.GetSelection()

		if k == tcell.KeyRune {
			switch event.Rune() {
			case 'a':
				modTree := modTreeMaker()
				actionPager.AddPage(pageModSelector, modTree, true, false)
				actionPager.SwitchToPage(pageModSelector)
				app.SetFocus(modTree)
				return nil
			case 'r':
				if r > 0 {
					mods := allGames[r-fixRows].Mods
					if len(mods) > 0 {
						allGames[r-fixRows].Mods = mods[:len(mods)-1]
						populateGamesTable()
						games.Persist()
					}
				}
				return nil
			// open dialog to insert new game
			case 'i':
				actionPager.SwitchToPage(pageNewForm)
				app.SetFocus(newForm)
				return nil
			}
		}

		if k == tcell.KeyDelete && r > 0 {
			if r == gamesTable.GetRowCount()-1 {
				gamesTable.Select(r-1, 0)
			}
			games.RemoveGameAt(r - fixRows)
			return nil
		}

		return event
	})
}

func populateCommandPreview(command string) {
	commandPreview.Clear()
	fmt.Fprintf(commandPreview, "preview $ %s", command)
}

func populateStats(g *games.Game) {
	statsTable.Clear()
	row := 0
	pts := float64(g.Playtime) / 1000 / 60
	saves := 0

	statsTable.SetCell(row, 0, tview.NewTableCell("# Savegames").SetTextColor(tview.Styles.SecondaryTextColor))
	statsTable.SetCell(row, 1, tview.NewTableCell(fmt.Sprintf("%v", saves)).SetAlign(tview.AlignLeft))
	row++
	statsTable.SetCell(row, 0, tview.NewTableCell("Playtime").SetTextColor(tview.Styles.SecondaryTextColor))
	statsTable.SetCell(row, 1, tview.NewTableCell(fmt.Sprintf("%.2f min", pts)).SetAlign(tview.AlignLeft))
	row++
	statsTable.SetCell(row, 0, tview.NewTableCell("Last Played").SetTextColor(tview.Styles.SecondaryTextColor))
	statsTable.SetCell(row, 1, tview.NewTableCell(fmt.Sprint(g.LastPlayed)).SetAlign(tview.AlignLeft))
	row++
	statsTable.SetCell(row, 0, tview.NewTableCell(""))
	//	statsTable.SetCell(row, 1, tview.NewTableCell(""))
	row++

	for k, v := range g.Stats {
		statsTable.SetCell(row, 0, tview.NewTableCell(strings.Title("# "+k)).SetTextColor(tview.Styles.SecondaryTextColor))
		statsTable.SetCell(row, 1, tview.NewTableCell(fmt.Sprintf("%v", v)).SetAlign(tview.AlignLeft))
		row++
	}

	statsTable.SetCell(row, 0, tview.NewTableCell("                    ").SetTextColor(tview.Styles.SecondaryTextColor))
	statsTable.SetCell(row, 1, tview.NewTableCell("                    ").SetAlign(tview.AlignLeft))
}

// helper functions

func filterExtensions(files []os.FileInfo) []os.FileInfo {
	tmp := files
	files = files[:0]
	for _, v := range tmp {
		ext := strings.ToLower(filepath.Ext(v.Name()))
		if _, found := config.ModExtensions[ext]; found || v.IsDir() {
			files = append(files, v)
		}
	}
	return files
}
