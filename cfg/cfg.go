package cfg

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/zmnpl/twad/helper"
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
	WadDir                 string         `json:"wad_dir"`
	WriteWadDirToEngineCfg bool           `json:"write_wad_dir_to_engine_cfg"`
	DontSetDoomwaddir      bool           `json:"dont_set_doomwaddir"`
	ModExtensions          map[string]int `json:"mod_extensions"`
	SourcePorts            []string       `json:"source_ports"`
	IWADs                  []string       `json:"iwa_ds"`
	Configured             bool           `json:"configured"`
	DeleteWithoutWarning   bool           `json:"delete_without_warning"`
	HideHeader             bool           `json:"hide_header"`
	GameListAbsoluteWidth  int            `json:"game_list_absolute_width"`
	GameListRelativeWidth  int            `json:"game_list_relative_width"`
}

func init() {
	firstStart()
	Instance()
	Persist() // just in case new settings made it into the programm
	EnableBasePath()
}

func defaultConfig() Cfg {
	config := Cfg{
		WadDir:                filepath.Join(helper.Home(), "/DOOM"),
		ModExtensions:         map[string]int{".wad": 1, ".pk3": 1, ".ipk3": 1},
		SourcePorts:           []string{"gzdoom", "zandronum", "lzdoom"},
		IWADs:                 []string{"doom2.wad", "doom.wad"},
		GameListRelativeWidth: 40,
		GameListAbsoluteWidth: 0,
	}

	// check if user has set DOOMWADDIR
	if dwd, exists := os.LookupEnv("DOOMWADDIR"); exists {
		config.WadDir = dwd
	}

	return config
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

	content, err := ioutil.ReadFile(configFullPath()) // TODO: Resolve simlinks
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

// Instance sets up and returns the singleton instance of config
func Instance() *Cfg {
	once.Do(func() {
		instance = &Cfg{}
		loadConfig()
	})
	return instance
}

// GetConfigFolder returns the folder where configuration is stored
func GetConfigFolder() string {
	return helper.Home() + configPath
}

// GetSavegameFolder returns the folder where savegames are stored
func GetSavegameFolder() string {
	return filepath.Join(GetConfigFolder(), "savegames")
}

// GetDemoFolder returns the folder where demos are stored
func GetDemoFolder() string {
	return filepath.Join(GetConfigFolder(), "demos")
}

// Persist writes all games into the according JSON file
func Persist() error {
	JSON, err := json.MarshalIndent(instance, "", "    ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filepath.Join(GetConfigFolder(), configName), JSON, 0755)
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
		go processSourcePortCfg(helper.Home() + "/.config/gzdoom/gzdoom.ini")
		go processSourcePortCfg(helper.Home() + "/.config/zandronum/zandronum.ini")
		go processSourcePortCfg(helper.Home() + "/.config/lzdoom/lzdoom.ini")
	}

	return nil
}

// ImportArchive imports given archive into a subfolder of the base path
func ImportArchive(zipPath, modName string) (err error) {
	_, err = helper.Unzip(zipPath, filepath.Join(instance.WadDir, modName))
	return
}

// Helper functions

func processSourcePortCfg(path string) {
	lines := sourcePortIniLines(path)
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

func sourcePortIniLines(path string) []string {
	lines := make([]string, 0, 1500)
	doomini, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer doomini.Close()

	scanner := bufio.NewScanner(doomini)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines
}

func configFullPath() string {
	return filepath.Join(GetConfigFolder(), configName)
}
