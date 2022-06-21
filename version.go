package main

import "fmt"

const VersionMajor = 0

const VersionMinor = 1

const VersionPatch = 0

func versionString() string {
	return fmt.Sprintf("%d.%d.%d", VersionMajor, VersionMinor, VersionPatch)
}
