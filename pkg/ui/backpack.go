// Package ui implements the Minecraft-like backpack inventory interface.
package ui

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"tesselbox/pkg/equipment"
	"tesselbox/pkg/health"
	"tesselbox/pkg/items"
)

// BackpackUI represents the Minecraft-like inventory interface
type BackpackUI struct {
	Open         bool
	ScreenWidth  int
	ScreenHeight int

	// Inventory references
	Inventory    *items.Inventory
	Equipment    *equipment.EquipmentSet
	HealthSystem *health.LocationalHealthSystem

	// UI Layout
	MainInventoryX float64
	MainInventoryY float64
	HotbarX        float64
	HotbarY        float64
	ArmorX         float64
	ArmorY         float64
	EquipmentX     float64
	EquipmentY     float64
	PlayerModelX   float64
	PlayerModelY   float64
	CraftingX      float64
	CraftingY      float64

	// Slot sizes
	SlotSize    float64
	SlotSpacing float64

	// Drag and drop
	DraggedItem     *items.Item
	DraggedQuantity int
	DragX, DragY    float64

	// Hover state
	HoveredSlot int
	HoveredType SlotType

	// Selected slot
	SelectedSlot int
}

// SlotType represents different slot areas
type SlotType int

const (
	SlotTypeInventory SlotType = iota
	SlotTypeHotbar
	SlotTypeArmor
	SlotTypeEquipment
	SlotTypeCrafting
	SlotTypeResult
)

// NewBackpackUI creates a new backpack UI
func NewBackpackUI(screenWidth, screenHeight int, inv *items.Inventory, eq *equipment.EquipmentSet, hs *health.LocationalHealthSystem) *BackpackUI {
	ui := &BackpackUI{
		ScreenWidth:  screenWidth,
		ScreenHeight: screenHeight,
		Inventory:    inv,
		Equipment:    eq,
		HealthSystem: hs,
		SlotSize:     48,
		SlotSpacing:  4,
		SelectedSlot: 0,
	}

	ui.calculateLayout()
	return ui
}

// calculateLayout positions UI elements
func (ui *BackpackUI) calculateLayout() {
	centerX := float64(ui.ScreenWidth) / 2
	centerY := float64(ui.ScreenHeight) / 2

	// Main inventory grid (9x3 = 27 slots below hotbar)
	ui.MainInventoryX = centerX - (9*ui.SlotSize+8*ui.SlotSpacing)/2
	ui.MainInventoryY = centerY + 20

	// Hotbar at bottom
	ui.HotbarX = centerX - (9*ui.SlotSize+8*ui.SlotSpacing)/2
	ui.HotbarY = ui.MainInventoryY + 3*ui.SlotSize + 3*ui.SlotSpacing + 20

	// Armor slots on left side
	ui.ArmorX = centerX - 250
	ui.ArmorY = centerY - 100

	// Equipment slots (wings, etc.)
	ui.EquipmentX = ui.ArmorX
	ui.EquipmentY = ui.ArmorY + 5*ui.SlotSize + 20

	// Player model in center-left
	ui.PlayerModelX = centerX - 80
	ui.PlayerModelY = centerY - 80

	// Crafting area on right
	ui.CraftingX = centerX + 150
	ui.CraftingY = centerY - 50
}

