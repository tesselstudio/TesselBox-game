// Package ui implements the chest interface
package ui

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"tesselbox/pkg/chest"
	"tesselbox/pkg/items"
)

// ChestUI represents the chest inventory interface
type ChestUI struct {
	Open          bool
	ScreenWidth   int
	ScreenHeight  int

	// References
	ChestManager  *chest.ChestManager
	PlayerInventory *items.Inventory

	// Current chest
	CurrentChestX float64
	CurrentChestY float64

	// UI Layout
	ChestGridX    float64
	ChestGridY    float64
	PlayerInvX    float64
	PlayerInvY    float64

	// Slot sizes
	SlotSize      float64
	SlotSpacing   float64

	// Drag and drop
	DraggedItem   *items.Item
	DraggedQuantity int
	DragX, DragY  float64

	// Hover state
	HoveredSlot   int
	HoveredArea   ChestArea // 'chest' or 'player'

	// Selected slots
	SelectedChestSlot   int
	SelectedPlayerSlot  int
}

// ChestArea represents which inventory area is being interacted with
type ChestArea int

const (
	AreaNone ChestArea = iota
	AreaChest
	AreaPlayer
)

// NewChestUI creates a new chest UI
func NewChestUI(screenWidth, screenHeight int, cm *chest.ChestManager, inv *items.Inventory) *ChestUI {
	ui := &ChestUI{
		ScreenWidth:     screenWidth,
		ScreenHeight:    screenHeight,
		ChestManager:    cm,
		PlayerInventory: inv,
		SlotSize:        48,
		SlotSpacing:     4,
		SelectedChestSlot: -1,
		SelectedPlayerSlot: -1,
	}

	ui.calculateLayout()
	return ui
}

// calculateLayout positions UI elements
func (ui *ChestUI) calculateLayout() {
	centerX := float64(ui.ScreenWidth) / 2

	// Chest grid at top (9x3)
	ui.ChestGridX = centerX - (9*ui.SlotSize+8*ui.SlotSpacing)/2
	ui.ChestGridY = 100

	// Player inventory below (9x3 main + 9 hotbar)
	ui.PlayerInvX = ui.ChestGridX
	ui.PlayerInvY = ui.ChestGridY + 3*ui.SlotSize + 3*ui.SlotSpacing + 40
}

// OpenChest opens the UI for a specific chest
func (ui *ChestUI) OpenChest(x, y float64) {
	ui.CurrentChestX = x
	ui.CurrentChestY = y
	ui.Open = true
	ui.SelectedChestSlot = -1
	ui.SelectedPlayerSlot = -1
}

// Close closes the chest UI
func (ui *ChestUI) Close() {
	ui.Open = false
	ui.DraggedItem = nil
	ui.DraggedQuantity = 0
}

// Update handles input and updates UI state
func (ui *ChestUI) Update() error {
	if !ui.Open {
		return nil
	}

	// Get mouse position
	mx, my := ebiten.CursorPosition()
	mouseX, mouseY := float64(mx), float64(my)

	// Update hover state
	ui.updateHover(mouseX, mouseY)

	// Handle clicks
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		ui.handleClick(mouseX, mouseY)
	}

	// Handle right clicks (quick transfer)
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		ui.handleRightClick(mouseX, mouseY)
	}

	// Update dragged item position
	if ui.DraggedItem != nil {
		ui.DragX = mouseX - ui.SlotSize/2
		ui.DragY = mouseY - ui.SlotSize/2
	}

	// Drop dragged item on mouse release
	if ui.DraggedItem != nil && inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		ui.handleDrop(mouseX, mouseY)
	}

	return nil
}

