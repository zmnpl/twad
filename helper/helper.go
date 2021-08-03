package helper

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// Home returns the users home directory
func Home() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr.HomeDir
}

// IsFileNameValid tests if the given file name can be used
func IsFileNameValid(fp string) bool {
	// Check if file already exists
	if _, err := os.Stat(fp); err == nil {
		return true
	}

	// Attempt to create it
	var d []byte
	if err := os.WriteFile(fp, d, 0644); err == nil {
		os.Remove(fp) // And delete it
		return true
	}

	return false
}

// FilterExtensions filters a given slice of FileInfo based on the passed extensions
// extensions can be a string containing multiple extensions; strings.Contains is used for comparison
func FilterExtensions(files []os.DirEntry, extensions string, includeDirs bool) []os.DirEntry {
	n := 0
	for _, f := range files {
		ext := strings.ToLower(filepath.Ext(f.Name()))
		if strings.Contains(extensions, ext) && !f.IsDir() {
			files[n] = f
			n++
		} else if includeDirs && f.IsDir() {
			files[n] = f
			n++
		}
	}
	files = files[:n]
	return files
}

// OpenBrowser tries to open the users default browser with given url
func Openbrowser(url string) {
	// test url
	timeout := time.Duration(1 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	resp, err := client.Get(url)
	// do nothing in case of error
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// dummy
	if resp.StatusCode == 404 {
	}

	// only when url is reachable
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		// TODO
		//log.Fatal(err)
	}
}
