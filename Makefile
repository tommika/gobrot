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

