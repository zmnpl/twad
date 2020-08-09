package tui

import (
	"strconv"
	"strings"
	"unicode"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"github.com/zmnpl/twad/games"
)

const (
	warpText = "Warp to episode level"
)

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
func makeWarp(game games.Game, onCancel func(), xOffset int, yOffset int) *tview.Flex {
	episode := 0
	level := 0

	warpTo := tview.NewInputField().SetLabel(warpText).
		SetAcceptanceFunc(warpStringAcceptance).
		SetFieldWidth(5)

	warpTo.SetDoneFunc(func(key tcell.Key) {
		episode, level = splitWarpString(warpTo.GetText())
		appModeNormal()
		game.Warp(episode, level)
	})

	// inner box
	// needed to set nice looking border + background color
	youSureBox := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(warpTo, 1, 0, true)
	youSureBox.SetBorder(true)
	youSureBox.SetBackgroundColor(tview.Styles.PrimitiveBackgroundColor)

	width := len([]rune(warpText)) + 10

	// surrounding layout
	youSureLayout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(nil, yOffset, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(nil, xOffset, 1, false).
			AddItem(youSureBox, width, 0, true).
			AddItem(nil, 0, 1, false),
			3, 1, true).
		AddItem(nil, 0, 1, false)

	return youSureLayout
}
