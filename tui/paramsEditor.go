package tui

import (
	"strings"

	"github.com/rivo/tview"
	"github.com/zmnpl/twad/games"
)

const (
	cpHeader            = "Custom Parameters"
	cpEnvironment       = "Environment Variables"
	cpEnvironmentDetail = `Provide environment variables here; To turn VSync off entirely for example:
"vblank_mode=1"`
	cpOthers       = "Others"
	cpOthersDetail = "Other parameters you want to pass to your ZDoom port"
)

func makeParamsEditor(g *games.Game) *tview.Flex {
	environmentLabel := tview.NewTextView().SetText(cpEnvironment).SetTextColor(tview.Styles.SecondaryTextColor)
	environmentDetailLabel := tview.NewTextView().SetText(cpEnvironmentDetail)
	environment := tview.NewInputField()
	environment.SetText(g.EnvironmentString())

	cpOthersLabel := tview.NewTextView().SetText(cpOthers).SetTextColor(tview.Styles.SecondaryTextColor)
	cpOthersDetailLabel := tview.NewTextView().SetText(cpOthersDetail)
	cpOthers := tview.NewInputField()
	cpOthers.SetText(g.ParamsString())

	okButton := tview.NewButton(optionsOkButtonLabel)
	okButton.SetBackgroundColorActivated(tview.Styles.PrimaryTextColor)
	okButton.SetLabelColorActivated(tview.Styles.ContrastBackgroundColor)

	okButton.SetSelectedFunc(func() {
		pfs := strings.Split(environment.GetText(), " ")
		for i := range pfs {
			pfs[i] = strings.TrimSpace(pfs[i])
		}
		g.Environment = pfs

		prms := strings.Split(cpOthers.GetText(), " ")
		for i := range prms {
			prms[i] = strings.TrimSpace(prms[i])
		}
		g.Parameters = prms

		games.Persist()
		appModeNormal()
	})

	// navigation path
	environment.SetInputCapture(optionMoveTo(cpOthers))
	cpOthers.SetInputCapture(optionMoveTo(okButton))
	okButton.SetInputCapture(optionMoveTo(environment))

	cps := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(environmentLabel, 1, 0, false).
		AddItem(environmentDetailLabel, 2, 0, false).
		AddItem(environment, 1, 0, true).
		AddItem(nil, 1, 0, false).
		AddItem(cpOthersLabel, 1, 0, false).
		AddItem(cpOthersDetailLabel, 1, 0, false).
		AddItem(cpOthers, 1, 0, false).
		AddItem(nil, 1, 0, false).
		AddItem(okButton, 1, 0, false)
	cps.SetBorder(true)
	cps.SetTitle(cpHeader)
	cps.SetBorderPadding(1, 1, 1, 1)

	settingsPage := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(cps, 90, 0, true).
		AddItem(tview.NewBox().SetBorder(false), 0, 1, false)

	return settingsPage
}
