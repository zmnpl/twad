# twad - terminal wad launcher

![demo](demo.gif)

If you love __DOOM__ and rather not leave your terminal like me, then you might be one of the few that might like **twad**. It is a terminal based WAD manager and launcher for doom source ports. At it's core twad lets you set up a multitude of WAD file combinations, store them and launch them with a couple of key strokes.

There are already great alternatives to manage and launch your WADs out there for many years and twad will probably never be as sophisticated. Though I figured: there are not so many for the terminal. Twad let's you stay in the terminal and on your keyboard as long as possible until you decide to rip and tear. Simple as that.

Needless to say, that this mostly was designed for *nix systems. However, with WSL this might as well be usable in Windows as well. Though quite some testing needs to go into this.

**Watch Out**: This tool is still in very early state and might contain bugs.

## Features

* Configure *games* you want to play (basically a combination of IWAD, Source Port and Mod Files)
* Run games
* Quickload (start from latest savegame)
* Warp to specific episodes / levels
* Record demos
* Watch demos
* Delete demos
* Collect some stats from the cames console output as well as playtime and so on (not sending it anywhere)
* Run games from rofi or dmenu
* Separate savegame / demo folders for games (in **~/.config/twad/...***)

## Installation Options

### a) Arch Linux: AUR

https://aur.archlinux.org/packages/twad-git

### b) Binary Download

I'll to add precompiled binaries to the [releases page](https://github.com/zmnpl/twad/releases). It comes without dependencies, just **download and run it** (on *nix systems).

### c) Compile yourself

```golang
go get -u github.com/zmnpl/twad
```

## Setup your environment

***twad*** assumes, you have **one folder**, where your IWADs are located. All your pwads (mapsets, gameplay mods, ...) need to be in the same folder or subfolder of this. The folder, where you put your IWADs is known to source ports as **DOOMWADDIR**.

An example could look like this:
```bash
~/DOOM ❯❯❯ tree
.
├── Ashes 2063
│   ├── Ashes2063Maps115.wad
│   └── Ashes2063Mod115.pk3
├── Back To Saturn X e1
│   ├── btsx_e1a.wad
│   ├── btsx_e1b.wad
│   └── btsx_e1.deh
├── D4T
│   └── D4Tv2.5.pk3
├── Sigil
│   ├── SIGIL_COMPAT.wad
│   └── SIGIL.wad
├── doom2.wad
└── doom.wad
```
1) Setup your **DOOMWADDIR** as described above
2) twad's first start will ask you for the path of **DOOMWADDIR**
3) Within twad create games
4) Add mods to your games
666)   __Rip and Tear__

### More on DOOMWADDIR

Your DOOM source port needs to know about the base folder of your mods and IWADs to work properly, since ***twad*** uses relative paths. ***twad**'s default method for this is to set the ***DOOMWADDIR*** environment variable when starting a game. This is only set for the current game session. (Should you already have set DOOMWADDIR, twad will shadow it with whatever is set in its configuration)

An alternative/additional method is to add paths to the respective source ports config. For *zdoom* ports it could look like this:
```bash
# in your doom engine .ini
[FileSearch.Directories]
PATH=/home/doomguy/Doom # path to DOOMWADDIR
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
- ~~Warp to map~~
- ~~Demo recording / viewing~~
- More statistics
- All the TODO flags
- Error Handling

## Credit where credit is due

### Doom logo

The use of the DOOM ASCII logo has been nicely permitted by Frans P. de Vries. Find it's history [here](http://www.gamers.org/~fpv/doomlogo.html)

DOOM and Quake are registered trademarks of id Software, Inc. The DOOM, Quake and id logos are trademarks of id Software, Inc. The ASCII version of the DOOM logo is Copyright © 1994 by F.P. de Vries.

### tview

[tview](https://github.com/rivo/tview) is used for the terminal ui elements.
