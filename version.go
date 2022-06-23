package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

const versionMajor = 0

const versionMinor = 1

const versionPatch = 0

func printVersionAndExit() {
	log.Println(fmt.Sprintf("%s version %s", filepath.Base(os.Args[0]), versionString()))
	os.Exit(0)
}

func versionString() string {
	return fmt.Sprintf("%d.%d.%d", versionMajor, versionMinor, versionPatch)
}
