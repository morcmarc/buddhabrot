buddhabrot
==========

Simple buddhabrot generator written in Go. 

![alt tag](https://raw.github.com/morcmarc/buddhabrot/master/example.png)

## Install

```
$ go get ./...
$ go install
```

## Usage

```
$ buddhabrot ~/buddha.png
```

The process will stop whenever you press control-C or send a SIGTERM to it.
Before shutting down it will render the image into the given file.

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