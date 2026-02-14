package crafting

import (
	"encoding/json"
	"fmt"
	"os"
	"tesselbox/pkg/items"
)

// RecipeInput represents an input item requirement
type RecipeInput struct {
	ItemType  items.ItemType `json:"item_type"`
	Quantity  int            `json:"quantity"`
}

// RecipeOutput represents an output item
type RecipeOutput struct {
	ItemType  items.ItemType `json:"item_type"`
	Quantity  int            `json:"quantity"`
}

// Recipe represents a crafting recipe
type Recipe struct {
	ID           string          `json:"id"`
	Name         string          `json:"name"`
	Description  string          `json:"description"`
	Inputs       []RecipeInput   `json:"inputs"`
	Outputs      []RecipeOutput  `json:"outputs"`
	CraftingTime float64         `json:"crafting_time"` // in seconds, 0 = instant
	RequiredTool items.ItemType  `json:"required_tool"` // NONE if no tool required
}

// CraftingSystem manages the crafting functionality
type CraftingSystem struct {
	recipes map[string]*Recipe
}

// NewCraftingSystem creates a new crafting system
func NewCraftingSystem() *CraftingSystem {
	return &CraftingSystem{
		recipes: make(map[string]*Recipe),
	}
}

// LoadRecipes loads recipes from a JSON file
func (cs *CraftingSystem) LoadRecipes(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read recipes file: %w", err)
	}
	
	var recipes []Recipe
	if err := json.Unmarshal(data, &recipes); err != nil {
		return fmt.Errorf("failed to parse recipes JSON: %w", err)
	}
	
	// Clear existing recipes
	cs.recipes = make(map[string]*Recipe)
	
	// Load recipes
	for i := range recipes {
		cs.recipes[recipes[i].ID] = &recipes[i]
	}
	
	return nil
}

// GetRecipe retrieves a recipe by ID
func (cs *CraftingSystem) GetRecipe(id string) (*Recipe, bool) {
	recipe, exists := cs.recipes[id]
	return recipe, exists
}

// GetAllRecipes returns all recipes
func (cs *CraftingSystem) GetAllRecipes() []*Recipe {
	recipes := make([]*Recipe, 0, len(cs.recipes))
	for _, recipe := range cs.recipes {
		recipes = append(recipes, recipe)
	}
	return recipes
}

// CanCraft checks if the player can craft a recipe
func (cs *CraftingSystem) CanCraft(recipe *Recipe, inventory *items.Inventory) bool {
	// Check if required tool is in selected slot
	if recipe.RequiredTool != items.NONE {
		selectedItem := inventory.GetSelectedItem()
		if selectedItem == nil || selectedItem.Type != recipe.RequiredTool {
			return false
		}
	}
	
	// Check if player has all required materials
	for _, input := range recipe.Inputs {
		if !inventory.HasItem(input.ItemType, input.Quantity) {
			return false
		}
	}
	
	return true
}

// GetMissingMaterials returns the materials needed to craft a recipe
func (cs *CraftingSystem) GetMissingMaterials(recipe *Recipe, inventory *items.Inventory) []RecipeInput {
	missing := []RecipeInput{}
	
	// Count available items
	available := make(map[items.ItemType]int)
	for _, slot := range inventory.Slots {
		if slot.Type != items.NONE {
			available[slot.Type] += slot.Quantity
		}
	}
	
	// Check each input
	for _, input := range recipe.Inputs {
		if available[input.ItemType] < input.Quantity {
			missing = append(missing, RecipeInput{
				ItemType: input.ItemType,
				Quantity: input.Quantity - available[input.ItemType],
			})
		}
	}
	
	return missing
}

// Craft attempts to craft a recipe
func (cs *CraftingSystem) Craft(recipeID string, inventory *items.Inventory) error {
	recipe, exists := cs.GetRecipe(recipeID)
	if !exists {
		return fmt.Errorf("recipe not found: %s", recipeID)
	}
	
	// Check if crafting is possible
	if !cs.CanCraft(recipe, inventory) {
		return fmt.Errorf("cannot craft %s: missing materials or tools", recipe.Name)
	}
	
	// Remove input materials
	for _, input := range recipe.Inputs {
		if !inventory.RemoveItemType(input.ItemType, input.Quantity) {
			return fmt.Errorf("failed to remove input materials")
		}
	}
	
	// Use tool durability if applicable
	if recipe.RequiredTool != items.NONE {
		inventory.UseItem()
	}
	
	// Add output items
	for _, output := range recipe.Outputs {
		if !inventory.AddItem(output.ItemType, output.Quantity) {
			// If we can't add the item, return it (inventory full)
			// In a real implementation, you might want to handle this differently
			return fmt.Errorf("inventory full")
		}
	}
	
	return nil
}

// GetAvailableRecipes returns recipes that can be crafted with current inventory
func (cs *CraftingSystem) GetAvailableRecipes(inventory *items.Inventory) []*Recipe {
	available := []*Recipe{}
	
	for _, recipe := range cs.recipes {
		if cs.CanCraft(recipe, inventory) {
			available = append(available, recipe)
		}
	}
	
	return available
}

// GetRecipeCount returns the number of loaded recipes
func (cs *CraftingSystem) GetRecipeCount() int {
	return len(cs.recipes)
}