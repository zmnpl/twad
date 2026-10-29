package base

import (
	"archive/zip"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
	// loadErr is set when an existing config file couldn't be loaded;
	// persisting is refused then, so the file doesn't get overwritten with defaults
	loadErr   error
	persistMu sync.Mutex
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
	CFG_VERSION      = 2
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
	UseTerminalColors      bool     `json:"use_terminal_colors"`
	CfgVersion             int      `json:"cfg_version"`
}

// Setup prepares everything twad needs before running games:
// it creates missing folders and the default config, loads and migrates the config,
// persists it and points the engines at the WAD dir
func Setup() error {
	if err := firstStart(); err != nil {
		return err
	}

	Config()
	if loadErr != nil {
		return fmt.Errorf("%w\nfix or remove the config file and start again", loadErr)
	}

	// just in case new settings made it into the programm
	if err := Persist(); err != nil {
		return err
	}

	return EnableBasePath()
}

func defaultConfig() Cfg {
	config := Cfg{
		WadDir:                filepath.Join(helper.Home(), "/DOOM"),
		ModExtensions:         ".wad.pk3.ipk3.pke.zip",
		Ports:                 []string{"uzdoom", "gzdoom", "zandronum", "lzdoom"},
		IWADs:                 []string{"doom2.wad", "doom.wad", "plutonia.wad", "tnt.wad", "heretic.wad", "boa.ipk3"},
		GameListRelativeWidth: 40,
		GameListAbsoluteWidth: 0,
		CfgVersion:            CFG_VERSION,
	}

	// check if user has set DOOMWADDIR
	if dwd := os.Getenv("DOOMWADDIR"); dwd != "" {
		config.WadDir = dwd
	}

	return config
}

// firstStart creates the folders for savegames and configs, as well as a default config if there is none
func firstStart() error {
	folders := []string{GetSavegameFolder(), GetGameConfigFolder(), GetSharedGameConfigFolder()}
	for _, canonical := range ports.PortCanonicalNames {
		folders = append(folders, PortSharedConfigPath(canonical))
	}
	for _, folder := range folders {
		if err := os.MkdirAll(folder, 0755); err != nil {
			return err
		}
	}

	if _, err := os.Stat(configFullPath()); !os.IsNotExist(err) {
		return nil
	}
	defaultConfigJSON, err := json.MarshalIndent(defaultConfig(), "", "    ")
	if err != nil {
		return err
	}
	return WriteFileAtomic(configFullPath(), defaultConfigJSON)
}

func loadConfig() error {
	dConf := defaultConfig()

	content, err := os.ReadFile(configFullPath())
	if err != nil {
		instance = &dConf
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	err = json.Unmarshal(content, instance)
	if err != nil {
		instance = &dConf
		return fmt.Errorf("%s: %w", configFullPath(), err)
	}

	// empty entries can't be selected in any meaningful way
	instance.Ports = withoutEmpty(instance.Ports)
	instance.IWADs = withoutEmpty(instance.IWADs)

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

	// a width of 0 would hide the games table
	if instance.GameListRelativeWidth < 1 || instance.GameListRelativeWidth > 100 {
		instance.GameListRelativeWidth = dConf.GameListRelativeWidth
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
		// additional known iwads
		instance.IWADs = append(instance.IWADs, "boa.ipk3", "plutonia.wad", "tnt.wad", "heretic.wad")
	}

	// v2
	if instance.CfgVersion < 2 {
		if !strings.Contains(instance.ModExtensions, ".zip") {
			instance.ModExtensions = instance.ModExtensions + ".zip"
		}
	}

	instance.CfgVersion = CFG_VERSION
}

// Exported functions

// Config returns the singleton instance of config
func Config() *Cfg {
	once.Do(func() {
		instance = &Cfg{}
		loadErr = loadConfig()
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

// Persist writes the configuration into its JSON file
func Persist() error {
	if loadErr != nil {
		return fmt.Errorf("not saving config, because it couldn't be loaded: %w", loadErr)
	}

	persistMu.Lock()
	defer persistMu.Unlock()

	JSON, err := json.MarshalIndent(instance, "", "    ")
	if err != nil {
		return err
	}

	return WriteFileAtomic(configFullPath(), JSON)
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
	if _, err := os.Stat(filepath.Join(instance.WadDir, mod)); os.IsNotExist(err) {
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

// WriteFileAtomic replaces the file at path by writing a temporary file next to it and renaming it,
// so the file is never left half written; symlinks are followed
// An unchanged file is left alone, which spares file syncing tools some work
func WriteFileAtomic(path string, data []byte) error {
	if target, err := filepath.EvalSymlinks(path); err == nil {
		path = target
	}

	if current, err := os.ReadFile(path); err == nil && bytes.Equal(current, data) {
		return nil
	}

	tmp, err := os.CreateTemp(filepath.Dir(path), filepath.Base(path)+".*.tmp")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name()) // fails harmlessly after the rename

	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Chmod(0644); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmp.Name(), path)
}

func withoutEmpty(list []string) []string {
	out := make([]string, 0, len(list))
	for _, s := range list {
		if s = strings.TrimSpace(s); s != "" {
			out = append(out, s)
		}
	}
	return out
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
			availableIwads = append(availableIwads, file.Name())
		}
	}

	return availableIwads, nil
}

// PortSharedConfigPath returns the path where common/shared configs for the given port should be stored
func PortSharedConfigPath(port string) string {
	return filepath.Join(GetSharedGameConfigFolder(), ports.CanonicalName(port))
}

// GetFileFromPK3 returns the content of one specific file from given pk3 file
// The file is matched by its name, regardless of the directory it is in
func GetFileFromPK3(pk3Path string, filename string) ([]byte, error) {
	r, err := zip.OpenReader(pk3Path)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	for _, f := range r.File {
		if path.Base(f.Name) != filename {
			continue
		}
		content, err := f.Open()
		if err != nil {
			return nil, err
		}
		defer content.Close()
		return io.ReadAll(content)
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

	contentString = string(content)
	return
}

// GetFileLinesFromPK3 is a wrapper for GetFileFromPK3
// It uses a bufio.Scanner to scan the file line by line and return them as slice
func GetFileLinesFromPK3(pk3path, filename string) (lines []string, err error) {
	content, err := GetFileFromPK3(pk3path, filename)
	if err != nil {
		return
	}

	scanner := bufio.NewScanner(bytes.NewReader(content))
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return
}

func DOWNLOAD_PATH() string {
	return filepath.Join(instance.WadDir, "twad_downloads")
}
