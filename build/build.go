package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	
	"tesselbox/assets"
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
	flag.StringVar(&config.Version, "version", "2.0.0", "Version string")
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
	outputDir := filepath.Dir(config.Output)
	if outputDir != "." {
		fullOutputDir := filepath.Join("..", outputDir)
		if err := os.MkdirAll(fullOutputDir, 0755); err != nil {
			log.Fatalf("Failed to create output directory: %v", err)
		}
	}

	// Prepare build command
	// Adjust output path for build context (running from parent dir)
	outputPath := config.Output
	if strings.HasPrefix(outputPath, "../") {
		outputPath = strings.TrimPrefix(outputPath, "../")
	}
	
	args := []string{
		"build",
		"-o", outputPath,
		"-ldflags", fmt.Sprintf("-X main.Version=%s -s -w", config.Version),
	}

	if config.Release {
		args = append(args, "-trimpath")
	}

	// Add source file (relative to project root)
	args = append(args, "cmd/main.go")

	// Set environment variables
	env := os.Environ()
	env = append(env, fmt.Sprintf("GOOS=%s", config.OS))
	env = append(env, fmt.Sprintf("GOARCH=%s", config.Arch))

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
	
	switch config.OS {
	case "windows":
		return embedWindowsIcon(config)
	case "darwin":
		return embedMacOSIcon(config)
	case "linux":
		return embedLinuxIcon(config)
	}
	
	return nil
}

func embedWindowsIcon(config BuildConfig) error {
	fmt.Printf("🎨 Embedding Windows icon...")
	
	// Check for rsrc tool
	rsrcPath := "rsrc"
	if _, err := exec.LookPath(rsrcPath); err != nil {
		// Try common Go bin paths
		homeDir, _ := os.UserHomeDir()
		possiblePaths := []string{
			filepath.Join(homeDir, "go", "bin", "rsrc"),
			filepath.Join(homeDir, ".local", "go", "bin", "rsrc"),
			"rsrc",
		}
		
		for _, path := range possiblePaths {
			if _, err := os.Stat(path); err == nil {
				rsrcPath = path
				break
			}
		}
		
		if _, err := os.Stat(rsrcPath); err != nil {
			fmt.Printf("⚠️  rsrc not found, skipping icon embedding\n")
			fmt.Printf("   Install with: go install github.com/akavel/rsrc@latest\n")
			return nil
		}
	}
	
	// Extract proper icon from embedded assets to temporary file
	iconData, err := extractEmbeddedIcon("tesselbox.ico")
	if err != nil {
		// Fallback to PNG if ICO not available
		iconData, err = extractEmbeddedIcon("icon.png")
		if err != nil {
			fmt.Printf("⚠️  Failed to extract embedded icon: %v\n", err)
			return nil
		}
	}
	
	tempIconPath := "temp_icon.png"
	if err := os.WriteFile(tempIconPath, iconData, 0644); err != nil {
		return fmt.Errorf("failed to write temporary icon: %v", err)
	}
	defer os.Remove(tempIconPath)
	
	// Create Windows manifest
	manifestPath := "windows.manifest"
	manifest := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<assembly xmlns="urn:schemas-microsoft-com:asm.v1" manifestVersion="1.0">
  <assemblyIdentity version="1.0.0.0" processorArchitecture="*" name="TesselBox" type="win32"/>
  <description>TesselBox - Hexagon Sandbox Game</description>
  <dependency>
    <dependentAssembly>
      <assemblyIdentity type="win32" name="Microsoft.Windows.Common-Controls" version="6.0.0.0" processorArchitecture="*" publicKeyToken="6595b64144ccf1df" language="*"/>
    </dependentAssembly>
  </dependency>
  <application xmlns="urn:schemas-microsoft-com:asm.v3">
    <windowsSettings>
      <dpiAware xmlns="http://schemas.microsoft.com/SMI/2005/WindowsSettings">true</dpiAware>
      <dpiAwareness xmlns="http://schemas.microsoft.com/SMI/2016/WindowsSettings">PerMonitorV2</dpiAwareness>
    </windowsSettings>
  </application>
</assembly>`
	
	if err := os.WriteFile(manifestPath, []byte(manifest), 0644); err != nil {
		return fmt.Errorf("failed to create manifest: %v", err)
	}
	defer os.Remove(manifestPath)
	
	// Run rsrc to generate .syso file using embedded icon
	cmd := exec.Command(rsrcPath, "-manifest", manifestPath, "-ico", tempIconPath, "-o", "rsrc.syso")
	if err := cmd.Run(); err != nil {
		fmt.Printf("⚠️  Failed to generate rsrc.syso: %v\n", err)
		return nil
	}
	defer os.Remove("rsrc.syso")
	
	fmt.Printf("✅ Windows icon embedded\n")
	return nil
}

func embedMacOSIcon(config BuildConfig) error {
	fmt.Printf("🎨 Setting up macOS icon...")
	
	// For macOS, we create an .app bundle for distribution
	appDir := fmt.Sprintf("%s.app", strings.TrimSuffix(filepath.Base(config.Output), filepath.Ext(config.Output)))
	appContents := filepath.Join(appDir, "Contents")
	appMacOS := filepath.Join(appContents, "MacOS")
	appResources := filepath.Join(appContents, "Resources")
	
	// Create app bundle structure
	if err := os.MkdirAll(appResources, 0755); err != nil {
		return fmt.Errorf("failed to create app bundle: %v", err)
	}
	
	// Copy binary to app bundle
	binaryPath := filepath.Join(appMacOS, "TesselBox")
	if err := os.Rename(config.Output, binaryPath); err != nil {
		return fmt.Errorf("failed to move binary to app bundle: %v", err)
	}
	
	// Extract proper icon from embedded assets
	iconData, err := extractEmbeddedIcon("TesselBox.icns")
	if err != nil {
		// Fallback to PNG if ICNS not available
		iconData, err = extractEmbeddedIcon("icon.png")
		if err != nil {
			fmt.Printf("⚠️  Failed to extract embedded icon: %v\n", err)
		} else {
			// Write embedded icon to app bundle
			iconDst := filepath.Join(appResources, "AppIcon.icns")
			if err := os.WriteFile(iconDst, iconData, 0644); err != nil {
				fmt.Printf("⚠️  Failed to write embedded icon: %v\n", err)
			}
		}
	} else {
		// Write embedded ICNS to app bundle
		iconDst := filepath.Join(appResources, "AppIcon.icns")
		if err := os.WriteFile(iconDst, iconData, 0644); err != nil {
			fmt.Printf("⚠️  Failed to write embedded icon: %v\n", err)
		}
	}
	
	// Create Info.plist
	infoPlist := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleExecutable</key>
    <string>TesselBox</string>
    <key>CFBundleIconFile</key>
    <string>AppIcon.icns</string>
    <key>CFBundleIdentifier</key>
    <string>com.tesselbox.game</string>
    <key>CFBundleName</key>
    <string>TesselBox</string>
    <key>CFBundleVersion</key>
    <string>%s</string>
    <key>CFBundleShortVersionString</key>
    <string>%s</string>
    <key>CFBundlePackageType</key>
    <string>APPL</string>
    <key>NSHighResolutionCapable</key>
    <true/>
    <key>NSSupportsAutomaticGraphicsSwitching</key>
    <true/>
</dict>
</plist>`, config.Version, config.Version)
	
	infoPath := filepath.Join(appContents, "Info.plist")
	if err := os.WriteFile(infoPath, []byte(infoPlist), 0644); err != nil {
		return fmt.Errorf("failed to create Info.plist: %v", err)
	}
	
	fmt.Printf("✅ macOS .app bundle created: %s\n", appDir)
	return nil
}

