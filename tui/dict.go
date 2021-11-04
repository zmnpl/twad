package tui

type dictionary struct {
	// demolist
	demosHeader string

	// add edit game
	addGame  string
	editGame string

	aeName       string
	aeSourcePort string
	aeIWAD       string
	aeOwnCfg     string
	aeNoDeh      string
	aeSharedCfgT string
	aeLink       string

	aeEnvironment       string
	aeEnvironmentDetail string

	aeOtherParams       string
	aeOtherParamsDetail string

	aeOkButton string

	// errorDisplay
	errTitleStart string
	errYolo       string
	errNotYet     string
	errAbort      string

	// gamestable
	gameTableHeaderRating     string
	gameTableHeaderName       string
	gameTableHeaderSourcePort string
	gameTableHeaderIwad       string

	deleteGameQuestion string
	deleteModQuestion  string

	// modlist
	overviewMods string

	modTreeTitle string

	// options
	optsErrPathDoesntExist string
	optsErrPathNoIWads     string
	optsLooksGood          string

	optsHeader                   string
	optsOkButtonLabel            string
	optsPathLabel                string
	optsDontDOOMWADDIR           string
	optsWriteBasePathToEngineCFG string
	optsDontWarn                 string
	optsSourcePortLabel          string
	optsIwadsLabel               string
	optsHideHeader               string
	optsGamesListRelativeWitdh   string

	// samegamelist
	savesHeader string

	// warprecord
	warpText          string
	mapSelectText     string
	skillText         string
	demoText          string
	demoTextOverwrite string
	warpOkButton      string

	// yousure
	confirmText string
	abortText   string

	// zipselect
	zipSelectTitle          string
	zipImportToLabel        string
	zipImportToExistsLabel  string
	zipImportToBadNameLabel string
	zipImportFormTitle      string
	zipImportFormOk         string
	zipImportCancel         string
	zipImportSecurityWarn   string
}

func defaultDict() dictionary {
	return dictionary{
		// demolist
		demosHeader: "Demos",

		// add edit game
		addGame:  "Add New Game",
		editGame: "Edit Game",

		aeName:       "Name",
		aeSourcePort: "Source Port",
		aeIWAD:       "IWAD",
		aeOwnCfg:     "Use Own Config",
		aeNoDeh:      "nodeh",
		aeSharedCfgT: "Use Shared Config [%v]",
		aeLink:       "Mod URL",

		aeEnvironment: "Environment Variables *",
		aeEnvironmentDetail: `* Provide environment variables here; To turn VSync off entirely for example:
	"vblank_mode:1"`,
		aeOtherParams:       "Others **",
		aeOtherParamsDetail: "** Other parameters you want to pass to your ZDoom port for this game",

		aeOkButton: "Ok",

		// errorDisplay
		errTitleStart: "ERROR",
		errYolo:       "I don't care. Go ahead.",
		errNotYet:     "Let me fix that first!",
		errAbort:      "Ok",

		// gamestable
		gameTableHeaderRating:     "Rating",
		gameTableHeaderName:       "Name",
		gameTableHeaderSourcePort: "SourcePort",
		gameTableHeaderIwad:       "Iwad",

		deleteGameQuestion: "Delete '%v'?",
		deleteModQuestion:  "Remove '%v' from '%v'?",

		// modlist
		overviewMods: "Mods in order",

		modTreeTitle: "Add new mod to game",

		// options
		optsErrPathDoesntExist: "doesn't exist",
		optsErrPathNoIWads:     "doesn't contain IWADs",
		optsLooksGood:          "looks good",

		optsHeader:                   "Options",
		optsOkButtonLabel:            "Save",
		optsPathLabel:                "WAD Dir",
		optsDontDOOMWADDIR:           "Do NOT set DOOMWADDIR for current session (use your shell's default)",
		optsWriteBasePathToEngineCFG: "Write the path into DOOM engines *.ini files",
		optsDontWarn:                 "Do NOT warn before deletion",
		optsSourcePortLabel:          "Source Port",
		optsIwadsLabel:               "IWADs",
		optsHideHeader:               "UI - Hide big DOOM logo",
		optsGamesListRelativeWitdh:   "UI - Game list relative width (1-100%)",

		// samegamelist
		savesHeader: "Savegames",

		// warprecord
		warpText:          "Warp ((Episode) Map)",
		mapSelectText:     "Select Map from Mod",
		skillText:         "Difficulty",
		demoText:          "Demo Name",
		demoTextOverwrite: "Overwriting",
		warpOkButton:      "Rip And Tear!",

		// yousure
		confirmText: "Delete",
		abortText:   "Hell No!",

		// zipselect
		zipSelectTitle:          "Select archive",
		zipImportToLabel:        "Folder name",
		zipImportToExistsLabel:  "exists already",
		zipImportToBadNameLabel: "cannot use that name",
		zipImportFormTitle:      "Import to",
		zipImportFormOk:         "Import",
		zipImportCancel:         "Back",
		zipImportSecurityWarn:   "SECURITY WARNING: Only import archives from trusted sources!",
	}
}

var dict dictionary

func init() {
	dict = defaultDict()
}
