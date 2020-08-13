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
	pageFirstSetup   = "firstsetup"
	pageHeader       = "header"
	pageMain         = "main"
	pageDetailGrid   = "detailgrid"
	pageContent      = "content"
	pageContentMain  = "maincontent"
	pageHelp         = "help"
	pageLicense      = "license"
	pageYouSure      = "yousure"
	pageMods         = "mods"
	pageDefaultRight = "right"
	pageWarp         = "warp"
	pageDemos        = "demos"

	tableBorders = false

	colorTagPrimaryText  = "[white]"
	colorTagContrast     = "[royalblue]"
	colorTagMoreContrast = "[orange]"
)

var (
	config *cfg.Cfg
	app    *tview.Application

	canvas       *tview.Pages
	mainFlex     *tview.Flex
	headerPages  *tview.Pages
	contentPages *tview.Pages
	footerPages  *tview.Pages

	detailPages         *tview.Pages
	detailSidePagesSub1 *tview.Pages
	detailSidePagesSub2 *tview.Pages

	gamesTable     *tview.Table
	commandPreview *tview.TextView
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
		contentPages.AddPage(pageFirstSetup, makeFirstTimeSetup(), true, true)
	}

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

	// build up layout
	canvas = tview.NewPages()
	headerPages = tview.NewPages()
	contentPages = tview.NewPages()
	footerPages = tview.NewPages()

	mainFlex = tview.NewFlex().SetDirection(tview.FlexRow)
	canvas.AddPage(pageMain, mainFlex, true, true)
	// header
	header, headerHeight := getHeader()
	headerPages.AddPage(pageHeader, header, true, true)
	mainFlex.AddItem(headerPages, headerHeight, 0, false)
	// content
	mainFlex.AddItem(contentPages, 0, 1, true)
	// footer
	helpPane, helpPaneHeight := makeHelpPane()
	footerPages.AddPage(pageHelp, helpPane, true, true)
	mainFlex.AddItem(footerPages, helpPaneHeight, 0, false)

	// command preview
	commandPreview = makeCommandPreview()
	// main view to select games
	gamesTable = makeGamesTable()
	// responsive detail grid
	detailGrid := tview.NewGrid()
	detailGrid.SetRows(-1, -1)
	detailGrid.SetColumns(-4, -6)
	detailSidePagesSub1 = tview.NewPages()
	detailSidePagesSub2 = tview.NewPages()
	// not so wide screens
	detailGrid.AddItem(detailSidePagesSub1, 0, 0, 1, 2, 0, 0, false)
	detailGrid.AddItem(detailSidePagesSub2, 1, 0, 1, 2, 0, 0, false)
	// wide screens
	detailGrid.AddItem(detailSidePagesSub1, 0, 0, 2, 1, 0, 75, false)
	detailGrid.AddItem(detailSidePagesSub2, 0, 1, 2, 1, 0, 75, false)

	detailPages = tview.NewPages()
	detailPages.AddPage(pageContentMain, detailGrid, true, true)

	contentFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	contentFlex.
		AddItem(commandPreview, 4, 0, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(gamesTable, 0, config.GameListRelativeWidth, true).
			AddItem(detailPages, 0, 100-config.GameListRelativeWidth, true), 0, 1, true)

	contentPages.AddPage(pageContent, contentFlex, true, true)
}

// small or big header
func getHeader() (tview.Primitive, int) {
	headerHeight := 20
	var header tview.Primitive
	header = makeHeader()
	if cfg.GetInstance().HideHeader {
		headerHeight = 1
		header = tview.NewTextView().SetDynamicColors(true).SetText(subtitle)
	}
	return header, headerHeight
}

// update functions
func selectedGameChanged(g *games.Game) {
	populateCommandPreview(g)
	detailSidePagesSub1.AddPage(pageMods, makeModList(g), true, true)
	frontPage, _ := detailSidePagesSub2.GetFrontPage()
	if frontPage != pageModSelector {
		detailSidePagesSub2.AddPage(pageStats, makeStatsTable(g), true, true)
	}
}

// redraw whole table
func whenGamesChanged() {
	populateGamesTable()
}

// reset ui
func appModeNormal() {
	// cleanup
	// clear bigMainPager
	contentPages.RemovePage(pageYouSure)
	contentPages.RemovePage(pageFirstSetup)
	contentPages.RemovePage(pageOptions)
	contentPages.RemovePage(pageWarp)
	// clear actionPager
	detailPages.RemovePage(pageAddEdit)
	detailSidePagesSub1.RemovePage(pageYouSure)
	detailSidePagesSub1.RemovePage(pageDemos)

	// set ui state
	detailPages.SwitchToPage(pageContentMain)
	detailSidePagesSub1.SwitchToPage(pageMods)
	detailSidePagesSub2.SwitchToPage(pageStats)
	contentPages.SwitchToPage(pageContent)

	app.SetFocus(gamesTable)
}

// used in options and such screens
//func tabNavigate(previous, next tview.Primitive) func(event *tcell.EventKey) *tcell.EventKey {
//	return func(event *tcell.EventKey) *tcell.EventKey {
//		k := event.Key()
//		switch k {
//		case tcell.KeyTab:
//			app.SetFocus(next)
//			return nil
//		case tcell.KeyBacktab:
//			app.SetFocus(previous)
//			return nil
//		}
//
//		return event
//	}
//}
