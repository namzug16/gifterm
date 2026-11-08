# Changelog

Each release's section is published as its GitHub release notes.

## 0.3.0 - 2026-11-08

- `-fit` controls how the GIF fills the terminal, using the CSS `object-fit` values:
  `contain` (default, the whole GIF with black bars), `cover` (fills the terminal, cropping the
  edges) and `fill` (stretches to the terminal size). Unknown values fail with a clear error.

## 0.2.0 - 2026-11-08

- `gifterm <giphy-url>` downloads and plays a GIF straight from Giphy, e.g.
  `gifterm https://giphy.com/gifs/anime-chainsaw-man-reze-IKFVtPf8jP6KJH16dB`. Giphy page links,
  `i.giphy.com/<id>.gif` and `media.giphy.com/.../giphy.gif` links are supported.
- URLs from other websites, unsupported Giphy pages, non-GIF files, GIFs that don't exist and
  downloads over 50 MB fail with a clear error.

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
