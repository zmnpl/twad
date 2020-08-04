package games

type runconfig struct {
	quickload   bool
	warp        bool
	warpEpisode int
	warpLevel   int
}

func NewRunConfig() *runconfig {
	var rcfg runconfig
	return &rcfg
}

func (rcfg *runconfig) Quickload(quickload bool) *runconfig {
	rcfg.quickload = quickload
	return rcfg
}

func (rcfg *runconfig) Warp(warp bool, episode, level int) *runconfig {
	rcfg.warp = warp
	rcfg.warpEpisode = episode
	rcfg.warpLevel = level
	return rcfg
}

func (rcfg *runconfig) Skill(skill int) *runconfig {

	return rcfg
}

func (rcfg *runconfig) RecordDemo(name string) *runconfig {

	return rcfg
}

func (rcfg *runconfig) PlayDemo(name string) *runconfig {

	return rcfg
}
