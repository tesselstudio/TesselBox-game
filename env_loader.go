package main

import (
	"bufio"
	"os"
	"strings"
	"log"
)

// Updated to accept a filename string to match your main.go call
func loadEnvFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Printf("Warning: Could not open %s: %v\n", filename, err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		
		// Skip empty lines or comments
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Split by the first "=" found
		if parts := strings.SplitN(line, "=", 2); len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			
			// Remove quotes if they exist around the value
			value = strings.Trim(value, `"'`)
			
			os.Setenv(key, value)
		}
	}
	log.Printf("Loaded %s successfully\n", filename)
}