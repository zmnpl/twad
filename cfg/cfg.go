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
)

const (
	configName = "twad.json"
	configPath = "/.config/twad"
)

// Cfg holds basic configuration settings
// Should only be instantiated via GetInstance
type Cfg struct {
	ModBasePath           string         `json:"mod_base_path,omitempty"`
	ModExtensions         map[string]int `json:"mod_extensions,omitempty"`
	SourcePorts           []string       `json:"source_ports,omitempty"`
	IWADs                 []string       `json:"iwads,omitempty"`
	Configured            bool           `json:"configured,omitempty"`
	SaveDirs              bool           `json:"save_dirs,omitempty"`
	WarnBeforeDelete      bool           `json:"warn_before_delete,omitempty"`
	PrintHeader           bool           `json:"print_header,omitempty"`
	NewModDisplay         bool           `json:"new_mod_display,omitempty"`
	GameListAbsoluteWidth int            `json:"game_list_absolute_width,omitempty"`
	GameListRelativeWidth int            `json:"game_list_relative_width,omitempty"`
}

func defaultConfig() Cfg {
	var dConf Cfg
	dConf.ModBasePath = home() + "/Games/Doom"
	dConf.ModExtensions = make(map[string]int)
	dConf.ModExtensions[".wad"] = 1
	dConf.ModExtensions[".pk3"] = 1
	dConf.ModExtensions[".ipk3"] = 1
	dConf.SourcePorts = []string{"gzdoom", "zandronum", "lzdoom"}
	dConf.IWADs = []string{"doom2.wad", "doom.wad"}
	dConf.Configured = false
	dConf.SaveDirs = true
	dConf.WarnBeforeDelete = true
	dConf.PrintHeader = true
	dConf.NewModDisplay = true
	dConf.GameListRelativeWidth = 5
	dConf.GameListAbsoluteWidth = 0
	return dConf
}

func init() {
	firstStart()
	GetInstance()
	loadConfig()
	Persist()
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
	content, err := ioutil.ReadFile(configFullPath())
	if err != nil {
		return err
	}

	err = json.Unmarshal(content, instance)
	if err != nil {
		return err
	}

	return nil
}

func configFullPath() string {
	return GetConfigFolder() + "/" + configName
}

// Exported functions

// GetInstance sets up and returns the singleton instance of config
func GetInstance() *Cfg {
	once.Do(func() {
		instance = &Cfg{}
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

// AddPathToCfgs adds the mod base path to the config ini files
// that enables the engine, to find mod files added with the -file parameter based on relative paths
func AddPathToCfgs() error {
	go processCfg(home() + "/.config/gzdoom/gzdoom.ini")
	go processCfg(home() + "/.config/zandronum/zandronum.ini")
	go processCfg(home() + "/.config/lzdoom/lzdoom.ini")
	return nil
}

func processCfg(path string) {
	lines := configLines(path)
	// if there are not lines for the respective config that is considered ok; maybe that config is not installed
	if lines == nil {
		return
	}

	entry := "PATH=" + instance.ModBasePath
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

// Helper functions

func home() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr.HomeDir
}
