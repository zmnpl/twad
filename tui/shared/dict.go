package shared

// Dictionary holds all texts of the user interfaces
type Dictionary struct {
	// demolist
	DemosHeader string

	// add edit game
	AddGame  string
	EditGame string

	AeName       string
	AeSourcePort string
	AeIWAD       string
	AeOwnCfg     string
	AeSharedCfgT string
	AeLink       string

	AeEnvironment       string
	AeEnvironmentDetail string

	AeOtherParams       string
	AeOtherParamsDetail string

	AeOkButton string

	// errorDisplay
	ErrTitleStart string
	ErrYolo       string
	ErrNotYet     string
	ErrAbort      string

	// gamestable
	GameTableHeaderRating     string
	GameTableHeaderName       string
	GameTableHeaderSourcePort string
	GameTableHeaderIwad       string

	DeleteGameQuestion string
	DeleteModQuestion  string

	// modlist
	OverviewMods string

	ModTreeTitle     string
	ModTreeErrWadDir string

	// options
	OptsErrPathDoesntExist string
	OptsErrPathNoIWads     string
	OptsLooksGood          string

	OptsHeader                   string
	OptsOkButtonLabel            string
	OptsPathLabel                string
	OptsDontDOOMWADDIR           string
	OptsWriteBasePathToEngineCFG string
	OptsDontWarn                 string
	OptsSourcePortLabel          string
	OptsIwadsLabel               string
	OptsHideHeader               string
	OptsGamesListRelativeWitdh   string
	OptsUseTerminalColors        string

	// samegamelist
	SavesHeader string

	// warprecord
	WarpText          string
	MapSelectText     string
	SkillText         string
	DemoText          string
	DemoTextOverwrite string
	WarpOkButton      string

	// yousure
	ConfirmText string
	AbortText   string

	// zipselect
	ZipSelectTitle              string
	ZipImportToLabel            string
	ZipImportToExistsLabel      string
	ZipImportToBadNameLabel     string
	ZipImportFormTitle          string
	ZipImportFormOk             string
	ZipImportCancel             string
	ZipImportTitle              string
	ZipImportSecurityWarn       string
	ZipImportNameInvalid        string
	ZipImportNameInvalidReasons string
	ZipImportArchiveNotFound    string
	ZipImportFailed             string
}

func defaultDict() Dictionary {
	return Dictionary{
		// demolist
		DemosHeader: "Demos",

		// add edit game
		AddGame:  "Add New Game",
		EditGame: "Edit Game",

		AeName:       "Name",
		AeSourcePort: "Source Port",
		AeIWAD:       "IWAD",
		AeOwnCfg:     "Use Own Config",
		AeSharedCfgT: "Use Shared Config [%v]",
		AeLink:       "Mod URL",

		AeEnvironment: "Environment Variables *",
		AeEnvironmentDetail: `* Provide environment variables here; To turn VSync off entirely for example:
	"vblank_mode:1"`,
		AeOtherParams:       "Others **",
		AeOtherParamsDetail: "** Other parameters you want to pass to your ZDoom port for this game",

		AeOkButton: "Ok",

		// errorDisplay
		ErrTitleStart: "ERROR",
		ErrYolo:       "I don't care. Go ahead.",
		ErrNotYet:     "Let me fix that first!",
		ErrAbort:      "Ok",

		// gamestable
		GameTableHeaderRating:     "Rating",
		GameTableHeaderName:       "Name",
		GameTableHeaderSourcePort: "SourcePort",
		GameTableHeaderIwad:       "Iwad",

		DeleteGameQuestion: "Delete '%v'?",
		DeleteModQuestion:  "Remove '%v' from '%v'?",

		// modlist
		OverviewMods: "Mods in order",

		ModTreeTitle:     "Add new mod to game",
		ModTreeErrWadDir: "WAD directory not available",

		// options
		OptsErrPathDoesntExist: "doesn't exist",
		OptsErrPathNoIWads:     "doesn't contain IWADs",
		OptsLooksGood:          "looks good",

		OptsHeader:                   "Options",
		OptsOkButtonLabel:            "Save",
		OptsPathLabel:                "WAD Dir",
		OptsDontDOOMWADDIR:           "Do NOT set DOOMWADDIR for current session (use your shell's default)",
		OptsWriteBasePathToEngineCFG: "Write the path into DOOM engines *.ini files",
		OptsDontWarn:                 "Do NOT warn before deletion",
		OptsSourcePortLabel:          "Source Port",
		OptsIwadsLabel:               "IWADs",
		OptsHideHeader:               "UI - Hide big DOOM logo",
		OptsGamesListRelativeWitdh:   "UI - Game list relative width (1-100%)",
		OptsUseTerminalColors:        "UI - Use terminal color theme",

		// samegamelist
		SavesHeader: "Savegames",

		// warprecord
		WarpText:          "Warp: Episode(optional) Map",
		MapSelectText:     "Select Map from Mod",
		SkillText:         "Difficulty",
		DemoText:          "Demo Name",
		DemoTextOverwrite: "Overwriting",
		WarpOkButton:      "Rip And Tear!",

		// yousure
		ConfirmText: "Ok",
		AbortText:   "Hell No!",

		// zipselect
		ZipSelectTitle:              "Select archive",
		ZipImportToLabel:            "Folder name",
		ZipImportToExistsLabel:      "exists already",
		ZipImportToBadNameLabel:     "cannot use that name",
		ZipImportFormTitle:          "Import archive to",
		ZipImportFormOk:             "Import",
		ZipImportCancel:             "No thanks!",
		ZipImportTitle:              "Import",
		ZipImportSecurityWarn:       "SECURITY WARNING: Only import archives from trusted sources!",
		ZipImportNameInvalid:        "Cannot use that name",
		ZipImportNameInvalidReasons: "Possible reasons:\n- File name contains forbidden characters\n- No permission to write this file/folder",
		ZipImportArchiveNotFound:    "Archive not found",
		ZipImportFailed:             "Could not import zip",
	}
}

// Dict holds the texts shown in the user interfaces
var Dict Dictionary

func init() {
	Dict = defaultDict()
}
