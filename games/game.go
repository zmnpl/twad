package games

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/zmnpl/twad/cfg"
)

// Game represents one game configuration
type Game struct {
	Name          string         `json:"name,omitempty"`
	SourcePort    string         `json:"source_port,omitempty"`
	Iwad          string         `json:"iwad,omitempty"`
	Environment   []string       `json:"environment,omitempty"`
	Mods          []string       `json:"mods,omitempty"`
	Parameters    []string       `json:"parameters,omitempty"`
	Stats         map[string]int `json:"stats,omitempty"`
	Playtime      int64          `json:"playtime,omitempty"`
	LastPlayed    string         `json:"last_played,omitempty"`
	SaveGameCount int            `json:"savegame_coutn,omitempty"`
	Rating        int            `json:"rating,omitempty"`
}

// NewGame creates new instance of a game
func NewGame(name, sourceport, iwad string) Game {
	config := cfg.GetInstance()
	var game Game
	game.Name = name

	// default source port
	game.SourcePort = "gzdoom"
	if len(config.SourcePorts) > 0 {
		game.SourcePort = config.SourcePorts[0]
	}

	// replace with given
	if sourceport != "" {
		game.SourcePort = sourceport
	}

	// default iwad
	game.Iwad = "doom2.wad"
	if len(config.IWADs) > 0 {
		game.Iwad = config.IWADs[0]
	}

	// replace with given
	if iwad != "" {
		game.Iwad = iwad
	}

	game.Environment = make([]string, 0)
	game.Parameters = make([]string, 0)
	game.Mods = make([]string, 0)
	game.Stats = make(map[string]int)

	return game
}

// AddMod adds mod
func (g *Game) AddMod(modFile string) {
	g.Mods = append(g.Mods, modFile)
	InformChangeListeners()
	Persist()
}

// Run executes given configuration and launches the mod
func (g *Game) Run() error {
	params := g.getLaunchParams()

	start := time.Now()

	// execute and capture output
	proc := exec.Command(g.SourcePort, params...)

	// add environment variables
	proc.Env = os.Environ()
	proc.Env = append(proc.Env, g.Environment...)

	// rip and tear!
	output, err := proc.CombinedOutput()
	if err != nil {
		return err
	}

	playtime := time.Since(start).Milliseconds()
	g.Playtime = g.Playtime + playtime
	g.LastPlayed = time.Now().Format("2006-01-02 15:04:05MST")

	go processOutput(string(output), g)

	return nil
}

// String returns the string which is run when running
func (g Game) String() string {
	params := g.getLaunchParams()

	return fmt.Sprintf("%s %s %s", g.EnvironmentString(), g.SourcePort, strings.TrimSpace(strings.Join(params, " ")))
}

// RatingString returns the string resulting from the games rating
func (g Game) RatingString() string {
	return strings.Repeat("*", g.Rating) + strings.Repeat("-", 5-g.Rating)
}

// EnvironmentString returns a join of all prefix parameters
func (g Game) EnvironmentString() string {
	return strings.TrimSpace(strings.Join(g.Environment, " "))
}

// ParamsString returns a join of all prefix parameters
func (g Game) ParamsString() string {
	return strings.TrimSpace(strings.Join(g.Parameters, " "))
}

// SaveCount returns the number of savegames existing for this game
func (g Game) SaveCount() int {
	if saves, err := ioutil.ReadDir(g.getSaveDir()); err == nil {
		return len(saves)
	}
	return 0
}

// Rate increases or decreases the games rating
func (g *Game) Rate(increment int) {
	g.Rating += increment
	switch {
	case g.Rating > 5:
		g.Rating = 5
	case g.Rating < 0:
		g.Rating = 0
	}

}

func (g Game) getLaunchParams() []string {
	params := make([]string, 1, 10)
	config := cfg.GetInstance()

	if g.Iwad != "" {
		params = append(params, "-iwad", g.Iwad)
	}

	if len(g.Mods) > 0 {
		params = append(params, "-file")
		params = append(params, g.Mods...)
	}

	if config.SaveDirs {
		err := os.MkdirAll(g.getSaveDir(), 0755)
		// only use separate save dir if directory has been craeted
		if err == nil {
			params = append(params, "-savedir")
			params = append(params, g.getSaveDir())
		}
	}

	// add custom parameters here
	params = append(params, g.Parameters...)

	return params
}

func (g Game) getSaveDir() string {
	return cfg.GetSavegameFolder() + "/" + g.cleansedName()
}

// cleansedName removes all but alphanumeric characters from name
// i.e. used for directory names
func (g Game) cleansedName() string {
	cleanser, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}

	return cleanser.ReplaceAllString(g.Name, "")
}

// processOutput processes the terminal output of the zdoom port
func processOutput(output string, g *Game) {
	for _, v := range strings.Split(output, "\n") {
		if stat, increment := parseStatline(v, g); stat != "" {
			g.Stats[stat] = g.Stats[stat] + increment
		}
	}

	Persist()
}

// parseStatLine receives each line from processOutput()
// if the line matches a known pattern it will be added to the games stats
func parseStatline(line string, g *Game) (string, int) {
	line = strings.TrimSpace(line)
	switch {

	case strings.HasPrefix(line, "Picked up a "):
		return strings.TrimSuffix(strings.TrimPrefix(line, "Picked up a "), "."), 1

	case strings.HasPrefix(line, "You got the "):
		return strings.TrimSuffix(strings.TrimPrefix(line, "You got the "), "!"), 1

	case strings.HasPrefix(line, "Level map01 - Kills: 10/19 - Items: 8/9 - Secrets: 0/5 - Time: 0:35"):
		return "", 1

	default:
		return "", 0
	}
}
