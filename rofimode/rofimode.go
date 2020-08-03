package rofimode

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"

	"github.com/zmnpl/twad/cfg"
	"github.com/zmnpl/twad/games"
)

// RunRofiMode starts rofi (or demenu) to select and run a already created game.
// It pipes all games as a list of names to the external program
func RunRofiMode(command string) {
	cfg.EnableBasePath()
	var params []string
	if command == "rofi" && commandExists("rofi") {
		params = []string{"rofi", "-dmenu", "-p", "'Rip & Tear'"}
	} else if command == "dmenu" && commandExists("dmenu") {
		params = []string{"-p", "'Rip & Tear'"}
	} else {
		return
	}

	rofi := exec.Command(command, params...)
	r, w := io.Pipe()
	rofi.Stdin = r
	var stdout bytes.Buffer
	rofi.Stdout = &stdout
	err := rofi.Start()
	if err != nil {
		//return err
	}

	rofiToGame := make(map[string]int)
	for i, v := range games.GetInstance() {
		displayName := fmt.Sprintf("%v: %s\n", i, v.Name)
		rofiToGame[displayName] = i
		w.Write([]byte(displayName)) // pipe game name to rofi
	}
	w.Close()

	rofi.Wait()

	result := string(stdout.Bytes())
	fmt.Println(result)

	// run selected game
	if i, exists := rofiToGame[result]; exists {
		games.GetInstance()[i].Run(false)
	}
}

func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}
