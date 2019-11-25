package main

import (
	"github.com/zmnpl/twad/cfg"
	"github.com/zmnpl/twad/tui"
)

func main() {
	cfg.GetInstance().Configured = false
	tui.Draw()
}