// Update handles input and updates UI state
func (ui *BackpackUI) Update() error {
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

	// Handle right clicks (split stack, etc.)
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
func (ui *BackpackUI) updateHover(mx, my float64) {
	ui.HoveredSlot = -1
	ui.HoveredType = SlotTypeInventory

	// Check hotbar slots
	for i := 0; i < 9; i++ {
		slotX := ui.HotbarX + float64(i)*(ui.SlotSize+ui.SlotSpacing)
		slotY := ui.HotbarY
		if mx >= slotX && mx <= slotX+ui.SlotSize && my >= slotY && my <= slotY+ui.SlotSize {
			ui.HoveredSlot = i
			ui.HoveredType = SlotTypeHotbar
			return
		}
	}

	// Check main inventory slots (9x3 grid, slots 9-35)
	for row := 0; row < 3; row++ {
		for col := 0; col < 9; col++ {
			slotIdx := 9 + row*9 + col
			slotX := ui.MainInventoryX + float64(col)*(ui.SlotSize+ui.SlotSpacing)
			slotY := ui.MainInventoryY + float64(row)*(ui.SlotSize+ui.SlotSpacing)
			if mx >= slotX && mx <= slotX+ui.SlotSize && my >= slotY && my <= slotY+ui.SlotSize {
				ui.HoveredSlot = slotIdx
				ui.HoveredType = SlotTypeInventory
				return
			}
		}
	}

	// Check armor slots
	armorSlots := []equipment.EquipmentSlot{
		equipment.SlotHelmet,
		equipment.SlotChestplate,
		equipment.SlotLeggings,
		equipment.SlotBoots,
	}
	for i, slot := range armorSlots {
		slotX := ui.ArmorX
		slotY := ui.ArmorY + float64(i)*(ui.SlotSize+ui.SlotSpacing)
		if mx >= slotX && mx <= slotX+ui.SlotSize && my >= slotY && my <= slotY+ui.SlotSize {
			ui.HoveredSlot = int(slot)
			ui.HoveredType = SlotTypeArmor
			return
		}
	}
}

// handleClick handles left mouse button clicks
func (ui *BackpackUI) handleClick(mx, my float64) {
	if ui.HoveredSlot < 0 {
		return
	}

	// Handle different slot types
	switch ui.HoveredType {
	case SlotTypeHotbar, SlotTypeInventory:
		// Get item from inventory
		if ui.HoveredSlot >= 0 && ui.HoveredSlot < len(ui.Inventory.Slots) {
			item := &ui.Inventory.Slots[ui.HoveredSlot]
			if item.Type != items.NONE {
				// Start dragging
				ui.DraggedItem = item
				ui.DraggedQuantity = item.Quantity
				// Clear the slot
				item.Type = items.NONE
				item.Quantity = 0
				item.Durability = -1
			}
		}
	case SlotTypeArmor:
		// Handle armor slot clicks
		slot := equipment.EquipmentSlot(ui.HoveredSlot)
		item := ui.Equipment.GetItem(slot)
		if item != nil {
			// Unequip
			ui.Equipment.UnequipItem(slot)
			// Add to inventory or start dragging
		}
	}
}

// handleRightClick handles right mouse button clicks
func (ui *BackpackUI) handleRightClick(mx, my float64) {
	if ui.HoveredSlot < 0 {
		return
	}

	// Get item and split stack if applicable
	if ui.HoveredSlot >= 0 && ui.HoveredSlot < len(ui.Inventory.Slots) {
		item := &ui.Inventory.Slots[ui.HoveredSlot]
		if item.Quantity > 1 {
			half := item.Quantity / 2
			ui.DraggedItem = item
			ui.DraggedQuantity = half
			ui.DraggedItem.Quantity = half
			item.Quantity -= half
		}
	}
}

// handleDrop handles dropping a dragged item
func (ui *BackpackUI) handleDrop(mx, my float64) {
	if ui.DraggedItem == nil {
		return
	}

	// Try to place in hovered slot
	if ui.HoveredSlot >= 0 {
		switch ui.HoveredType {
		case SlotTypeHotbar, SlotTypeInventory:
			// Place in inventory
			if ui.HoveredSlot >= 0 && ui.HoveredSlot < len(ui.Inventory.Slots) {
				ui.DraggedItem.Quantity = ui.DraggedQuantity
				ui.Inventory.Slots[ui.HoveredSlot] = *ui.DraggedItem
			}
		case SlotTypeArmor:
			// Try to equip armor
			// Check if item is armor type
			// ...
		}
	}

	// Clear drag state
	ui.DraggedItem = nil
	ui.DraggedQuantity = 0
}

// Draw renders the backpack UI
func (ui *BackpackUI) Draw(screen *ebiten.Image) {
	if !ui.Open {
		return
	}

	// Draw semi-transparent background
	bgColor := color.RGBA{20, 20, 30, 230}
	ebitenutil.DrawRect(screen, 0, 0, float64(ui.ScreenWidth), float64(ui.ScreenHeight), bgColor)

	// Draw title
	ebitenutil.DebugPrintAt(screen, "INVENTORY", ui.ScreenWidth/2-40, 20)

	// Draw player model/health visualization
	ui.drawPlayerModel(screen)

	// Draw armor slots
	ui.drawArmorSlots(screen)

	// Draw equipment slots
	ui.drawEquipmentSlots(screen)

	// Draw main inventory
	ui.drawMainInventory(screen)

	// Draw hotbar
	ui.drawHotbar(screen)

	// Draw crafting area
	ui.drawCraftingArea(screen)

	// Draw dragged item
	if ui.DraggedItem != nil {
		ui.drawItem(screen, ui.DraggedItem, ui.DragX, ui.DragY)
	}

	// Draw hover tooltip
	if ui.HoveredSlot >= 0 {
		mx, my := ebiten.CursorPosition()
		ui.drawTooltip(screen, float64(mx), float64(my))
	}
}

// drawPlayerModel draws the player character with body part health visualization
func (ui *BackpackUI) drawPlayerModel(screen *ebiten.Image) {
	// Draw player silhouette as a simple rectangle representation
	playerWidth := 80.0
	playerHeight := 120.0

	// Get health for each body part and color accordingly
	bodyParts := []health.BodyPart{
		health.PartHead,
		health.PartTorso,
		health.PartLeftArm,
		health.PartRightArm,
		health.PartLeftLeg,
		health.PartRightLeg,
	}

	for _, part := range bodyParts {
		x, y, w, h := health.BodyPartPosition(part)
		healthPct := ui.HealthSystem.GetPartHealthPercentage(part)
		col := health.GetBodyPartColor(healthPct)

		// Scale to player model size
		scaleX := playerWidth / 140.0
		scaleY := playerHeight / 200.0

		screenX := ui.PlayerModelX + x*scaleX
		screenY := ui.PlayerModelY + y*scaleY
		screenW := w * scaleX
		screenH := h * scaleY

		// Draw body part
		ebitenutil.DrawRect(screen, screenX, screenY, screenW, screenH, col)

		// Draw border
		borderColor := color.RGBA{100, 100, 100, 255}
		ebitenutil.DrawRect(screen, screenX, screenY, screenW, 2, borderColor)
		ebitenutil.DrawRect(screen, screenX, screenY+screenH-2, screenW, 2, borderColor)
		ebitenutil.DrawRect(screen, screenX, screenY, 2, screenH, borderColor)
		ebitenutil.DrawRect(screen, screenX+screenW-2, screenY, 2, screenH, borderColor)
	}
}

// drawArmorSlots draws the armor equipment slots
func (ui *BackpackUI) drawArmorSlots(screen *ebiten.Image) {
	slotNames := []string{"Helmet", "Chest", "Legs", "Boots"}
	armorSlots := []equipment.EquipmentSlot{
		equipment.SlotHelmet,
		equipment.SlotChestplate,
		equipment.SlotLeggings,
		equipment.SlotBoots,
	}

	for i, slot := range armorSlots {
		slotX := ui.ArmorX
		slotY := ui.ArmorY + float64(i)*(ui.SlotSize+ui.SlotSpacing)

		// Draw slot background
		bgColor := color.RGBA{60, 60, 70, 255}
		if ui.HoveredSlot == int(slot) && ui.HoveredType == SlotTypeArmor {
			bgColor = color.RGBA{80, 80, 100, 255}
		}
		ebitenutil.DrawRect(screen, slotX, slotY, ui.SlotSize, ui.SlotSize, bgColor)

		// Draw border
		borderColor := color.RGBA{100, 100, 120, 255}
		ebitenutil.DrawRect(screen, slotX, slotY, ui.SlotSize, 2, borderColor)
		ebitenutil.DrawRect(screen, slotX, slotY+ui.SlotSize-2, ui.SlotSize, 2, borderColor)
		ebitenutil.DrawRect(screen, slotX, slotY, 2, ui.SlotSize, borderColor)
		ebitenutil.DrawRect(screen, slotX+ui.SlotSize-2, slotY, 2, ui.SlotSize, borderColor)

		// Draw slot label
		ebitenutil.DebugPrintAt(screen, slotNames[i], int(slotX)+5, int(slotY-15))

		// Draw equipped item if any
		item := ui.Equipment.GetItem(slot)
		if item != nil {
			ui.drawEquipmentItem(screen, item, slotX+4, slotY+4, ui.SlotSize-8)
		}
	}
}

// drawEquipmentSlots draws special equipment slots (wings, etc.)
func (ui *BackpackUI) drawEquipmentSlots(screen *ebiten.Image) {
	// Wings slot
	slotX := ui.EquipmentX
	slotY := ui.EquipmentY

	bgColor := color.RGBA{60, 60, 70, 255}
	if ui.HoveredSlot == int(equipment.SlotWings) {
		bgColor = color.RGBA{80, 80, 100, 255}
	}
	ebitenutil.DrawRect(screen, slotX, slotY, ui.SlotSize, ui.SlotSize, bgColor)

	borderColor := color.RGBA{100, 100, 120, 255}
	ebitenutil.DrawRect(screen, slotX, slotY, ui.SlotSize, 2, borderColor)
	ebitenutil.DrawRect(screen, slotX, slotY+ui.SlotSize-2, ui.SlotSize, 2, borderColor)
	ebitenutil.DrawRect(screen, slotX, slotY, 2, ui.SlotSize, borderColor)
	ebitenutil.DrawRect(screen, slotX+ui.SlotSize-2, slotY, 2, ui.SlotSize, borderColor)

	ebitenutil.DebugPrintAt(screen, "Wings", int(slotX)+5, int(slotY-15))

	wings := ui.Equipment.GetItem(equipment.SlotWings)
	if wings != nil {
		ui.drawEquipmentItem(screen, wings, slotX+4, slotY+4, ui.SlotSize-8)
	}
}

// drawMainInventory draws the main inventory grid
func (ui *BackpackUI) drawMainInventory(screen *ebiten.Image) {
	for row := 0; row < 3; row++ {
		for col := 0; col < 9; col++ {
			slotIdx := 9 + row*9 + col
			slotX := ui.MainInventoryX + float64(col)*(ui.SlotSize+ui.SlotSpacing)
			slotY := ui.MainInventoryY + float64(row)*(ui.SlotSize+ui.SlotSpacing)

			// Draw slot background
			bgColor := color.RGBA{50, 50, 60, 255}
			if ui.HoveredSlot == slotIdx && ui.HoveredType == SlotTypeInventory {
				bgColor = color.RGBA{70, 70, 90, 255}
			}
			ebitenutil.DrawRect(screen, slotX, slotY, ui.SlotSize, ui.SlotSize, bgColor)

			// Draw border
			borderColor := color.RGBA{80, 80, 100, 255}
			ebitenutil.DrawRect(screen, slotX, slotY, ui.SlotSize, 2, borderColor)
			ebitenutil.DrawRect(screen, slotX, slotY+ui.SlotSize-2, ui.SlotSize, 2, borderColor)
			ebitenutil.DrawRect(screen, slotX, slotY, 2, ui.SlotSize, borderColor)
			ebitenutil.DrawRect(screen, slotX+ui.SlotSize-2, slotY, 2, ui.SlotSize, borderColor)

			// Draw item if present
			if slotIdx >= 0 && slotIdx < len(ui.Inventory.Slots) {
				item := &ui.Inventory.Slots[slotIdx]
				if item.Type != items.NONE {
					ui.drawItem(screen, item, slotX+4, slotY+4)
				}
			}
		}
	}
}

// drawHotbar draws the hotbar
func (ui *BackpackUI) drawHotbar(screen *ebiten.Image) {
	for i := 0; i < 9; i++ {
		slotX := ui.HotbarX + float64(i)*(ui.SlotSize+ui.SlotSpacing)
		slotY := ui.HotbarY

		// Draw slot background
		bgColor := color.RGBA{50, 50, 60, 255}
		if i == ui.SelectedSlot {
			bgColor = color.RGBA{100, 200, 100, 255} // Highlight selected
		} else if ui.HoveredSlot == i && ui.HoveredType == SlotTypeHotbar {
			bgColor = color.RGBA{70, 70, 90, 255}
		}
		ebitenutil.DrawRect(screen, slotX, slotY, ui.SlotSize, ui.SlotSize, bgColor)

		// Draw border
		borderColor := color.RGBA{80, 80, 100, 255}
		if i == ui.SelectedSlot {
			borderColor = color.RGBA{150, 255, 150, 255}
		}
		ebitenutil.DrawRect(screen, slotX, slotY, ui.SlotSize, 2, borderColor)
		ebitenutil.DrawRect(screen, slotX, slotY+ui.SlotSize-2, ui.SlotSize, 2, borderColor)
		ebitenutil.DrawRect(screen, slotX, slotY, 2, ui.SlotSize, borderColor)
		ebitenutil.DrawRect(screen, slotX+ui.SlotSize-2, slotY, 2, ui.SlotSize, borderColor)

		// Draw slot number
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d", i+1), int(slotX)+20, int(slotY-15))

		// Draw item if present
		if i >= 0 && i < len(ui.Inventory.Slots) {
			item := &ui.Inventory.Slots[i]
			if item.Type != items.NONE {
				ui.drawItem(screen, item, slotX+4, slotY+4)
			}
		}
	}
}

// drawCraftingArea draws the crafting interface
func (ui *BackpackUI) drawCraftingArea(screen *ebiten.Image) {
	// Draw crafting grid label
	ebitenutil.DebugPrintAt(screen, "Crafting", int(ui.CraftingX), int(ui.CraftingY-25))

	// Draw 2x2 crafting grid
	for row := 0; row < 2; row++ {
		for col := 0; col < 2; col++ {
			slotX := ui.CraftingX + float64(col)*(ui.SlotSize+ui.SlotSpacing)
			slotY := ui.CraftingY + float64(row)*(ui.SlotSize+ui.SlotSpacing)

			bgColor := color.RGBA{50, 50, 60, 255}
			ebitenutil.DrawRect(screen, slotX, slotY, ui.SlotSize, ui.SlotSize, bgColor)

			borderColor := color.RGBA{80, 80, 100, 255}
			ebitenutil.DrawRect(screen, slotX, slotY, ui.SlotSize, 2, borderColor)
			ebitenutil.DrawRect(screen, slotX, slotY+ui.SlotSize-2, ui.SlotSize, 2, borderColor)
			ebitenutil.DrawRect(screen, slotX, slotY, 2, ui.SlotSize, borderColor)
			ebitenutil.DrawRect(screen, slotX+ui.SlotSize-2, slotY, 2, ui.SlotSize, borderColor)
		}
	}

	// Draw result slot
	resultX := ui.CraftingX + 2*ui.SlotSize + 30
	resultY := ui.CraftingY + ui.SlotSize/2 - ui.SlotSize/2

	bgColor := color.RGBA{60, 60, 80, 255}
	ebitenutil.DrawRect(screen, resultX, resultY, ui.SlotSize, ui.SlotSize, bgColor)

	borderColor := color.RGBA{120, 120, 150, 255}
	ebitenutil.DrawRect(screen, resultX, resultY, ui.SlotSize, 2, borderColor)
	ebitenutil.DrawRect(screen, resultX, resultY+ui.SlotSize-2, ui.SlotSize, 2, borderColor)
	ebitenutil.DrawRect(screen, resultX, resultY, 2, ui.SlotSize, borderColor)
	ebitenutil.DrawRect(screen, resultX+ui.SlotSize-2, resultY, 2, ui.SlotSize, borderColor)
}

// drawItem draws an item at the specified position
func (ui *BackpackUI) drawItem(screen *ebiten.Image, item *items.Item, x, y float64) {
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

// drawEquipmentItem draws an equipment item
func (ui *BackpackUI) drawEquipmentItem(screen *ebiten.Image, item *equipment.EquipmentItem, x, y, size float64) {
	if item == nil {
		return
	}

	// Draw item background
	ebitenutil.DrawRect(screen, x, y, size, size, item.IconColor)

	// Draw durability bar if applicable
	if item.Durability > 0 && item.MaxDurability > 0 {
		durabilityPct := float64(item.Durability) / float64(item.MaxDurability)
		barWidth := size * durabilityPct
		var barColor color.RGBA
		if durabilityPct > 0.5 {
			barColor = color.RGBA{0, 255, 0, 255}
		} else if durabilityPct > 0.25 {
			barColor = color.RGBA{255, 255, 0, 255}
		} else {
			barColor = color.RGBA{255, 0, 0, 255}
		}
		ebitenutil.DrawRect(screen, x, y+size-4, barWidth, 4, barColor)
	}
}

// drawTooltip draws a tooltip for hovered items
func (ui *BackpackUI) drawTooltip(screen *ebiten.Image, mx, my float64) {
	// Get hovered item info
	var itemName string
	var itemDesc string

	if ui.HoveredType == SlotTypeInventory || ui.HoveredType == SlotTypeHotbar {
		if ui.HoveredSlot >= 0 && ui.HoveredSlot < len(ui.Inventory.Slots) {
			item := &ui.Inventory.Slots[ui.HoveredSlot]
			if item.Type != items.NONE {
				itemProps := items.ItemDefinitions[item.Type]
				if itemProps != nil {
					itemName = itemProps.Name
					itemDesc = itemProps.Description
				}
			}
		}
	}

	if itemName != "" {
		// Draw tooltip background
		tooltipX := mx + 15
		tooltipY := my + 15
		tooltipW := 150.0
		tooltipH := 50.0

		ebitenutil.DrawRect(screen, tooltipX, tooltipY, tooltipW, tooltipH, color.RGBA{30, 30, 40, 230})
		ebitenutil.DrawRect(screen, tooltipX, tooltipY, tooltipW, 2, color.RGBA{100, 100, 120, 255})
		ebitenutil.DrawRect(screen, tooltipX, tooltipY+tooltipH-2, tooltipW, 2, color.RGBA{100, 100, 120, 255})
		ebitenutil.DrawRect(screen, tooltipX, tooltipY, 2, tooltipH, color.RGBA{100, 100, 120, 255})
		ebitenutil.DrawRect(screen, tooltipX+tooltipW-2, tooltipY, 2, tooltipH, color.RGBA{100, 100, 120, 255})

		// Draw text
		ebitenutil.DebugPrintAt(screen, itemName, int(tooltipX)+5, int(tooltipY)+5)
		ebitenutil.DebugPrintAt(screen, itemDesc, int(tooltipX)+5, int(tooltipY)+25)
	}
}

// Toggle opens/closes the backpack UI
func (ui *BackpackUI) Toggle() {
	ui.Open = !ui.Open
	if !ui.Open {
		// Cancel any drag operation when closing
		ui.DraggedItem = nil
	}
}

// SetSelectedSlot sets the selected hotbar slot
func (ui *BackpackUI) SetSelectedSlot(slot int) {
	if slot >= 0 && slot < 9 {
		ui.SelectedSlot = slot
	}
}

// IsOpen returns whether the UI is open
func (ui *BackpackUI) IsOpen() bool {
	return ui.Open
}
