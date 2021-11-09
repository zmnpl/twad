package games

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/zmnpl/twad/base"
	"github.com/zmnpl/twad/helper"
	"github.com/zmnpl/twad/ports"
	st "github.com/zmnpl/twad/ports/savesStats"
)

// Game represents one game configuration
type Game struct {
	Name             string         `json:"name"`
	Port             string         `json:"source_port"`
	Iwad             string         `json:"iwad"`
	Environment      []string       `json:"environment"`
	Mods             []string       `json:"mods"`
	CustomParameters []string       `json:"custom_parameters"`
	ConsoleStats     map[string]int `json:"stats"`
	Playtime         int64          `json:"playtime"`
	LastPlayed       string         `json:"last_played"`
	SaveGameCount    int            `json:"save_game_count"`
	Rating           int            `json:"rating"`
	AddEdit          time.Time      `json:"added"`
	Link             string         `json:"link"`
	PersonalPortCfg  bool           `json:"own_source_port_cfg"`
	NoDeh            bool           `json:"no_deh"`
	SharedConfig     string         `json:"shared_config"`
	Stats            []st.MapStats
	StatsTotal       st.MapStats
	Savegames        []*st.Savegame
}

// NewGame creates new instance of a game
func NewGame(name, sourceport, sharedConfig, iwad string) Game {
	game := Game{
		Name:             name,
		Port:             "gzdoom",
		Iwad:             "doom2.wad",
		Environment:      make([]string, 0),
		CustomParameters: make([]string, 0),
		Mods:             make([]string, 0),
		ConsoleStats:     make(map[string]int),
		AddEdit:          time.Now(),
	}

	// replace with given
	if sourceport != "" {
		game.Port = sourceport
	}
	if iwad != "" {
		game.Iwad = iwad
	}

	return game
}

// Checks if FILE is a DeHacked file
// by comparing its last three letters to "deh" or "DEH"
func isDehFile(file string) bool {
	s := file[len(file)-3:]

	// Lowercase or uppercase are both common
	if (s == "deh" || s == "DEH") {
		return true
	}

	return false
}

// Run executes given configuration and launches the mod
// Just a wrapper for game.run
func (g *Game) Run() (err error) {
	g.run(newRunConfig())
	return
}

// Quickload starts the game from it's last savegame
// Just a wrapper for game.run
func (g *Game) Quickload() (err error) {
	g.run(newRunConfig(quickload()))
	return
}

// Warp lets you select episode and level to start in
// Just a wrapper for game.run
func (g *Game) Warp(episode, level, skill int) (err error) {
	g.run(newRunConfig(
		warp(episode, level),
		setSkill(ports.AdjustedSkill(g.Port, skill))))
	return
}

// WarpRecord lets you select episode and level to start in
// Just a wrapper for game.run
func (g *Game) WarpRecord(episode, level, skill int, demoName string) (err error) {
	g.run(newRunConfig(
		warp(episode, level),
		setSkill(ports.AdjustedSkill(g.Port, skill)),
		recordDemo(demoName)))
	return
}

// GoToMap lets you select a specific map from a mod based on it's name
// Just a wrapper for game.run
func (g *Game) GoToMap(mapName string, skill int) (err error) {
	g.run(newRunConfig(
		goToMap(mapName),
		setSkill(ports.AdjustedSkill(g.Port, skill))))
	return
}

// GoToMapRecord lets you select a specific map from a mod based on it's name
// Just a wrapper for game.run
func (g *Game) GoToMapRecord(mapName string, skill int, demoName string) (err error) {
	g.run(newRunConfig(
		goToMap(mapName),
		setSkill(ports.AdjustedSkill(g.Port, skill)),
		recordDemo(demoName)))
	return
}

// PlayDemo replays the given demo file
// Wrapper for game.run
func (g *Game) PlayDemo(name string) {
	g.run(newRunConfig(playDemo(name)))
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

func (g *Game) run(rcfg runOptionSet) (err error) {
	start := time.Now()

	// change working directory to redirect stat file output
	// for boom ports only
	wd, wdChangeError := os.Getwd()
	if ports.Family(g.Port) == ports.Boom {
		os.Chdir(g.getSaveDir())
	}

	// rip and tear!
	doom := g.composeProcess(g.getLaunchParams(rcfg))
	output, err := doom.CombinedOutput()
	if err != nil {
		os.WriteFile("twad.log", []byte(fmt.Sprintf("%v\n\n%v\n\n%v\n\n%v", string(output), err.Error(), g.getLaunchParams(rcfg), doom)), 0755)
		return err
	}

	// change back working directory to where it was
	wdNow, _ := os.Getwd()
	if wd != wdNow && wdChangeError == nil {
		os.Chdir(wd)
	}

	playtime := time.Since(start).Milliseconds()
	g.Playtime = g.Playtime + playtime
	g.LastPlayed = time.Now().Format("2006-01-02 15:04:05MST")

	// could take a while ...
	go g.processOutput(string(output))
	go g.ReadLatestStats()

	return
}

func (g *Game) composeProcess(params []string) (cmd *exec.Cmd) {
	// create process object
	cmd = exec.Command(g.Port, params...)
	// add environment variables; use os environment as basis
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, g.Environment...)
	return
}

