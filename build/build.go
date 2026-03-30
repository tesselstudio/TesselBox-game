package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
)

// BuildConfig holds build configuration
type BuildConfig struct {
	OS       string
	Arch     string
	Output   string
	Release  bool
	Version  string
}

func main() {
	// Parse command line flags
	config := BuildConfig{}
	flag.StringVar(&config.OS, "os", runtime.GOOS, "Target operating system")
	flag.StringVar(&config.Arch, "arch", runtime.GOARCH, "Target architecture")
	flag.StringVar(&config.Output, "output", "", "Output file path")
	flag.BoolVar(&config.Release, "release", false, "Release build (optimized)")
	flag.StringVar(&config.Version, "version", "v0.3-alpha", "Version string")
	flag.Parse()

	// Set default output if not specified
	if config.Output == "" {
		ext := ""
		if config.OS == "windows" {
			ext = ".exe"
		}
		config.Output = fmt.Sprintf("tesselbox-%s-%s%s", config.OS, config.Arch, ext)
	}

	fmt.Printf("Building TesselBox v%s for %s/%s\n", config.Version, config.OS, config.Arch)
	fmt.Printf("Output: %s\n", config.Output)
	fmt.Printf("Current dir: %s\n", func() string { dir, _ := os.Getwd(); return dir }())

	// Create output directory
	outputDir := "bin"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Build arguments
	args := []string{"build", "-o", config.Output}
	
	// Add version info via ldflags
	ldflags := fmt.Sprintf("-X main.Version=%s", config.Version)
	
	// Add optimization flags for release builds
	if config.Release {
		ldflags += " -s -w"
		args = append(args, "-ldflags", ldflags)
		// Add trimpath as a separate build flag
		args = append(args, "-trimpath")
	} else {
		args = append(args, "-ldflags", ldflags)
	}
	
	// Add source file
	args = append(args, "cmd/main.go")

	// Set environment variables for cross-compilation
	env := os.Environ()
	env = append(env, fmt.Sprintf("GOOS=%s", config.OS))
	env = append(env, fmt.Sprintf("GOARCH=%s", config.Arch))
	
	// Only disable CGO for actual cross-compilation (not for current platform)
	currentOS := runtime.GOOS
	currentArch := runtime.GOARCH
	if config.OS != currentOS || config.Arch != currentArch {
		env = append(env, "CGO_ENABLED=0")
	}

	// Execute build from project root
	cmd := exec.Command("go", args...)
	cmd.Env = env
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("Build command: go %v\n", args)
	fmt.Printf("Build dir: %s\n", cmd.Dir)

	if err := cmd.Run(); err != nil {
		log.Fatalf("Build failed: %v", err)
	}

	fmt.Printf("✅ Build completed successfully!\n")

	// Add version info for release builds
	if config.Release {
		if err := addReleaseInfo(config); err != nil {
			fmt.Printf("Warning: Failed to add release info: %v\n", err)
		}
	}

	// Show file size
	if info, err := os.Stat(config.Output); err == nil {
		size := info.Size()
		if size < 1024 {
			fmt.Printf("Size: %d bytes\n", size)
		} else if size < 1024*1024 {
			fmt.Printf("Size: %.1f KB\n", float64(size)/1024)
		} else {
			fmt.Printf("Size: %.1f MB\n", float64(size)/(1024*1024))
		}
	}
}

func addReleaseInfo(config BuildConfig) error {
	fmt.Printf("📦 Release build completed for %s/%s\n", config.OS, config.Arch)
	return nil
}
