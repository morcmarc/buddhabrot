buddhabrot
==========

Simple buddhabrot generator written in Go. 

![alt tag](https://raw.github.com/morcmarc/buddhabrot/master/example.png)

## Requirements

The renderer uses SDL. On Mac you can install the library via homebrew:

```
$ brew install sdl 
$ brew install sdl_image
$ brew install sdl_ttf
```

## Install

```
$ go get ./...
$ go install
```

## Usage

```
$ buddhabrot ~/buddha.png
```

This will open a window and display the render progress more or less real-time.
When you're happy with the end result just close the window and the image will
be saved into the given file.

```
NAME:
   Buddhabrot - generate buddahbrot fractals

USAGE:
   Buddhabrot [global options] command [command options] [arguments...]

VERSION:
   0.1.0

COMMANDS:
   help, h  Shows a list of commands or help for one command
   
GLOBAL OPTIONS:
   --width, -w '600'    
   --height, -t '600'   
   --color, -c '50,200,500' color palette, if only one value given (i.e.: '-c 40') it will render a greyscale image
   --help, -h     show help
   --version, -v    print the version
```