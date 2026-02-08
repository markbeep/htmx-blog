# https://just.systems

alias tw := tailwindcss

tailwindcss:
    tailwindcss -i static/tw.css -o static/main.css --watch

new PATH:
    ./post.sh {{ PATH }}

alias d := dev

dev:
    air
