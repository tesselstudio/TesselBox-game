package main

import (
	"fmt"
	"os"
	"path/filepath"
)

//go:generate go run github.com/josephspurrier/goversioninfo
func main() {
	// This file is used for generating version info for Windows builds
	// The go:generate directive above will create versioninfo.go
	fmt.Println("Version info generator for TesselBox")
}
