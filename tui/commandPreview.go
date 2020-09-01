package tui

import (
	"fmt"
	"strings"

	"github.com/rivo/tview"
	"github.com/zmnpl/twad/games"
)

const (
	previewText = "~$"
)

// command preview
func makeCommandPreview() *tview.TextView {
	commandPreview = tview.NewTextView().SetDynamicColors(true)
	commandPreview.SetBorder(true)
	fmt.Fprintf(commandPreview, "")

	return commandPreview
}

func populateCommandPreview(g *games.Game) {
	commandPreview.Clear()
	fmt.Fprintf(commandPreview, fmt.Sprintf("%s %s", previewText, stylizeCommandList(g.CommandList())))
}

func stylizeCommandList(params []string) string {
	keywords := map[string]int{
		"-iwad":          1,
		"-file":          1,
		"-savedir":       1,
		"-save":          1,
		"zdoom":          1,
		"gzdoom":         1,
		"lzdoom":         1,
		"chocolate-doom": 1,
		"crispy-doom":    1,
	}

	optionals := map[string]int{
		"-loadgame": 1,
	}

	var command strings.Builder
	for _, s := range params {
		if s == "" {
			continue
		}
		if _, isKnown := keywords[s]; isKnown {
			command.WriteString(fmt.Sprintf(" %s%s%s", colorTagMoreContrast, s, colorTagPrimaryText))
			continue
		}
		if _, isKnownOptional := optionals[s]; isKnownOptional {
			command.WriteString(fmt.Sprintf(" %s%s%s", colorTagContrast, s, colorTagPrimaryText))
			continue
		}
		command.WriteString(fmt.Sprintf(" %s", s))
	}
	return strings.TrimSpace(command.String())
}
