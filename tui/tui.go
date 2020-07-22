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
	config              *cfg.Cfg
	app                 *tview.Application
	gamesTable          *tview.Table
	commandPreview      *tview.TextView
	mainApplicationPage *tview.Flex
	rightSidePager      *tview.Pages
	rightSidePagerSub1  *tview.Pages
	rightSidePagerSub2  *tview.Pages
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
	initUIElements()

	// settings - only when first start of app
	if !config.Configured {
		settingsPage := makeSettingsPage()
		bigMainPager.AddPage(pageSettings, settingsPage, true, true)
	}

	// main layout
	header, headerHeight := getHeader()
	helpPane, helpPaneHeight := makeHelpPane()
	canvas := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(header, headerHeight, 0, false).
		AddItem(bigMainPager, 0, 1, true).
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

		// runes are not really to handle at app leve -.-
		if k == tcell.KeyRune {
			//switch event.Rune() {
			//case 'q': // get out
			//app.Stop()
			//}
		}

		return event
	})

	// run app
	if err := app.SetRoot(canvas, true).SetFocus(bigMainPager).Run(); err != nil {
		panic(err)
	}
}

func initUIElements() {
	// init basic primitives
	app = tview.NewApplication()
	gamesTable = makeGamesTable()
	commandPreview = makeCommandPreview()

	// main page containing all the content
	mainApplicationPage = tview.NewFlex().SetDirection(tview.FlexRow)

	// right side
	rightSidePager = tview.NewPages()
	rightSidePagerSub1 = tview.NewPages()
	rightSidePagerSub2 = tview.NewPages()

	// TODO: make layout a bit more flexible
	defaultRightPage := tview.NewFlex().SetDirection(tview.FlexColumn)
	defaultRightPage.
		AddItem(tview.NewTextView().SetBackgroundColor(tview.Styles.PrimaryTextColor), 2, 0, false).
		AddItem(rightSidePagerSub1, 0, 5, false).
		AddItem(tview.NewTextView(), 2, 0, false).
		AddItem(rightSidePagerSub2, 0, 5, false)
	rightSidePager.AddPage(pageDefaultRight, defaultRightPage, true, true)

	mainApplicationPage.AddItem(commandPreview, 1, 0, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(gamesTable, 0, config.GameListRelativeWidth, true).
			AddItem(rightSidePager, 0, 10-config.GameListRelativeWidth, false), 0, 2, true)

	// center with main content
	bigMainPager = tview.NewPages()
	bigMainPager.AddPage(pageMain, mainApplicationPage, true, true)
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
	rightSidePagerSub1.AddPage(pageMods, makeModList(g), true, true)
	frontPage, _ := rightSidePagerSub2.GetFrontPage()
	if frontPage != pageModSelector {
		rightSidePagerSub2.AddPage(pageStats, makeStatsTable(g), true, true)
	}
}

// redraw whole table
func whenGamesChanged() {
	populateGamesTable()
}

// reset ui
func appModeNormal() {
	rightSidePager.SwitchToPage(pageDefaultRight)
	rightSidePagerSub1.SwitchToPage(pageMods)
	rightSidePagerSub2.SwitchToPage(pageStats)
	bigMainPager.SwitchToPage(pageMain)

	// clear bigMainPager
	if bigMainPager.HasPage(pageYouSure) {
		bigMainPager.RemovePage(pageYouSure)
	}
	if bigMainPager.HasPage(pageSettings) {
		bigMainPager.RemovePage(pageSettings)
	}
	if bigMainPager.HasPage(pageOptions) {
		bigMainPager.RemovePage(pageOptions)
	}

	// clear actionPager
	if rightSidePager.HasPage(pageAddEdit) {
		rightSidePager.RemovePage(pageAddEdit)
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
