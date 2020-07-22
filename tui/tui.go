package tui

import (
	"github.com/zmnpl/twad/cfg"
	"github.com/zmnpl/twad/games"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

const (
	previewBackgroundColor = tcell.ColorRoyalBlue
	accentColor            = tcell.ColorOrange

	pageOptions      = "options"
	pageStats        = "stats"
	pageNewForm      = "newform"
	pageModSelector  = "modselector"
	pageSettings     = "settings"
	pageMain         = "main"
	pageHelp         = "help"
	pageLicense      = "license"
	pageYouSure      = "yousure"
	pageParamsEdit   = "paramseditor"
	pageGameOverview = "gameoverview"
	pageMods         = "mods"

	tableBorders = false
)

var (
	config         *cfg.Cfg
	app            *tview.Application
	gamesTable     *tview.Table
	commandPreview *tview.TextView
	actionPager    *tview.Pages
	modPager       *tview.Pages
	modTree        *tview.TreeView
	licensePage    *tview.TextView

	bigMainPager *tview.Pages
)

func init() {
	config = cfg.GetInstance()
	games.RegisterChangeListener(whenGamesChanged)

	// ui stylepageSettings
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
	modPager = makeModListPager()
	selectedGameChanged(&games.Game{})
	populateGamesTable()

	// center with main content
	bigMainPager = tview.NewPages()

	// main page with games table and stats
	mainPage := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(commandPreview, 1, 0, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(gamesTable, 0, 5, true).
			AddItem(tview.NewTextView(), 2, 0, false).
			AddItem(modPager, 0, 2, false).
			AddItem(tview.NewTextView(), 2, 0, false).
			AddItem(actionPager, 0, 3, false), 0, 2, true)

	bigMainPager.AddPage(pageMain, mainPage, true, true)

	// settings - only when first start of app
	if !config.Configured {
		settingsPage := makeSettingsPage()
		bigMainPager.AddPage(pageSettings, settingsPage, true, true)
	}

	// main layout
	headerHeight := 20
	var header tview.Primitive
	header = makeHeader()
	if !cfg.GetInstance().PrintHeader {
		headerHeight = 1
		header = tview.NewTextView().SetDynamicColors(true).SetText(subtitle)
	}
	canvas := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(header, headerHeight, 0, false).
		AddItem(bigMainPager, 0, 1, true).
		AddItem(makeHelpPane(), 5, 0, false)

	// capture input
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		k := event.Key()

		if k == tcell.KeyRune {
			switch event.Rune() {
			case 'q':
				app.Stop()
			}
		}

		// show help at bottom of screen
		// if k == tcell.KeyF1 {
		// 	frontPage, _ := bigMainPager.GetFrontPage()
		// 	if frontPage == pageHelp {
		// 		appModeNormal()
		// 		return nil
		// 	}
		// 	help := makeHelpPane()
		// 	app.SetFocus(help)
		// 	bigMainPager.AddPage(pageHelp, help, true, true)
		// 	return nil
		// }

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
	actionPager.AddPage(pageStats, makeStatsTable(g), true, true)
	modPager.AddPage(pageMods, makeModList(g), true, true)
	if actionPager.HasPage(pageGameOverview) {
		actionPager.AddPage(pageGameOverview, makeModList(g), true, true)
	}
}

func whenGamesChanged() {
	populateGamesTable()
}

func appModeNormal() {
	actionPager.SwitchToPage(pageStats)
	bigMainPager.SwitchToPage(pageMain)

	// clear bigMainPager
	if bigMainPager.HasPage(pageYouSure) {
		bigMainPager.RemovePage(pageYouSure)
	}
	if bigMainPager.HasPage(pageHelp) {
		bigMainPager.RemovePage(pageHelp)
	}
	if bigMainPager.HasPage(pageSettings) {
		bigMainPager.RemovePage(pageSettings)
	}
	if bigMainPager.HasPage(pageOptions) {
		bigMainPager.RemovePage(pageOptions)
	}

	// clear actionPager
	if actionPager.HasPage(pageNewForm) {
		actionPager.RemovePage(pageNewForm)
	}
	if actionPager.HasPage(pageParamsEdit) {
		actionPager.RemovePage(pageParamsEdit)
	}
	if actionPager.HasPage(pageGameOverview) {
		actionPager.RemovePage(pageGameOverview)
	}
	app.SetFocus(gamesTable)
}

func tabNavigate(previous, next tview.Primitive) func(event *tcell.EventKey) *tcell.EventKey {
	return func(event *tcell.EventKey) *tcell.EventKey {
		k := event.Key()
		switch k {
		case tcell.KeyTab:
			app.SetFocus(next)
			return nil
		case tcell.KeyBacktab:
			app.SetFocus(previous)
			return nil
		}

		return event
	}
}