// updateHover updates the hovered slot
func (ui *ChestUI) updateHover(mx, my float64) {
	ui.HoveredSlot = -1
	ui.HoveredArea = AreaNone

	// Check chest grid
	for row := 0; row < 3; row++ {
		for col := 0; col < 9; col++ {
			slotIdx := row*9 + col
			slotX := ui.ChestGridX + float64(col)*(ui.SlotSize+ui.SlotSpacing)
			slotY := ui.ChestGridY + float64(row)*(ui.SlotSize+ui.SlotSpacing)
			if mx >= slotX && mx <= slotX+ui.SlotSize && my >= slotY && my <= slotY+ui.SlotSize {
				ui.HoveredSlot = slotIdx
				ui.HoveredArea = AreaChest
				return
			}
		}
	}

	// Check player inventory
	for row := 0; row < 3; row++ {
		for col := 0; col < 9; col++ {
			slotIdx := 9 + row*9 + col // Skip hotbar
			slotX := ui.PlayerInvX + float64(col)*(ui.SlotSize+ui.SlotSpacing)
			slotY := ui.PlayerInvY + float64(row)*(ui.SlotSize+ui.SlotSpacing)
			if mx >= slotX && mx <= slotX+ui.SlotSize && my >= slotY && my <= slotY+ui.SlotSize {
				ui.HoveredSlot = slotIdx
				ui.HoveredArea = AreaPlayer
				return
			}
		}
	}

	// Check hotbar
	for i := 0; i < 9; i++ {
		slotX := ui.PlayerInvX + float64(i)*(ui.SlotSize+ui.SlotSpacing)
		slotY := ui.PlayerInvY + 3*ui.SlotSize + 3*ui.SlotSpacing + 20
		if mx >= slotX && mx <= slotX+ui.SlotSize && my >= slotY && my <= slotY+ui.SlotSize {
			ui.HoveredSlot = i
			ui.HoveredArea = AreaPlayer
			return
		}
	}
}

// handleClick handles left mouse button clicks
func (ui *ChestUI) handleClick(mx, my float64) {
	if ui.HoveredSlot < 0 {
		return
	}

	switch ui.HoveredArea {
	case AreaChest:
		// Get item from chest
		contents := ui.ChestManager.GetChestContents(ui.CurrentChestX, ui.CurrentChestY)
		if ui.HoveredSlot < len(contents) && contents[ui.HoveredSlot].Type != items.NONE {
			item := contents[ui.HoveredSlot]
			ui.DraggedItem = &item
			ui.DraggedQuantity = item.Quantity
			ui.SelectedChestSlot = ui.HoveredSlot

			// Remove from chest
			ui.ChestManager.RemoveItemFromChest(ui.CurrentChestX, ui.CurrentChestY, ui.HoveredSlot, item.Quantity)
		}

	case AreaPlayer:
		// Get item from player inventory
		if ui.HoveredSlot >= 0 && ui.HoveredSlot < len(ui.PlayerInventory.Slots) {
			item := &ui.PlayerInventory.Slots[ui.HoveredSlot]
			if item.Type != items.NONE {
				ui.DraggedItem = item
				ui.DraggedQuantity = item.Quantity
				ui.SelectedPlayerSlot = ui.HoveredSlot

				// Clear the slot
				item.Type = items.NONE
				item.Quantity = 0
				item.Durability = -1
			}
		}
	}
}

// handleRightClick handles right mouse button clicks (quick transfer)
func (ui *ChestUI) handleRightClick(mx, my float64) {
	if ui.HoveredSlot < 0 {
		return
	}

	switch ui.HoveredArea {
	case AreaChest:
		// Transfer from chest to player
		contents := ui.ChestManager.GetChestContents(ui.CurrentChestX, ui.CurrentChestY)
		if ui.HoveredSlot < len(contents) && contents[ui.HoveredSlot].Type != items.NONE {
			item := contents[ui.HoveredSlot]
			// Add to player inventory
			if ui.PlayerInventory.AddItem(item.Type, item.Quantity) {
				// Remove from chest
				ui.ChestManager.RemoveItemFromChest(ui.CurrentChestX, ui.CurrentChestY, ui.HoveredSlot, item.Quantity)
			}
		}

	case AreaPlayer:
		// Transfer from player to chest
		if ui.HoveredSlot >= 0 && ui.HoveredSlot < len(ui.PlayerInventory.Slots) {
			item := ui.PlayerInventory.Slots[ui.HoveredSlot]
			if item.Type != items.NONE {
				// Add to chest
				if ui.ChestManager.AddItemToChest(ui.CurrentChestX, ui.CurrentChestY, item.Type, item.Quantity) {
					// Remove from player
					ui.PlayerInventory.Slots[ui.HoveredSlot] = items.Item{Type: items.NONE, Quantity: 0, Durability: -1}
				}
			}
		}
	}
}

