###HOW TO DOWNLOAD A VIDEO

youtube-dl -f worstvideo

ffmpeg -i input_video.mp4 -vf "fps=24" frame_%04d.png

ffmpeg -i water.mp4 -vf "fps=24" input2/frame_%04d.png

