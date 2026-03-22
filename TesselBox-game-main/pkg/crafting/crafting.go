package crafting

import (
	"fmt"
	"os"
	"tesselbox/pkg/items"

	"gopkg.in/yaml.v3"
)

// CraftingStation represents different crafting stations
type CraftingStation int

const (
	STATION_NONE CraftingStation = iota
	STATION_WORKBENCH
	STATION_FURNACE
	STATION_ANVIL
)

// RecipeInput represents an input item requirement
type RecipeInput struct {
	ItemType items.ItemType `yaml:"item_type"`
	Quantity int            `yaml:"quantity"`
}

// RecipeOutput represents an output item
type RecipeOutput struct {
	ItemType items.ItemType `yaml:"item_type"`
	Quantity int            `yaml:"quantity"`
}

// Recipe represents a crafting recipe
type Recipe struct {
	ID              string          `yaml:"id"`
	Name            string          `yaml:"name"`
	Description     string          `yaml:"description"`
	Inputs          []RecipeInput   `yaml:"inputs"`
	Outputs         []RecipeOutput  `yaml:"outputs"`
	CraftingTime    float64         `yaml:"crafting_time"`    // in seconds, 0 = instant
	RequiredTool    items.ItemType  `yaml:"required_tool"`    // NONE if no tool required
	RequiredStation CraftingStation `yaml:"required_station"` // STATION_NONE if no station required
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

// LoadRecipes loads recipes from a YAML file
func (cs *CraftingSystem) LoadRecipes(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read recipes file: %w", err)
	}

	var recipes []Recipe
	if err := yaml.Unmarshal(data, &recipes); err != nil {
		return fmt.Errorf("failed to parse recipes YAML: %w", err)
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

// CanCraft checks if the player can craft a recipe at the given station
func (cs *CraftingSystem) CanCraft(recipe *Recipe, inventory *items.Inventory, station CraftingStation) bool {
	// Check if required station matches
	if recipe.RequiredStation != STATION_NONE && recipe.RequiredStation != station {
		return false
	}

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

// GetMissingMaterials returns the materials needed to craft a recipe at the given station
func (cs *CraftingSystem) GetMissingMaterials(recipe *Recipe, inventory *items.Inventory, station CraftingStation) []RecipeInput {
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

// Craft attempts to craft a recipe at the given station
func (cs *CraftingSystem) Craft(recipeID string, inventory *items.Inventory, station CraftingStation) error {
	recipe, exists := cs.GetRecipe(recipeID)
	if !exists {
		return fmt.Errorf("recipe not found: %s", recipeID)
	}

	// Check if crafting is possible
	if !cs.CanCraft(recipe, inventory, station) {
		return fmt.Errorf("cannot craft %s: missing materials, tools, or station", recipe.Name)
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

// GetAvailableRecipes returns recipes that can be crafted with current inventory at the given station
func (cs *CraftingSystem) GetAvailableRecipes(inventory *items.Inventory, station CraftingStation) []*Recipe {
	available := []*Recipe{}

	for _, recipe := range cs.recipes {
		if cs.CanCraft(recipe, inventory, station) {
			available = append(available, recipe)
		}
	}

	return available
}

// GetRecipeCount returns the number of loaded recipes
func (cs *CraftingSystem) GetRecipeCount() int {
	return len(cs.recipes)
}
