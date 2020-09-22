package games

import (
	"archive/zip"
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/bits"
	"os"
	"regexp"
	"strconv"
	"unsafe"
)

// nested structs match json structure of zdoom json format
type sg struct {
	Stats Savegame `json:"statistics"`
}

// Savegame saves a dual purpose
// Is a nested struct piece to easily parse zdoom savegame stats
// Represents a games savegame and some of it's properties
type Savegame struct {
	FI        os.FileInfo
	Directory string
	Name      string
	Slot      int
	Levels    []MapStats `json:"levels"`
}

// MapStats contains the stats for one single level read from a savegame
type MapStats struct {
	//RecordTime   time.Time
	TotalKills   uint32 `json:"totalkills"`
	KillCount    uint32 `json:"killcount"`
	TotalSecrets uint32 `json:"totalsecrets"`
	SecretCount  uint32 `json:"secretcount"`
	LevelTime    uint32 `json:"leveltime"`

	//
	TotalItems uint32 `json:"totalitems"`
	ItemCount  uint32 `json:"itemcount"`
	LevelName  string `json:"levelname"`
}

// NewSavegame initializes a new Savegame struct
func NewSavegame(fi os.FileInfo, dir string) Savegame {
	savegame := Savegame{
		Directory: dir,
	}
	if fi != nil {
		savegame.FI = fi
	}
	return savegame
}

func getZDoomStats(path string) []MapStats {
	stats, err := zdoomStatsFromJSON(path)
	if err == nil {
		return stats
	}

	stats, err = zdoomStatsFromBinary(path)
	if err == nil {
		return stats
	}

	return make([]MapStats, 0)
}

// ZDOOM

func zdoomStatsFromJSON(path string) ([]MapStats, error) {
	jsonContent, err := getFileContentFromZip(path, "globals.json")
	if err != nil {
		return nil, err
	}

	save := sg{
		Stats: Savegame{
			Directory: path,
		},
	}

	if err := json.Unmarshal(jsonContent, &save); err != nil {
		return nil, err
	}

	return save.Stats.Levels, nil
}

func zdoomStatsFromBinary(path string) ([]MapStats, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// reader to read from
	contentReader := bytes.NewReader(content)

	// statistics start from here
	magicSeries := []byte("sTat")
	readFrom := BinaryStartPosition(content, magicSeries)
	if readFrom == -1 {
		return nil, fmt.Errorf("Could not find magic series: %v", magicSeries)
	}
	contentReader.Seek(int64(readFrom), 0)

	// count of levels in savegame
	levelCountBytes := make([]byte, 4)
	contentReader.Read(levelCountBytes)
	levelCount := binary.BigEndian.Uint32(levelCountBytes)
	levelCount = reverseBitsIfNeeded(levelCount)

	// golang has buildin functionality for this; though it doesn't use reader
	// size, bs := binary.Uvarint(content[readFrom+4:])
	size := ReadSize(contentReader)
	skip := make([]byte, size-1)
	contentReader.Read(skip)

	maps := make([]MapStats, 0)
	for i := 0; uint32(i) < levelCount; i++ {
		currentMap := readLevelStats(contentReader)

		// position + 1; skip NEW_NAME (27) (like DOOMLAUNCHER)
		contentReader.ReadByte()

		// level name is the last piece
		size := ReadSize(contentReader)
		levelNameBytes := make([]byte, size-1)
		contentReader.Read(levelNameBytes)
		currentMap.LevelName = string(levelNameBytes)

		maps = append(maps, currentMap)
	}

	return maps, nil
}

// BinaryStartPosition returns the position after the search series has been found
func BinaryStartPosition(binaryData []byte, startAfterSeries []byte) int {
	seriesLength := len(startAfterSeries)
	readFrom := -1
	for index, _ := range binaryData {
		if bytes.Equal(startAfterSeries, binaryData[index:index+seriesLength]) {
			readFrom = index + seriesLength
			break
		}
	}
	return readFrom
}

// ReadSize gets the size of coming string in variable length encoding
// Here: looks like reversed order; lowest bits are first bytes
//
// https://en.wikipedia.org/wiki/Variable-length_quantity
func ReadSize(reader io.ReadSeeker) int {
	b := make([]byte, 1)
	count := 0
	ofset := 0

	for {
		reader.Read(b)
		count = count | (int(b[0])&0x7f)<<ofset
		ofset = 7

		// Checks if the MSB is 0
		if (int(b[0]) & 0x80) == 0 {
			break
		}
	}

	return count
}

func readLevelStats(reader io.Reader) MapStats {
	lvlStats := MapStats{}

	// make byte slice of length 4
	// is used for the other reads as well...
	b := make([]byte, unsafe.Sizeof(lvlStats.KillCount))

	reader.Read(b)
	lvlStats.TotalKills = binary.BigEndian.Uint32(b)
	lvlStats.TotalKills = reverseBitsIfNeeded(lvlStats.TotalKills)

	reader.Read(b)
	lvlStats.KillCount = binary.BigEndian.Uint32(b)
	lvlStats.KillCount = reverseBitsIfNeeded(lvlStats.KillCount)

	reader.Read(b)
	lvlStats.TotalSecrets = binary.BigEndian.Uint32(b)
	lvlStats.TotalSecrets = reverseBitsIfNeeded(lvlStats.TotalSecrets)

	reader.Read(b)
	lvlStats.SecretCount = binary.BigEndian.Uint32(b)
	lvlStats.SecretCount = reverseBitsIfNeeded(lvlStats.SecretCount)

	reader.Read(b)
	lvlStats.LevelTime = binary.BigEndian.Uint32(b)
	lvlStats.LevelTime = reverseBitsIfNeeded(lvlStats.LevelTime)

	return lvlStats
}

