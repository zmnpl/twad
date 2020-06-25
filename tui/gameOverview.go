package tui

import (
	"github.com/rivo/tview"
	"github.com/zmnpl/twad/games"
)

const (
	overviewHeaderText      = "Overview"
	overviewEnvironmentVars = "Environment Variables"
	overviewIwad            = "IWAD"
	overviewSourcePort      = "Source Port"
	overviewGameName        = "Name"
	overviewOtherParams     = "Others"
	overviewMods            = "Mods"
)

func makeGameOverview(g *games.Game) *tview.Flex {

	overviewWindow := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewTextView().SetText(overviewGameName).SetTextColor(tview.Styles.SecondaryTextColor), 1, 0, false).
		AddItem(tview.NewTextView().SetText(g.Name), 1, 0, false).
		AddItem(nil, 1, 0, false).
		AddItem(tview.NewTextView().SetText(overviewSourcePort).SetTextColor(tview.Styles.SecondaryTextColor), 1, 0, false).
		AddItem(tview.NewTextView().SetText(g.SourcePort), 1, 0, false).
		AddItem(nil, 1, 0, false).
		AddItem(tview.NewTextView().SetText(overviewIwad).SetTextColor(tview.Styles.SecondaryTextColor), 1, 0, false).
		AddItem(tview.NewTextView().SetText(g.Iwad), 1, 0, false).
		AddItem(nil, 1, 0, false)

	if len(g.Environment) > 0 {
		if g.Environment[0] != "" {
			overviewWindow.AddItem(tview.NewTextView().SetText(overviewEnvironmentVars).SetTextColor(tview.Styles.SecondaryTextColor), 1, 0, false).
				AddItem(tview.NewTextView().SetText(g.EnvironmentString()), 2, 0, false).
				AddItem(nil, 1, 0, false)
		}
	}

	if len(g.Parameters) > 0 {
		if g.Parameters[0] != "" {
			overviewWindow.AddItem(tview.NewTextView().SetText(overviewOtherParams).SetTextColor(tview.Styles.SecondaryTextColor), 1, 0, false).
				AddItem(tview.NewTextView().SetText(g.ParamsString()), 1, 0, false).
				AddItem(nil, 1, 0, false)
		}
	}

	overviewWindow.AddItem(tview.NewTextView().SetText(overviewMods).SetTextColor(tview.Styles.SecondaryTextColor), 1, 0, false)
	for _, mod := range g.Mods {
		overviewWindow.AddItem(tview.NewTextView().SetText(mod), 1, 0, false)
	}

	overviewWindow.SetBorder(true)
	overviewWindow.SetTitle(overviewHeaderText)
	overviewWindow.SetBorderPadding(1, 1, 1, 1)

	return overviewWindow
}
