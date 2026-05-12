# WeazlTunes

```text
 __      __          _______________._____________           ________         
/  \    /  \ ____   /  |  \____    /|  \__    ___/_ __  ____ \_____  \  ______
\   \/\/   // __ \ /   |  |_/     / |  | |    | |  |  \/    \  _(__  < /  ___/
 \        /\  ___//    ^   /     /_ |  |_|    | |  |  /   |  \/       \\___ \ 
  \__/\  /  \___  >____   /_______ \|____/____| |____/|___|  /______  /____  >
       \/       \/     |__|       \/                       \/       \/     \/ 
```

WeazlTunes is a terminal radio tuner for Icecast streams, SomaFM, Xiph, and
whatever weird little URL you copied from a corner of the internet at 2:13 AM.
Think Winamp muscle memory, public radio directory energy, and a TUI that still
knows how to wear a gradient. No browser tabs, no autoplay sludge, just stations
and `mpv`.

## Defaults

On first launch, WeazlTunes drops a fresh `config.json` into
`~/.config/weazltunes/` with five SomaFM presets:

- Groove Salad
- Drone Zone
- DEF CON Radio
- Indie Pop Rocks
- Deep Space One

Custom URLs are saved under `my_stations`. Presets stay intentionally tight at
five slots so they feel like radio buttons, not a junk drawer.

## Run

```sh
go run ./cmd/weazltunes
```

## Install

```sh
./scripts/install.sh
```

The installer builds `weazltunes`, tucks it into `~/.weazltunes/bin`, and adds
that directory to your shell `PATH` if it is not already there. It uses local Go
caches inside the repo by default, checks your Go version against `go.mod`, and
warns if `mpv` is not on `PATH`.

Set `WEAZLTUNES_SKIP_LAUNCH=1` if you want setup without the first boot:

```sh
WEAZLTUNES_SKIP_LAUNCH=1 ./scripts/install.sh
```

## Requirements

- Go 1.25 or newer
- `mpv` in `PATH`

WeazlTunes is the tuner. `mpv` is the amp. That means MP3, AAC, Ogg, playlist
URLs, direct streams, and YouTube live streams work as well as your local `mpv`
build supports them.

## Build From Source

```sh
go build -o weazltunes ./cmd/weazltunes
```

No CGO, no database, no native extension puzzle box. If Go can build Bubble Tea
apps on your machine, it can build this.

## Keys

- `1`: top five presets
- `2`: SomaFM channels
- `3`: Xiph search mode
- `4`: my stations
- `/`: focus the tune box
- `v`: switch between sphere and bars visualizers
- `enter`: search, play selected station, or save/play a pasted URL
- `ctrl+p`: add pasted, selected, or playing station to presets
- `ctrl+r`: rename selected preset or saved station
- `space`: pause/resume playback
- `s`: stop playback
- `esc`: leave the tune box
- `q` / `ctrl+c`: quit

## Tuning

Paste a direct stream URL into the tune box and press `enter`:

```text
https://radio.prendie.io/radio.mp3
```

It gets saved to `my stations` and starts playing immediately. Paste a YouTube
live URL and WeazlTunes hands it to `mpv` the same way. If `mpv` can resolve it,
WeazlTunes can tune it.

Press `ctrl+p` on a saved station when it deserves one of the five preset slots.
Newest promotions move to the top, and duplicate URLs are folded together rather
than copied around like mystery mixtapes.

## Directories

SomaFM comes from:

```text
https://api.somafm.com/channels.json
```

Xiph search uses the public Icecast directory. Genre pages are tried first, then
WeazlTunes falls back to the directory XML feed for broader text matching.

## Visuals

The colors come from WeazlChat, inverted for the radio sibling: a yellow
WeazlTunes wordmark with purple diagonal rails, dark panels, mint status text,
and a Harmonica-powered `sphere` visualizer inspired by the Vegas Sphere happy
face. The old bar visualizer is still available with `v`.

The sphere is playback-reactive today. True transient sync is the next audio
backend step: WeazlTunes will need an energy signal from decoded audio, while
Harmonica handles the smooth motion once that signal exists.

## Config

Edit `~/.config/weazltunes/config.json` when you want to hand-tune presets or
saved stations:

```json
{
  "presets": [
    {
      "name": "Prendie Radio",
      "url": "https://radio.prendie.io/radio.mp3"
    }
  ],
  "my_stations": []
}
```

WeazlTunes keeps the file human-readable because a radio preset should not need
a migration ceremony.
