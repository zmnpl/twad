# twad - terminal wad launcher

![demo](demo.gif)

If you love __DOOM__ and you love your terminal, then twad might be for you. It is a terminal based WAD manager and launcher for ZDoom engine games. At it's core twad lets you set up a multitude of WAD file combinations, store them and launch them with the press of a button.

There are already great alternatives to manage and launch your WADs out there for many years and twad will probably never be as sophisticated. Though I figured: there are not so many for the terminal. Twad let's you stay in the terminal and on your keyboard as long as possible. Simple as that.

__ALPHA__
This tool is still in very early state.

## Installation

### AUR

An AUR package is planned.

### Manual

```golang
go get -u github.com/zmnpl/twad
```

## Usage

1) Set up your base directory
2) Create games
3) Add mods to your games
4) __Rip and Tear__

## Rofi Mode

You can use [***rofi***](https://github.com/davatorium/rofi) or [***dmenu***](https://tools.suckless.org/dmenu/) to launch your games. Run twad like this to use the respective programm. This will open rofi/dmenu and show a list of all games you already have. Select one you want to play and hit enter. Of course this will also track your statistics.
```bash
twad --rofi
# or
twad --dmenu
```
**For instant Rip & Tear:** Bind this to a keyboard shortcut


## Plans

- ~~Separate savegames folders per game~~
- ~~AUR package~~
- ~~Rofi mode~~
- ~~Help area~~
- ~~Savegame Count~~
- Import for downloaded Zips
- More statistics
- Fading popup highlight where somthing was added
- All the TODO flags

## Credit where credit is due

### Doom logo

The use of the DOOM ASCII logo has been nicely permitted by Frans P. de Vries. Find it's history [here](http://www.gamers.org/~fpv/doomlogo.html)

DOOM and Quake are registered trademarks of id Software, Inc. The DOOM, Quake and id logos are trademarks of id Software, Inc. The ASCII version of the DOOM logo is Copyright © 1994 by F.P. de Vries.

### tview

[tview](https://github.com/rivo/tview) is used for the terminal ui elements.
