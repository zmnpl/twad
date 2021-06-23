package portspec

import "strings"

const (
	Zdoom = iota
	Chocolate
	Boom

	ZdoomSaveExtension     = ".zds.json"
	BoomSaveExtension      = ".dsg"
	ChocolateSaveExtension = ".dsg"
)

var (
	PortCanonicalNames = map[string]string{
		"gzdoom":     "gzdoom",
		"zandronum":  "zandronum",
		"lzdoom":     "lzdoom",
		"crispy":     "crispydoom",
		"chocolate":  "chocolatedoom",
		"prboomplus": "prboomplus",
		"boom":       "boom",
		"na":         "unknown_port",
	}
)

// PortFamily checks the games engine type by inspecting the string
// known keyphrases will be interpreted as a certain source port family
func PortFamily(sourcePort string) (t int) {
	t = Zdoom

	sp := strings.ToLower(sourcePort)

	if strings.Contains(sp, "crispy") || strings.Contains(sp, "chocolate") {
		t = Chocolate
		return
	}

	if strings.Contains(sp, "boom") {
		t = Boom
		return
	}

	return
}

// PortConfigFileExtension returns the file extension of config files for the give port
func PortConfigFileExtension(port string) string {
	switch PortFamily(port) {
	case Chocolate, Boom:
		return ".cfg"
	default:
		return ".ini"
	}
}

// PortSaveDirParam returns the right paramter key for specifying the savegame directory
// accounts for zdoom-, chocolate-doom and boom ports at the moments
func PortSaveDirParam(port string) string {
	switch PortFamily(port) {
	case Boom:
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
	switch PortFamily(port) {
	case Chocolate:
		return skill + 1
	case Boom:
		return skill + 1
	default:
		return skill + 1
	}
}

// PortSaveFileExtension gives the appropriate file extension
// adjusted for the games source port
func PortSaveFileExtension(port string) string {
	switch PortFamily(port) {
	case Chocolate, Boom:
		return ".dsg"
	default:
		return ".zds"
	}
}

// PortSaveGameName gives the appropriate syntax for save names
// adjusted for the games source port
func PortSaveGameName(port, save string) string {
	switch PortFamily(port) {
	case Chocolate, Boom:
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