func embedLinuxIcon(config BuildConfig) error {
	fmt.Printf("🎨 Setting up Linux icon...")
	
	// For Linux, create a .desktop file
	desktopFile := fmt.Sprintf(`[Desktop Entry]
Version=1.0
Type=Application
Name=TesselBox
Comment=Hexagon Sandbox Game
Exec=%s
Icon=tesselbox
Terminal=false
Categories=Game;Simulation;
StartupNotify=true
Keywords=game;sandbox;hexagon;blocks;crafting;`, config.Output)
	
	desktopPath := "tesselbox.desktop"
	if err := os.WriteFile(desktopPath, []byte(desktopFile), 0644); err != nil {
		return fmt.Errorf("failed to create desktop file: %v", err)
	}
	defer os.Remove(desktopPath)
	
	// Extract proper icon from embedded assets for Linux integration
	iconData, err := extractEmbeddedIcon("tesselbox.png")
	if err != nil {
		// Fallback to generic icon.png
		iconData, err = extractEmbeddedIcon("icon.png")
		if err != nil {
			fmt.Printf("⚠️  Failed to extract embedded icon: %v\n", err)
		} else {
			iconDst := "tesselbox.png"
			if err := os.WriteFile(iconDst, iconData, 0644); err != nil {
				fmt.Printf("⚠️  Failed to write embedded icon: %v\n", err)
			} else {
				defer os.Remove(iconDst)
			}
		}
	} else {
		iconDst := "tesselbox.png"
		if err := os.WriteFile(iconDst, iconData, 0644); err != nil {
			fmt.Printf("⚠️  Failed to write embedded icon: %v\n", err)
		} else {
			defer os.Remove(iconDst)
		}
	}
	
	fmt.Printf("✅ Linux desktop integration created\n")
	fmt.Printf("   Install with: cp tesselbox.desktop ~/.local/share/applications/\n")
	fmt.Printf("   Copy icon to: ~/.local/share/icons/tesselbox.png\n")
	
	return nil
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0644)
}

// extractEmbeddedIcon extracts an icon from the embedded assets
func extractEmbeddedIcon(filename string) ([]byte, error) {
	return assets.GetIconFile(filename)
}
