package games

import (
	"archive/zip"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/bits"
	"unsafe"
)

// nested structs match json structure of zdoom json format

type saveGame struct {
	Stats SaveStats `json:"statistics"`
}

// SaveStats contains stats for all levels of a savegame
type SaveStats struct {
	Savegame string
	Levels   []LevelStats `json:"levels"`
}

// LevelStats contains the stats for one single level read from a savegame
type LevelStats struct {
	//RecordTime   time.Time
	TotalKills   uint32 `json:"totalkills"`
	KillCount    uint32 `json:"killcount"`
	TotalSecrets uint32 `json:"totalsecrets"`
	SecretCount  uint32 `json:"secretcount"`
	LevelTime    uint32 `json:"leveltime"`

	//
	TotalItems int    `json:"totalitems"`
	ItemCount  int    `json:"itemcount"`
	LevelName  string `json:"levelname"`
}

func getStatsFromSavegame(path string) SaveStats {
	sls, err := zdoomStatsFromJSON(path)
	if err == nil {
		return sls
	}

	sls, err = zdoomStatsFromBinary(path)
	if err == nil {
		return sls
	}

	return SaveStats{}
}

// ZDOOM

func zdoomStatsFromJSON(path string) (SaveStats, error) {
	jsonContent, err := getFileContentFromZip(path, "globals.json")
	if err != nil {
		return SaveStats{}, err
	}

	save := saveGame{
		Stats: SaveStats{
			Savegame: path,
		},
	}

	save.Stats.Savegame = path
	json.Unmarshal(jsonContent, &save)

	return save.Stats, nil
}

func zdoomStatsFromBinary(path string) (SaveStats, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return SaveStats{}, err
	}

	levelStats := SaveStats{
		Savegame: path,
	}

	// reader to read from
	contentReader := bytes.NewReader(content)

	// statistics start from here
	magicSeries := []byte("sTat")
	readFrom := BinaryStartPosition(content, magicSeries)
	if readFrom == -1 {
		return SaveStats{}, fmt.Errorf("Could not find magic series: %v", magicSeries)
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

	for i := 0; uint32(i) < levelCount; i++ {
		levelStats := readLevelStats(contentReader)

		// position + 1; skip NEW_NAME (27) (like DOOMLAUNCHER)
		contentReader.ReadByte()

		// level name is the last piece
		size := ReadSize(contentReader)
		levelNameBytes := make([]byte, size-1)
		contentReader.Read(levelNameBytes)
		levelStats.LevelName = string(levelNameBytes)
	}

	return levelStats, nil
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

func readLevelStats(reader io.Reader) LevelStats {
	lvlStats := LevelStats{}

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

// TODO
// BOOM
// Crispy
