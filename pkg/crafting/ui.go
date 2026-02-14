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

	// Recipe display
	visibleRecipes []*Recipe

	// Quantity selector
	craftQuantity int

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
		craftQuantity:     1,
		animationProgress: 0.0,
	}
}

// Toggle opens or closes the crafting UI
func (ui *CraftingUI) Toggle() {
	ui.Open = !ui.Open
	if ui.Open {
		ui.visibleRecipes = ui.craftingSystem.GetAvailableRecipes(ui.inventory)
		if len(ui.visibleRecipes) > 0 {
			ui.SelectedRecipe = 0
		} else {
			ui.SelectedRecipe = -1
		}
	}
}

// Update handles input updates for the crafting UI
func (ui *CraftingUI) Update() error {
	if !ui.Open {
		return nil
	}

	// Refresh available recipes
	ui.visibleRecipes = ui.craftingSystem.GetAvailableRecipes(ui.inventory)

	// Keyboard navigation
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		ui.Toggle()
		return nil
	}

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
				if err := ui.craftingSystem.Craft(recipe.ID, ui.inventory); err != nil {
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
				if err := ui.craftingSystem.Craft(recipe.ID, ui.inventory); err != nil {
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

	// Draw semi-transparent background
	bgColor := color.RGBA{30, 30, 40, 230}
	ebitenutil.DrawRect(screen, 0, 0, 1280, 720, bgColor)

	// Draw title
	ui.drawText(screen, "CRAFTING MENU", 50, 40, color.RGBA{255, 255, 255, 255})
	ui.drawText(screen, "Press ESC to close", 50, 70, color.RGBA{180, 180, 180, 255})

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
		ui.drawText(screen, recipe.Name, recipeListX+10, y+10, color.RGBA{255, 255, 255, 255})

		// Recipe description
		ui.drawText(screen, recipe.Description, recipeListX+10, y+35, color.RGBA{180, 180, 180, 255})
	}

	// No recipes available message
	if len(ui.visibleRecipes) == 0 {
		ui.drawText(screen, "No recipes available!", recipeListX, recipeListY, color.RGBA{200, 100, 100, 255})
	}
}

// drawRecipeDetails draws the details of the selected recipe
func (ui *CraftingUI) drawRecipeDetails(screen *ebiten.Image) {
	recipe := ui.visibleRecipes[ui.SelectedRecipe]
	detailsX := 500
	detailsY := 100

	// Recipe name
	ui.drawText(screen, "Recipe: "+recipe.Name, detailsX, detailsY, color.RGBA{255, 255, 255, 255})

	// Inputs section
	ui.drawText(screen, "Required Materials:", detailsX, detailsY+50, color.RGBA{200, 200, 255, 255})

	inputY := detailsY + 80
	for _, input := range recipe.Inputs {
		itemName := items.ItemNameByID(input.ItemType)
		itemColor := items.ItemColorByID(input.ItemType)

		// Draw item color indicator
		ebitenutil.DrawRect(screen, float64(detailsX), float64(inputY), 20, 20, itemColor)

		// Draw item name and quantity
		have := ui.inventory.HasItem(input.ItemType, input.Quantity)
		textColor := color.RGBA{100, 255, 100, 255} // Green if have enough
		if !have {
			textColor = color.RGBA{255, 100, 100, 255} // Red if missing
		}

		ui.drawText(screen, fmt.Sprintf("%s x%d", itemName, input.Quantity), detailsX+30, inputY+2, textColor)
		inputY += 30
	}

	// Outputs section
	outputY := inputY + 20
	ui.drawText(screen, "Results:", detailsX, outputY, color.RGBA{200, 200, 255, 255})
	outputY += 30

	for _, output := range recipe.Outputs {
		itemName := items.ItemNameByID(output.ItemType)
		itemColor := items.ItemColorByID(output.ItemType)

		// Draw item color indicator
		ebitenutil.DrawRect(screen, float64(detailsX), float64(outputY), 20, 20, itemColor)

		// Draw item name and quantity
		ui.drawText(screen, fmt.Sprintf("%s x%d", itemName, output.Quantity), detailsX+30, outputY+2, color.RGBA{255, 255, 255, 255})
		outputY += 30
	}

	// Quantity selector
	quantityY := outputY + 30
	ui.drawText(screen, "Quantity:", detailsX, quantityY, color.RGBA{200, 200, 200, 255})
	ui.drawText(screen, fmt.Sprintf("%d", ui.craftQuantity), detailsX+200, quantityY, color.RGBA{255, 255, 255, 255})

	// Craft button
	craftButtonX := detailsX
	craftButtonY := quantityY + 50
	craftButtonWidth := 180
	craftButtonHeight := 60

	// Check if can craft
	canCraft := ui.craftingSystem.CanCraft(recipe, ui.inventory)
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
			ui.drawText(screen, buttonText, textX+dx, textY+dy, color.RGBA{255, 255, 255, 255})
		}
	}

	// Instructions
	ui.drawText(screen, "Use arrow keys to navigate, +/- or left/right for quantity, ENTER to craft", detailsX, 650, color.RGBA{150, 150, 150, 255})
}

// drawText draws text on the screen with larger, more readable text
func (ui *CraftingUI) drawText(screen *ebiten.Image, text string, x, y int, col color.RGBA) {
	// Draw text multiple times with slight offsets for thicker, more readable text
	for dx := 0; dx < 2; dx++ {
		for dy := 0; dy < 2; dy++ {
			ebitenutil.DebugPrintAt(screen, text, x+dx, y+dy)
		}
	}
}