func (g *Game) getLaunchParams(rcfg runOptionSet) []string {
	params := make([]string, 0, 10)

	// IWAD
	if g.Iwad != "" {
		params = append(params, "-iwad", g.Iwad) // -iwad seems to be universal across zdoom, boom and chocolate doom
	}

	// some mapsets need both .wad files (wad, pk3, pke, etc.) and DeHacked files (.deh) to work on some ports
	// but the former needs a -file parameter, while the latter needs a -deh parameter
	// so we split g.Mods in two by using the isDehFile() function
	modsWads := []string{}
	modsDeh  := []string{}

	if len(g.Mods) > 0 {
		for _, file := range g.Mods {
			if isDehFile(file) {
				// some ports (e.g. Woof!) don't auto check for .deh files on DOOMWADDIR
				modsDeh = append(modsDeh, os.Getenv("DOOMWADDIR") + "/" + file)
			} else {
				modsWads = append(modsWads, file)
			}
		}
	}

	// mods 
	if len(modsWads) > 0 {
		params = append(params, "-file") // -file seems to be universal across zdoom, boom and chocolate doom
		params = append(params, modsWads...)
	}

	// some mods need this
	if g.NoDeh {
		params = append(params, "-nodeh")
	}

	// DeHacked files
	if len(modsDeh) > 0 {
		params = append(params, "-deh")
		params = append(params, modsDeh...)
	}

	// custom game save directory
	// making dir seems to be redundant, since engines do that already
	// still keeping it to possibly keep track of it / handle errors
	// only use separate save dir if directory has been craeted or path exists already
	if err := os.MkdirAll(g.getSaveDir(), 0755); err == nil {
		params = append(params, ports.SaveDirParam(g.Port))
		params = append(params, g.getSaveDir())
	}

	// custom config directory
	if g.PersonalPortCfg {
		if err := os.MkdirAll(g.getConfigDir(), 0755); err == nil {
			params = append(params, "-config")
			params = append(params, filepath.Join(g.getConfigDir(), ports.CanonicalName(g.Port)+ports.ConfigFileExtension(g.Port)))
		}
	} else if g.SharedConfig != "" {
		if err := os.MkdirAll(base.PortSharedConfigPath(g.Port), 0755); err == nil {
			params = append(params, "-config")
			params = append(params, filepath.Join(base.PortSharedConfigPath(g.Port), g.SharedConfig))
		}
	}

	// stats for zdoom on windows
	if ports.Family(g.Port) == ports.Zdoom && runtime.GOOS == "windows" {
		params = append(params, "-stdout")
	}

	// stats for chocolate doom and ports
	if ports.Family(g.Port) == ports.Chocolate {
		params = append(params, "-statdump")
		params = append(params, path.Join(g.getSaveDir(), "statdump.txt"))
	}

	// stats for chocolate doom and ports
	if ports.Family(g.Port) == ports.Boom {
		params = append(params, "-levelstat")
	}

	// quickload
	if rcfg.loadLastSave {
		params = append(params, g.getLastSaveLaunchParams()...)
	}

	// start a specific map
	// either by warping to a episode / level
	// or by
	if ports.Family(g.Port) == ports.Zdoom && rcfg.shouldGoToMap && len(rcfg.goToMap) > 0 {
		// first try mod map by name (only zdoom for now)
		params = append(params, "+map")
		params = append(params, rcfg.goToMap)

		// add skill
		params = append(params, "-skill")
		params = append(params, strconv.Itoa(rcfg.skill))
	} else if rcfg.shouldWarp && (rcfg.warpEpisode > 0 || rcfg.warpLevel > 0) {
		// warp
		// only warp if no specific map has been selected
		params = append(params, "-warp")
		if rcfg.warpEpisode > 0 {
			params = append(params, strconv.Itoa(rcfg.warpEpisode))
		}
		if rcfg.warpLevel > 0 {
			params = append(params, strconv.Itoa(rcfg.warpLevel))
		}

		// add skill
		params = append(params, "-skill")
		params = append(params, strconv.Itoa(rcfg.skill))
	}

	// demo recording
	if rcfg.recDemo {
		if err := os.MkdirAll(g.getDemoDir(), 0755); err == nil {
			params = append(params, "-record") // TODO: Does -record behave equally across ports?
			params = append(params, g.getDemoDir()+"/"+rcfg.demoName)
		}
	}

	// play demo
	if rcfg.plyDemo {
		params = append(params, "-playdemo")
		params = append(params, g.getDemoDir()+"/"+rcfg.demoName)
	}

	return append(params, g.CustomParameters...)
}

