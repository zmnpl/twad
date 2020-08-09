package games

type runconfig struct {
	loadLastSave bool
	beam         bool
	warpEpisode  int
	warpLevel    int
}

// newRunConfig returns an insance of *runconfig
func newRunConfig() *runconfig {
	var rcfg runconfig
	return &rcfg
}

func (rcfg *runconfig) quickload() *runconfig {
	rcfg.loadLastSave = true
	return rcfg
}

func (rcfg *runconfig) warp(episode, level int) *runconfig {
	rcfg.beam = true
	rcfg.warpEpisode = episode
	rcfg.warpLevel = level
	return rcfg
}

func (rcfg *runconfig) setSkill(skillLevel int) *runconfig {

	return rcfg
}

func (rcfg *runconfig) recordDemo(name string) *runconfig {

	return rcfg
}

func (rcfg *runconfig) palyDemo(name string) *runconfig {

	return rcfg
}
