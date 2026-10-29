package games

import (
	"os"
	"sync"
	"testing"

	"github.com/zmnpl/twad/base"
)

// newHome points twad at an empty home directory and resets the loaded games
func newHome(t *testing.T) {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)

	instance = nil
	once = sync.Once{}
	loadErr = nil

	if err := os.MkdirAll(base.GetConfigFolder(), 0755); err != nil {
		t.Fatal(err)
	}
}

func TestPersistAndLoadGames(t *testing.T) {
	newHome(t)

	if err := LoadError(); err != nil {
		t.Fatalf("missing games file should be fine: %v", err)
	}
	AddGame(NewGame("Aeon DM", "gzdoom", "", "doom2.wad"))

	// load again from disk
	once = sync.Once{}
	if err := LoadError(); err != nil {
		t.Fatal(err)
	}
	if len(Games()) != 1 || Games()[0].Name != "Aeon DM" {
		t.Errorf("loaded games %v", Games())
	}
}

func TestBrokenGamesFileIsKept(t *testing.T) {
	newHome(t)
	broken := `[{ "name": "Aeon DM", oops`
	if err := os.WriteFile(gamesJSONPath(), []byte(broken), 0644); err != nil {
		t.Fatal(err)
	}

	if err := LoadError(); err == nil {
		t.Fatal("expected an error for a broken games file")
	}

	// adding a game must not replace the file with a list holding only that game
	AddGame(NewGame("New", "gzdoom", "", "doom2.wad"))
	if err := Persist(); err == nil {
		t.Error("expected persisting to be refused")
	}

	content, _ := os.ReadFile(gamesJSONPath())
	if string(content) != broken {
		t.Errorf("broken games file got overwritten with %q", content)
	}
}