func reverseBitsIfNeeded(i uint32) uint32 {
	if i > 0x0000FFFF {
		i = bits.Reverse32(i)
	}
	return i
}

func getFileContentFromZip(src string, fileName string) ([]byte, error) {
	r, err := zip.OpenReader(src)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	for _, f := range r.File {
		if f.Name == fileName {
			content := make([]byte, 0, f.FileInfo().Size())
			contentBuffer := bytes.NewBuffer(content)

			rc, err := f.Open()
			if err != nil {
				return nil, err
			}

			_, err = io.Copy(contentBuffer, rc)

			rc.Close()

			if err != nil {
				return nil, err
			}
			return contentBuffer.Bytes(), nil
		}
	}
	return nil, err
}

// BOOM

func getBoomStats(path string) ([]MapStats, error) {
	lines, err := fileLines(path)
	if err != nil {
		return nil, err
	}

	boomStatsRege := regexp.MustCompile(`(?P<mapName>.*?)\s+-\s+` +
		`(?P<timeMinutes>\d+):(?P<timeSeconds>\d+)\.(?P<timeMilliSeconds>\d+) \(` +
		`(?P<timeParMinutes>\d+):(?P<timeParSeconds>\d+)\)\s+` +
		`K:\s+(?P<kills>\d+)/(?P<killsTotal>\d+)\s+` +
		`I:\s+(?P<items>\d+)/(?P<itemsTotal>\d+)\s+` +
		`S:\s+(?P<secrets>\d+)/(?P<secretsTotal>\d+)`)

	maps := make([]MapStats, 0)
	for _, l := range lines {
		lvlMap := reSubMatchMap(boomStatsRege, l)
		currentMap := MapStats{}

		if mapName, ok := lvlMap["mapName"]; ok {
			currentMap.LevelName = mapName
		}
		if killcount, err := strconv.Atoi(lvlMap["kills"]); err == nil {
			currentMap.KillCount = uint32(killcount)
		}
		if totalkills, err := strconv.Atoi(lvlMap["killsTotal"]); err == nil {
			currentMap.TotalKills = uint32(totalkills)
		}
		if itemCount, err := strconv.Atoi(lvlMap["items"]); err == nil {
			currentMap.ItemCount = uint32(itemCount)
		}
		if totalItems, err := strconv.Atoi(lvlMap["itemsTotal"]); err == nil {
			currentMap.TotalItems = uint32(totalItems)
		}
		if secretCount, err := strconv.Atoi(lvlMap["secrets"]); err == nil {
			currentMap.SecretCount = uint32(secretCount)
		}
		if totalSecrets, err := strconv.Atoi(lvlMap["secretsTotal"]); err == nil {
			currentMap.TotalSecrets = uint32(totalSecrets)
		}

		maps = append(maps, currentMap)
	}

	return maps, nil
}

func fileLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, nil
}

// Crispy / Chocolate
func getChocolateStats(path string) ([]MapStats, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	chocolateStatsRegex := regexp.MustCompile(`(?s)=+.*?(?P<mapName>\w+).*?=+\s*?` +
		`Time: (?P<timeMinutes>\d+):(?P<timeSeconds>\d+) \(par: (?P<timeParMinutes>\d+):(?P<timeParSeconds>\d+)\)\s*?` +
		`Player.*?` +
		`Kills: (?P<kills>\d+) \/ (?P<killsTotal>\d+) .*?` +
		`Items: (?P<items>\d+) \/ (?P<itemsTotal>\d+).*?` +
		`Secrets: (?P<secrets>\d+) \/ (?P<secretsTotal>\d+) .*?\)`)

	maps := make([]MapStats, 0)
	matchedMaps := chocolateStatsRegex.FindAllString(string(content), -1)
	for _, v := range matchedMaps {
		lvlMap := reSubMatchMap(chocolateStatsRegex, v)
		currentMap := MapStats{}

		if mapName, ok := lvlMap["mapName"]; ok {
			currentMap.LevelName = mapName
		}
		if killcount, err := strconv.Atoi(lvlMap["kills"]); err == nil {
			currentMap.KillCount = uint32(killcount)
		}
		if totalkills, err := strconv.Atoi(lvlMap["killsTotal"]); err == nil {
			currentMap.TotalKills = uint32(totalkills)
		}
		if itemCount, err := strconv.Atoi(lvlMap["items"]); err == nil {
			currentMap.ItemCount = uint32(itemCount)
		}
		if totalItems, err := strconv.Atoi(lvlMap["itemsTotal"]); err == nil {
			currentMap.TotalItems = uint32(totalItems)
		}
		if secretCount, err := strconv.Atoi(lvlMap["secrets"]); err == nil {
			currentMap.SecretCount = uint32(secretCount)
		}
		if totalSecrets, err := strconv.Atoi(lvlMap["secretsTotal"]); err == nil {
			currentMap.TotalSecrets = uint32(totalSecrets)
		}

		maps = append(maps, currentMap)
	}

	return maps, nil
}

func reSubMatchMap(r *regexp.Regexp, str string) map[string]string {
	match := r.FindStringSubmatch(str)
	subMatchMap := make(map[string]string)
	for i, name := range r.SubexpNames() {
		if i != 0 {
			subMatchMap[name] = match[i]
		}
	}
	return subMatchMap
}
