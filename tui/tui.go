package tui

import (
	"fmt"
	"path/filepath"

	"github.com/zmnpl/twad/games"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/zmnpl/goidgames"
	"github.com/zmnpl/twad/base"
)

type Foo *tview.TextView

type Bar struct {
	*tview.TextView
}

const (
	colorTagPrimaryText  = "[white]"
	colorTagContrast     = "[royalblue]"
	colorTagMoreContrast = "[orange]"

	warnColor  = "[red]"
	warnColorO = tcell.ColorRed
	goodColor  = "[green]"
	goodColorO = tcell.ColorGreen

	pageOptions        = "options"
	pageStats          = "stats"
	pageAddEdit        = "addEdit"
	pageModSelector    = "modselector"
	pageFirstSetup     = "firstsetup"
	pageHeader         = "header"
	pageMain           = "main"
	pageDetailGrid     = "detailgrid"
	pageContent        = "content"
	pageContentMain    = "maincontent"
	pageHelp           = "help"
	pageLicense        = "license"
	pageYouSure        = "yousure"
	pageMods           = "mods"
	pageDefaultRight   = "right"
	pageWarp           = "warp"
	pageDemos          = "demos"
	pageSaves          = "saves"
	pageError          = "error"
	pageZipImport      = "zipselect"
	pageHello          = "hello"
	pageHelpKeymap     = "helpkeymap"
	pageIdgamesBrowser = "idgamesbrowse"
	pageDLConfirm      = "downloadConfirm"

	tableBorders = false
)

var (
	config *base.Cfg
	app    *tview.Application

	canvas       *tview.Pages
	headerPages  *tview.Pages
	contentPages *tview.Pages
	footerPages  *tview.Pages

	detailPages         *tview.Pages
	detailSidePagesSub1 *tview.Pages
	detailSidePagesSub2 *tview.Pages

	gamesTable     *tview.Table
	commandPreview *tview.TextView

	zipInput       *zipSelect
	idgamesBrowser *goidgames.IdgamesBrowser

	statusline *tview.TextView
)

func init() {
	config = base.Config()
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
		hello := makeFirstStartHello()
		contentPages.AddPage(pageHello, hello, true, false)
		contentPages.SwitchToPage(pageHello)
		app.SetFocus(hello)
		config.Configured = true
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

	statusline = tview.NewTextView()
	statusline.SetChangedFunc(func() {
		app.Draw()
	})

	// build up layout
	canvas = tview.NewPages()
	headerPages = tview.NewPages()
	contentPages = tview.NewPages()
	footerPages = tview.NewPages()

	// create header
	header, headerHeight := getHeader()
	headerPages.AddPage(pageHeader, header, true, true)
	// create footer
	helpPane, helpPaneHeight := makeKeyMap()
	footerPages.AddPage(pageHelp, helpPane, true, true)

	// set up main grid layout
	mainGrid := tview.NewGrid()
	mainGrid.SetRows(headerHeight, -1, helpPaneHeight, 1)
	mainGrid.SetColumns(-1)
	canvas.AddPage(pageMain, mainGrid, true, true)

	// add to main grid
	// header
	mainGrid.AddItem(headerPages, 0, 0, 1, 1, headerHeight+20, 0, false)
	// content
	mainGrid.AddItem(contentPages, 0, 0, 2, 1, 0, 0, true)
	mainGrid.AddItem(contentPages, 1, 0, 1, 1, headerHeight+20, 0, true)
	// footer
	mainGrid.AddItem(footerPages, 2, 0, 1, 1, 0, 0, false)

	// status line
	mainGrid.AddItem(statusline, 3, 0, 1, 1, 0, 0, false)

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

	zipInput = newZipImportUI()
	contentPages.AddPage(pageZipImport, zipInput.layout, true, true)

	dlPath := filepath.Join(base.Config().WadDir, "twad_downloads")
	idgamesBrowser = goidgames.NewIdgamesBrowser(app)
	idgamesBrowser.SetDownloadPath(dlPath)

	idgamesBrowser.SetConfirmCallback(func(g goidgames.Idgame) {
		youSure := makeYouSureBox(fmt.Sprintf("Download %v?", g.Title),
			func() {
				path, err := DownloadIdGame(g, dlPath)
				if err != nil {
					showError("Download Failed", err.Error(), tview.NewInputField(), nil)
					contentPages.RemovePage(pageDLConfirm)
					return
				}

				contentPages.RemovePage(pageDLConfirm)
				runZipImport(path, g.Title, 0, 5, idgamesBrowser.GetRootLayout())
				//app.SetFocus(idgamesBrowser.GetRootLayout())
			},
			func() {
				contentPages.RemovePage(pageDLConfirm)
				app.SetFocus(idgamesBrowser.GetRootLayout())
			},
			0,
			5,
			idgamesBrowser.GetRootLayout().Box)

		contentPages.AddPage(pageDLConfirm, youSure,
			true, true)
	})
	contentPages.AddPage(pageIdgamesBrowser, idgamesBrowser.GetRootLayout(), true, true)

	contentPages.AddPage(pageContent, contentFlex, true, true)
}

// small or big header
func getHeader() (tview.Primitive, int) {
	headerHeight := 19
	var header tview.Primitive
	header = makeHeader()
	if base.Config().HideHeader {
		headerHeight = 1
		header = tview.NewTextView().SetDynamicColors(true).SetText(subtitle)
	}
	return header, headerHeight
}

// update functions
func selectedGameChanged(g *games.Game) {
	populateCommandPreview(g)
	detailSidePagesSub1.AddPage(pageMods, makeModList(g), true, true)
	detailSidePagesSub2.AddPage(pageStats, makeStatsTable(g), true, true)
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
	contentPages.RemovePage(pageError)

	// clear actionPager
	detailPages.RemovePage(pageAddEdit)
	detailSidePagesSub1.RemovePage(pageYouSure)
	detailSidePagesSub1.RemovePage(pageDemos)
	detailSidePagesSub1.RemovePage(pageSaves)

	// reset import
	zipInput.reset()

	// set ui state
	detailPages.SwitchToPage(pageContentMain)
	detailSidePagesSub1.SwitchToPage(pageMods)
	detailSidePagesSub2.SwitchToPage(pageStats)
	contentPages.SwitchToPage(pageContent)

	statusline.Clear()

	app.SetFocus(gamesTable)
}
