package helper

import (
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
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
	if err := ioutil.WriteFile(fp, d, 0644); err == nil {
		os.Remove(fp) // And delete it
		return true
	}

	return false
}

// FilterExtensions filters a given slice of FileInfo based on the passed extensions
// extensions can be a string containing multiple extensions; strings.Contains is used for comparison
func FilterExtensions(files []os.FileInfo, extensions string, includeDirs bool) []os.FileInfo {
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
