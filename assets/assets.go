package assets

import (
	"embed"
	"io/fs"
)

//go:embed config/*.yaml
var ConfigFS embed.FS

//go:embed icons/*
var IconFS embed.FS

// GetConfigFile reads a config file from the embedded filesystem
func GetConfigFile(filename string) ([]byte, error) {
	return ConfigFS.ReadFile("config/" + filename)
}

// GetIconFile reads an icon file from the embedded filesystem
func GetIconFile(filename string) ([]byte, error) {
	return IconFS.ReadFile("icons/" + filename)
}

// ListIcons returns a list of all embedded icon files
func ListIcons() ([]string, error) {
	var icons []string
	err := fs.WalkDir(IconFS, "icons", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			// Remove "icons/" prefix for cleaner names
			if len(path) > 6 {
				icons = append(icons, path[6:])
			}
		}
		return nil
	})
	return icons, err
}
