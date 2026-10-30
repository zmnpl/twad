package shared

import (
	"encoding/csv"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"unicode"
)

// SkillLevels are the difficulties, as named by the game
var SkillLevels = []string{"I'M TOO YOUNG TO DIE.", "HEY, NOT TOO ROUGH.", "HURT ME PLENTY.", "ULTRA-VIOLENCE.", "NIGHTMARE!"}

// SplitParams splits space separated parameters; quotes keep parameters with spaces together
func SplitParams(params string) []string {
	paramsReader := strings.NewReader(params)
	c := csv.NewReader(paramsReader)
	c.Comma = ' '

	p, err := c.ReadAll()
	if err == nil && len(p) > 0 {
		for i, v := range p[0] {
			p[0][i] = strings.TrimSpace(v)
		}
		return p[0]
	}
	return []string{}
}

// WarpStringAcceptance tells if lastChar may be added to a warp string
// warp strings are expected to be of form
// "e l"
// e=episode (number)
// space-character
// l=level (number)
// if one is ommited, the other one just works as "level" for doom ii and the like
func WarpStringAcceptance(warp string, lastChar rune) (warpable bool) {
	return (unicode.IsDigit(lastChar) || unicode.IsSpace(lastChar)) && len([]rune(warp)) <= 5
}

// SplitWarpString returns episode and level of a warp string
func SplitWarpString(warp string) (episode, level int) {
	parts := strings.Split(warp, " ")
	// episode
	if len(parts) > 0 {
		episode, _ = strconv.Atoi(parts[0])
	}
	// level
	if len(parts) > 1 {
		level, _ = strconv.Atoi(parts[1])
	}

	return
}

// CommandExists tells if cmd can be found in path
func CommandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// StylizeCommandList joins the parameters of a command to a line
// known keywords and optional parameters get highlighted with the given color tags
func StylizeCommandList(params []string, keywordTag, optionalTag, textTag string) string {
	keywords := map[string]bool{
		"-iwad":          true,
		"-file":          true,
		"-savedir":       true,
		"-save":          true,
		"-config":        true,
		"-statdump":      true,
		"-levelstat":     true,
		"-fast":          true,
		"zdoom":          true,
		"gzdoom":         true,
		"lzdoom":         true,
		"zandronum":      true,
		"chocolate-doom": true,
		"crispy-doom":    true,
		"prboom":         true,
		"prboom-plus":    true,
		"woof":           true,
		"dsda":           true,
	}

	optionals := map[string]int{
		"-loadgame": 1,
	}

	var command strings.Builder
	for _, s := range params {
		if s == "" {
			continue
		}
		if _, isKnown := keywords[s]; isKnown {
			command.WriteString(fmt.Sprintf(" %s%s%s", keywordTag, s, textTag))
			continue
		}
		if _, isKnownOptional := optionals[s]; isKnownOptional {
			command.WriteString(fmt.Sprintf(" %s%s%s", optionalTag, s, textTag))
			continue
		}
		command.WriteString(fmt.Sprintf(" %s", s))
	}
	return strings.TrimSpace(command.String())
}
