package games

// just to avoid a truckload of parameters to the composing method...
type runOptionSet struct {
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

var defaultRunConfig = runOptionSet{
	//loadLastSave: false,
	//shouldWarp: false,
	//shouldBeam: false,
	//beamToMap 2,
	//recDemo: false,
	//plyDemo: false,
	//demoName: "",
	//warpEpisode: 1,
	//warpLevel: 1,
	//skill: 2,
}

// use function currying approach to set options instead of working on a pointer
// results in nice syntac for using newRunConfig(...)
type runOption func(r runOptionSet) runOptionSet

func newRunConfig(os ...runOption) runOptionSet {
	ros := defaultRunConfig
	for _, o := range os {
		ros = o(ros)
	}
	return ros
}

func quickload() runOption {
	return func(ros runOptionSet) runOptionSet {
		ros.loadLastSave = true
		return ros
	}
}

func beam(beamToMap string) runOption {
	return func(ros runOptionSet) runOptionSet {
		ros.shouldBeam = true
		ros.beamToMap = beamToMap
		return ros
	}
}

func warp(episode, level int) runOption {
	return func(ros runOptionSet) runOptionSet {
		ros.shouldWarp = true
		ros.warpEpisode = episode
		ros.warpLevel = level
		return ros
	}
}

// skill is taken for gzdoom (0-4)
// game will remap if other engine is used
func setSkill(skillLevel int) runOption {
	return func(ros runOptionSet) runOptionSet {
		ros.skill = skillLevel
		return ros
	}
}

func recordDemo(name string) runOption {
	return func(ros runOptionSet) runOptionSet {
		ros.demoName = name
		ros.recDemo = true
		return ros
	}
}

func playDemo(name string) runOption {
	return func(ros runOptionSet) runOptionSet {
		ros.demoName = name
		ros.plyDemo = true
		return ros
	}
}
