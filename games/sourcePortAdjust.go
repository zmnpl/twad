package games

import "strings"

const (
	zdoom = iota
	chocolate
	boom
)

// sourcePortFamily checks the games engine type by inspecting the string
// known keyphrases will be interpreted as a certain source port family
func sourcePortFamily(sourcePort string) (t int) {
	t = zdoom

	sp := strings.ToLower(sourcePort)

	if strings.Contains(sp, "crispy") || strings.Contains(sp, "chocolate") {
		t = chocolate
		return
	}

	if strings.Contains(sp, "boom") {
		t = boom
		return
	}

	return
}

// spSaveDirParam returns the right paramter key for specifying the savegame directory
// accounts for zdoom-, chocolate-doom and boom ports at the moments
func (g Game) spSaveDirParam() string {
	switch sourcePortFamily(g.SourcePort) {
	case boom:
		return "-save"
	default:
		return "-savedir"
	}
}

// adjust skill for source port
// default(zdoom): 0-4 (documenation seems wrong?, so 1-5)
// chocolate: 1-5
// boom: 1-5
func (g Game) spAdjustedSkill(inSkill int) int {
	switch sourcePortFamily(g.SourcePort) {
	case chocolate:
		return inSkill + 1
	case boom:
		return inSkill + 1
	default:
		return inSkill + 1
	}
}

// spSaveFileExtension gives the appropriate file extension
// adjusted for the games source port
func (g Game) spSaveFileExtension() string {
	switch sourcePortFamily(g.SourcePort) {
	case chocolate, boom:
		return ".dsg"
	default:
		return ".zds"
	}
}

// spSaveGameName gives the appropriate syntax for save names
// adjusted for the games source port
func (g Game) spSaveGameName(save string) string {
	switch sourcePortFamily(g.SourcePort) {
	case chocolate, boom:
		if save != "" {
			tmp := []rune(save)
			save = string(tmp[len(tmp)-5 : len(tmp)-4])
			return save
		}
		return save
	default:
		return save
	}
}
