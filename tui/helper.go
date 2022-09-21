package tui

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gdamore/tcell/v2"
	"github.com/zmnpl/goidgames"
)

type WriteCounter struct {
	Total uint64
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.PrintProgress()
	return n, nil
}

func (wc WriteCounter) PrintProgress() {
	statusline.Clear()
	statusline.Write([]byte(fmt.Sprintf(" Downloading... %v complete", wc.Total)))
}

// DownloadTo tries to download the game to given path and returns the full path of the downloaded file
func DownloadIdGame(g goidgames.Idgame, path string) (filePath string, err error) {
	success := false
	if err = os.MkdirAll(path, 0755); err != nil {
		return "", err
	}

	filePath = filepath.Join(path, g.Filename)
	// try for all mirrors
	for _, mirror := range goidgames.Mirrors {
		resp, err := http.Get(fmt.Sprintf("%s/%s/%s", mirror, g.Dir, g.Filename))
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		out, err := os.Create(filePath)
		if err != nil {
			continue
		}
		defer out.Close()

		counter := &WriteCounter{}
		_, err = io.Copy(out, io.TeeReader(resp.Body, counter))
		if err == nil {
			success = true
			break
		}
	}
	//statusline.Clear()
	if !success {
		return "", fmt.Errorf("%s", "Unable to download.")
	}
	return filePath, nil
}

func hexStringFromColor(c tcell.Color) string {
	r, g, b := c.RGB()
	return fmt.Sprintf("[#%02x%02x%02x]", r, g, b)
}
