#!/bin/sh

name="$1"
basename=$(basename "$name")
date=$(date +"%Y-%m-%dT%H:%M:%S%:z")

echo "---" > "$name.md"
echo "title: \"$basename\"" >> "$name.md"
echo "date: $date" >> "$name.md"
echo "showDate: true" >> "$name.md"
echo "tags: [ ]" >> "$name.md"
echo "draft: true" >> "$name.md"
echo "---" >> "$name.md"
