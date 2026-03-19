package crafting

import (
	"fmt"

	"image/color"
	"tesselbox/pkg/items"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// CraftingUI represents the crafting interface
type CraftingUI struct {
	craftingSystem *CraftingSystem
	inventory      *items.Inventory

	// UI state
	Open           bool
	SelectedRecipe int
	CreativeMode   bool

	// Recipe display
	visibleRecipes []*Recipe

	// Creative mode items
	allItems     []items.ItemType
	scrollOffset int

	// Quantity selector
	craftQuantity int

	// Current crafting station
	currentStation CraftingStation

	// Animation
	animationProgress float64
}

// NewCraftingUI creates a new crafting UI
func NewCraftingUI(craftingSystem *CraftingSystem, inventory *items.Inventory) *CraftingUI {
	return &CraftingUI{
		craftingSystem:    craftingSystem,
		inventory:         inventory,
		Open:              false,
		SelectedRecipe:    -1,
		CreativeMode:      true, // Enable creative mode
		scrollOffset:      0,
		craftQuantity:     1,
		currentStation:    STATION_NONE,
		animationProgress: 0.0,
	}
}

// Toggle opens or closes the crafting UI
func (ui *CraftingUI) Toggle() {
	ui.Open = !ui.Open
	if ui.Open {
		if ui.CreativeMode {
			ui.allItems = []items.ItemType{}
			for itemType := range items.ItemDefinitions {
				ui.allItems = append(ui.allItems, itemType)
			}
		} else {
			ui.visibleRecipes = ui.craftingSystem.GetAvailableRecipes(ui.inventory, ui.currentStation)
			if len(ui.visibleRecipes) > 0 {
				ui.SelectedRecipe = 0
			} else {
				ui.SelectedRecipe = -1
			}
		}
	}
}

// SetStation sets the current crafting station and refreshes available recipes
func (ui *CraftingUI) SetStation(station CraftingStation) {
	ui.currentStation = station
	if ui.Open {
		ui.visibleRecipes = ui.craftingSystem.GetAvailableRecipes(ui.inventory, ui.currentStation)
		if len(ui.visibleRecipes) == 0 {
			ui.SelectedRecipe = -1
		} else if ui.SelectedRecipe >= len(ui.visibleRecipes) {
			ui.SelectedRecipe = len(ui.visibleRecipes) - 1
		}
	}
}

// GetCurrentStation returns the current crafting station
func (ui *CraftingUI) GetCurrentStation() CraftingStation {
	return ui.currentStation
}

// Update handles input updates for the crafting UI
func (ui *CraftingUI) Update() error {
	if !ui.Open {
		return nil
	}

	// Keyboard navigation
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		ui.Toggle()
		return nil
	}

	if ui.CreativeMode {
		// Handle scroll
		_, scrollY := ebiten.Wheel()
		if scrollY > 0 {
			ui.scrollOffset = max(0, ui.scrollOffset-1)
		} else if scrollY < 0 {
			ui.scrollOffset++
		}

		// Handle mouse click to select item
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			mx, my := ebiten.CursorPosition()
			x := 50
			y := 100
			itemSize := 50
			perRow := 10

			if mx >= x && my >= y {
				col := (mx - x) / itemSize
				row := (my - y) / itemSize
				localIndex := row*perRow + col
				globalIndex := ui.scrollOffset*perRow + localIndex

				if globalIndex >= 0 && globalIndex < len(ui.allItems) {
					itemType := ui.allItems[globalIndex]
					// Add to selected inventory slot (assume slot 0 for now, or find empty slot)
					for i := 0; i < len(ui.inventory.Slots); i++ {
						if ui.inventory.Slots[i].Type == items.NONE {
							ui.inventory.Slots[i] = items.Item{Type: itemType, Quantity: 1, Durability: -1}
							break
						}
					}
				}
			}
		}
	} else {
		// Normal crafting mode
		// Refresh available recipes
		ui.visibleRecipes = ui.craftingSystem.GetAvailableRecipes(ui.inventory, ui.currentStation)

		if ebiten.IsKeyPressed(ebiten.KeyUp) {
			ui.SelectedRecipe--
			if ui.SelectedRecipe < 0 {
				ui.SelectedRecipe = len(ui.visibleRecipes) - 1
			}
		}

		if ebiten.IsKeyPressed(ebiten.KeyDown) {
			ui.SelectedRecipe++
			if ui.SelectedRecipe >= len(ui.visibleRecipes) {
				ui.SelectedRecipe = 0
			}
		}

		// Quantity selection
		if ebiten.IsKeyPressed(ebiten.KeyLeft) {
			ui.craftQuantity--
			if ui.craftQuantity < 1 {
				ui.craftQuantity = 1
			}
		}

		if ebiten.IsKeyPressed(ebiten.KeyRight) {
			ui.craftQuantity++
		}

		// Craft on Enter or Space
		if ebiten.IsKeyPressed(ebiten.KeyEnter) || ebiten.IsKeyPressed(ebiten.KeySpace) {
			if ui.SelectedRecipe >= 0 && ui.SelectedRecipe < len(ui.visibleRecipes) {
				recipe := ui.visibleRecipes[ui.SelectedRecipe]
				for i := 0; i < ui.craftQuantity; i++ {
					if err := ui.craftingSystem.Craft(recipe.ID, ui.inventory, ui.currentStation); err != nil {
						break // Stop if crafting fails (e.g., inventory full)
					}
				}
			}
		}

		// Mouse click handling (simplified)
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			mx, my := ebiten.CursorPosition()
			ui.handleClick(mx, my)
		}
	}

	return nil
}

