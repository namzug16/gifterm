# Changelog

Each release's section is published as its GitHub release notes.

## 0.1.0 - 2026-11-08

The first release, for macOS and Linux on arm64 and amd64.

- Install with `brew install namzug16/tap/gifterm`.
- `gifterm <input.gif>` plays a GIF in the terminal as colored ASCII art, refitting it when the
  terminal is resized without stretching or cropping.
- `-fps` sets the playback speed, `-cd` the character density string (multibyte characters
  allowed), `-randomBlank` fills blank pixels with random characters and `-ofg` colors only the
  foreground.
- `-car` overrides the character aspect ratio detected from the terminal, which falls back to 2.0
  when the terminal reports no pixel size.
- GIF disposal methods and canvas size are honored, so partial-frame GIFs render correctly.
- Flags are accepted before or after the file, and `-version` prints the version.
- Quit with `q`, `esc` or `ctrl+c`.
