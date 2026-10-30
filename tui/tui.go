package tui

import (
	"fmt"

	"github.com/zmnpl/twad/games"
	"github.com/zmnpl/twad/tui/shared"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/zmnpl/goidgames"
	"github.com/zmnpl/twad/base"
)

const (
	pageOptions        = "options"
	pageStats          = "stats"
	pageSaveStats      = "savestats"
	pageAddEdit        = "addEdit"
	pageModSelector    = "modselector"
	pageHeader         = "header"
	pageMain           = "main"
	pageContent        = "content"
	pageContentMain    = "maincontent"
	pageHelp           = "help"
	pageLicense        = "license"
	pageYouSure        = "yousure"
	pageMods           = "mods"
	pageWarp           = "warp"
	pageDemos          = "demos"
	pageSaves          = "saves"
	pageError          = "error"
	pageZipSelect      = "zipselect"
	pageZipImport      = "zipimport"
	pageHello          = "hello"
	pageHelpKeymap     = "helpkeymap"
	pageIdgamesBrowser = "idgamesbrowse"

	tableBorders = false
)

// ui holds the state of the terminal user interface
type ui struct {
	config *base.Cfg
	app    *tview.Application

	canvas       *tview.Pages
	contentPages *layer

	detailPages         *layer
	detailSidePagesSub1 *layer
	detailSidePagesSub2 *layer

	gamesTable     *tview.Table
	commandPreview *tview.TextView

	zipSelector    *zipSelect
	idgamesBrowser *goidgames.IdgamesBrowser

	statusline *tview.TextView

	keyGroups []keyGroup
}

// Draw performs all necessary steps to start the ui
func Draw() {
	u := &ui{config: base.Config()}
	u.selectTheme()
	u.buildKeyGroups()
	games.RegisterChangeListener(u.whenGamesChanged)
	u.run()
}

// run performs all necessary steps to start the ui
func (u *ui) run() {
	u.initUIElements()

	// settings - only when first start of app
	if !u.config.Configured {
		u.contentPages.switchTo(pageHello, u.makeFirstStartHello())
		u.config.Configured = true
	}

	// populate
	u.selectedGameChanged(&games.Game{})
	u.populateGamesTable()

	// capture input
	u.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		k := event.Key()
		// switch back to nowmal mode
		if k == tcell.KeyESC {
			u.appModeNormal()
			return nil
		}
		return event
	})

	// run app
	if err := u.app.SetRoot(u.canvas, true).SetFocus(u.contentPages).Run(); err != nil {
		panic(err)
	}
}