// handleClick handles mouse clicks on the crafting UI
func (ui *CraftingUI) handleClick(mx, my int) {
	// Recipe list area
	recipeListX := 50
	recipeListY := 100
	recipeHeight := 60

	for i := range ui.visibleRecipes {
		recipeY := recipeListY + i*recipeHeight

		// Check if click is on this recipe
		if mx >= recipeListX && mx <= recipeListX+400 &&
			my >= recipeY && my <= recipeY+recipeHeight-10 {
			ui.SelectedRecipe = i
			ui.craftQuantity = 1
			return
		}
	}

	// Craft button area
	craftButtonX := 500
	craftButtonY := 500
	craftButtonWidth := 180
	craftButtonHeight := 60
	if ui.SelectedRecipe >= 0 && ui.SelectedRecipe < len(ui.visibleRecipes) {
		if mx >= craftButtonX && mx <= craftButtonX+craftButtonWidth &&
			my >= craftButtonY && my <= craftButtonY+craftButtonHeight {
			recipe := ui.visibleRecipes[ui.SelectedRecipe]
			for i := 0; i < ui.craftQuantity; i++ {
				if err := ui.craftingSystem.Craft(recipe.ID, ui.inventory, ui.currentStation); err != nil {
					break
				}
			}
		}
	}
}

// Draw renders the crafting UI
func (ui *CraftingUI) Draw(screen *ebiten.Image) {
	if !ui.Open {
		return
	}

	if ui.CreativeMode {
		// Creative mode: draw item grid
		ui.drawCreativeGrid(screen)
	} else {
		// Normal crafting mode
		ui.drawCraftingUI(screen)
	}
}

// drawCreativeGrid draws the creative item selection grid
func (ui *CraftingUI) drawCreativeGrid(screen *ebiten.Image) {
	// Draw semi-transparent background
	bgColor := color.RGBA{30, 30, 40, 230}
	ebitenutil.DrawRect(screen, 0, 0, 1280, 720, bgColor)

	// Draw title
	ui.drawText(screen, "CREATIVE INVENTORY", 50, 40)
	ui.drawText(screen, "Press ESC to close, mouse wheel to scroll, click to select item", 50, 70)

	// Draw grid of items
	x := 50
	y := 100
	itemSize := 50
	perRow := 10
	visibleRows := 10

	for i := ui.scrollOffset * perRow; i < len(ui.allItems) && (i-ui.scrollOffset*perRow) < perRow*visibleRows; i++ {
		localIndex := i - ui.scrollOffset*perRow
		row := localIndex / perRow
		col := localIndex % perRow
		ix := x + col*itemSize
		iy := y + row*itemSize

		itemType := ui.allItems[i]
		itemColor := items.ItemColorByID(itemType)
		ebitenutil.DrawRect(screen, float64(ix), float64(iy), float64(itemSize-2), float64(itemSize-2), itemColor)
	}
}

// drawCraftingUI draws the normal crafting interface
func (ui *CraftingUI) drawCraftingUI(screen *ebiten.Image) {
	// Draw semi-transparent background
	bgColor := color.RGBA{30, 30, 40, 230}
	ebitenutil.DrawRect(screen, 0, 0, 1280, 720, bgColor)

	// Draw title
	ui.drawText(screen, "CRAFTING MENU", 50, 40)
	ui.drawText(screen, "Press ESC to close", 50, 70)

	// Draw recipe list
	ui.drawRecipeList(screen)

	// Draw selected recipe details
	if ui.SelectedRecipe >= 0 && ui.SelectedRecipe < len(ui.visibleRecipes) {
		ui.drawRecipeDetails(screen)
	}
}

