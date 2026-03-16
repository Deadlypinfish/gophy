#!/bin/bash
set -e

echo "Building for Linux..."
GOOS=linux GOARCH=amd64 go build -ldflags "-X 'main.buildTime=$(date '+%Y-%m-%d %H:%M:%S')'" -o gophy .

echo "Copying binary to server..."
scp gophy forge:~/apps/gophy/gophy.new

echo "Swapping binary and restarting..."
ssh forge 'export XDG_RUNTIME_DIR=/run/user/$(id -u) && mv ~/apps/gophy/gophy.new ~/apps/gophy/gophy && systemctl --user restart gophy'

rm gophy
echo "Deployed! https://gif.randomthings.org"
