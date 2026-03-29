package assets

import (
	"embed"
)

//go:embed config/*.yaml
var ConfigFS embed.FS

// GetConfigFile reads a config file from the embedded filesystem
func GetConfigFile(filename string) ([]byte, error) {
	return ConfigFS.ReadFile("config/" + filename)
}