// drawRecipeList draws the list of available recipes
func (ui *CraftingUI) drawRecipeList(screen *ebiten.Image) {
	recipeListX := 50
	recipeListY := 100
	recipeWidth := 400
	recipeHeight := 60

	for i, recipe := range ui.visibleRecipes {
		y := recipeListY + i*recipeHeight

		// Recipe background
		bgColor := color.RGBA{50, 50, 60, 255}
		if i == ui.SelectedRecipe {
			bgColor = color.RGBA{80, 80, 100, 255}
		}
		ebitenutil.DrawRect(screen, float64(recipeListX), float64(y), float64(recipeWidth), float64(recipeHeight-10), bgColor)

		// Recipe border
		borderColor := color.RGBA{100, 100, 120, 255}
		ebitenutil.DrawRect(screen, float64(recipeListX), float64(y), float64(recipeWidth), 2, borderColor)
		ebitenutil.DrawRect(screen, float64(recipeListX), float64(y+recipeHeight-12), float64(recipeWidth), 2, borderColor)

		// Recipe name
		ui.drawText(screen, recipe.Name, recipeListX+10, y+10)

		// Recipe description
		ui.drawText(screen, recipe.Description, recipeListX+10, y+35)
	}

	// No recipes available message
	if len(ui.visibleRecipes) == 0 {
		ui.drawText(screen, "No recipes available!", recipeListX, recipeListY)
	}
}

// drawRecipeDetails draws the details of the selected recipe
func (ui *CraftingUI) drawRecipeDetails(screen *ebiten.Image) {
	recipe := ui.visibleRecipes[ui.SelectedRecipe]
	detailsX := 500
	detailsY := 100

	// Recipe name
	ui.drawText(screen, "Recipe: "+recipe.Name, detailsX, detailsY)

	// Inputs section
	ui.drawText(screen, "Required Materials:", detailsX, detailsY+50)

	inputY := detailsY + 80
	for _, input := range recipe.Inputs {
		itemName := items.ItemNameByID(input.ItemType)
		itemColor := items.ItemColorByID(input.ItemType)

		// Draw item color indicator
		ebitenutil.DrawRect(screen, float64(detailsX), float64(inputY), 20, 20, itemColor)

		// Draw item name and quantity
		ui.drawText(screen, fmt.Sprintf("%s x%d", itemName, input.Quantity), detailsX+30, inputY+2)
		inputY += 30
	}

	// Outputs section
	outputY := inputY + 20
	ui.drawText(screen, "Results:", detailsX, outputY)
	outputY += 30

	for _, output := range recipe.Outputs {
		itemName := items.ItemNameByID(output.ItemType)
		itemColor := items.ItemColorByID(output.ItemType)

		// Draw item color indicator
		ebitenutil.DrawRect(screen, float64(detailsX), float64(outputY), 20, 20, itemColor)

		// Draw item name and quantity
		ui.drawText(screen, fmt.Sprintf("%s x%d", itemName, output.Quantity), detailsX+30, outputY+2)
		outputY += 30
	}

	// Quantity selector
	quantityY := outputY + 30
	ui.drawText(screen, "Quantity:", detailsX, quantityY)
	ui.drawText(screen, fmt.Sprintf("%d", ui.craftQuantity), detailsX+200, quantityY)

	// Craft button
	craftButtonX := detailsX
	craftButtonY := quantityY + 50
	craftButtonWidth := 180
	craftButtonHeight := 60

	// Check if can craft
	canCraft := ui.craftingSystem.CanCraft(recipe, ui.inventory, ui.currentStation)
	buttonColor := color.RGBA{100, 200, 100, 255}
	if !canCraft {
		buttonColor = color.RGBA{150, 150, 150, 255}
	}

	ebitenutil.DrawRect(screen, float64(craftButtonX), float64(craftButtonY), float64(craftButtonWidth), float64(craftButtonHeight), buttonColor)
	ebitenutil.DrawRect(screen, float64(craftButtonX), float64(craftButtonY), float64(craftButtonWidth), 3, color.RGBA{255, 255, 255, 255})

	buttonText := "CRAFT"
	if !canCraft {
		buttonText = "MISSING ITEMS"
	}

	// Center text in button with larger, thicker text
	textWidth := len(buttonText) * 8
	textX := craftButtonX + (craftButtonWidth-textWidth)/2
	textY := craftButtonY + 25

	// Draw text multiple times for thicker appearance
	for dx := 0; dx < 2; dx++ {
		for dy := 0; dy < 2; dy++ {
			ui.drawText(screen, buttonText, textX+dx, textY+dy)
		}
	}

	// Instructions
	ui.drawText(screen, "Use arrow keys to navigate, +/- or left/right for quantity, ENTER to craft", detailsX, 650)
}

// drawText draws text on the screen with larger, more readable text
func (ui *CraftingUI) drawText(screen *ebiten.Image, text string, x, y int) {
	// Draw text multiple times with slight offsets for thicker, more readable text
	for dx := 0; dx < 2; dx++ {
		for dy := 0; dy < 2; dy++ {
			ebitenutil.DebugPrintAt(screen, text, x+dx, y+dy)
		}
	}
}
