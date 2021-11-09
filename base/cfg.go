package base

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/mholt/archiver/v3"
	"github.com/zmnpl/twad/helper"
	"github.com/zmnpl/twad/ports"
)

var (
	instance *Cfg
	once     sync.Once
	// KnownIwads contains names of the most common iwads
	KnownIwads = map[string]bool{
		"doom.wad":      true,
		"doom2.wad":     true,
		"tnt.wad":       true,
		"plutonia.wad":  true,
		"heretic.wad":   true,
		"hexen.wad":     true,
		"strive1.wad":   true,
		"sve.wad":       true,
		"chex.wad":      true,
		"strife0.wad":   true,
		"freedoom1.wad": true,
		"freedoom2.wad": true,
		"freedm.wad":    true,
		"chex3.wad":     true,
		"action2.wad":   true,
		"harm1.wad":     true,
		"hacx.wad":      true,
		"boa.ipk3":      true,
	}
)

const (
	CFG_VERSION      = 1
	configName       = "twad.json"
	configPath       = ".config/twad"
	MAX_SOURCE_PORTS = 6
)

// Cfg holds basic configuration settings
// Should only be instantiated via GetInstance
type Cfg struct {
	WadDir                 string   `json:"wad_dir"`
	WriteWadDirToEngineCfg bool     `json:"write_wad_dir_to_engine_cfg"`
	DontSetDoomwaddir      bool     `json:"dont_set_doomwaddir"`
	ModExtensions          string   `json:"mod_extensions"`
	Ports                  []string `json:"source_ports"`
	IWADs                  []string `json:"iwa_ds"`
	Configured             bool     `json:"configured"`
	DeleteWithoutWarning   bool     `json:"delete_without_warning"`
	HideHeader             bool     `json:"hide_header"`
	GameListAbsoluteWidth  int      `json:"game_list_absolute_width"`
	GameListRelativeWidth  int      `json:"game_list_relative_width"`
	CfgVersion             int      `json:"cfg_version"`
}

func init() {
	firstStart()
	Config()
	Persist() // just in case new settings made it into the programm
	EnableBasePath()
}

func defaultConfig() Cfg {
	config := Cfg{
		WadDir:                filepath.Join(helper.Home(), "/DOOM"),
		ModExtensions:         ".wad.pk3.ipk3.pke.deh",
		Ports:                 []string{"gzdoom", "zandronum", "lzdoom"},
		IWADs:                 []string{"doom2.wad", "doom.wad", "plutonia.wad", "tnt.wad", "heretic.wad", "boa.ipk3"},
		GameListRelativeWidth: 40,
		GameListAbsoluteWidth: 0,
		CfgVersion:            CFG_VERSION,
	}

	// check if user has set DOOMWADDIR
	if dwd, exists := os.LookupEnv("DOOMWADDIR"); exists {
		config.WadDir = dwd
	}

	return config
}

