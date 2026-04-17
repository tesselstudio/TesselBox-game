// Package ui implements on-screen notification system for user feedback
package ui

import (
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// NotificationType represents the type of notification
type NotificationType int

const (
	NotificationTypeInfo NotificationType = iota
	NotificationTypeWarning
	NotificationTypeError
	NotificationTypeSuccess
)

// Notification represents a single on-screen notification
type Notification struct {
	Message   string
	Type      NotificationType
	StartTime time.Time
	Duration  time.Duration
}

// NotificationManager manages and displays on-screen notifications
type NotificationManager struct {
	notifications []*Notification
	maxVisible    int
	screenWidth   int
	screenHeight  int
}

// NewNotificationManager creates a new notification manager
func NewNotificationManager(screenWidth, screenHeight int) *NotificationManager {
	return &NotificationManager{
		notifications: make([]*Notification, 0),
		maxVisible:    5,
		screenWidth:   screenWidth,
		screenHeight:  screenHeight,
	}
}

// Add adds a new notification
func (nm *NotificationManager) Add(message string, notifType NotificationType, duration time.Duration) {
	notif := &Notification{
		Message:   message,
		Type:      notifType,
		StartTime: time.Now(),
		Duration:  duration,
	}
	nm.notifications = append(nm.notifications, notif)
}

// AddInfo adds an info notification
func (nm *NotificationManager) AddInfo(message string) {
	nm.Add(message, NotificationTypeInfo, 3*time.Second)
}

// AddWarning adds a warning notification
func (nm *NotificationManager) AddWarning(message string) {
	nm.Add(message, NotificationTypeWarning, 5*time.Second)
}

// AddError adds an error notification
func (nm *NotificationManager) AddError(message string) {
	nm.Add(message, NotificationTypeError, 8*time.Second)
}

// AddSuccess adds a success notification
func (nm *NotificationManager) AddSuccess(message string) {
	nm.Add(message, NotificationTypeSuccess, 3*time.Second)
}

// Update removes expired notifications
func (nm *NotificationManager) Update() {
	now := time.Now()
	active := make([]*Notification, 0)
	for _, notif := range nm.notifications {
		if now.Sub(notif.StartTime) < notif.Duration {
			active = append(active, notif)
		}
	}
	nm.notifications = active
}

// Draw renders the notifications
func (nm *NotificationManager) Draw(screen *ebiten.Image) {
	if len(nm.notifications) == 0 {
		return
	}

	// Show only the most recent notifications
	visibleCount := len(nm.notifications)
	if visibleCount > nm.maxVisible {
		visibleCount = nm.maxVisible
	}

	startY := float64(nm.screenHeight) - 100
	for i := visibleCount - 1; i >= 0; i-- {
		notif := nm.notifications[i]
		y := startY - float64(visibleCount-1-i)*40

		// Determine color based on type
		var bgColor color.RGBA
		switch notif.Type {
		case NotificationTypeInfo:
			bgColor = color.RGBA{50, 100, 200, 230}
		case NotificationTypeWarning:
			bgColor = color.RGBA{200, 150, 50, 230}
		case NotificationTypeError:
			bgColor = color.RGBA{200, 50, 50, 230}
		case NotificationTypeSuccess:
			bgColor = color.RGBA{50, 200, 100, 230}
		}

		// Draw background
		bgWidth := 400.0
		bgHeight := 35.0
		bgX := float64(nm.screenWidth)/2 - bgWidth/2
		ebitenutil.DrawRect(screen, bgX, y, bgWidth, bgHeight, bgColor)

		// Draw border
		borderColor := color.RGBA{255, 255, 255, 200}
		ebitenutil.DrawRect(screen, bgX, y, bgWidth, 2, borderColor)
		ebitenutil.DrawRect(screen, bgX, y+bgHeight-2, bgWidth, 2, borderColor)
		ebitenutil.DrawRect(screen, bgX, y, 2, bgHeight, borderColor)
		ebitenutil.DrawRect(screen, bgX+bgWidth-2, y, 2, bgHeight, borderColor)

		// Draw text
		ebitenutil.DebugPrintAt(screen, notif.Message, int(bgX)+10, int(y)+10)
	}
}

// Clear removes all notifications
func (nm *NotificationManager) Clear() {
	nm.notifications = make([]*Notification, 0)
}
