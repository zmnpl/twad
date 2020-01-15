# twad - terminal wad launcher

![demo](demo.gif)

If you love __DOOM__ and rather not leave your terminal like me, then you might be one of the few that might like **twad**. It is a terminal based WAD manager and launcher for ZDoom engine games. At it's core twad lets you set up a multitude of WAD file combinations, store them and launch them with the press of a button.

There are already great alternatives to manage and launch your WADs out there for many years and twad will probably never be as sophisticated. Though I figured: there are not so many for the terminal. Twad let's you stay in the terminal and on your keyboard as long as possible. Simple as that.

As a little bonus,  twad collects some statistics for you (of course doesn't send them anywhere!) and organizes savegames for each mod combination in a separate folder.

## Watch Out

This tool is still in very early state and might contain bugs.

## Installation

### AUR

https://aur.archlinux.org/packages/twad-git

### Manually

```golang
go get -u github.com/zmnpl/twad
```

### Binary Download

I'll to add precompiled binaries to the [releases page](https://github.com/zmnpl/twad/releases). It comes without dependencies, just download and run it.

## Initial Config

1) Create a base directory which holds gets to hold all your DOOM files
2) Put your **doom.wad** and **doom2.wad** in that base dir
3) Drop all your mods in here (Subdirectories are of course possible)
4) Within twad create games
5) Add mods to your games
666) __Rip and Tear__


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
