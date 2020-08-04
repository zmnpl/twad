# twad - terminal wad launcher

![demo](demo.gif)

If you love __DOOM__ and rather not leave your terminal like me, then you might be one of the few that might like **twad**. It is a terminal based WAD manager and launcher for ZDoom engine games. At it's core twad lets you set up a multitude of WAD file combinations, store them and launch them with a couple of key strokes.

There are already great alternatives to manage and launch your WADs out there for many years and twad will probably never be as sophisticated. Though I figured: there are not so many for the terminal. Twad let's you stay in the terminal and on your keyboard as long as possible. Simple as that.

As a little bonus, twad collects some statistics for you (of course doesn't send them anywhere!) and organizes savegames for each mod combination in a separate folder.

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

I'll to add precompiled binaries to the [releases page](https://github.com/zmnpl/twad/releases). It comes without dependencies, just download and run it (on *nix systems atm).

## Initial Config

1) Create a base directory which gets to hold all your DOOM files (long time ago, this was known as DOOMWADDIR)
2) Put your **doom.wad** and **doom2.wad** in that base dir
3) Drop all your mods in here (Subdirectories are of course possible)
4) Within twad create games
5) Add mods to your games
666) __Rip and Tear__

## DOOMWADDIR

Your DOOM engine needs to know about the base folder of your mods and IWADs to work properly. Twad's default method for this is to set the ***DOOMWADDIR*** environment variable when starting a game. This is only set for the current game session.

An alternative method is to add paths to the respective source ports config. For *zdoom* ports it could look like this:
```
# in your doom engine .ini
[FileSearch.Directories]
PATH=/home/doomguy/Doom
```

There is flag in the options which lets Twad try to do this automatically for these engines if it finds the respective config:
- **Zandronum** *(~/.config/zandronum/zandronum.ini)*
- **LZDoom** *(~/.config/lzdoom/lzdoom.ini)*
- **GZDoom** *(~/.config/gzdoom/gzdoom.ini)*

If you are using something different, please configure it accoridingly or send in an issue or pull request ;)

## Rofi Mode

You can use [***rofi***](https://github.com/davatorium/rofi) or [***dmenu***](https://tools.suckless.org/dmenu/) to launch your games. Run twad like this to use the respective programm. This will open rofi/dmenu and show a list of all games you already have. Select one you want to play and hit enter. Of course this will also track your statistics.
```bash
twad --rofi
# or
twad --dmenu
```
**For instant Rip & Tear:** Bind this to a keyboard shortcut

![rofimode](rofimode.png)

## Plans / Ideas

- ~~Separate savegames folders per game~~
- ~~AUR package~~
- ~~Rofi mode~~
- ~~Help area~~
- ~~Savegame Count~~
- ~~Unified Add/Edit dialog~~
- ~~Opions scren~~
- ~~Ability to hide the header for screens with few rows~~
- ~~Add button for path setup~~
- ~~Quickload~~
- WSL support
- Warp to map
- Demo recording / viewing
- More statistics
- All the TODO flags

## Credit where credit is due

### Doom logo

The use of the DOOM ASCII logo has been nicely permitted by Frans P. de Vries. Find it's history [here](http://www.gamers.org/~fpv/doomlogo.html)

DOOM and Quake are registered trademarks of id Software, Inc. The DOOM, Quake and id logos are trademarks of id Software, Inc. The ASCII version of the DOOM logo is Copyright © 1994 by F.P. de Vries.

### tview

[tview](https://github.com/rivo/tview) is used for the terminal ui elements.
