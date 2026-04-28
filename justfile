# https://just.systems

alias tw := tailwindcss

tailwindcss:
    tailwindcss -i static/tw.css -o static/main.css --watch

# Create a new empty blogpost template under the path.
[script("bash")]
new PATH:
    name="{{ PATH }}"
    basename=$(basename "$name")
    date=$(date +"%Y-%m-%dT%H:%M:%S%:z")

    echo "---" > "$name.md"
    echo "title: \"$basename\"" >> "$name.md"
    echo "description: \"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAH\"" >> "$name.md"
    echo "date: $date" >> "$name.md"
    echo "showDate: true" >> "$name.md"
    echo "tags: [ ]" >> "$name.md"
    echo "draft: true" >> "$name.md"
    echo "---" >> "$name.md"

generate:
    templ generate

alias d := dev

dev: generate
    air