func (g *Game) getLastSaveLaunchParams() (params []string) {
	params = []string{}

	if lastSave, err := g.lastSave(); err == nil {
		params = append(params, []string{"-loadgame", lastSave}...) // -loadgame seems to be universal across zdoom, boom and chocolate doom
	}
	return
}

// CommandList returns the full slice of strings in order to launch the game
func (g *Game) CommandList() (command []string) {
	command = g.Environment
	command = append(command, g.Port)
	command = append(command, g.getLaunchParams(newRunConfig(quickload()))...)
	return
}

// SaveCount returns the number of savegames existing for this game
func (g *Game) SaveCount() int {
	if saves, err := g.savegameFiles(); err == nil {
		return len(saves)
	}
	return 0
}

// savegameFiles returns a slice of os.FileInfo with all savegmes for this game
func (g *Game) savegameFiles() ([]os.DirEntry, error) {
	saves, err := os.ReadDir(g.getSaveDir())
	if err != nil {
		return nil, err
	}
	saves = helper.FilterExtensions(saves, ports.SaveFileExtension(g.Port), false)

	sort.Slice(saves, func(i, j int) bool {
		foo, _ := saves[i].Info()
		bar, _ := saves[j].Info()
		return foo.ModTime().After(bar.ModTime())
	})

	return saves, nil
}

// LoadSavegames returns a slice of Savegames for the game
func (g *Game) LoadSavegames() []*st.Savegame {
	saveDir := g.getSaveDir()
	savegames := make([]*st.Savegame, 0)
	savegameFiles, _ := g.savegameFiles()

	for i, s := range savegameFiles {
		savegame := st.NewSavegame(s, saveDir)
		g.loadSaveMeta(&savegame)

		// load stats
		// after 3 load parallel
		if i <= 2 {
			g.loadSaveStats(&savegame)
		} else {
			go g.loadSaveStats(&savegame)
		}
		savegames = append(savegames, &savegame)
	}
	g.Savegames = savegames
	return savegames
}

func (g *Game) loadSaveMeta(s *st.Savegame) {
	s.Meta = g.GetSaveMeta(path.Join(s.Directory, s.FI.Name()))
}

func (g *Game) loadSaveStats(s *st.Savegame) {
	s.Levels = g.GetStats(path.Join(s.Directory, s.FI.Name()))
}

// GetSaveMeta reads meta information for the given savegame
func (g *Game) GetSaveMeta(savePath string) st.SaveMeta {
	if ports.Family(g.Port) == ports.Chocolate {
		meta, _ := st.ChocolateMetaFromBinary(savePath)
		return meta
	} else if ports.Family(g.Port) == ports.Boom {
		meta, _ := st.ChocolateMetaFromBinary(savePath)
		return meta
	}

	return st.GetZDoomSaveMeta(savePath)
}

// GetStats reads stats from the given savegame path for zdoom ports
// If the port is boom or chocolate, their respective dump-files are used
func (g *Game) GetStats(savePath string) []st.MapStats {
	var stats []st.MapStats
	if ports.Family(g.Port) == ports.Chocolate {
		stats, _ = st.GetChocolateStats(path.Join(g.getSaveDir(), "statdump.txt"))
	} else if ports.Family(g.Port) == ports.Boom {
		stats, _ = st.GetBoomStats(path.Join(g.getSaveDir(), "levelstat.txt"))
	} else {
		stats = st.GetZDoomStats(savePath)
	}

	return stats
}

// ReadLatestStats tries to read stats from the newest existing savegame
func (g *Game) ReadLatestStats() {
	lastSavePath, _ := g.lastSave()
	g.Stats = g.GetStats(lastSavePath)
	g.StatsTotal = st.SummarizeStats(g.Stats)
}

