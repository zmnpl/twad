package tui

import (
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/rivo/tview"
	"github.com/zmnpl/twad/games"
)

var (
	skillLevels []string
)

func init() {
	skillLevels = []string{"I'M TOO YOUNG TO DIE.", "HEY, NOT TOO ROUGH.", "HURT ME PLENTY.", "ULTRA-VIOLENCE.", "NIGHTMARE!"}
}

// warp strings are expected to be of form
// "e l"
// e=episode (number)
// space-character
// l=level (number)
// if one is ommited, the other one just works as "level" for doom ii and the like
func warpStringAcceptance(warp string, lastChar rune) (warpable bool) {
	return (unicode.IsDigit(lastChar) || unicode.IsSpace(lastChar)) && len([]rune(warp)) <= 5
}

func splitWarpString(warp string) (episode, level int) {
	parts := strings.Split(warp, " ")
	// episode
	if len(parts) > 0 {
		episode, _ = strconv.Atoi(parts[0])
	}
	// level
	if len(parts) > 1 {
		level, _ = strconv.Atoi(parts[1])
	}

	return
}

// warp dialog
func makeWarpRecord(game games.Game, onCancel func(), xOffset int, yOffset int, container *tview.Box) *tview.Flex {
	episode := 0
	level := 0

	warpRecordForm := tview.NewForm()
	warpRecordForm.
		SetBorder(true).
		SetTitle(game.Name).
		SetBackgroundColor(tview.Styles.PrimitiveBackgroundColor)

	// warp
	warpTo := tview.NewInputField().SetLabel(dict.warpText).SetAcceptanceFunc(warpStringAcceptance).SetFieldWidth(5)
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
	skl := tview.NewDropDown().SetOptions(skillLevels, nil).SetCurrentOption(2).SetLabel(dict.skillText)
	warpRecordForm.AddFormItem(skl)

	// to record a demo, specify a name
	demoName := tview.NewInputField().SetLabel(dict.demoText).SetFieldWidth(21)

	demoName.SetChangedFunc(func(text string) {
		demoName.SetLabel(dict.demoText)
		if game.DemoExists(demoName.GetText()) {
			demoName.SetLabel(warnColor + dict.demoTextOverwrite)
		}
	})
	warpRecordForm.AddFormItem(demoName)

	// confirm button
	warpRecordForm.AddButton(dict.warpOkButton, func() {
		episode, level = splitWarpString(warpTo.GetText())
		difficulty, _ := skl.GetCurrentOption()
		demo := demoName.GetText()

		appModeNormal() // TODO: looks like this is only executed after the game closed; not sure why

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

	// surrounding layout
	helpHeight := 5
	width := 50
	_, _, _, height := warpRecordForm.GetRect()
	_, _, _, containerHeight := container.GetRect()

	// though, if it flows out of the screen, then on top of the game
	if yOffset+height > containerHeight+helpHeight {
		yOffset = yOffset - height - 1
	}

	warpWindowLayout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(nil, yOffset, 0, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(nil, xOffset, 1, false).
			AddItem(warpRecordForm, width, 0, true).
			AddItem(nil, 0, 1, false),
			height+3, 0, true). // + 3 because default box size in tview is 15x10 (width x height)
		AddItem(nil, 0, 1, false)

	return warpWindowLayout
}
