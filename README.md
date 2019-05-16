# F1_trackSaver

> A track information saver for Codemasters F1 2018 game for PC, XBOX, and Playstation

> F1 does not provide information on their tracks so map creation is halted. With this the user will be able to easily create a map and easily use the M_worldPosition elements to position themself on the created map!

F1_trackSaver is Written in Go and Utilizes Redis and MySQL

---------------------------------------
  * [Features](#features)
  * [Requirements](#requirements)
  * [Installation](#installation)
  * [Usage](#usage)

---------------------------------------

## Features
  * Saves track information from F1
  * Optimizes for usage as a map, in this case for F1_GO
  * Splits map into sectors and turns for easy refrencing
  * Uses fatih/color to pretty output to the command line

## Requirements
  * Some version of go so you know, stuff actually works
  * Redigo, a Go client for the Redis database. *Make sure this is up and running before starting F1_trackSaver*
  * Go-MySQL-Driver, A MySQL-Driver for Go's database/sql package
  * Color, a ANSI color package to output colorized or SGR defined output to the standard output.

---------------------------------------

## Installation
Simple install the package to your [$GOPATH](https://github.com/golang/go/wiki/GOPATH "GOPATH") with the [go tool](https://golang.org/cmd/go/ "go command") from shell:
```bash
$ go get -u github.com/crocotelementry/F1_trackSaver
```
Make sure [Git is installed](https://git-scm.com/downloads) on your machine and in your system's `PATH`.

Until we find a way to have our requirements including in the F1_trackSaver package, We will also need to install four more items into your go path.

**redigo:**
```bash
$ go get github.com/gomodule/redigo/redis
```

**go-sql-driver:**
```bash
$ go get -u github.com/go-sql-driver/mysql
```

**Color**
```bash
$ go get github.com/fatih/color
```

## Usage
*F1_trackSaver* is ran by running the main executable. Some features that are critical to *F1_trackSaver's* usability are able to be ran from the terminal window in which you start *F1_trackSaver*, but it is not recommended. Before starting *F1_trackSaver*, make sure your Redis database is up and running, if it isn't, start it before starting *F1_trackSaver*.

To check if your Redis database is up and running:
```bash
$ redis-cli ping
```

If this returns PONG like below, then continue to starting *F1_trackSaver*:
```bash
$ redis-cli ping
PONG
```
If it does not return PONG, then start up your Redis database:
```bash
$ redis-server
```

To run *F1_trackSaver*:
```bash
go run *.go
```

---------------------------------------
