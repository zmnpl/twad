package games

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/zmnpl/twad/cfg"
)

const (
	zdoom = iota
	chocolate
	boom
)

// Game represents one game configuration
type Game struct {
	Name          string         `json:"name"`
	SourcePort    string         `json:"source_port"`
	Iwad          string         `json:"iwad"`
	Environment   []string       `json:"environment"`
	Mods          []string       `json:"mods"`
	Parameters    []string       `json:"parameters"`
	Stats         map[string]int `json:"stats"`
	Playtime      int64          `json:"playtime"`
	LastPlayed    string         `json:"last_played"`
	SaveGameCount int            `json:"save_game_count"`
	Rating        int            `json:"rating"`
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
// Just a wrapper for game.run
func (g *Game) Run() (err error) {
	g.run(*newRunConfig())
	return
}

// Quickload starts the game from it's last savegame
// Just a wrapper for game.run
func (g *Game) Quickload() (err error) {
	g.run(*newRunConfig().quickload())
	return
}

// Warp lets you select episode and level to start in
// Just a wrapper for game.run
func (g *Game) Warp(episode, level int) (err error) {
	g.run(*newRunConfig().warp(episode, level))
	return
}

func (g *Game) run(rcfg runconfig) (err error) {
	start := time.Now()

	// rip and tear!
	doom := g.composeProcess(g.getLaunchParams(rcfg))

	output, err := doom.CombinedOutput()
	if err != nil {
		return err
	}

	playtime := time.Since(start).Milliseconds()
	g.Playtime = g.Playtime + playtime
	g.LastPlayed = time.Now().Format("2006-01-02 15:04:05MST")

	go processOutput(string(output), g)

	return
}

func (g Game) composeProcess(params []string) (cmd *exec.Cmd) {
	// execute and capture output
	cmd = exec.Command(g.SourcePort, params...)
	// add environment variables; use os environment as basis
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, g.Environment...)
	return
}

func (g Game) getLaunchParams(rcfg runconfig) []string {
	params := make([]string, 1, 10)
	config := cfg.GetInstance()

	// IWAD
	if g.Iwad != "" {
		params = append(params, "-iwad", g.Iwad) // -iwad seems to be universal across zdoom, boom and chocolate doom
	}

	// mods
	if len(g.Mods) > 0 {
		params = append(params, "-file") // -file seems to be universal across zdoom, boom and chocolate doom
		params = append(params, g.Mods...)
	}

	// custom game save directory
	if config.DefaultSaveDir == false {
		// making dir seems to be redundant, since engines do that already
		// still keeping it to possibly keep track of it / handle errors
		err := os.MkdirAll(g.getSaveDir(), 0755)

		// only use separate save dir if directory has been craeted or path exists already
		if err == nil {
			params = append(params, g.saveDirParam())
			params = append(params, g.getSaveDir())
		}
	}

	// quickload
	if rcfg.loadLastSave {
		params = append(params, g.getLastSaveLaunchParams()...)
	}

	// warp
	if rcfg.beam {
		params = append(params, "-warp")
		if rcfg.warpEpisode > 0 {
			params = append(params, strconv.Itoa(rcfg.warpEpisode))
		}
		if rcfg.warpLevel > 0 {
			params = append(params, strconv.Itoa(rcfg.warpLevel))
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
		params = append(params, []string{"-loadgame", lastSave}...) // -loadgame seems to be universal across zdoom, boom and chocolate doom
	}
	return
}

// String returns the string which is run when running
func (g Game) String() string {
	return fmt.Sprintf("%s", strings.TrimSpace(strings.Join(g.CommandList(), " ")))
}

// CommandList returns the full slice of strings in order to launch the game
//
func (g Game) CommandList() (command []string) {
	command = g.Environment
	command = append(command, g.SourcePort)
	command = append(command, g.getLaunchParams(*newRunConfig().quickload())...)
	return
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

// saveDirParam returns the right paramter key for specifying the savegame directory
// accounts for zdoom-, chocolate-doom and boom ports at the moments
func (g Game) saveDirParam() (parameter string) {
	sp := g.sourcePortFamily()

	// assume zdoom port
	// works for chocolate-doom family as well
	parameter = "-savedir"

	// boom derivates
	if sp == boom {
		parameter = "-save"
	}
	return
}

// lastSave returns the the file name or slotnumber (depending on source port) for the game
func (g Game) lastSave() (save string, err error) {
	saveDir := g.getSaveDir()
	saves, err := ioutil.ReadDir(saveDir)
	if err != nil {
		return
	}

	// retrieve source port family; do once since it is not the fastet function
	sp := g.sourcePortFamily()

	// assume zdoom
	portSaveFileExtension := ".zds"
	// chocolate doom + ports
	if sp == chocolate {
		portSaveFileExtension = ".dsg"
	}

	// find the newest file
	newestTime, _ := time.Parse(time.RFC3339, "1900-01-01T00:00:00+00:00")
	for _, file := range saves {
		extension := strings.ToLower(filepath.Ext(file.Name()))
		if file.Mode().IsRegular() && file.ModTime().After(newestTime) && extension == portSaveFileExtension {
			save = filepath.Join(saveDir, file.Name())
			newestTime = file.ModTime()
		}
	}

	// chocolate doom and boom take a slot from 0-7 as parameter
	// the filenames usually look like this: doomsav0.dsg
	// the zero is just cut out and returned
	if (sp == chocolate || sp == boom) && save != "" {
		tmp := []rune(save)
		save = string(tmp[len(tmp)-5 : len(tmp)-4])
		return
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
	if g.Stats == nil {
		g.Stats = make(map[string]int)
	}
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

// checks the games engine type by inspecting the string
// known keyphrases will be interpreted as a certain source port family
func (g Game) sourcePortFamily() (t int) {
	// assume zdoom family
	t = zdoom

	sp := strings.ToLower(g.SourcePort)

	// chocolate doom family
	if strings.Contains(sp, "crispy") || strings.Contains(sp, "chocolate") {
		t = chocolate
		return
	}

	// boom family
	if strings.Contains(sp, "boom") {
		t = boom
		return
	}

	return
}