func firstStart() {
	// create directory for games and configs
	err := os.MkdirAll(GetSavegameFolder(), 0755)
	if err != nil {
		log.Fatal(err)
	}

	err = os.MkdirAll(GetGameConfigFolder(), 0755)
	if err != nil {
		log.Fatal(err)
	}

	err = os.MkdirAll(GetSharedGameConfigFolder(), 0755)
	if err != nil {
		log.Fatal(err)
	}

	// shared config paths
	for _, canonical := range ports.PortCanonicalNames {
		err = os.MkdirAll(PortSharedConfigPath(canonical), 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	configPath := configFullPath()
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
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

	content, err := os.ReadFile(configFullPath()) // TODO: Resolve simlinks
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

	if instance.ModExtensions == "" {
		instance.ModExtensions = dConf.ModExtensions
	}

	if len(instance.Ports) == 0 {
		instance.Ports = dConf.Ports
	}

	if len(instance.IWADs) == 0 {
		instance.IWADs = dConf.IWADs
	}

	if instance.GameListAbsoluteWidth == 0 {
		instance.GameListAbsoluteWidth = 40
	}

	updateConfig()

	return nil
}

// if something needs to be added to the config this can be done here
// existing configs otherwise would not get additions
func updateConfig() {
	// going to v1 apply these changes
	if instance.CfgVersion < 1 {
		// new known mod extension
		if !strings.Contains(instance.ModExtensions, ".pke") {
			instance.ModExtensions = instance.ModExtensions + ".pke"
		}

		if !strings.Contains(instance.ModExtensions, ".deh") {
			instance.ModExtensions = instance.ModExtensions + ".deh"
		}
		// additional known iwads
		instance.IWADs = append(instance.IWADs, "boa.ipk3", "plutonia.wad", "tnt.wad", "heretic.wad")
	}

	instance.CfgVersion = CFG_VERSION
	go Persist()
}

// Exported functions

// Config returns the singleton instance of config
func Config() *Cfg {
	once.Do(func() {
		instance = &Cfg{}
		loadConfig()
	})
	return instance
}

// GetConfigFolder returns the folder where configuration is stored
func GetConfigFolder() string {
	return filepath.Join(helper.Home(), configPath)
}

// GetSavegameFolder returns the folder where savegames are stored
func GetSavegameFolder() string {
	return filepath.Join(GetConfigFolder(), "savegames")
}

// GetGameConfigFolder returns the folder where savegames are stored
func GetGameConfigFolder() string {
	return filepath.Join(GetConfigFolder(), "configs")
}

// GetSharedGameConfigFolder returns the folder where savegames are stored
func GetSharedGameConfigFolder() string {
	return filepath.Join(GetConfigFolder(), "configs_shared")
}

// GetSharedGameConfigs returns a list with configs for given port name in the according shared subfolder
func GetSharedGameConfigs(port string) []string {
	files, err := os.ReadDir(PortSharedConfigPath(port))

	if err != nil {
		return nil
	}

	cfgs := make([]string, len(files))
	for i, file := range files {
		cfgs[i] = file.Name()
	}

	return cfgs
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

	err = os.WriteFile(filepath.Join(GetConfigFolder(), configName), JSON, 0755)
	if err != nil {
		return err
	}

	return nil
}

// WadDirIsSane checks if the configured DOOMWADDIR is something useful
// The root directory "/" for example is a bad idea, especially considering zip import functionality
func WadDirIsSane() bool {
	if instance.WadDir == "/" {
		return false
	}
	// TODO: other bad ideas?
	// contains no wad

	return true
}

// ModOk checks, if the given mod (relative file path) exists within the configured DOOMWADDIR
// TODO: Maybe even check with MD5
func ModOk(mod string) bool {
	if _, err := os.Stat(path.Join(instance.WadDir, mod)); os.IsNotExist(err) {
		return false
	}
	return true
}

// EnableBasePath adds the mod base path to the config ini files and/or sets it as DOOMWADDIR
// that enables the engine, to find mod files added with the -file parameter based on relative paths
func EnableBasePath() error {
	// DOOMWADDIR
	if !instance.DontSetDoomwaddir {
		os.Setenv("DOOMWADDIR", instance.WadDir)
	}

	// Engine-Configs
	if instance.WriteWadDirToEngineCfg {
		go processSourcePortCfg(filepath.Join(helper.Home(), ".config/gzdoom/gzdoom.ini"))
		go processSourcePortCfg(filepath.Join(helper.Home(), ".config/zandronum/zandronum.ini"))
		go processSourcePortCfg(filepath.Join(helper.Home(), ".config/lzdoom/lzdoom.ini"))
	}

	return nil
}

// ImportArchive imports given archive into a subfolder of the base path
func ImportArchive(zipPath, modName string) (err error) {
	err = archiver.Unarchive(zipPath, filepath.Join(instance.WadDir, modName))
	return
}

func SourcePorts() []string {
	sourceports := make([]string, len(instance.Ports))
	for i := range instance.Ports {
		sourceports[i] = filepath.Base(instance.Ports[i])
	}
	return sourceports
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

	err := os.WriteFile(path, configData.Bytes(), 0755)
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

// General Helper
func PathHasIwads(path string) (bool, error) {
	files, err := os.ReadDir(path)

	if err != nil {
		return false, err
	}

	for _, file := range files {
		_, isIwad := KnownIwads[strings.ToLower(file.Name())]
		if isIwad {
			return true, nil
		}
	}

	return false, nil
}

func GePathIwads(path string) ([]string, error) {
	files, err := os.ReadDir(path)

	if err != nil {
		return nil, err
	}

	availableIwads := make([]string, 0, len(KnownIwads))
	for _, file := range files {
		iwadName := strings.ToLower(file.Name())
		_, isIwad := KnownIwads[iwadName]
		if isIwad {
			availableIwads = append(availableIwads, iwadName)
		}
	}

	return availableIwads, nil
}

// PortSharedConfigPath returns the path where common/shared configs for the given port should be stored
func PortSharedConfigPath(port string) string {
	return filepath.Join(GetSharedGameConfigFolder(), ports.CanonicalName(port))
}

// GetFileFromPK3 returns one specific file from given wad or pk3 file
// Don't forget to close the ReadCloser when done reading
// Can return nil if it didn't work out
func GetFileFromPK3(pk3Path string, filename string) (io.ReadCloser, error) {
	f, err := os.Open(pk3Path)
	if err != nil {
		return nil, err
	}

	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}

	z := archiver.NewZip()
	z.Open(f, fi.Size())
	for {
		f, err := z.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			break
		}

		if f.Name() == filename {
			fmt.Println(f.Name())
			return f, nil
		}

		err = f.Close()
		if err != nil {
			// TODO: Does that interest me?
		}
	}
	return nil, fmt.Errorf("couldn't find %v in %v", filename, pk3Path)
}

// GetFileContentStringFromPK3 is a wrapper for GetFileFromPK3
// It returns the files contents as string
func GetFileContentStringFromPK3(pk3path, filename string) (contentString string, err error) {
	content, err := GetFileFromPK3(pk3path, filename)
	if err != nil {
		return
	}
	defer content.Close()

	contentBytes, err := io.ReadAll(content)
	if err != nil {
		return
	}

	contentString = string(contentBytes)
	return
}

// GetFileLinesFromPK3 is a wrapper for GetFileFromPK3
// It uses a bufio.Scanner to scan the file line by line and return them as slice
func GetFileLinesFromPK3(pk3path, filename string) (lines []string, err error) {
	content, err := GetFileFromPK3(pk3path, filename)
	if err != nil {
		return
	}
	defer content.Close()

	scanner := bufio.NewScanner(content)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return
}
