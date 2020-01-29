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

	pageOptions     = "options"
	pageStats       = "stats"
	pageNewForm     = "newform"
	pageModSelector = "modselector"
	pageSettings    = "settings"
	pageMain        = "main"
	pageHelp        = "help"
	pageLicense     = "license"
	pageYouSure     = "yousure"

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
	if bigMainPager.HasPage(pageYouSure) {
		bigMainPager.RemovePage(pageYouSure)
	}
	if bigMainPager.HasPage(pageHelp) {
		bigMainPager.RemovePage(pageHelp)
	}
	app.SetFocus(gamesTable)
}
