package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

const VersionMajor = 0

const VersionMinor = 1

const VersionPatch = 0

func printVersionAndExit() {
	log.Println(fmt.Sprintf("%s version %s", filepath.Base(os.Args[0]), versionString()))
	os.Exit(0)
}

func versionString() string {
	return fmt.Sprintf("%d.%d.%d", VersionMajor, VersionMinor, VersionPatch)
}
