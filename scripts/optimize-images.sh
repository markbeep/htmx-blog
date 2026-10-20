#!/bin/sh
set -eu

dir="${1:-content}"

if ! command -v cwebp >/dev/null 2>&1 && ! command -v magick >/dev/null 2>&1; then
  echo "cwebp or magick is required to convert images" >&2
  exit 1
fi

encode_static() {
  src="$1"
  dst="$2"
  if command -v cwebp >/dev/null 2>&1; then
    cwebp -quiet -q 82 -m 6 -sharp_yuv -mt "$src" -o "$dst"
  else
    magick "$src" -quality 82 -define webp:method=6 "$dst"
  fi
}

encode_animation() {
  src="$1"
  dst="$2"
  if command -v gif2webp >/dev/null 2>&1; then
    gif2webp -quiet -q 80 -m 6 -mt "$src" -o "$dst"
  else
    magick "$src" -quality 80 -define webp:method=6 "$dst"
  fi
}

find "$dir" -type f \( -iname '*.png' -o -iname '*.jpg' -o -iname '*.jpeg' \) | while IFS= read -r src; do
  dst="${src%.*}.webp"
  if [ -f "$dst" ] && [ "$dst" -nt "$src" ]; then
    continue
  fi
  if ! encode_static "$src" "$dst"; then
    echo "failed to convert $src" >&2
    rm -f "$dst"
  fi
done

find "$dir" -type f -iname '*.gif' | while IFS= read -r src; do
  dst="${src%.*}.webp"
  if [ -f "$dst" ] && [ "$dst" -nt "$src" ]; then
    continue
  fi
  if ! encode_animation "$src" "$dst"; then
    echo "failed to convert $src" >&2
    rm -f "$dst"
  fi
done
