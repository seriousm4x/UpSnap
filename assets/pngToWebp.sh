#!/usr/bin/env bash

for i in *.png; do
    magick -quality 100 "$i" "$(basename "$i" .png)".webp
    rm "$i"
done
