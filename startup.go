// startup.go
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"tesselbox/pkg/render"
)

// App holds shared services and state for the application lifecycle
type App struct {
	Game *render.Game // the main game instance

	ctx    context.Context
	cancel context.CancelFunc
}

// NewApp initializes the application with proper order and error handling
func NewApp() (*App, error) {
	ctx, cancel := context.WithCancel(context.Background())

	// No more OAuth / auth service — running fully local / guest mode by default

	// Create the game instance (no auth parameter needed anymore)
	game := render.NewGame() // ← adjust this call if your NewGame still expects arguments

	app := &App{
		Game:  game,
		ctx:   ctx,
		cancel: cancel,
	}

	// Basic validation
	if err := app.validate(); err != nil {
		return nil, err
	}

	return app, nil
}

// validate performs basic startup checks (expand later)
func (a *App) validate() error {
	if a.Game == nil {
		return errors.New("game instance not created")
	}
	return nil
}

// Run starts the application and blocks until shutdown
func (a *App) Run() error {
	// Set up OS signal channel for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Channel to capture errors from the game loop
	errChan := make(chan error, 1)

	// Run Ebitengine game loop in a separate goroutine (RunGame blocks)
	go func() {
		ebiten.SetWindowSize(render.ScreenWidth, render.ScreenHeight)
		ebiten.SetWindowTitle("Tesselbox - Go/Ebiten Version")
		ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

		if err := ebiten.RunGame(a.Game); err != nil {
			errChan <- err
		} else {
			errChan <- nil
		}
	}()

	// Wait for either shutdown signal or game error
	var shutdownErr error
	select {
	case sig := <-sigChan:
		log.Printf("Received shutdown signal: %v", sig)
	case err := <-errChan:
		if err != nil {
			shutdownErr = fmt.Errorf("game loop error: %w", err)
		}
		// normal game exit (uncommon)
	}

	// Begin graceful shutdown with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(a.ctx, 15*time.Second)
	defer shutdownCancel()

	// Future cleanup points can go here:
	// - world.SaveIfDirty(shutdownCtx)
	// - close any open files / connections
	// - wait for background goroutines

	// Wait until timeout expires (gives time for cleanup)
	<-shutdownCtx.Done()

	log.Println("Graceful shutdown complete")
	return shutdownErr
}
