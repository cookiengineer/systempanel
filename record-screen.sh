#!/bin/bash

ffmpeg -f x11grab -video_size 960x1080 -i :0.0+960,0 demo.mp4;
sleep 3;
ffmpeg -i demo.mp4 -vf "fps=10,scale=480:-1:flags=lanczos,palettegen" demo-palette.png;
ffmpeg -i demo.mp4 -i palette.png -lavfi "fps=10,scale=480:-1:flags=lanczos[x];[x][1:v]paletteuse" docs/demo.gif;