// handleDrop handles dropping a dragged item
func (ui *ChestUI) handleDrop(mx, my float64) {
	if ui.DraggedItem == nil {
		return
	}

	// Try to place in hovered slot
	if ui.HoveredSlot >= 0 {
		switch ui.HoveredArea {
		case AreaChest:
			// Place in chest
			contents := ui.ChestManager.GetChestContents(ui.CurrentChestX, ui.CurrentChestY)
			if ui.HoveredSlot < len(contents) {
				ui.DraggedItem.Quantity = ui.DraggedQuantity
				ui.ChestManager.GetChest(ui.CurrentChestX, ui.CurrentChestY).Slots[ui.HoveredSlot] = *ui.DraggedItem
			}

		case AreaPlayer:
			// Place in player inventory
			if ui.HoveredSlot >= 0 && ui.HoveredSlot < len(ui.PlayerInventory.Slots) {
				ui.DraggedItem.Quantity = ui.DraggedQuantity
				ui.PlayerInventory.Slots[ui.HoveredSlot] = *ui.DraggedItem
			}
		}
	}

	// Clear drag state
	ui.DraggedItem = nil
	ui.DraggedQuantity = 0
}

// Draw renders the chest UI
func (ui *ChestUI) Draw(screen *ebiten.Image) {
	if !ui.Open {
		return
	}

	// Draw semi-transparent background
	bgColor := color.RGBA{20, 20, 30, 230}
	ebitenutil.DrawRect(screen, 0, 0, float64(ui.ScreenWidth), float64(ui.ScreenHeight), bgColor)

	// Draw chest label
	ebitenutil.DebugPrintAt(screen, "CHEST", int(ui.ChestGridX)+180, int(ui.ChestGridY-30))

	// Draw chest grid
	ui.drawChestGrid(screen)

	// Draw player inventory label
	ebitenutil.DebugPrintAt(screen, "INVENTORY", int(ui.PlayerInvX)+160, int(ui.PlayerInvY-30))

	// Draw player inventory
	ui.drawPlayerInventory(screen)

	// Draw dragged item
	if ui.DraggedItem != nil {
		ui.drawItem(screen, ui.DraggedItem, ui.DragX, ui.DragY)
	}
}

// drawChestGrid draws the chest inventory grid
func (ui *ChestUI) drawChestGrid(screen *ebiten.Image) {
	contents := ui.ChestManager.GetChestContents(ui.CurrentChestX, ui.CurrentChestY)

	for row := 0; row < 3; row++ {
		for col := 0; col < 9; col++ {
			slotIdx := row*9 + col
			slotX := ui.ChestGridX + float64(col)*(ui.SlotSize+ui.SlotSpacing)
			slotY := ui.ChestGridY + float64(row)*(ui.SlotSize+ui.SlotSpacing)

			// Draw slot background
			bgColor := color.RGBA{60, 60, 70, 255}
			if ui.HoveredSlot == slotIdx && ui.HoveredArea == AreaChest {
				bgColor = color.RGBA{80, 80, 100, 255}
			}
			if slotIdx == ui.SelectedChestSlot {
				bgColor = color.RGBA{100, 100, 130, 255}
			}
			ebitenutil.DrawRect(screen, slotX, slotY, ui.SlotSize, ui.SlotSize, bgColor)

			// Draw border
			borderColor := color.RGBA{100, 100, 120, 255}
			ebitenutil.DrawRect(screen, slotX, slotY, ui.SlotSize, 2, borderColor)
			ebitenutil.DrawRect(screen, slotX, slotY+ui.SlotSize-2, ui.SlotSize, 2, borderColor)
			ebitenutil.DrawRect(screen, slotX, slotY, 2, ui.SlotSize, borderColor)
			ebitenutil.DrawRect(screen, slotX+ui.SlotSize-2, slotY, 2, ui.SlotSize, borderColor)

			// Draw item if present
			if slotIdx < len(contents) && contents[slotIdx].Type != items.NONE {
				ui.drawItem(screen, &contents[slotIdx], slotX+4, slotY+4)
			}
		}
	}
}

