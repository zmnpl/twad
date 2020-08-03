package games

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
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

// RemoveMod removes mod at the given index
func (g *Game) RemoveMod(i int) {
	g.Mods = append(g.Mods[0:i], g.Mods[i+1:]...)
}

// Run executes given configuration and launches the mod
func (g *Game) Run(loadLastSave bool) error {
	params := g.getLaunchParams()
	if loadLastSave {
		params = append(params, g.getLastSaveLaunchParams()...)
	}

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
	return fmt.Sprintf("%s", strings.TrimSpace(strings.Join(g.CommandList(), " ")))
}

// CommandList returns the full slice of strings in order to launch the game
func (g Game) CommandList() []string {
	result := g.Environment
	result = append(result, g.SourcePort)
	result = append(result, g.getLaunchParams()...)
	result = append(result, g.getLastSaveLaunchParams()...)
	return result
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

// SwitchMods switches both entries within the mod slice
func (g *Game) SwitchMods(a, b int) {
	if a < len(g.Mods) && b < len(g.Mods) {
		modA := g.Mods[a]
		modB := g.Mods[b]
		g.Mods[a] = modB
		g.Mods[b] = modA
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

	if config.DefaultSaveDir == false {
		// making dir seems to be redundant, since engines do that already
		// still keeping it to possibly keep track of it / handle errors
		err := os.MkdirAll(g.getSaveDir(), 0755)

		// only use separate save dir if directory has been craeted or path exists already
		if err == nil {
			params = append(params, "-savedir") // -savedir works for zdoom and chocolate-doom derivates
			params = append(params, g.getSaveDir())
		}
	}

	// add custom parameters here
	params = append(params, g.Parameters...)

	return params
}

func (g Game) getLastSaveLaunchParams() (params []string) {
	params = []string{}

	// if the default savedir is used, it can not be made sure
	// that the savegame belongs to the game/mod combindation
	// therefore it doesn't make sense to do
	if config.DefaultSaveDir {
		return
	}

	// only if the last savegame could be determined successfully
	// otherwise params will stay empty
	if lastSave, err := g.lastSave(); err == nil {
		params = append(params, []string{"-loadgame", lastSave}...)
	}
	return
}

func (g Game) lastSave() (save string, err error) {
	saveDir := g.getSaveDir()
	saves, err := ioutil.ReadDir(saveDir)
	if err != nil {
		return
	}

	newestTime, _ := time.Parse(time.RFC3339, "1900-01-01T00:00:00+00:00")
	for _, file := range saves {
		extension := strings.ToLower(filepath.Ext(file.Name()))
		if file.Mode().IsRegular() && file.ModTime().After(newestTime) && (extension == ".zds" || extension == ".dsg") {
			save = filepath.Join(saveDir, file.Name())
			newestTime = file.ModTime()
		}
	}

	if save == "" {
		err = os.ErrNotExist
	}

	return
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
