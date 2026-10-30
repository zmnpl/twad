package shared

import (
	"github.com/gdamore/tcell/v2"
)

// Theme holds the colors of the user interfaces
type Theme struct {
	PrimitiveBackground    tcell.Color
	ContrastBackground     tcell.Color
	MoreContrastBackground tcell.Color
	Border                 tcell.Color
	Title                  tcell.Color
	Graphics               tcell.Color
	PrimaryText            tcell.Color
	SecondaryText          tcell.Color
	TertiaryText           tcell.Color
	InverseText            tcell.Color
	ContrastSecondaryText  tcell.Color

	Good tcell.Color
	Warn tcell.Color

	// Terminal is set for the theme based on terminal colors;
	// its colors are best referred to by name then
	Terminal bool
}

// "Hard coded" theme
var twadTheme = Theme{
	PrimitiveBackground:    tcell.NewHexColor(0x000000), // black
	ContrastBackground:     tcell.NewHexColor(0x4169E1), // royal blue
	MoreContrastBackground: tcell.NewHexColor(0xFFA500), // organge
	Border:                 tcell.NewHexColor(0x4169E1), // royal blue
	Title:                  tcell.NewHexColor(0x4169E1), // royal blue
	Graphics:               tcell.NewHexColor(0x4169E1), // royal blue
	PrimaryText:            tcell.NewHexColor(0xFFFFFF), // white
	SecondaryText:          tcell.NewHexColor(0xFFA500), // organge
	TertiaryText:           tcell.NewHexColor(0xFF69B4), // hot pink
	InverseText:            tcell.NewHexColor(0xFFFACD), // lemon chiffon
	ContrastSecondaryText:  tcell.NewHexColor(0xFFDAB9), // peach puff

	Good: tcell.NewHexColor(0x008000), // green
	Warn: tcell.NewHexColor(0xFF0000), // red
}

// Theme based on terminal colors
var terminalTheme = Theme{
	PrimitiveBackground:    tcell.ColorBlack,
	ContrastBackground:     tcell.ColorGray,
	MoreContrastBackground: tcell.ColorMaroon,
	Border:                 tcell.ColorNavy,
	Title:                  tcell.ColorMaroon,
	Graphics:               tcell.ColorDefault,
	PrimaryText:            tcell.ColorWhite,
	SecondaryText:          tcell.ColorPurple,
	TertiaryText:           tcell.ColorTeal,
	InverseText:            tcell.ColorBlack,
	ContrastSecondaryText:  tcell.ColorBlack,

	Good:     tcell.ColorGreen,
	Warn:     tcell.ColorRed,
	Terminal: true,
}

// SelectTheme returns the theme based on terminal colors or the twad theme
func SelectTheme(useTerminalColors bool) Theme {
	if useTerminalColors {
		return terminalTheme
	}
	return twadTheme
}

// ColorName returns the name of the color for terminal themes, its hex value otherwise
func (t Theme) ColorName(c tcell.Color) string {
	if t.Terminal {
		return c.Name()
	}
	return c.CSS()
}

// Tag returns the tview color tag for the given color of the theme
func (t Theme) Tag(c tcell.Color) string {
	return "[" + t.ColorName(c) + "]"
}
