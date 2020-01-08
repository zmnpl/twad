package tui

const (
	inputLabelName       = "Name"
	inputLabelSourcePort = "Source Port"
	inputLabelIWad       = "IWad"

	subtitle  = "[orange]twad[white] - [orange]t[white]erminal [orange]wad[white] manager and launcher[orange]"
	subtitle2 = "twad - terminal wad manager and launcher"

	setupOkHint      = "Hit [red]Ctrl+O[white] when you are done."
	setupPathExplain = `For [orange]twad[white] to function correctly, you should have all your DOOM mod files organized in one central directory. Subdirectories per mod are possible of course.
Navigate with arrow keys or Vim bindings. [red]Enter[white] or [red]Space[white] expand the directory. Highlight the righ one and hit [red]Ctrl+O[white]`
	setupPathExample = `[red]->[white]/home/slayer/games/DOOMmods            [red]# i need this folder
  [white]/home/slayer/games/DOOMmods[orange]/BrutalDoom [grey]# sub dir for Brutal Doom
  [white]/home/slayer/games/DOOMmods[orange]/QCDE       [grey]# sub dir for QCDE`

	tviewHeader = "[orange]tview"
	creditTview = `The terminal user interface is build with tview:
https://github.com/rivo/tview`

	doomLogoCreditHeader = "[orange]DOOM Logo"
	creditDoomLogo       = `DOOM and Quake are registered trademarks of id Software, Inc. The DOOM, Quake and id logos are trademarks of id Software, Inc.

The ASCII version of the DOOM logo is Copyright © 1994 by F.P. de Vries.

This logo is work from Frans P. de Vries who originally made it and nicely granted me permission to use it here

Details can be found in this little piece of video game history:
http://www.gamers.org/~fpv/doomlogo.html`

	licenseHeader = "[orange]License"
	mitLicense    = `MIT License

Copyright (c) 2019 Simon Paul

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.`
)
