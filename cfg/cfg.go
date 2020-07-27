package cfg

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"strings"
	"sync"
)

var (
	instance *Cfg
	once     sync.Once
	// KnownIwads contains names of the most common iwads
	KnownIwads = [...]string{"doom.wad", "doom2.wad", "tnt.wad", "plutonia.wad", "heretic.wad", "hexen.wad", "strive1.wad", "sve.wad", "chex.wad", "strife0.wad", "freedoom1.wad", "freedoom2.wad", "freedm.wad", "chex3.wad", "action2.wad", "harm1.wad", "hacx.wad"}
)

const (
	configName = "twad.json"
	configPath = "/.config/twad"
)

// Cfg holds basic configuration settings
// Should only be instantiated via GetInstance
type Cfg struct {
	WadDir                  string         `json:"wad_dir"`
	WriteWadDirToEngineCfg  bool           `json:"write_wad_dir_to_engine_cfg"`
	DontSetDoomwaddir       bool           `json:"dont_set_doomwaddir"`
	ModExtensions           map[string]int `json:"mod_extensions"`
	SourcePorts             []string       `json:"source_ports"`
	IWADs                   []string       `json:"iwa_ds"`
	Configured              bool           `json:"configured"`
	DefaultSaveDir          bool           `json:"default_save_dir"`
	DeleteWithoutWarning    bool           `json:"delete_without_warning"`
	HideHeader              bool           `json:"hide_header"`
	GameListAbsoluteWidth   int            `json:"game_list_absolute_width"`
	GameListRelativeWidth   int            `json:"game_list_relative_width"`
	DetailPaneSplitVertical bool           `json:"detail_pane_split_vertical"`
}

func init() {
	firstStart()
	GetInstance()
	Persist() // just in case new settings made it into the programm
	EnableBasePath()
}

func defaultConfig() Cfg {
	var dConf Cfg
	dConf.WadDir = home() + "/DOOM"
	if dwd, exists := os.LookupEnv("DOOMWADDIR"); exists {
		dConf.WadDir = dwd
	}
	dConf.ModExtensions = make(map[string]int)
	dConf.ModExtensions[".wad"] = 1
	dConf.ModExtensions[".pk3"] = 1
	dConf.ModExtensions[".ipk3"] = 1
	dConf.SourcePorts = []string{"gzdoom", "zandronum", "lzdoom"}
	dConf.IWADs = []string{"doom2.wad", "doom.wad"}
	dConf.GameListRelativeWidth = 40
	dConf.GameListAbsoluteWidth = 0

	return dConf
}

func firstStart() {
	// create directory for games and configs
	configFolder := GetSavegameFolder()
	configPath := configFullPath()
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		err := os.MkdirAll(configFolder, 0755)
		if err != nil {
			log.Fatal(err)
		}

		f, err := os.Create(configPath)
		if err != nil {
			log.Fatal(err)
		}

		defaulConfigJSON, _ := json.MarshalIndent(defaultConfig(), "", "    ")
		if _, err = f.Write([]byte(defaulConfigJSON)); err != nil {
			log.Fatal(err)
		}
		defer f.Close()
	}
}

func loadConfig() error {
	dConf := defaultConfig()

	content, err := ioutil.ReadFile(configFullPath())
	if err != nil {
		instance = &dConf
		return err
	}

	err = json.Unmarshal(content, instance)
	if err != nil {
		instance = &dConf
		return err
	}

	// check zero values for certain variables
	// empty ones do not really make sense
	// so set them to the defaults
	if instance.WadDir == "" {
		instance.WadDir = dConf.WadDir
	}

	if len(instance.ModExtensions) == 0 {
		instance.ModExtensions = dConf.ModExtensions
	}

	if len(instance.SourcePorts) == 0 {
		instance.SourcePorts = dConf.SourcePorts
	}

	if len(instance.IWADs) == 0 {
		instance.IWADs = dConf.IWADs
	}

	return nil
}

// Exported functions

// GetInstance sets up and returns the singleton instance of config
func GetInstance() *Cfg {
	once.Do(func() {
		instance = &Cfg{}
		loadConfig()
	})
	return instance
}

// GetConfigFolder returns the folder where configuration is stored
func GetConfigFolder() string {
	return home() + configPath
}

// GetSavegameFolder returns the folder where savegames are stored
func GetSavegameFolder() string {
	return GetConfigFolder() + "/savegames"
}

// Persist writes all games into the according JSON file
func Persist() error {
	JSON, err := json.MarshalIndent(instance, "", "    ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(GetConfigFolder()+"/"+configName, JSON, 0755)
	if err != nil {
		return err
	}

	return nil
}

// EnableBasePath adds the mod base path to the config ini files and/or sets it as DOOMWADDIR
// that enables the engine, to find mod files added with the -file parameter based on relative paths
func EnableBasePath() error {
	// DOOMWADDIR
	if instance.DontSetDoomwaddir == false {
		os.Setenv("DOOMWADDIR", instance.WadDir)
	}

	// Engine-Configs
	if instance.WriteWadDirToEngineCfg {
		go processEngineCfg(home() + "/.config/gzdoom/gzdoom.ini")
		go processEngineCfg(home() + "/.config/zandronum/zandronum.ini")
		go processEngineCfg(home() + "/.config/lzdoom/lzdoom.ini")
	}

	return nil
}

// Helper functions

func processEngineCfg(path string) {
	lines := configLines(path)
	// if there are not lines for the respective config that is considered ok; maybe that config is not installed
	if lines == nil {
		return
	}

	entry := "PATH=" + instance.WadDir
	// if the config already has the set path, there is nothing more to do here
	for _, l := range lines {
		if strings.Contains(l, entry) {
			return
		}
	}

	var configData bytes.Buffer

	for _, v := range lines {
		v = strings.TrimSpace(v)
		configData.WriteString(v + "\n")
		if v == "[FileSearch.Directories]" {
			configData.WriteString(entry + "\n")
		}
	}

	err := ioutil.WriteFile(path, configData.Bytes(), 0755)
	if err != nil {
		// TODO - do we want to see that error?
	}
}

func configLines(path string) []string {
	lines := make([]string, 0, 1500)
	gzdoomini, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer gzdoomini.Close()

	scanner := bufio.NewScanner(gzdoomini)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines
}

func configFullPath() string {
	return GetConfigFolder() + "/" + configName
}

func home() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr.HomeDir
}
