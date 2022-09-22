package main

import (
	"flag"

	"github.com/zmnpl/twad/base"
	"github.com/zmnpl/twad/rofimode"
	"github.com/zmnpl/twad/tui"
)

func main() {
	rofi := flag.Bool("rofi", false, "Run rofi mode.")
	dmenu := flag.Bool("dmenu", false, "Run dmenu mode.")
	flag.Parse()

	base.Config()

	if *rofi {
		rofimode.RunRofiMode("rofi")
		return
	}

	if *dmenu {
		rofimode.RunRofiMode("rofi")
		return
	}

	//cfg.GetInstance().Configured = false
	tui.Draw()
}
