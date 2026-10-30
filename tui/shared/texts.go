package shared

import (
	"bytes"
	"text/template"
)

// texts shown in both user interfaces; colors are given as tview color tags
const (
	HelloText = "Hi, it looks like you are starting [orange]twad[white] for the first time. Hit [red]o[white] to bring up the options menu and set your DOOM directory ($DOOMWADDIR)"
	Subtitle  = "[orange]twad[white] - [orange]t[white]erminal [orange]wad[white] launcher[orange]"

	DoomLogoCreditHeader = "[orange]DOOM Logo"
	CreditDoomLogo       = `DOOM and Quake are registered trademarks of id Software, Inc. The DOOM, Quake and id logos are trademarks of id Software, Inc.

The ASCII version of the DOOM logo is Copyright © 1994 by F.P. de Vries.

This logo is work from Frans P. de Vries who originally made it and nicely granted me permission to use it here

Details can be found in this little piece of video game history:
http://www.gamers.org/~fpv/doomlogo.html`

	LicenseHeader = "[orange]License"
	MITLicense    = `MIT License

Copyright (c) 2020 Simon Paul

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

	// PreviewText is the prompt in front of the command preview
	PreviewText = "~$"
)

// DOOM and Quake are registered trademarks of id Software, Inc. The DOOM, Quake and id logos are trademarks of id Software, Inc.
// The ASCII version of the DOOM logo is Copyright © 1994 by F.P. de Vries.

// this logo is work from Frans P. de Vries who originally made it and nicely granted me permission to use it here
// details can be found in this little piece of video game history:
// http://www.gamers.org/~fpv/doomlogo.html
const doomLogoTemplate = "[{{.Title}}]=================     ===============     ===============   ========  ========\n[{{.Title}}]\\\\ . . . . . . .\\\\   //. . . . . . .\\\\   //. . . . . . .\\\\  \\\\. . .\\\\// . . //\n[{{.Secondary}}]||[{{.Title}}]. . .[{{.Secondary}}]_____[{{.Title}}]. . .[{{.Secondary}}]|| ||[{{.Title}}]. . .[{{.Secondary}}]_____[{{.Title}}]. . .[{{.Secondary}}]|| ||[{{.Title}}]. . .[{{.Secondary}}]_____[{{.Title}}]. . .[{{.Secondary}}]|| ||[{{.Title}}] . . .[{{.Secondary}}]\\/[{{.Title}}] . . .[{{.Secondary}}]||\n[{{.Secondary}}]||[{{.Title}}] . .[{{.Secondary}}]||   ||[{{.Title}}]. . [{{.Secondary}}]|| ||[{{.Title}}] . .[{{.Secondary}}]||   ||[{{.Title}}]. . [{{.Secondary}}]|| ||[{{.Title}}] . .[{{.Secondary}}]||   ||[{{.Title}}]. . [{{.Secondary}}]|| ||[{{.Title}}]. . . . . . . [{{.Secondary}}]||\n[{{.Secondary}}]||[{{.Title}}]. . [{{.Secondary}}]||   ||[{{.Title}}] . .[{{.Secondary}}]|| ||[{{.Title}}]. . [{{.Secondary}}]||   ||[{{.Title}}] . .[{{.Secondary}}]|| ||[{{.Title}}]. . [{{.Secondary}}]||   ||[{{.Title}}] . .[{{.Secondary}}]|| ||[{{.Title}}] . [{{.Secondary}}]|[{{.Title}}] . . . . .[{{.Secondary}}]||\n[{{.Secondary}}]||[{{.Title}}] . .[{{.Secondary}}]||   ||[{{.Title}}]. [{{.Secondary}}]_-|| ||-_[{{.Title}}] .[{{.Secondary}}]||   ||[{{.Title}}]. . [{{.Secondary}}]|| ||[{{.Title}}] . .[{{.Secondary}}]||   ||[{{.Title}}]. [{{.Secondary}}]_-|| ||-_[{{.Title}}].[{{.Secondary}}]|\\[{{.Title}}] . . . . [{{.Secondary}}]||\n[{{.Secondary}}]||[{{.Title}}]. . [{{.Secondary}}]||   ||-'  || ||  `-||   || [{{.Title}}]. .[{{.Secondary}}]|| ||[{{.Title}}]. . [{{.Secondary}}]||   ||-'  || ||  `|\\_[{{.Title}}] . .[{{.Secondary}}]|[{{.Title}}]. .[{{.Secondary}}]||\n[{{.Secondary}}]||[{{.Title}}] . [{{.Secondary}}]_||   ||    || ||    ||   ||_[{{.Title}}] . [{{.Secondary}}]|| ||[{{.Title}}] . [{{.Secondary}}]_||   ||    || ||   |\\ `-_/|[{{.Title}}] . [{{.Secondary}}]||\n[{{.Secondary}}]||_-' ||  .|/    || ||    \\|.  || `-_|| ||_-' ||  .|/    || ||   | \\  / |-_[{{.Title}}].[{{.Secondary}}]||\n[{{.Secondary}}]||    ||_-'      || ||      `-_||    || ||    ||_-'      || ||   | \\  / |  `||\n[{{.Secondary}}]||    `'         || ||         `'    || ||    `'         || ||   | \\  / |   ||\n[{{.Secondary}}]||            .===' `===.         .==='.`===.         .===' /==. |  \\/  |   ||\n[{{.Secondary}}]||         .=='   \\_|-_ `===. .==='   _|_   `===. .===' _-|/   `==  \\/  |   ||\n[{{.Secondary}}]||      .=='    _-'    `-_  `='    _-'   `-_    `='  _-'   `-_  /|  \\/  |   ||\n[{{.Secondary}}]||   .=='    _-'          `-__\\._-'         `-_./__-'         `' |. /|  |   ||\n[{{.Secondary}}]||.=='    _-'                                                     `' |  /==.||\n[{{.Secondary}}]=='    _-'                                                            \\/   `==\n[{{.Secondary}}]\\   _-'                                                                `-_   /\n[{{.Secondary}}] `''                     [{{.Secondary}}]twad[white] - [{{.Secondary}}]t[white]erminal [{{.Secondary}}]wad[white] launcher[{{.Secondary}}]                     ``'"

var logoTemplate = template.Must(template.New("logo").Parse(doomLogoTemplate))

// Logo returns the DOOM logo with tview color tags; colors can be names or hex values
func Logo(titleColor, secondaryColor string) string {
	buf := new(bytes.Buffer)
	logoTemplate.Execute(buf, struct{ Title, Secondary string }{titleColor, secondaryColor})
	return buf.String()
}
