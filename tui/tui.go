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
	config              *cfg.Cfg
	app                 *tview.Application
	gamesTable          *tview.Table
	commandPreview      *tview.TextView
	mainApplicationPage *tview.Flex
	actionPager         *tview.Pages
	actionPagerSub1     *tview.Pages
	actionPagerSub2     *tview.Pages
	modTree             *tview.TreeView
	licensePage         *tview.TextView

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
	initUiElements()

	mainApplicationPage.AddItem(commandPreview, 1, 0, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(gamesTable, 0, config.GameListRelativeWidth, true).
			AddItem(tview.NewTextView(), 2, 0, false).
			AddItem(actionPagerSub1, 0, 2, false).
			AddItem(tview.NewTextView(), 2, 0, false).
			AddItem(actionPager, 0, 3, false), 0, 2, true)

	bigMainPager.AddPage(pageMain, mainApplicationPage, true, true)

	// settings - only when first start of app
	if !config.Configured {
		settingsPage := makeSettingsPage()
		bigMainPager.AddPage(pageSettings, settingsPage, true, true)
	}

	// main layout
	header, headerHeight := getHeader()
	canvas := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(header, headerHeight, 0, false).
		AddItem(bigMainPager, 0, 1, true).
		AddItem(makeHelpPane(), 5, 0, false)

	// capture input
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		k := event.Key()

		if k == tcell.KeyRune {
			switch event.Rune() {
			// get out
			case 'q':
				app.Stop()
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
			}

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

func initUiElements() {
	// init basic primitives
	app = tview.NewApplication()
	gamesTable = makeGamesTable()
	commandPreview = makeCommandPreview()
	actionPager = makeActionPager()
	actionPagerSub1 = makeModListPager()
	selectedGameChanged(&games.Game{})
	populateGamesTable()

	// center with main content
	bigMainPager = tview.NewPages()

	// main page containing all the content
	mainApplicationPage = tview.NewFlex().SetDirection(tview.FlexRow)

	// right side
	actionPager = tview.NewPages()
	actionPagerSub1 = tview.NewPages()
	actionPagerSub2 = tview.NewPages()
}

func getHeader() (tview.Primitive, int) {
	headerHeight := 20
	var header tview.Primitive
	header = makeHeader()
	if !cfg.GetInstance().PrintHeader {
		headerHeight = 1
		header = tview.NewTextView().SetDynamicColors(true).SetText(subtitle)
	}
	return header, headerHeight
}

// update functions
func selectedGameChanged(g *games.Game) {
	populateCommandPreview(g.String())
	actionPager.AddPage(pageStats, makeStatsTable(g), true, true)
	actionPagerSub1.AddPage(pageMods, makeModList(g), true, true)
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
