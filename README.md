# WeazlTunes

```text
 __      __          _______________._____________           ________         
/  \    /  \ ____   /  |  \____    /|  \__    ___/_ __  ____ \_____  \  ______
\   \/\/   // __ \ /   |  |_/     / |  | |    | |  |  \/    \  _(__  < /  ___/
 \        /\  ___//    ^   /     /_ |  |_|    | |  |  /   |  \/       \\___ \ 
  \__/\  /  \___  >____   /_______ \|____/____| |____/|___|  /______  /____  >
       \/       \/     |__|       \/                       \/       \/     \/ 
```

WeazlTunes is a terminal Icecast player written in Go with Bubble Tea, Lip Gloss,
and Harmonica. It ships with five presets, can browse SomaFM, and can search the
Xiph/Icecast directory by genre or station text. Direct stream URLs and YouTube
live URLs can be saved under `my stations`; playback is delegated to `mpv`.

## Requirements

- Go 1.25 or newer
- `mpv` in `PATH`

The TUI controls search and playback; `mpv` handles MP3, AAC, Ogg, and any
network source your local mpv build supports.

## Run

```sh
go run ./cmd/weazltunes
```

## Install

```sh
./scripts/install.sh
```

## Keys

- `1`: top five presets
- `2`: SomaFM channels
- `3`: Xiph search mode
- `4`: my stations
- `/`: focus tune box
- `enter`: search, play selected station, or save/play a pasted URL
- `ctrl+p`: add pasted, selected, or playing station to presets
- `ctrl+r`: rename selected preset or saved station
- `space`: pause/resume playback
- `s`: stop playback
- `q` / `ctrl+c`: quit

On first launch, WeazlTunes writes `~/.config/weazltunes/config.json`. Edit that
file to replace presets or saved stations.
