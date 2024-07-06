package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// first 16 tcell colors
// ColorBlack
// ColorMaroon
// ColorGreen
// ColorOlive
// ColorNavy
// ColorPurple
// ColorTeal
// ColorSilver
// ColorGray
// ColorRed
// ColorLime
// ColorYellow
// ColorBlue
// ColorFuchsia
// ColorAqua
// ColorWhite

var (
	colorTagWarnColor = "[#FF0000]"
	warnColor         = tcell.NewHexColor(0xFF0000) // colorred

	colorTagGoodColor = "[#008000]"
	goodColor         = tcell.NewHexColor(0x008000) // colorgreen

	// these get changed/set up with the seme selection
	colorTagPrimaryText  = "[white]"
	colorTagContrast     = "[royalblue]"
	colorTagMoreContrast = "[orange]"
)

func init() {
	//selectTheme()
}

func selectTheme() {
	// "Hard coded" theme
	twadTheme := tview.Theme{
		// ui stylepageSettings
		PrimitiveBackgroundColor:    tcell.NewHexColor(0x000000), // black
		ContrastBackgroundColor:     tcell.NewHexColor(0x4169E1), // royal blue
		MoreContrastBackgroundColor: tcell.NewHexColor(0xFFA500), // organge
		BorderColor:                 tcell.NewHexColor(0x4169E1), // royal blue
		TitleColor:                  tcell.NewHexColor(0x4169E1), // royal blue
		GraphicsColor:               tcell.NewHexColor(0x4169E1), // royal blue
		PrimaryTextColor:            tcell.NewHexColor(0xFFFFFF), // white
		SecondaryTextColor:          tcell.NewHexColor(0xFFA500), // organge
		TertiaryTextColor:           tcell.NewHexColor(0xFF69B4), // hot pink
		InverseTextColor:            tcell.NewHexColor(0xFFFACD), // lemon chiffon
		ContrastSecondaryTextColor:  tcell.NewHexColor(0xFFDAB9), // peach puff
	}

	// Theme based on terminal colors
	terminalTheme := tview.Theme{
		PrimitiveBackgroundColor:    tcell.ColorBlack,
		ContrastBackgroundColor:     tcell.ColorGray,
		MoreContrastBackgroundColor: tcell.ColorMaroon,
		BorderColor:                 tcell.ColorNavy,
		TitleColor:                  tcell.ColorMaroon,
		GraphicsColor:               tcell.ColorDefault,
		PrimaryTextColor:            tcell.ColorWhite,
		SecondaryTextColor:          tcell.ColorPurple,
		TertiaryTextColor:           tcell.ColorTeal,
		InverseTextColor:            tcell.ColorBlack,
		ContrastSecondaryTextColor:  tcell.ColorBlack,
	}

	// Select theme
	tview.Styles = twadTheme

	colorTagPrimaryText = "[" + "#FFFFFF" + "]"  // like primary text
	colorTagContrast = "[" + "#4169E1" + "]"     // like title
	colorTagMoreContrast = "[" + "#FFA500" + "]" // like secondary text

	if config.UseTerminalColors {
		tview.Styles = terminalTheme

		colorTagPrimaryText = "[" + tview.Styles.PrimaryTextColor.Name() + "]"
		colorTagContrast = "[" + tview.Styles.TitleColor.Name() + "]"
		colorTagMoreContrast = "[" + tview.Styles.SecondaryTextColor.Name() + "]"

		colorTagGoodColor = "[green]"
		colorTagWarnColor = "[red]"
	}

}
