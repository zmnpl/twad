package games

import (
	"strings"

	"github.com/zmnpl/twad/portspec"
)

// PortSaveDirParam returns the right paramter key for specifying the savegame directory
// accounts for zdoom-, chocolate-doom and boom ports at the moments
func PortSaveDirParam(port string) string {
	switch portspec.PortFamily(port) {
	case portspec.Boom:
		return "-save"
	default:
		return "-savedir"
	}
}

// PortAdjustedSkill for source port
// default(zdoom): 0-4 (documenation seems wrong?, so 1-5)
// chocolate: 1-5
// boom: 1-5
func PortAdjustedSkill(port string, skill int) int {
	switch portspec.PortFamily(port) {
	case portspec.Chocolate:
		return skill + 1
	case portspec.Boom:
		return skill + 1
	default:
		return skill + 1
	}
}

// PortSaveFileExtension gives the appropriate file extension
// adjusted for the games source port
func PortSaveFileExtension(port string) string {
	switch portspec.PortFamily(port) {
	case portspec.Chocolate, portspec.Boom:
		return ".dsg"
	default:
		return ".zds"
	}
}

// PortSaveGameName gives the appropriate syntax for save names
// adjusted for the games source port
func PortSaveGameName(port, save string) string {
	switch portspec.PortFamily(port) {
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
