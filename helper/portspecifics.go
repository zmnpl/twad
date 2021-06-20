package helper

import "strings"

const (
	zdoom = iota
	chocolate
	boom
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

// PortConfigFileExtension returns the file extension of config files for the give port
func PortConfigFileExtension(port string) string {
	switch PortFamily(port) {
	case chocolate, boom:
		return ".cfg"
	default:
		return ".ini"
	}
}
