package savesStats

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"unsafe"
)

// GetZDoomStats returns a slice of MapStats
// Tries to parse json first
// If that doesn't work falls back to old binary mode (zandronum)
func GetZDoomStats(path string) []MapStats {
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

// GetZDoomSaveMeta returns a slice of SaveMeta
// Tries to parse json first
// If that doesn't work falls back to old binary mode (zandronum)
func GetZDoomSaveMeta(path string) SaveMeta {
	meta, err := zdoomMetaFromJSON(path)
	if err == nil {
		return meta
	}

	meta, err = zdoomMetaFromBinary(path)
	if err == nil {
		return meta
	}

	return SaveMeta{
		Title: "FROM INCOMPATIBLE SOURCE PORT",
	}
}

func zdoomMetaFromJSON(path string) (SaveMeta, error) {
	meta := SaveMeta{}

	jsonContent, err := getFileContentFromZip(path, "info.json")
	if err != nil {
		return meta, err
	}

	if err := json.Unmarshal(jsonContent, &meta); err != nil {
		return meta, err
	}

	return meta, nil
}

func zdoomMetaFromBinary(path string) (SaveMeta, error) {
	meta := SaveMeta{}
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return meta, err
	}
	contentReader := bytes.NewReader(content)

	magicSeries := []byte("tEXtTitle")
	readStart := binaryStartPosition(content, magicSeries, 0)

	if -1 != readStart {
		magicSeriesEnd := []byte{0, 0, 0}
		readEnd := binaryStartPosition(content, magicSeriesEnd, readStart)
		if -1 != readEnd {
			stringLength := readEnd - readStart - 7
			if stringLength > 0 {
				if stringLength > 24 {
					stringLength = 24
				}
				contentReader.Seek(int64(readStart), io.SeekStart)
				result := make([]byte, stringLength)
				contentReader.Read(result)

				// cast to runes (utf8)
				buf := make([]rune, len(result))
				for i, b := range result {
					buf[i] = rune(b)
				}

				meta.Title = string(buf)

				return meta, nil
			}
		}
		return meta, fmt.Errorf("Could not find name end position in binay")
	}
	return meta, fmt.Errorf("Could not find name start position in binary")
}

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
	readFrom := binaryStartPosition(content, magicSeries, 0)
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
	size := readSize(contentReader)
	skip := make([]byte, size-1)
	contentReader.Read(skip)

	maps := make([]MapStats, 0)
	for i := 0; uint32(i) < levelCount; i++ {
		currentMap := readLevelStats(contentReader)

		// position + 1; skip NEW_NAME (27) (like DOOMLAUNCHER)
		contentReader.ReadByte()

		// level name is the last piece
		size := readSize(contentReader)
		levelNameBytes := make([]byte, size-1)
		contentReader.Read(levelNameBytes)
		currentMap.LevelName = string(levelNameBytes)

		maps = append(maps, currentMap)
	}

	return maps, nil
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
