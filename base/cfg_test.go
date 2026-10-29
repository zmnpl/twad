package base

import (
	"archive/zip"
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"testing"
)

// newHome points twad at an empty home directory and resets the loaded config
func newHome(t *testing.T) string {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	t.Setenv("DOOMWADDIR", "")

	instance = nil
	once = sync.Once{}
	loadErr = nil

	return home
}

func writeConfig(t *testing.T, content string) {
	t.Helper()
	if err := os.MkdirAll(GetConfigFolder(), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(configFullPath(), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

func readConfig(t *testing.T) Cfg {
	t.Helper()
	content, err := os.ReadFile(configFullPath())
	if err != nil {
		t.Fatal(err)
	}
	var c Cfg
	if err := json.Unmarshal(content, &c); err != nil {
		t.Fatal(err)
	}
	return c
}

func TestSetupFirstStart(t *testing.T) {
	home := newHome(t)

	if err := Setup(); err != nil {
		t.Fatal(err)
	}

	for _, folder := range []string{GetSavegameFolder(), GetGameConfigFolder(), PortSharedConfigPath("gzdoom")} {
		if fi, err := os.Stat(folder); err != nil || !fi.IsDir() {
			t.Errorf("folder %s not created: %v", folder, err)
		}
	}

	c := readConfig(t)
	if want := filepath.Join(home, "DOOM"); c.WadDir != want {
		t.Errorf("wad dir %q, want %q", c.WadDir, want)
	}
	if c.CfgVersion != CFG_VERSION {
		t.Errorf("config version %d, want %d", c.CfgVersion, CFG_VERSION)
	}
	if got := os.Getenv("DOOMWADDIR"); got != c.WadDir {
		t.Errorf("DOOMWADDIR %q, want %q", got, c.WadDir)
	}
}

func TestSetupKeepsBrokenConfig(t *testing.T) {
	newHome(t)
	broken := `{ "wad_dir": "/doom", oops`
	writeConfig(t, broken)

	if err := Setup(); err == nil {
		t.Fatal("expected an error for a broken config")
	}
	if err := Persist(); err == nil {
		t.Fatal("expected persisting to be refused")
	}

	content, _ := os.ReadFile(configFullPath())
	if string(content) != broken {
		t.Errorf("broken config got overwritten with %q", content)
	}
}

func TestLoadMigratesAndNormalizes(t *testing.T) {
	newHome(t)
	writeConfig(t, `{
		"wad_dir": "/doom",
		"mod_extensions": ".wad.pk3",
		"source_ports": ["gzdoom", "", "  "],
		"iwa_ds": ["doom2.wad", ""],
		"game_list_relative_width": 0
	}`)

	if err := Setup(); err != nil {
		t.Fatal(err)
	}

	for _, c := range []Cfg{*Config(), readConfig(t)} {
		if !slices.Equal(c.Ports, []string{"gzdoom"}) {
			t.Errorf("ports %q", c.Ports)
		}
		if c.IWADs[0] != "doom2.wad" || slices.Contains(c.IWADs, "") {
			t.Errorf("iwads %q", c.IWADs)
		}
		if c.GameListRelativeWidth != defaultConfig().GameListRelativeWidth {
			t.Errorf("relative width %d", c.GameListRelativeWidth)
		}
		if !strings.Contains(c.ModExtensions, ".pke") || !strings.Contains(c.ModExtensions, ".zip") {
			t.Errorf("mod extensions not migrated: %q", c.ModExtensions)
		}
		if c.CfgVersion != CFG_VERSION {
			t.Errorf("config version %d", c.CfgVersion)
		}
	}
}

func TestPersistLeavesUnchangedFileAlone(t *testing.T) {
	newHome(t)
	if err := Setup(); err != nil {
		t.Fatal(err)
	}
	before, _ := os.Stat(configFullPath())

	if err := Persist(); err != nil {
		t.Fatal(err)
	}
	after, _ := os.Stat(configFullPath())
	if !os.SameFile(before, after) || !before.ModTime().Equal(after.ModTime()) {
		t.Error("unchanged config was rewritten")
	}

	Config().HideHeader = !Config().HideHeader
	if err := Persist(); err != nil {
		t.Fatal(err)
	}
	if readConfig(t).HideHeader != Config().HideHeader {
		t.Error("changed config wasn't written")
	}

	if tmps, _ := filepath.Glob(filepath.Join(GetConfigFolder(), "*.tmp")); len(tmps) > 0 {
		t.Errorf("temporary files left: %v", tmps)
	}
}

func TestPersistFollowsSymlink(t *testing.T) {
	newHome(t)
	if err := Setup(); err != nil {
		t.Fatal(err)
	}

	// config managed elsewhere, e.g. in a dotfiles repository
	target := filepath.Join(t.TempDir(), "twad.json")
	if err := os.Rename(configFullPath(), target); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(target, configFullPath()); err != nil {
		t.Skip("symlinks not available:", err)
	}

	Config().HideHeader = !Config().HideHeader
	if err := Persist(); err != nil {
		t.Fatal(err)
	}

	if fi, err := os.Lstat(configFullPath()); err != nil || fi.Mode()&os.ModeSymlink == 0 {
		t.Error("symlink got replaced")
	}
	if readConfig(t).HideHeader != Config().HideHeader {
		t.Error("change didn't reach the symlink target")
	}
}

func TestGetFileFromPK3(t *testing.T) {
	pk3 := filepath.Join(t.TempDir(), "mod.pk3")
	f, err := os.Create(pk3)
	if err != nil {
		t.Fatal(err)
	}
	w := zip.NewWriter(f)
	for name, content := range map[string]string{
		"maps/map01.wad":  "not a map",
		"content/mapinfo": "map aeon01 \"Vanish\"\nmap aeon02 \"Solitaire\"",
	} {
		fw, _ := w.Create(name)
		fw.Write([]byte(content))
	}
	w.Close()
	f.Close()

	lines, err := GetFileLinesFromPK3(pk3, "mapinfo")
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) != 2 || lines[0] != `map aeon01 "Vanish"` {
		t.Errorf("lines %q", lines)
	}

	if _, err := GetFileFromPK3(pk3, "missing"); err == nil {
		t.Error("expected an error for a missing file")
	}
}
