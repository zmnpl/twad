package savesStats

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Crispy / Chocolate
func GetChocolateStats(path string) ([]MapStats, error) {
	content, err := os.ReadFile(path)
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

func ChocolateMetaFromBinary(path string) (SaveMeta, error) {
	fallbackName := "NA"
	if len(path) > 5 && strings.HasSuffix(strings.ToLower(path), ".dsg") {
		runes := []rune(path)
		fallbackName = fmt.Sprintf("Slot %v", string(runes[len(runes)-5:len(runes)-4]))
	}

	meta := SaveMeta{}
	meta.Title = fallbackName

	content, err := os.ReadFile(path)
	if err != nil {
		return meta, err
	}
	contentReader := bytes.NewReader(content)
	result := make([]byte, 24)

	read, err := contentReader.Read(result)
	if err != nil || read <= 0 {
		return meta, fmt.Errorf("could not read name from save file")
	}

	// cast to runes (utf8)
	buf := make([]rune, len(result))
	for i, b := range result {
		buf[i] = rune(b)
	}
	meta.Title = string(buf)

	return meta, nil
}
