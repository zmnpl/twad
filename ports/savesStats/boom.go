package savesStats

import (
	"regexp"
	"strconv"
)

func GetBoomStats(path string) ([]MapStats, error) {
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
