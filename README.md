<!--
Copyright (c) 2024 Thomas Mikalsen. Subject to the MIT License
-->

![alt gobrot](./doc/gobrot.jpg "Gobrot")


GoBrot
======

An interactive Mandelbrot Set explorer written in Go using the SDL2
library for platform-independent graphics. A multi-threaded (goroutine)
rendering process provides efficient and realtime exploration of this
fascinating mathematical object.

Pre-reqs:
---------

## Go

Install Go 1.13+

## SDL

Install SDL2 development libs and headers. Make sure that the `SDL2` header
files are in a directory called `SDL2` and specify that directory in the
`CGO_CFLAGS` env variable:  
```
export CGO_CFLAGS=-I/path/to/SDL2
```

Specify the location of the `SDL2` libraries in the `CGO_LDFLAGS` environment
variable:
```
export CGO_LDFLAGS=-L/path/to/sdl2libs
```

All of that should be straightforward on Mac and Ubuntu.
On Mac, use
```
brew install pkg-config
brew install sdl2
```

On Ubuntu, use `apt-get install libsdl2-dev`

On Windows, ensure that `SDL2.dll` is on your path.

On Windows, you may need to do a little work to ensure that the headers are
under a directory called `SDL2`. I copied all of the SDL2 includes to
`%GOPATH%\include\SDL2`, and copied the SDL2 libs to `%GOPATH%\lib\x64`, and copied `SDL2.dll` to
`%GOPATH%\bin`. I then setup my environment (in .bashrc) like so:

```
export GOBIN=${GOPATH}\\bin
export CGO_CFLAGS=-I${GOPATH}\\include
export CGO_LDFLAGS=-L${GOPATH}\\lib\\SDL2\\x64
export PATH="${PATH}:${GOBIN}"
```


Dependencies
------------

Get required Go dependencies

```
go get -v github.com/veandco/go-sdl2/sdl@master
```

Build
-----


```
make
```

for a clean build,
```
make clean all
```

Run It
------

```
make run
```

### Mouse Controls:

 * *Left Click-and-Drag*: Pan/scroll
 * *Left Double-click*: Center and zoom in
 * *Shift + Left Double-click*: Center and zoom out
 * *Mouse Wheel*: Zoom in and out
 * *Shift + Wheel*: Rotate
 * *Ctrl + Wheel*: Increase/decrease max number of iterations
 * *Alt + Wheel*: Increase/decrease number of entries in palette

### Keyboard Controls:
 * *Esc*: Exit the program
 * *Home*: Return to initial/home waypoint and settings
 * *Arrows*: Pan/scroll
 * *PgUp/PgDn*: Zoom in/out
 * *Shift+PgUp/PgDown*: zoom in/out more
 * *+/-*: Increase/decrease max number of iterations

