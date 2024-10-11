###HOW TO DOWNLOAD A VIDEO

youtube-dl -f worstvideo

ffmpeg -i input_video.mp4 -vf "fps=24" frame_%04d.png

ffmpeg -i water.mp4 -vf "fps=24" input2/frame_%04d.png

ffmpeg -i br.mp4 -vf "fps=24" input/frame_%04d.png










# This converts a video to a GIF

ffmpeg -i test.mp4 -vf "fps=12" -c:v pam -f image2pipe - | convert - output.gif