// drawPlayerInventory draws the player inventory
func (ui *ChestUI) drawPlayerInventory(screen *ebiten.Image) {
	// Draw main inventory (9x3)
	for row := 0; row < 3; row++ {
		for col := 0; col < 9; col++ {
			slotIdx := 9 + row*9 + col
			slotX := ui.PlayerInvX + float64(col)*(ui.SlotSize+ui.SlotSpacing)
			slotY := ui.PlayerInvY + float64(row)*(ui.SlotSize+ui.SlotSpacing)

			ui.drawPlayerSlot(screen, slotIdx, slotX, slotY)
		}
	}

	// Draw hotbar
	hotbarY := ui.PlayerInvY + 3*ui.SlotSize + 3*ui.SlotSpacing + 20
	for i := 0; i < 9; i++ {
		slotX := ui.PlayerInvX + float64(i)*(ui.SlotSize+ui.SlotSpacing)
		ui.drawPlayerSlot(screen, i, slotX, hotbarY)
	}
}

// drawPlayerSlot draws a single player inventory slot
func (ui *ChestUI) drawPlayerSlot(screen *ebiten.Image, slotIdx int, x, y float64) {
	// Draw slot background
	bgColor := color.RGBA{50, 50, 60, 255}
	if ui.HoveredSlot == slotIdx && ui.HoveredArea == AreaPlayer {
		bgColor = color.RGBA{70, 70, 90, 255}
	}
	if slotIdx == ui.SelectedPlayerSlot {
		bgColor = color.RGBA{90, 90, 120, 255}
	}
	if slotIdx == ui.PlayerInventory.Selected {
		bgColor = color.RGBA{100, 200, 100, 255} // Highlight selected hotbar slot
	}
	ebitenutil.DrawRect(screen, x, y, ui.SlotSize, ui.SlotSize, bgColor)

	// Draw border
	borderColor := color.RGBA{80, 80, 100, 255}
	if slotIdx == ui.PlayerInventory.Selected {
		borderColor = color.RGBA{150, 255, 150, 255}
	}
	ebitenutil.DrawRect(screen, x, y, ui.SlotSize, 2, borderColor)
	ebitenutil.DrawRect(screen, x, y+ui.SlotSize-2, ui.SlotSize, 2, borderColor)
	ebitenutil.DrawRect(screen, x, y, 2, ui.SlotSize, borderColor)
	ebitenutil.DrawRect(screen, x+ui.SlotSize-2, y, 2, ui.SlotSize, borderColor)

	// Draw item if present
	if slotIdx >= 0 && slotIdx < len(ui.PlayerInventory.Slots) {
		item := &ui.PlayerInventory.Slots[slotIdx]
		if item.Type != items.NONE {
			ui.drawItem(screen, item, x+4, y+4)
		}
	}
}

// drawItem draws an item at the specified position
func (ui *ChestUI) drawItem(screen *ebiten.Image, item *items.Item, x, y float64) {
	if item == nil {
		return
	}

	// Draw item background based on item type color
	itemProps := items.ItemDefinitions[item.Type]
	if itemProps != nil {
		ebitenutil.DrawRect(screen, x, y, ui.SlotSize-8, ui.SlotSize-8, itemProps.IconColor)
	}

	// Draw quantity
	if item.Quantity > 1 {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d", item.Quantity), int(x)+25, int(y)+25)
	}
}

// IsOpen returns whether the UI is open
func (ui *ChestUI) IsOpen() bool {
	return ui.Open
}

// GetCurrentChestPosition returns the position of the currently open chest
func (ui *ChestUI) GetCurrentChestPosition() (float64, float64) {
	return ui.CurrentChestX, ui.CurrentChestY
}
