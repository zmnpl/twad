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
	pageAddEdit      = "addEdit"
	pageModSelector  = "modselector"
	pageSettings     = "settings"
	pageMain         = "main"
	pageHelp         = "help"
	pageLicense      = "license"
	pageYouSure      = "yousure"
	pageMods         = "mods"
	pageDefaultRight = "right"

	tableBorders = false
)

var (
	config             *cfg.Cfg
	app                *tview.Application
	mainContentPage    *tview.Flex
	contentPages       *tview.Pages
	gamesTable         *tview.Table
	commandPreview     *tview.TextView
	rightSidePages     *tview.Pages
	rightSidePagesSub1 *tview.Pages
	rightSidePagesSub2 *tview.Pages
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
	initUIElements()

	// settings - only when first start of app
	if !config.Configured {
		settingsPage := makeSettingsPage()
		contentPages.AddPage(pageSettings, settingsPage, true, true)
	}

	// main layout
	header, headerHeight := getHeader()
	helpPane, helpPaneHeight := makeHelpPane()
	canvas := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(header, headerHeight, 0, false).
		AddItem(contentPages, 0, 1, true).
		AddItem(helpPane, helpPaneHeight, 0, false)

	// populate
	selectedGameChanged(&games.Game{})
	populateGamesTable()

	// capture input
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		k := event.Key()
		// switch back to nowmal mode
		if k == tcell.KeyESC {
			appModeNormal()
			return nil
		}
		return event
	})

	// run app
	if err := app.SetRoot(canvas, true).SetFocus(contentPages).Run(); err != nil {
		panic(err)
	}
}

func initUIElements() {
	// init basic primitives
	app = tview.NewApplication()

	// command preview
	commandPreview = makeCommandPreview()

	// main view to select games
	gamesTable = makeGamesTable()

	// main page containing all the content
	mainContentPage = tview.NewFlex().SetDirection(tview.FlexRow)

	// right side
	rightSidePages = tview.NewPages()
	rightSidePagesSub1 = tview.NewPages()
	rightSidePagesSub2 = tview.NewPages()

	// TODO: make layout a bit more flexible
	defaultRightPage := tview.NewFlex().SetDirection(tview.FlexColumn)
	defaultRightPage.
		AddItem(tview.NewTextView().SetBackgroundColor(tview.Styles.PrimaryTextColor), 2, 0, false).
		AddItem(rightSidePagesSub1, 0, 5, false).
		AddItem(tview.NewTextView(), 2, 0, false).
		AddItem(rightSidePagesSub2, 0, 5, false)
	rightSidePages.AddPage(pageDefaultRight, defaultRightPage, true, true)

	mainContentPage.AddItem(commandPreview, 1, 0, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(gamesTable, 0, config.GameListRelativeWidth, true).
			AddItem(rightSidePages, 0, 10-config.GameListRelativeWidth, false), 0, 2, true)

	// center with main content
	contentPages = tview.NewPages()
	contentPages.AddPage(pageMain, mainContentPage, true, true)
}

// small or big header
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
	rightSidePagesSub1.AddPage(pageMods, makeModList(g), true, true)
	frontPage, _ := rightSidePagesSub2.GetFrontPage()
	if frontPage != pageModSelector {
		rightSidePagesSub2.AddPage(pageStats, makeStatsTable(g), true, true)
	}
}

// redraw whole table
func whenGamesChanged() {
	populateGamesTable()
}

// reset ui
func appModeNormal() {
	rightSidePages.SwitchToPage(pageDefaultRight)
	rightSidePagesSub1.SwitchToPage(pageMods)
	rightSidePagesSub2.SwitchToPage(pageStats)
	contentPages.SwitchToPage(pageMain)

	// clear bigMainPager
	if contentPages.HasPage(pageYouSure) {
		contentPages.RemovePage(pageYouSure)
	}
	if contentPages.HasPage(pageSettings) {
		contentPages.RemovePage(pageSettings)
	}
	if contentPages.HasPage(pageOptions) {
		contentPages.RemovePage(pageOptions)
	}

	// clear actionPager
	if rightSidePages.HasPage(pageAddEdit) {
		rightSidePages.RemovePage(pageAddEdit)
	}

	app.SetFocus(gamesTable)
}

// used in options and such screens
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
