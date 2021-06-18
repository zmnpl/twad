package tui

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/zmnpl/twad/cfg"
	"github.com/zmnpl/twad/games"
)

const (
	demosHeader = "Demos"
)

func makeDemoList(g *games.Game) (*tview.Flex, error) {
	// surrounding container
	demoFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	demoFlex.SetBorderPadding(0, 0, 1, 1)
	demoFlex.AddItem(tview.NewTextView().
		SetText(demosHeader).
		SetTextColor(tview.Styles.SecondaryTextColor), 1, 0, false)

	// list
	demoList := tview.NewList()
	demoFlex.AddItem(demoList, 0, 1, true)
	demoList.SetSecondaryTextColor(tview.Styles.TitleColor).SetSelectedFocusOnly(true)

	// get demos
	demos, err := g.Demos()
	if err != nil {
		return nil, err
	}

	if len(demos) == 0 {
		return nil, fmt.Errorf("no demos in demo dir")
	}

	// how to populate the list
	populate := func() {
		demoList.Clear()
		for _, demo := range demos {
			di, err := demo.Info()
			var t time.Time
			if err != nil {
				t = time.Now()
			} else {
				t = di.ModTime()
			}
			demoList.AddItem(demo.Name(), fmt.Sprintf("%v (%.2f KiB)", t.Format("2006-01-02 15:04"), float32(di.Size())/1024), '|', nil)
		}
	}

	// do it
	populate()

	// hit enter plays demo
	demoList.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		g.PlayDemo(demos[index].Name())
	})

	// removes demo at given index and focuses app properly
	removeDemo := func(i int) {
		// TODO: bug in tview; remove when fixed
		demoList.SetChangedFunc(nil) // BUG WORKAROUND

		demos, err = g.RemoveDemo(demos[i].Name())
		if err != nil {
			showError("could not remove demo", err.Error(), nil, nil)
			return
		}

		if len(demos) > 0 {
			populate()
			app.SetFocus(demoList)
		} else {
			appModeNormal()
		}
		games.Persist()

		//demoList.SetChangedFunc(changeFunc) // BUG WORKAROUND
	}

	// tab navigates back to games table; tab navigation on list is redundant
	demoList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		k := event.Key()

		// delete mod
		if k == tcell.KeyDelete {
			if demoList.GetItemCount() > 0 {
				// when in edit mode, this is only confusing
				ci := demoList.GetCurrentItem()
				if cfg.Config().DeleteWithoutWarning {
					removeDemo(ci)
					return nil
				}

				youSure := makeYouSureBox(demos[ci].Name(),
					func() {
						removeDemo(ci)
						detailSidePagesSub1.RemovePage(pageYouSure)
					},
					func() {
						detailSidePagesSub1.RemovePage(pageYouSure)
						app.SetFocus(demoList)
					},
					2, 2, demoList.Box)
				detailSidePagesSub1.AddPage(pageYouSure,
					youSure, true, true) // TODO: calculate offsets
				app.SetFocus(youSure)
			}
			return nil
		}

		if k == tcell.KeyRune {
			switch event.Rune() {
			// quit app from here as well
			case 'q':
				app.Stop()
				return nil

			// start zip import from here
			case 'i':
				contentPages.SwitchToPage(pageZipImport)
				app.SetFocus(zipInput.selectTree)
				return nil

			}

		}

		return event
	})

	return demoFlex, nil
}
