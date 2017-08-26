#!/bin/bash
if [ "$1" = "linux" ] || [ "$1" = "" ]; then
    echo "Generating Linux binary"
    env GOOS=linux GOARCH=amd64 go build -o bin/grawler-linux
fi
if [ "$1" = "windows" ] || [ "$1" = "" ]; then
    echo "Generating Windows binary"
    env GOOS=windows GOARCH=amd64 go build -o bin/grawler-windows.exe
fi
if [ "$1" = "mac" ] || [ "$1" = "" ]; then
    echo "Generating MacOS binary"
    env GOOS=darwin GOARCH=amd64 go build -o bin/grawler-darwin
fi
echo "Done!"
ls -al bin/
