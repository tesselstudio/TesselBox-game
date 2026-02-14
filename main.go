// main.go
package main

import "log"

func main() {
    app, err := NewApp()
    if err != nil {
        log.Fatalf("Startup failed: %v", err)
    }

    if err := app.Run(); err != nil {
        log.Printf("Application exited with error: %v", err)
    }
}
