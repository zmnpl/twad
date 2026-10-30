package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/zmnpl/twad/tui/shared"
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

func (u *ui) selectTheme() {
	t := shared.SelectTheme(u.config.UseTerminalColors)

	tview.Styles = tview.Theme{
		PrimitiveBackgroundColor:    t.PrimitiveBackground,
		ContrastBackgroundColor:     t.ContrastBackground,
		MoreContrastBackgroundColor: t.MoreContrastBackground,
		BorderColor:                 t.Border,
		TitleColor:                  t.Title,
		GraphicsColor:               t.Graphics,
		PrimaryTextColor:            t.PrimaryText,
		SecondaryTextColor:          t.SecondaryText,
		TertiaryTextColor:           t.TertiaryText,
		InverseTextColor:            t.InverseText,
		ContrastSecondaryTextColor:  t.ContrastSecondaryText,
	}

	colorTagPrimaryText = t.Tag(t.PrimaryText)    // like primary text
	colorTagContrast = t.Tag(t.Title)             // like title
	colorTagMoreContrast = t.Tag(t.SecondaryText) // like secondary text
	colorTagGoodColor = t.Tag(t.Good)
	colorTagWarnColor = t.Tag(t.Warn)
}
