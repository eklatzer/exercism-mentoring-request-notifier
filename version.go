package main

import (
	"fmt"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

const versionMajor = 0

const versionMinor = 2

const versionPatch = 0

func printVersionAndExit() {
	log.Println(fmt.Sprintf("%s version %s", filepath.Base(os.Args[0]), versionString()))
	os.Exit(0)
}

func versionString() string {
	return fmt.Sprintf("%d.%d.%d", versionMajor, versionMinor, versionPatch)
}
