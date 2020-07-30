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
	commandPreview = tview.NewTextView().
		SetDynamicColors(true)

		//	commandPreview.SetBackgroundColor(previewBackgroundColor)
	fmt.Fprintf(commandPreview, "")

	return commandPreview
}

func populateCommandPreview(g *games.Game) {
	commandPreview.Clear()
	fmt.Fprintf(commandPreview, fmt.Sprintf("%s %s", previewText, stylizeCommandList(g.CommandList())))
}

func stylizeCommandList(params []string) string {
	keywords := make(map[string]int)
	keywords["-iwad"] = 1
	keywords["-file"] = 1
	keywords["-savedir"] = 1
	keywords["zdoom"] = 1
	keywords["gzdoom"] = 1
	keywords["lzdoom"] = 1

	var command strings.Builder
	for _, s := range params {
		if s == "" {
			continue
		}
		if _, isKnown := keywords[s]; isKnown {
			command.WriteString(fmt.Sprintf(" %s%s%s", colorTagMoreContrast, s, colorTagPrimaryText))
			continue
		}
		command.WriteString(fmt.Sprintf(" %s", s))
	}
	return strings.TrimSpace(command.String())
}
