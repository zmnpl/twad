package games

import (
	"encoding/json"
	"fmt"
	"github.com/zmnpl/twad/cfg"
	"io/ioutil"
	"os/exec"
	"strings"
	"sync"
	"time"
)

type GameList []Game

// Game represents one game configuration
type Game struct {
	Name       string         `json:"name"`
	SourcePort string         `json:"sourceport"`
	Iwad       string         `json:"iwad"`
	Mods       []string       `json:"mods"`
	Params     []string       `json:"params"`
	Stats      map[string]int `json:"stats"`
	Playtime   int64          `json:"playtime"`
}

// Games holds TODO
var (
	once            sync.Once
	gamesInstance   GameList
	config          *cfg.Cfg
	changeListeners []func()
	gamesJSONName   = "games.json"
)

// TODO: make singleton GetInstance

func init() {
	config = cfg.GetInstance()
	GetInstance()
	changeListeners = make([]func(), 0)
}

// GetInstance sets up and returns the singleton instance of games
func GetInstance() GameList {
	once.Do(func() {
		gamesInstance = make(GameList, 0, 0)
		loadGames()
	})
	return gamesInstance
}

// c := make(chan func() error, 10)
// go foo(c)
// QueuePersist = func() {
// 	c <- Persist
// }
// var QueuePersist func()
// func foo(bar chan func() error) {
// 	for f := range bar {
// 		f()
// 	}
//}

func informChangeListeners() {
	for _, f := range changeListeners {
		f()
	}
}

func RegisterChangeListener(f func()) {
	changeListeners = append(changeListeners, f)
}

// AddGame adds a game to the list
// this triggers the list to be written to disk as well
func AddGame(g Game) {
	gamesInstance = append(gamesInstance, g)
	informChangeListeners()
	Persist() // TODO: Could be done in a goroutine; Maybe queue via channel
}

// RemoveGameAt removes the game at the given index
func RemoveGameAt(i int) {
	gamesInstance = append(gamesInstance[:i], gamesInstance[i+1:]...)
	informChangeListeners()
	Persist()
}

// NewGame creates new instance of a game
func NewGame(name, sourceport, iwad string) Game {
	var mod Game
	mod.Name = name
	// TODO: default from config
	mod.SourcePort = "/usr/bin/gzdoom"
	if sourceport != "" {
		mod.SourcePort = sourceport
	}
	mod.Iwad = "doom2.wad"
	if iwad != "" {
		mod.Iwad = iwad
	}
	mod.Mods = make([]string, 0)
	mod.Params = make([]string, 0)
	mod.Stats = make(map[string]int)

	return mod
}

// MaxModCount returns the biggest number of mods for a single games
// this is useful for table creation, to know how many colums one needs
func MaxModCount() int {
	maxCnt := 0
	for _, g := range gamesInstance {
		if len(g.Mods) > maxCnt {
			maxCnt = len(g.Mods)
		}
	}
	return maxCnt + 1
}

// GameCount returns the number of games available
func GameCount() int {
	return len(gamesInstance)
}

// AddMod adds mod
func (g *Game) AddMod(modFile string) {
	g.Mods = append(g.Mods, modFile)
	informChangeListeners()
	Persist()
}

// Run executes given configuration and launches the mod
func (g *Game) Run() error {
	var params []string

	if g.Iwad != "" {
		params = append(params, "-iwad", g.Iwad)
	}

	if len(g.Mods) > 0 {
		params = append(params, "-file")
		params = append(params, g.Mods...)
	}

	start := time.Now()

	// execute and capture output
	proc := exec.Command(g.SourcePort, params...)
	output, err := proc.CombinedOutput()
	if err != nil {
		return err
	}

	playtime := time.Since(start).Milliseconds()
	g.Playtime = g.Playtime + playtime

	processOutput(string(output), g)

	// Call Persist to write stats
	Persist()

	return nil
}

func processOutput(output string, g *Game) {
	for _, v := range strings.Split(output, "\n") {
		if stat, increment := parseStatline(v, g); stat != "" {
			g.Stats[stat] = g.Stats[stat] + increment
		}
	}
}

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

// String returns the string which is run when running
func (g Game) String() string {
	var iwad string
	if g.Iwad != "" {
		iwad = fmt.Sprintf(" -iwad %s", g.Iwad)
	}

	var mods string
	if len(g.Mods) > 0 {
		mods = fmt.Sprintf(" -file %s", strings.Trim(strings.Join(g.Mods, " "), " "))
	}

	return fmt.Sprintf("%s%s%s", g.SourcePort, iwad, mods)
}

// Persist writes all games into the according JSON file
func Persist() error {
	gamesJSON, err := json.MarshalIndent(gamesInstance, "", "    ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(cfg.GetConfigFolder()+"/"+gamesJSONName, gamesJSON, 0755)
	if err != nil {
		return err
	}
	return nil
}

func loadGames() error {
	content, err := ioutil.ReadFile(cfg.GetConfigFolder() + "/" + gamesJSONName)
	if err != nil {
		return err
	}

	err = json.Unmarshal(content, &gamesInstance)
	if err != nil {
		return err
	}

	return nil
}
