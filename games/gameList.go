package games

import (
	"encoding/json"
	"os"
	"sort"
	"sync"

	"github.com/zmnpl/twad/cfg"
)

// GameList holds all the games configured
type GameList []Game

var (
	once            sync.Once
	instance        GameList
	config          *cfg.Cfg
	changeListeners []func()
	gamesJSONName   = "games.json"
)

func init() {
	config = cfg.Instance()
	GetInstance()
	changeListeners = make([]func(), 0)
}

// GetInstance sets up and returns the singleton instance of games
func GetInstance() GameList {
	once.Do(func() {
		instance = make(GameList, 0, 0)
		loadGames()
	})
	return instance
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

// InformChangeListeners triggers the given function for each registered listener
func InformChangeListeners() {
	for _, f := range changeListeners {
		f()
	}
}

// RegisterChangeListener takes functions, that should get executed once
// configuration change
func RegisterChangeListener(f func()) {
	changeListeners = append(changeListeners, f)
}

// AddGame adds a game to the list
// this triggers the list to be written to disk as well
func AddGame(g Game) {
	instance = append(instance, g)
	InformChangeListeners()
	Persist() // TODO: Could be done in a goroutine; Maybe queue via channel
}

// RemoveGameAt removes the game at the given index
func RemoveGameAt(i int) {
	instance = append(instance[:i], instance[i+1:]...)
	InformChangeListeners()
	Persist()
}

// SortAlph sorts games alphabetically
func SortAlph() {
	sort.Slice(instance, func(i, j int) bool {
		return instance[i].Name < instance[j].Name
	})
	Persist()
}

// MaxModCount returns the biggest number of mods for a single game
// this is useful for table creation, to know how many colums one needs
func MaxModCount() int {
	maxCnt := 0
	for _, g := range instance {
		if len(g.Mods) > maxCnt {
			maxCnt = len(g.Mods)
		}
	}
	return maxCnt + 1
}

// GameCount returns the number of games available
func GameCount() int {
	return len(instance)
}

// Persist writes all games into the according JSON file
func Persist() error {
	gamesJSON, err := json.MarshalIndent(instance, "", "    ")
	if err != nil {
		return err
	}

	err = os.WriteFile(cfg.GetConfigFolder()+"/"+gamesJSONName, gamesJSON, 0755)
	if err != nil {
		return err
	}
	return nil
}

func loadGames() error {
	content, err := os.ReadFile(cfg.GetConfigFolder() + "/" + gamesJSONName) // TODO: Resolve simlinks
	if err != nil {
		return err
	}

	err = json.Unmarshal(content, &instance)
	if err != nil {
		return err
	}

	for i, _ := range instance {
		go instance[i].ReadLatestStats()
	}

	return nil
}
