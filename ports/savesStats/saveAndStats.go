package savesStats

import (
	"os"
)

// Nested structs match json structure of zdoom json format
type sg struct {
	Stats   Savegame `json:"statistics"`
	Ticrate float64  `json:"ticrate"`
}

// Savegame saves a dual purpose
// Is a nested struct piece to easily parse zdoom savegame stats
// Represents a games savegame and some of it's properties
type Savegame struct {
	FI        os.FileInfo
	Directory string
	Name      string
	Slot      int
	Meta      SaveMeta
	Levels    []MapStats `json:"levels"`
}

// SaveMeta is metainformatino such as given name for a savegame
type SaveMeta struct {
	Software    string `json:"Software"`
	Engine      string `json:"Engine"`
	SaveVersion int    `json:"Save Version"`
	Title       string `json:"Title"`
	CurrentMap  string `json:"Current Map"`
	GameWAD     string `json:"Game WAD"`
	MapWAD      string `json:"Map WAD"`
	Comment     string `json:"Comment"`
	//Creation Time": "2020-09-02 21:59:18",
}

// MapStats contains the stats for one single level read from a savegame
type MapStats struct {
	//RecordTime   time.Time
	TotalKills   uint32 `json:"totalkills"`
	KillCount    uint32 `json:"killcount"`
	TotalSecrets uint32 `json:"totalsecrets"`
	SecretCount  uint32 `json:"secretcount"`
	LevelTime    uint32 `json:"leveltime"`

	TotalItems uint32 `json:"totalitems"`
	ItemCount  uint32 `json:"itemcount"`
	LevelName  string `json:"levelname"`
}

// NewSavegame initializes a new Savegame struct
func NewSavegame(fi os.DirEntry, dir string) Savegame {
	savegame := Savegame{
		Directory: dir,
	}
	if fi != nil {
		savegame.FI, _ = fi.Info()
	}
	return savegame
}

// SummarizeStats sums up a slice of stats to a total one
func SummarizeStats(stats []MapStats) (total MapStats) {
	total.LevelName = "TOTALS"
	for _, s := range stats {
		total.KillCount += s.KillCount
		total.TotalKills += s.TotalKills
		total.ItemCount += s.ItemCount
		total.TotalItems += s.TotalItems
		total.SecretCount += s.SecretCount
		total.TotalSecrets += s.TotalSecrets
		total.LevelTime += s.LevelTime
	}
	return
}

func (save Savegame) ReversedLevels() (r []MapStats) {
	r = make([]MapStats, len(save.Levels))
	copy(r, save.Levels)

	for i := 0; i < len(r)/2; i++ {
		j := len(r) - i - 1
		r[i], r[j] = r[j], r[i]
	}

	return
}
