package base

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/zmnpl/goidgames"
)

func TestDownloadIdGame(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/good/levels/foo.zip":
			w.Write([]byte(strings.Repeat("x", 100000)))
		case "/broken/levels/foo.zip":
			w.Header().Set("Content-Length", "100000")
			w.Write([]byte("short")) // connection ends early
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	mirrors := goidgames.Mirrors
	t.Cleanup(func() { goidgames.Mirrors = mirrors })
	g := goidgames.Idgame{Dir: "levels", Filename: "foo.zip"}

	t.Run("skips failing mirror", func(t *testing.T) {
		goidgames.Mirrors = []string{srv.URL + "/missing", srv.URL + "/good"}
		var total uint64
		p, err := DownloadIdGame(g, t.TempDir(), func(n uint64) { total = n })
		if err != nil {
			t.Fatal(err)
		}
		if total != 100000 {
			t.Errorf("progress reported %d bytes", total)
		}
		if fi, err := os.Stat(p); err != nil || fi.Size() != 100000 {
			t.Errorf("downloaded file: %v", err)
		}
	})

	t.Run("no partial file on failure", func(t *testing.T) {
		goidgames.Mirrors = []string{srv.URL + "/missing", srv.URL + "/broken"}
		dir := t.TempDir()
		if _, err := DownloadIdGame(g, dir, nil); err == nil {
			t.Fatal("expected an error")
		}
		if _, err := os.Stat(filepath.Join(dir, g.Filename)); !os.IsNotExist(err) {
			t.Error("partial file left behind")
		}
	})
}
