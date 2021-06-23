package games

import (
	"strings"

	"github.com/zmnpl/twad/portspec"
)

// spSaveDirParam returns the right paramter key for specifying the savegame directory
// accounts for zdoom-, chocolate-doom and boom ports at the moments
func (g Game) spSaveDirParam() string {
	switch portspec.PortFamily(g.SourcePort) {
	case portspec.Boom:
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
	switch portspec.PortFamily(g.SourcePort) {
	case portspec.Chocolate:
		return inSkill + 1
	case portspec.Boom:
		return inSkill + 1
	default:
		return inSkill + 1
	}
}

// spSaveFileExtension gives the appropriate file extension
// adjusted for the games source port
func (g Game) spSaveFileExtension() string {
	switch portspec.PortFamily(g.SourcePort) {
	case portspec.Chocolate, portspec.Boom:
		return ".dsg"
	default:
		return ".zds"
	}
}

// spSaveGameName gives the appropriate syntax for save names
// adjusted for the games source port
func (g Game) spSaveGameName(save string) string {
	switch portspec.PortFamily(g.SourcePort) {
	case portspec.Chocolate, portspec.Boom:
		if save != "" {
			//tmp := []rune(save)
			//save = string(tmp[len(tmp)-5 : len(tmp)-4])
			save = strings.TrimSuffix(strings.TrimPrefix(save, "doomsav"), ".dsg")
			return save
		}
		return save
	default:
		return save
	}
}
