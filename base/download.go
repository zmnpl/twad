package base

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/zmnpl/goidgames"
)

// progressWriter reports the number of bytes written so far
type progressWriter struct {
	total    uint64
	progress func(total uint64)
}

func (pw *progressWriter) Write(p []byte) (int, error) {
	pw.total += uint64(len(p))
	if pw.progress != nil {
		pw.progress(pw.total)
	}
	return len(p), nil
}

// DownloadIdGame tries to download the game to the given directory from one of the idgames mirrors
// and returns the full path of the downloaded file
// progress gets called with the number of bytes downloaded so far; it may be nil
func DownloadIdGame(g goidgames.Idgame, dir string, progress func(total uint64)) (string, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	filePath := filepath.Join(dir, g.Filename)
	var errs []error
	for _, mirror := range goidgames.Mirrors {
		url := fmt.Sprintf("%s/%s/%s", mirror, g.Dir, g.Filename)
		err := download(url, filePath, progress)
		if err == nil {
			return filePath, nil
		}
		errs = append(errs, err)
	}

	return "", fmt.Errorf("unable to download %s: %w", g.Filename, errors.Join(errs...))
}

// download fetches url into filePath; a partial file is removed on failure
func download(url, filePath string, progress func(total uint64)) (err error) {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s: %s", url, resp.Status)
	}

	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := out.Close(); err == nil {
			err = cerr
		}
		if err != nil {
			os.Remove(filePath)
		}
	}()

	_, err = io.Copy(out, io.TeeReader(resp.Body, &progressWriter{progress: progress}))
	return err
}
