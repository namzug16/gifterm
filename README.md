# GifTerm
GIFs on the terminal

# Build

You can build and install the tool with the Makefile

```$ make```

# Usage

```$ gifterm input.gif```


# How to convert a video into a GIF

Make sure to have installed ffmpeg and imagemagick and ten run

```$ ffmpeg -i <input.mp4> -vf "fps=12" -c:v pam -f image2pipe - | convert - output.gif```
