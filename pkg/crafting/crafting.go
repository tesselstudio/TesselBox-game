package crafting

import (
	"fmt"
	"tesselbox/pkg/items"

	"tesselbox/assets"

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

// CraftingSystem manages recipes and crafting operations
type CraftingSystem struct {
	recipes       map[string]*Recipe
	OnItemCrafted func(recipeID string) // Callback for when an item is crafted
}

// NewCraftingSystem creates a new crafting system
func NewCraftingSystem() *CraftingSystem {
	return &CraftingSystem{
		recipes: make(map[string]*Recipe),
	}
}

// LoadRecipes loads recipes from a YAML file with fallback support
func (cs *CraftingSystem) LoadRecipes(filePath string) error {
	data, err := assets.GetConfigFile("crafting_recipes.yaml")
	if err != nil {
		// If embedded file fails, try to load default recipes
		return cs.loadDefaultRecipes()
	}

	var recipes []Recipe
	if err := yaml.Unmarshal(data, &recipes); err != nil {
		// If parsing fails, fall back to default recipes
		return cs.loadDefaultRecipes()
	}

	// Clear existing recipes
	cs.recipes = make(map[string]*Recipe)

	// Load recipes
	loadedCount := 0
	for i := range recipes {
		if recipes[i].ID != "" {
			cs.recipes[recipes[i].ID] = &recipes[i]
			loadedCount++
		}
	}

	// If no valid recipes were loaded, fall back to defaults
	if loadedCount == 0 {
		return cs.loadDefaultRecipes()
	}

	return nil
}

// loadDefaultRecipes loads a minimal set of default recipes for basic gameplay
func (cs *CraftingSystem) loadDefaultRecipes() error {
	// Clear existing recipes
	cs.recipes = make(map[string]*Recipe)

	// Basic crafting recipes
	defaultRecipes := []Recipe{
		{
			ID:          "wooden_pickaxe",
			Name:        "Wooden Pickaxe",
			Description: "A basic pickaxe for mining stone",
			Inputs: []RecipeInput{
				{ItemType: items.PLANKS, Quantity: 3},
				{ItemType: items.STICK, Quantity: 2},
			},
			Outputs: []RecipeOutput{
				{ItemType: items.WOODEN_PICKAXE, Quantity: 1},
			},
			CraftingTime:    2.0,
			RequiredTool:    items.NONE,
			RequiredStation: STATION_WORKBENCH,
		},
		{
			ID:          "stone_pickaxe",
			Name:        "Stone Pickaxe",
			Description: "A sturdy pickaxe for mining ores",
			Inputs: []RecipeInput{
				{ItemType: items.STONE_BLOCK, Quantity: 3},
				{ItemType: items.STICK, Quantity: 2},
			},
			Outputs: []RecipeOutput{
				{ItemType: items.STONE_PICKAXE, Quantity: 1},
			},
			CraftingTime:    3.0,
			RequiredTool:    items.WOODEN_PICKAXE,
			RequiredStation: STATION_WORKBENCH,
		},
		{
			ID:          "planks",
			Name:        "Wooden Planks",
			Description: "Basic building material",
			Inputs: []RecipeInput{
				{ItemType: items.LOG_BLOCK, Quantity: 1},
			},
			Outputs: []RecipeOutput{
				{ItemType: items.PLANKS, Quantity: 4},
			},
			CraftingTime:    1.0,
			RequiredTool:    items.NONE,
			RequiredStation: STATION_NONE,
		},
		{
			ID:          "sticks",
			Name:        "Sticks",
			Description: "Basic crafting component",
			Inputs: []RecipeInput{
				{ItemType: items.PLANKS, Quantity: 2},
			},
			Outputs: []RecipeOutput{
				{ItemType: items.STICK, Quantity: 4},
			},
			CraftingTime:    0.5,
			RequiredTool:    items.NONE,
			RequiredStation: STATION_NONE,
		},
		{
			ID:          "workbench",
			Name:        "Workbench",
			Description: "Basic crafting station",
			Inputs: []RecipeInput{
				{ItemType: items.PLANKS, Quantity: 4},
			},
			Outputs: []RecipeOutput{
				{ItemType: items.WORKBENCH, Quantity: 1},
			},
			CraftingTime:    2.0,
			RequiredTool:    items.NONE,
			RequiredStation: STATION_NONE,
		},
		{
			ID:          "furnace",
			Name:        "Furnace",
			Description: "Smelting and cooking station",
			Inputs: []RecipeInput{
				{ItemType: items.STONE_BLOCK, Quantity: 8},
			},
			Outputs: []RecipeOutput{
				{ItemType: items.FURNACE, Quantity: 1},
			},
			CraftingTime:    3.0,
			RequiredTool:    items.WOODEN_PICKAXE,
			RequiredStation: STATION_WORKBENCH,
		},
		{
			ID:          "iron_pickaxe",
			Name:        "Iron Pickaxe",
			Description: "A durable pickaxe for mining tough materials",
			Inputs: []RecipeInput{
				{ItemType: items.IRON_INGOT, Quantity: 3},
				{ItemType: items.STICK, Quantity: 2},
			},
			Outputs: []RecipeOutput{
				{ItemType: items.IRON_PICKAXE, Quantity: 1},
			},
			CraftingTime:    4.0,
			RequiredTool:    items.STONE_PICKAXE,
			RequiredStation: STATION_WORKBENCH,
		},
		{
			ID:          "wooden_sword",
			Name:        "Wooden Sword",
			Description: "A basic weapon for defense",
			Inputs: []RecipeInput{
				{ItemType: items.PLANKS, Quantity: 2},
				{ItemType: items.STICK, Quantity: 1},
			},
			Outputs: []RecipeOutput{
				{ItemType: items.WOODEN_SWORD, Quantity: 1},
			},
			CraftingTime:    2.0,
			RequiredTool:    items.NONE,
			RequiredStation: STATION_WORKBENCH,
		},
		{
			ID:          "stone_sword",
			Name:        "Stone Sword",
			Description: "A sturdy weapon for combat",
			Inputs: []RecipeInput{
				{ItemType: items.STONE_BLOCK, Quantity: 2},
				{ItemType: items.STICK, Quantity: 1},
			},
			Outputs: []RecipeOutput{
				{ItemType: items.STONE_SWORD, Quantity: 1},
			},
			CraftingTime:    3.0,
			RequiredTool:    items.WOODEN_PICKAXE,
			RequiredStation: STATION_WORKBENCH,
		},
		{
			ID:          "iron_ingot",
			Name:        "Iron Ingot",
			Description: "Smelted iron for crafting",
			Inputs: []RecipeInput{
				{ItemType: items.COAL, Quantity: 1},
			},
			Outputs: []RecipeOutput{
				{ItemType: items.IRON_INGOT, Quantity: 1},
			},
			CraftingTime:    5.0,
			RequiredTool:    items.NONE,
			RequiredStation: STATION_FURNACE,
		},
		{
			ID:          "gold_ingot",
			Name:        "Gold Ingot",
			Description: "Smelted gold for crafting",
			Inputs: []RecipeInput{
				{ItemType: items.DIAMOND, Quantity: 1},
			},
			Outputs: []RecipeOutput{
				{ItemType: items.GOLD_INGOT, Quantity: 1},
			},
			CraftingTime:    5.0,
			RequiredTool:    items.NONE,
			RequiredStation: STATION_FURNACE,
		},
		{
			ID:          "iron_sword",
			Name:        "Iron Sword",
			Description: "A strong weapon for serious combat",
			Inputs: []RecipeInput{
				{ItemType: items.IRON_INGOT, Quantity: 2},
				{ItemType: items.STICK, Quantity: 1},
			},
			Outputs: []RecipeOutput{
				{ItemType: items.IRON_SWORD, Quantity: 1},
			},
			CraftingTime:    5.0,
			RequiredTool:    items.STONE_PICKAXE,
			RequiredStation: STATION_WORKBENCH,
		},
		{
			ID:          "anvil",
			Name:        "Anvil",
			Description: "Advanced crafting station for metalworking",
			Inputs: []RecipeInput{
				{ItemType: items.IRON_INGOT, Quantity: 5},
			},
			Outputs: []RecipeOutput{
				{ItemType: items.ANVIL, Quantity: 1},
			},
			CraftingTime:    6.0,
			RequiredTool:    items.IRON_PICKAXE,
			RequiredStation: STATION_WORKBENCH,
		},
		// Armor recipes - using wool as leather substitute
		{
			ID:          "leather_helmet",
			Name:        "Leather Cap",
			Description: "Basic head protection",
			Inputs: []RecipeInput{
				{ItemType: items.WOOL, Quantity: 5},
			},
			Outputs: []RecipeOutput{
				{ItemType: items.LEATHER_HELMET, Quantity: 1},
			},
			CraftingTime:    3.0,
			RequiredTool:    items.NONE,
			RequiredStation: STATION_WORKBENCH,
		},
		{
			ID:          "leather_chestplate",
			Name:        "Leather Tunic",
			Description: "Basic body protection",
			Inputs: []RecipeInput{
				{ItemType: items.WOOL, Quantity: 8},
			},
			Outputs: []RecipeOutput{
				{ItemType: items.LEATHER_CHESTPLATE, Quantity: 1},
			},
			CraftingTime:    4.0,
			RequiredTool:    items.NONE,
			RequiredStation: STATION_WORKBENCH,
		},
		{
			ID:          "leather_leggings",
			Name:        "Leather Pants",
			Description: "Basic leg protection",
			Inputs: []RecipeInput{
				{ItemType: items.WOOL, Quantity: 7},
			},
			Outputs: []RecipeOutput{
				{ItemType: items.LEATHER_LEGGINGS, Quantity: 1},
			},
			CraftingTime:    4.0,
			RequiredTool:    items.NONE,
			RequiredStation: STATION_WORKBENCH,
		},
		{
			ID:          "leather_boots",
			Name:        "Leather Boots",
			Description: "Basic foot protection",
			Inputs: []RecipeInput{
				{ItemType: items.WOOL, Quantity: 4},
			},
			Outputs: []RecipeOutput{
				{ItemType: items.LEATHER_BOOTS, Quantity: 1},
			},
			CraftingTime:    3.0,
			RequiredTool:    items.NONE,
			RequiredStation: STATION_WORKBENCH,
		},
		{
			ID:          "iron_helmet",
			Name:        "Iron Helmet",
			Description: "Metal head protection",
			Inputs: []RecipeInput{
				{ItemType: items.IRON_INGOT, Quantity: 5},
			},
			Outputs: []RecipeOutput{
				{ItemType: items.IRON_HELMET, Quantity: 1},
			},
			CraftingTime:    5.0,
			RequiredTool:    items.STONE_PICKAXE,
			RequiredStation: STATION_ANVIL,
		},
		{
			ID:          "iron_chestplate",
			Name:        "Iron Chestplate",
			Description: "Metal body protection",
			Inputs: []RecipeInput{
				{ItemType: items.IRON_INGOT, Quantity: 8},
			},
			Outputs: []RecipeOutput{
				{ItemType: items.IRON_CHESTPLATE, Quantity: 1},
			},
			CraftingTime:    6.0,
			RequiredTool:    items.STONE_PICKAXE,
			RequiredStation: STATION_ANVIL,
		},
		{
			ID:          "iron_leggings",
			Name:        "Iron Leggings",
			Description: "Metal leg protection",
			Inputs: []RecipeInput{
				{ItemType: items.IRON_INGOT, Quantity: 7},
			},
			Outputs: []RecipeOutput{
				{ItemType: items.IRON_LEGGINGS, Quantity: 1},
			},
			CraftingTime:    6.0,
			RequiredTool:    items.STONE_PICKAXE,
			RequiredStation: STATION_ANVIL,
		},
		{
			ID:          "iron_boots",
			Name:        "Iron Boots",
			Description: "Metal foot protection",
			Inputs: []RecipeInput{
				{ItemType: items.IRON_INGOT, Quantity: 4},
			},
			Outputs: []RecipeOutput{
				{ItemType: items.IRON_BOOTS, Quantity: 1},
			},
			CraftingTime:    5.0,
			RequiredTool:    items.STONE_PICKAXE,
			RequiredStation: STATION_ANVIL,
		},
		// Chest recipe
		{
			ID:          "chest",
			Name:        "Chest",
			Description: "Storage container",
			Inputs: []RecipeInput{
				{ItemType: items.PLANKS, Quantity: 8},
			},
			Outputs: []RecipeOutput{
				{ItemType: items.CHEST, Quantity: 1},
			},
			CraftingTime:    4.0,
			RequiredTool:    items.NONE,
			RequiredStation: STATION_WORKBENCH,
		},
		// Wings recipe - using string as feather substitute
		{
			ID:          "wings",
			Name:        "Wings",
			Description: "Allows flight",
			Inputs: []RecipeInput{
				{ItemType: items.WOOL, Quantity: 6},
				{ItemType: items.STRING, Quantity: 20},
				{ItemType: items.IRON_INGOT, Quantity: 2},
			},
			Outputs: []RecipeOutput{
				{ItemType: items.BOW, Quantity: 1}, // Using BOW as wings placeholder
			},
			CraftingTime:    8.0,
			RequiredTool:    items.IRON_PICKAXE,
			RequiredStation: STATION_ANVIL,
		},
		// Advanced tools
		{
			ID:          "diamond_pickaxe",
			Name:        "Diamond Pickaxe",
			Description: "The ultimate mining tool",
			Inputs: []RecipeInput{
				{ItemType: items.DIAMOND, Quantity: 3},
				{ItemType: items.STICK, Quantity: 2},
			},
			Outputs: []RecipeOutput{
				{ItemType: items.DIAMOND_SWORD, Quantity: 1}, // Using DIAMOND_SWORD as placeholder
			},
			CraftingTime:    6.0,
			RequiredTool:    items.IRON_PICKAXE,
			RequiredStation: STATION_ANVIL,
		},
		{
			ID:          "diamond_sword",
			Name:        "Diamond Sword",
			Description: "The ultimate weapon",
			Inputs: []RecipeInput{
				{ItemType: items.DIAMOND, Quantity: 2},
				{ItemType: items.STICK, Quantity: 1},
			},
			Outputs: []RecipeOutput{
				{ItemType: items.DIAMOND_SWORD, Quantity: 1},
			},
			CraftingTime:    6.0,
			RequiredTool:    items.IRON_PICKAXE,
			RequiredStation: STATION_ANVIL,
		},
	}

	// Load default recipes
	for i := range defaultRecipes {
		cs.recipes[defaultRecipes[i].ID] = &defaultRecipes[i]
	}

	return nil
}

// LoadRecipesFromAssets loads recipes from embedded assets
func (cs *CraftingSystem) LoadRecipesFromAssets() error {
	return cs.LoadRecipes("crafting_recipes.yaml")
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

	// Call crafting callback if set
	if cs.OnItemCrafted != nil {
		cs.OnItemCrafted(recipeID)
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
