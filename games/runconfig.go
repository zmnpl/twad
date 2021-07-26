package games

// just to avoid a truckload of parameters to the composing method...
type runconfig struct {
	loadLastSave bool
	shouldWarp   bool
	shouldBeam   bool
	beamToMap    string
	recDemo      bool
	plyDemo      bool
	demoName     string
	warpEpisode  int
	warpLevel    int
	skill        int
}

func newRunConfig() *runconfig {
	return &runconfig{}
}

func (rcfg *runconfig) quickload() *runconfig {
	rcfg.loadLastSave = true
	return rcfg
}

func (rcfg *runconfig) beam(beamToMap string) *runconfig {
	rcfg.shouldBeam = true
	rcfg.beamToMap = beamToMap
	return rcfg
}

func (rcfg *runconfig) warp(episode, level int) *runconfig {
	rcfg.shouldWarp = true
	rcfg.warpEpisode = episode
	rcfg.warpLevel = level
	return rcfg
}

// skill is taken for gzdoom (0-4)
// game will remap if other engine is used
func (rcfg *runconfig) setSkill(skillLevel int) *runconfig {
	rcfg.skill = skillLevel
	return rcfg
}

func (rcfg *runconfig) recordDemo(name string) *runconfig {
	rcfg.demoName = name
	rcfg.recDemo = true
	return rcfg
}

func (rcfg *runconfig) playDemo(name string) *runconfig {
	rcfg.demoName = name
	rcfg.plyDemo = true
	return rcfg
}
