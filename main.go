package main

import (
	"os"

	"github.com/zmnpl/twad/cfg"
	"github.com/zmnpl/twad/rofimode"
	"github.com/zmnpl/twad/tui"
)

func main() {
	args := os.Args[1:]

	cfg.Config()

	for _, v := range args {
		switch v {
		case "--rofi":
			rofimode.RunRofiMode("rofi")
			return
		case "--dmenu":
			rofimode.RunRofiMode("dmenu")
			return
		}

	}
	//cfg.GetInstance().Configured = false
	tui.Draw()
}
