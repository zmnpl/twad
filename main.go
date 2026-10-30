package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/zmnpl/twad/base"
	"github.com/zmnpl/twad/bubbletui"
	"github.com/zmnpl/twad/games"
	"github.com/zmnpl/twad/rofimode"
	"github.com/zmnpl/twad/tui"
)

func main() {
	rofi := flag.Bool("rofi", false, "Run rofi mode.")
	dmenu := flag.Bool("dmenu", false, "Run dmenu mode.")
	bubbletea := flag.Bool("bubbletea", false, "Use the user interface built with bubbletea.")
	flag.Parse()

	if err := base.Setup(); err != nil {
		fmt.Fprintf(os.Stderr, "twad: %v\n", err)
		os.Exit(1)
	}
	if err := games.LoadError(); err != nil {
		fmt.Fprintf(os.Stderr, "twad: %v\nfix or remove the games file and start again\n", err)
		os.Exit(1)
	}

	if *rofi {
		rofimode.RunRofiMode("rofi")
		return
	}

	if *dmenu {
		rofimode.RunRofiMode("rofi")
		return
	}

	if *bubbletea {
		if err := bubbletui.Draw(); err != nil {
			fmt.Fprintf(os.Stderr, "twad: %v\n", err)
			os.Exit(1)
		}
		return
	}

	//cfg.GetInstance().Configured = false
	tui.Draw()
}