// DemoCount returns the number of demos existing for this game
func (g *Game) DemoCount() int {
	if demos, err := os.ReadDir(g.getDemoDir()); err == nil {
		return len(demos)
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

func (g *Game) getSaveDir() string {
	return filepath.Join(base.GetSavegameFolder(), g.cleansedName())
}

func (g *Game) getConfigDir() string {
	return filepath.Join(base.GetGameConfigFolder(), g.cleansedName())
}

// lastSave returns the the file name or slotnumber (depending on source port) for the game
func (g *Game) lastSave() (save string, err error) {
	saveDir := g.getSaveDir()
	saves, err := os.ReadDir(saveDir)
	if err != nil {
		return
	}

	// assume zdoom
	portSaveFileExtension := ports.SaveFileExtension(g.Port)

	// find the newest file
	newestTime, _ := time.Parse(time.RFC3339, "1900-01-01T00:00:00+00:00")
	for _, file := range saves {
		extension := strings.ToLower(filepath.Ext(file.Name()))
		fi, err := file.Info()
		if err != nil {
			continue
		}
		if fi.Mode().IsRegular() && fi.ModTime().After(newestTime) && extension == portSaveFileExtension {
			save = filepath.Join(saveDir, file.Name())
			newestTime = fi.ModTime()
		}
	}

	// adjust for different souce ports
	save = ports.SaveGameName(g.Port, save)

	if save == "" {
		err = os.ErrNotExist
	}

	return
}

func (g Game) ModMaps() map[string]string {
	maps := make(map[string]string)

	// check all mods
	for _, v := range g.Mods {
		// pk3s
		if strings.HasSuffix(strings.ToLower(v), ".pk3") {
			mapCounter := 1
			lines, _ := base.GetFileLinesFromPK3(filepath.Join(base.Config().WadDir, v), "mapinfo")
			for _, l := range lines {
				// example line
				// map aeon22 "Decayed and Conquered"
				if strings.HasPrefix(l, "map") {
					fields := strings.Split(l, " ")
					if len(fields) >= 3 {
						maps[fmt.Sprintf("%02d %v", mapCounter, strings.Trim(fields[2], "\""))] = fields[1]
						mapCounter += 1
					}
				}
			}
		}
	}

	return maps
}

// Demos returns the demo files existing for the game
func (g *Game) Demos() ([]os.DirEntry, error) {
	demos, err := os.ReadDir(g.getDemoDir())
	if err != nil {
		return nil, err
	}
	sort.Slice(demos, func(i, j int) bool {
		foo, err := demos[i].Info()
		if err != nil {
			return false
		}
		bar, err := demos[j].Info()
		if err != nil {
			return true
		}
		return foo.ModTime().After(bar.ModTime())
	})
	return demos, err
}

// DemoExists checks if a file with the same name already exists in the default demo dir
// Doesn't use standard library to ignore file ending; design decision
func (g *Game) DemoExists(name string) bool {
	if files, err := os.ReadDir(g.getDemoDir()); err == nil {
		for _, f := range files {
			nameWithouthExt := strings.TrimSuffix(f.Name(), filepath.Ext(f.Name()))
			if nameWithouthExt == name {
				return true
			}
		}
	}
	return false
}

// RemoveDemo removes the demo file with the given name
// and returns the new set of demos
func (g *Game) RemoveDemo(name string) ([]os.DirEntry, error) {
	err := os.Remove(filepath.Join(g.getDemoDir(), name))
	if err != nil {
		return nil, err
	}
	return g.Demos()
}

func (g *Game) getDemoDir() string {
	return filepath.Join(base.GetDemoFolder(), g.cleansedName())
}

// cleansedName removes all but alphanumeric characters from name
// used for directory names
func (g *Game) cleansedName() string {
	cleanser, _ := regexp.Compile("[^a-zA-Z0-9]+")
	return cleanser.ReplaceAllString(g.Name, "")
}

// processOutput processes the terminal output of the zdoom port
func (g *Game) processOutput(output string) {
	if g.ConsoleStats == nil {
		g.ConsoleStats = make(map[string]int)
	}
	for _, v := range strings.Split(output, "\n") {
		if stat, increment := parseStatline(v, g); stat != "" {
			g.ConsoleStats[stat] = g.ConsoleStats[stat] + increment
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

// Printing Methods
// String returns the string which is run when running

// RatingString returns the string resulting from the games rating
func (g *Game) RatingString() string {
	return strings.Repeat("*", g.Rating) + strings.Repeat("-", 5-g.Rating)
}

// EnvironmentString returns a join of all prefix parameters
func (g *Game) EnvironmentString() string {
	return strings.TrimSpace(strings.Join(g.Environment, " "))
}

// ParamsString returns a join of all prefix parameters
func (g *Game) ParamsString() string {
	return strings.TrimSpace(strings.Join(g.CustomParameters, " "))
}
