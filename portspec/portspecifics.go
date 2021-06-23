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
