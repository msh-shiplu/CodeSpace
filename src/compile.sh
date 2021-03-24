echo "compiling osx executable"
go build -o gem
echo "done compiling for osx. Output: gem"
echo "compiling win executable"
env GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc go build -o gem.exe
echo "done compiling for windows. Output: gem.exe"
echo "compiling linux executable"
env GOOS=linux GOARCH=amd64 go build -o gem_linux_amd64
