# GifTerm

GIFs on the terminal

<img src="./assets/gifterm.gif" alt="Gifterm demo" style="width: 100%;">

*The terminal output is actual text*

## Contents

- [Installation](#installation)
- [Usage](#usage)
- [Video to Gif](#how-to-convert-a-video-into-a-gif)

## Installation

You can simply install it by running

```
$ make
```

Just make sure to have installed and setup the minimum required version of Go. You can find it in the go.mod file 

## Usage

### Basic usage

```
$ gifterm <input.gif>
```

### Flags

- **cd**: Character Density. Default ".,:-=i|%O#@$X"
- **fps**: FPS. Default 12
- **randomBlank**: Set if a random character from CD should be picked for a blank pixel
- **ofg**: Only Foreground. Set if only the foregroud color should be set

> The Font Aspect Ratio changes from terminal to terminal, this value works in mine, so make sure to play a round to find yours

```
$ gifterm -cd ".,:-=i|%O#@$X" -fps 12 -randomBlank <input.gif> 
```
Will produce something like this

<img src="./assets/gifterm_random_blank.gif" alt="Gifterm demo" style="width: 100%;">

## How to convert a video into a GIF

Make sure to have installed ffmpeg and imagemagick and then run

```
$ ffmpeg -i <input.mp4> -vf "fps=12" -c:v pam -f image2pipe - | convert - input.gif
```