func (u *ui) initUIElements() {
	// init basic primitives
	u.app = tview.NewApplication()

	u.statusline = tview.NewTextView()

	// build up layout
	u.canvas = tview.NewPages()
	headerPages := tview.NewPages()
	u.contentPages = newLayer(u, pageContent)
	footerPages := tview.NewPages()

	// create header
	header, headerHeight := u.getHeader()
	headerPages.AddPage(pageHeader, header, true, true)
	// create footer
	helpPane, helpPaneHeight := u.makeKeyMap()
	footerPages.AddPage(pageHelp, helpPane, true, true)

	// set up main grid layout
	mainGrid := tview.NewGrid()
	mainGrid.SetRows(headerHeight, -1, helpPaneHeight, 1)
	mainGrid.SetColumns(-1)
	u.canvas.AddPage(pageMain, mainGrid, true, true)

	// add to main grid
	// header
	mainGrid.AddItem(headerPages, 0, 0, 1, 1, headerHeight+20, 0, false)
	// content
	mainGrid.AddItem(u.contentPages, 0, 0, 2, 1, 0, 0, true)
	mainGrid.AddItem(u.contentPages, 1, 0, 1, 1, headerHeight+20, 0, true)
	// footer
	mainGrid.AddItem(footerPages, 2, 0, 1, 1, 0, 0, false)

	// status line
	mainGrid.AddItem(u.statusline, 3, 0, 1, 1, 0, 0, false)

	// command preview
	u.commandPreview = u.makeCommandPreview()
	// main view to select games
	u.gamesTable = u.makeGamesTable()
	// responsive detail grid
	detailGrid := tview.NewGrid()
	detailGrid.SetRows(-1, -1)
	detailGrid.SetColumns(-4, -6)
	u.detailSidePagesSub1 = newLayer(u, pageMods)
	u.detailSidePagesSub2 = newLayer(u, pageStats)
	// first not so wide screens, then add wide screen layouts
	detailGrid.AddItem(u.detailSidePagesSub1, 0, 0, 1, 2, 0, 0, false).AddItem(u.detailSidePagesSub1, 0, 0, 2, 1, 0, 75, false)
	detailGrid.AddItem(u.detailSidePagesSub2, 1, 0, 1, 2, 0, 0, false).AddItem(u.detailSidePagesSub2, 0, 1, 2, 1, 0, 75, false)

	u.detailPages = newLayer(u, pageContentMain)
	u.detailPages.setBase(detailGrid)

	u.zipSelector = newZipImportUI(u)
	u.contentPages.AddPage(pageZipSelect, u.zipSelector.layout, true, false)

	// id games browser
	u.idgamesBrowser = goidgames.NewIdgamesBrowser(u.app)
	u.idgamesBrowser.SetDownloadPath(base.DOWNLOAD_PATH())
	u.idgamesBrowser.SetConfirmCallback(func(g goidgames.Idgame) {
		confirm(u.contentPages, fmt.Sprintf("Download %v?", g.Title), 0, 5, u.idgamesBrowser.GetRootLayout(), func() {
			// keep input away from the browser while downloading
			u.app.SetFocus(u.statusline)

			go func() {
				path, err := base.DownloadIdGame(g, base.DOWNLOAD_PATH(), func(total uint64) {
					u.app.QueueUpdateDraw(func() {
						u.statusline.SetText(fmt.Sprintf(" Downloading... %v complete", total))
					})
				})

				// back to the ui goroutine
				u.app.QueueUpdateDraw(func() {
					u.statusline.Clear()
					if err != nil {
						u.showError("Download Failed", err.Error(), u.idgamesBrowser.GetRootLayout(), nil)
						return
					}

					// suggest import of downloaded archive
					u.runZipImport(path, g.Title, 0, 5, u.idgamesBrowser.GetRootLayout())
				})
			}()
		})
	})
	u.contentPages.AddPage(pageIdgamesBrowser, u.idgamesBrowser.GetRootLayout(), true, false)

	// add central content page
	contentFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	contentFlex.
		AddItem(u.commandPreview, 4, 0, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(u.gamesTable, 0, u.config.GameListRelativeWidth, true).
			AddItem(u.detailPages, 0, 100-u.config.GameListRelativeWidth, true), 0, 1, true)
	u.contentPages.setBase(contentFlex)
}

// small or big header
func (u *ui) getHeader() (tview.Primitive, int) {
	headerHeight := 19
	var header tview.Primitive
	header = u.makeHeader()
	if base.Config().HideHeader {
		headerHeight = 1
		header = tview.NewTextView().SetDynamicColors(true).SetText(shared.Subtitle)
	}
	return header, headerHeight
}

// update functions
func (u *ui) selectedGameChanged(g *games.Game) {
	u.populateCommandPreview(g)
	u.detailSidePagesSub1.setBase(u.makeModList(g))
	u.detailSidePagesSub2.setBase(makeStatsTable(g))
}

// launch runs a game in the background, so the ui stays responsive while playing
func (u *ui) launch(g *games.Game, run func()) {
	u.statusline.SetText(fmt.Sprintf(" Playing %v ...", g.Name))

	go func() {
		run()

		// back to the ui goroutine
		u.app.QueueUpdateDraw(func() {
			u.statusline.Clear()
			// playtime and stats changed
			if selected, ok := u.selectedGame(); ok && selected == g {
				u.selectedGameChanged(g)
			}
		})
	}()
}

// redraw whole table
func (u *ui) whenGamesChanged() {
	u.populateGamesTable()
}

// reset ui
func (u *ui) appModeNormal() {
	u.contentPages.reset()
	u.detailPages.reset()
	u.detailSidePagesSub1.reset()
	u.detailSidePagesSub2.reset()

	// reset import
	u.zipSelector.reset()

	u.statusline.Clear()

	u.app.SetFocus(u.gamesTable)
}
