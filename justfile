# https://just.systems

alias tw := tailwindcss

tailwindcss:
    tailwindcss -i static/tw.css -o static/main.css --watch

new PATH:
    ./post.sh {{ PATH }}

generate:
    templ generate

alias d := dev

dev: generate
    air
