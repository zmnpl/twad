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
	keywords := map[string]bool{
		"-iwad":          true,
		"-file":          true,
		"-nodeh":         true,
		"-deh":           true,
		"-savedir":       true,
		"-save":          true,
		"-config":        true,
		"-statdump":      true,
		"-levelstat":     true,
		"zdoom":          true,
		"gzdoom":         true,
		"lzdoom":         true,
		"zandronum":      true,
		"chocolate-doom": true,
		"crispy-doom":    true,
		"prboom":         true,
		"prboom-plus":    true,
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
