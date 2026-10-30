package tui

import (
	"sort"

	"github.com/rivo/tview"
	"github.com/zmnpl/twad/games"
	"github.com/zmnpl/twad/tui/shared"
)

const (
	warpDialogWidth  = 50
	warpDialogHeight = 13
)

// warp dialog
func (u *ui) makeWarpRecord(game games.Game) *tview.Form {
	episode := 0
	level := 0

	warpRecordForm := tview.NewForm()
	warpRecordForm.
		SetBorder(true).
		SetTitle(game.Name).
		SetBackgroundColor(tview.Styles.PrimitiveBackgroundColor)

	// warp
	warpTo := tview.NewInputField().SetLabel(shared.Dict.WarpText).SetAcceptanceFunc(shared.WarpStringAcceptance).SetFieldWidth(5)
	warpRecordForm.AddFormItem(warpTo)

	// map select
	maps := game.ModMaps()
	displayMaps := make([]string, 0, len(maps)+1)
	displayMaps = append(displayMaps, "")
	for k := range maps {
		displayMaps = append(displayMaps, k)
	}
	sort.Strings(displayMaps)
	mapSelect := tview.NewDropDown().SetOptions(displayMaps, nil).SetCurrentOption(0).SetLabel("maps")
	warpRecordForm.AddFormItem(mapSelect)

	// skill level
	skl := tview.NewDropDown().SetOptions(shared.SkillLevels, nil).SetCurrentOption(2).SetLabel(shared.Dict.SkillText)
	warpRecordForm.AddFormItem(skl)

	// to record a demo, specify a name
	demoName := tview.NewInputField().SetLabel(shared.Dict.DemoText).SetFieldWidth(21)

	demoName.SetChangedFunc(func(text string) {
		demoName.SetLabel(shared.Dict.DemoText)
		if game.DemoExists(demoName.GetText()) {
			demoName.SetLabel(colorTagWarnColor + shared.Dict.DemoTextOverwrite)
		}
	})
	warpRecordForm.AddFormItem(demoName)

	// confirm button
	warpRecordForm.AddButton(shared.Dict.WarpOkButton, func() {
		episode, level = shared.SplitWarpString(warpTo.GetText())
		difficulty, _ := skl.GetCurrentOption()
		demo := demoName.GetText()

		u.appModeNormal() // TODO: looks like this is only executed after the game closed; not sure why

		i, beamToMapDisplayName := mapSelect.GetCurrentOption()
		if i > 0 {
			if len(demo) > 0 {
				game.GoToMapRecord(maps[beamToMapDisplayName], difficulty, demo)
				return
			} else {
				game.GoToMap(maps[beamToMapDisplayName], difficulty)
				return
			}
		}

		// supplying a demoname automatically starts recording
		if len(demo) > 0 {
			game.WarpRecord(episode, level, difficulty, demo)
		} else {
			game.Warp(episode, level, difficulty)
		}
	})

	return warpRecordForm
}
