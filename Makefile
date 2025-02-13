# vim: ts=4 sw=4
# Copyright (c) 2019, 2024 Thomas Mikalsen. Subject to the MIT License 
ifeq ($(OS),Windows_NT)
EXE_X=.exe
else
EXE_X=
endif

EXE_NAME=gomandelbrot$(EXE_X)

all: build

build:
	mkdir -p bin
	go mod tidy
	go vet ./...
	go build -o bin/$(EXE_NAME)

run:
	./bin/$(EXE_NAME)

clean:
	rm -rf bin

