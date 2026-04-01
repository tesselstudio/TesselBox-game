package entities

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"tesselbox/assets"
	"time"

	"gopkg.in/yaml.v3"
)

// DataLoader handles loading entity definitions from various sources
type DataLoader struct {
	entityManager     *EntityManager
	componentRegistry map[string]reflect.Type
	loadedTemplates   map[string]bool
	loadedConfigs     map[string]bool
}

// NewDataLoader creates a new data loader
func NewDataLoader(entityManager *EntityManager) *DataLoader {
	return &DataLoader{
		entityManager:     entityManager,
		componentRegistry: ComponentRegistry,
		loadedTemplates:   make(map[string]bool),
		loadedConfigs:     make(map[string]bool),
	}
}

// LoadAll loads all entity data from embedded assets
func (dl *DataLoader) LoadAll() error {
	log.Println("Loading entity data...")

	// Load entity templates
	if err := dl.LoadEntityTemplates(); err != nil {
		return fmt.Errorf("failed to load entity templates: %v", err)
	}

	// Load blocks configuration
	if err := dl.LoadBlocksConfig(); err != nil {
		return fmt.Errorf("failed to load blocks config: %v", err)
	}

	// Load items configuration
	if err := dl.LoadItemsConfig(); err != nil {
		return fmt.Errorf("failed to load items config: %v", err)
	}

	// Load organisms configuration
	if err := dl.LoadOrganismsConfig(); err != nil {
		return fmt.Errorf("failed to load organisms config: %v", err)
	}

	// Load crafting recipes
	if err := dl.LoadCraftingConfig(); err != nil {
		return fmt.Errorf("failed to load crafting config: %v", err)
	}

	log.Println("Entity data loading completed")
	return nil
}

// LoadEntityTemplates loads entity templates from YAML
func (dl *DataLoader) LoadEntityTemplates() error {
	data, err := assets.GetConfigFile("entities.yaml")
	if err != nil {
		// File not found is not an error - entities may be defined in other files
		log.Printf("No entities.yaml found, skipping entity templates")
		return nil
	}

	return dl.entityManager.LoadTemplates(data)
}

// LoadBlocksConfig loads block definitions and converts to entity templates
func (dl *DataLoader) LoadBlocksConfig() error {
	if dl.loadedConfigs["blocks"] {
		return nil
	}

	data, err := assets.GetConfigFile("blocks.yaml")
	if err != nil {
		log.Printf("No blocks.yaml found, skipping blocks config")
		return nil
	}

	var blocks map[string]*BlockConfig
	err = yaml.Unmarshal(data, &blocks)
	if err != nil {
		return fmt.Errorf("failed to parse blocks.yaml: %v", err)
	}

	// Convert blocks to entity templates
	templates := make(map[string]*EntityTemplate)
	for blockID, block := range blocks {
		template := dl.convertBlockToTemplate(blockID, block)
		templates[blockID] = template
	}

	// Load templates into entity manager
	templatesData, err := yaml.Marshal(templates)
	if err != nil {
		return fmt.Errorf("failed to marshal block templates: %v", err)
	}

	err = dl.entityManager.LoadTemplates(templatesData)
	if err != nil {
		return fmt.Errorf("failed to load block templates: %v", err)
	}

	dl.loadedConfigs["blocks"] = true
	log.Printf("Loaded %d block templates", len(blocks))
	return nil
}

// LoadItemsConfig loads item definitions and converts to entity templates
func (dl *DataLoader) LoadItemsConfig() error {
	if dl.loadedConfigs["items"] {
		return nil
	}

	data, err := assets.GetConfigFile("items.yaml")
	if err != nil {
		log.Printf("No items.yaml found, skipping items config")
		return nil
	}

	var items map[string]*ItemConfig
	err = yaml.Unmarshal(data, &items)
	if err != nil {
		return fmt.Errorf("failed to parse items.yaml: %v", err)
	}

	// Convert items to entity templates
	templates := make(map[string]*EntityTemplate)
	for itemID, item := range items {
		template := dl.convertItemToTemplate(itemID, item)
		templates[itemID] = template
	}

	// Load templates into entity manager
	templatesData, err := yaml.Marshal(templates)
	if err != nil {
		return fmt.Errorf("failed to marshal item templates: %v", err)
	}

	err = dl.entityManager.LoadTemplates(templatesData)
	if err != nil {
		return fmt.Errorf("failed to load item templates: %v", err)
	}

	dl.loadedConfigs["items"] = true
	log.Printf("Loaded %d item templates", len(items))
	return nil
}

// LoadOrganismsConfig loads organism definitions and converts to entity templates
func (dl *DataLoader) LoadOrganismsConfig() error {
	if dl.loadedConfigs["organisms"] {
		return nil
	}

	data, err := assets.GetConfigFile("organisms.yaml")
	if err != nil {
		log.Printf("No organisms.yaml found, skipping organisms config")
		return nil
	}

	var organisms map[string]*OrganismConfig
	err = yaml.Unmarshal(data, &organisms)
	if err != nil {
		return fmt.Errorf("failed to parse organisms.yaml: %v", err)
	}

	// Convert organisms to entity templates
	templates := make(map[string]*EntityTemplate)
	for organismID, organism := range organisms {
		template := dl.convertOrganismToTemplate(organismID, organism)
		templates[organismID] = template
	}

	// Load templates into entity manager
	templatesData, err := yaml.Marshal(templates)
	if err != nil {
		return fmt.Errorf("failed to marshal organism templates: %v", err)
	}

	err = dl.entityManager.LoadTemplates(templatesData)
	if err != nil {
		return fmt.Errorf("failed to load organism templates: %v", err)
	}

	dl.loadedConfigs["organisms"] = true
	log.Printf("Loaded %d organism templates", len(organisms))
	return nil
}

// LoadCraftingConfig loads crafting recipes
func (dl *DataLoader) LoadCraftingConfig() error {
	if dl.loadedConfigs["crafting"] {
		return nil
	}

	data, err := assets.GetConfigFile("crafting_recipes.yaml")
	if err != nil {
		log.Printf("No crafting_recipes.yaml found, skipping crafting config")
		return nil
	}

	var recipes map[string]*CraftingRecipeConfig
	err = yaml.Unmarshal(data, &recipes)
	if err != nil {
		// If that fails, try to parse as list (legacy format)
		var recipeList []*CraftingRecipeConfig
		err = yaml.Unmarshal(data, &recipeList)
		if err != nil {
			return fmt.Errorf("failed to parse crafting_recipes.yaml: %v", err)
		}

		// Convert list to map
		recipes = make(map[string]*CraftingRecipeConfig)
		for _, recipe := range recipeList {
			recipes[recipe.ID] = recipe
		}
	}

	// Convert recipes to crafting system format
	recipeMap := make(map[string]*CraftingRecipe)
	for recipeID, recipe := range recipes {
		craftingRecipe := dl.convertCraftingRecipe(recipeID, recipe)
		recipeMap[recipeID] = craftingRecipe
	}

	// Add recipes to crafting system
	for _, recipe := range recipeMap {
		// This would be added to the crafting system
		log.Printf("Loaded crafting recipe: %s", recipe.ID)
	}

	dl.loadedConfigs["crafting"] = true
	log.Printf("Loaded %d crafting recipes", len(recipes))
	return nil
}

// ============================================================================
// Configuration Structures
// ============================================================================

// BlockConfig represents a block configuration from YAML
type BlockConfig struct {
	ID          string                 `yaml:"id"`
	Name        string                 `yaml:"name"`
	Color       []uint8                `yaml:"color"`
	TopColor    []uint8                `yaml:"topColor,omitempty"`
	SideColor   []uint8                `yaml:"sideColor,omitempty"`
	Colors      [][]uint8              `yaml:"colors,omitempty"`
	Hardness    float64                `yaml:"hardness"`
	Transparent bool                   `yaml:"transparent"`
	Solid       bool                   `yaml:"solid"`
	Collectible bool                   `yaml:"collectible"`
	Flammable   bool                   `yaml:"flammable"`
	LightLevel  int                    `yaml:"lightLevel"`
	Gravity     bool                   `yaml:"gravity"`
	Viscosity   float64                `yaml:"viscosity"`
	Pattern     string                 `yaml:"pattern"`
	Texture     string                 `yaml:"texture,omitempty"`
	UI          map[string]interface{} `yaml:"ui,omitempty"`
	Function    map[string]interface{} `yaml:"function,omitempty"`
}

// ItemConfig represents an item configuration from YAML
type ItemConfig struct {
	ID           string  `yaml:"id"`
	Name         string  `yaml:"name"`
	IconColor    []uint8 `yaml:"iconColor"`
	Description  string  `yaml:"description"`
	StackSize    int     `yaml:"stackSize"`
	Durability   int     `yaml:"durability"`
	IsTool       bool    `yaml:"isTool"`
	ToolPower    float64 `yaml:"toolPower"`
	IsPlaceable  bool    `yaml:"isPlaceable"`
	BlockType    string  `yaml:"blockType"`
	IsWeapon     bool    `yaml:"isWeapon"`
	WeaponDamage float64 `yaml:"weaponDamage"`
	WeaponRange  float64 `yaml:"weaponRange"`
	WeaponSpeed  float64 `yaml:"weaponSpeed"`
	WeaponType   string  `yaml:"weaponType"`
	IsArmor      bool    `yaml:"isArmor"`
	ArmorType    string  `yaml:"armorType"`
	ArmorDefense float64 `yaml:"armorDefense"`
}

// OrganismConfig represents an organism configuration from YAML
type OrganismConfig struct {
	ID         string                 `yaml:"id"`
	Name       string                 `yaml:"name"`
	Type       string                 `yaml:"type"`
	Appearance map[string]interface{} `yaml:"appearance"`
	Properties map[string]interface{} `yaml:"properties"`
	Behavior   map[string]interface{} `yaml:"behavior"`
	Function   map[string]interface{} `yaml:"function"`
	Drops      []string               `yaml:"drops"`
}

// CraftingRecipeConfig represents a crafting recipe configuration from YAML
type CraftingRecipeConfig struct {
	ID           string                   `yaml:"id"`
	Name         string                   `yaml:"name"`
	Inputs       []map[string]interface{} `yaml:"inputs"`
	Outputs      []map[string]interface{} `yaml:"outputs"`
	RequiredTool string                   `yaml:"requiredTool,omitempty"`
	CraftingTime string                   `yaml:"craftingTime"`
	Category     string                   `yaml:"category"`
	Tier         int                      `yaml:"tier"`
	Properties   map[string]interface{}   `yaml:"properties,omitempty"`
}

// ============================================================================
// Conversion Functions
// ============================================================================

// convertBlockToTemplate converts a block config to an entity template
func (dl *DataLoader) convertBlockToTemplate(blockID string, block *BlockConfig) *EntityTemplate {
	template := &EntityTemplate{
		ID:         blockID,
		Type:       "block",
		Name:       block.Name,
		Tags:       []string{"block", "static"},
		Components: make(map[string]interface{}),
		Inherits:   []string{},
	}

	// Render component
	renderComp := map[string]interface{}{
		"type":        "render",
		"color":       block.Color,
		"pattern":     block.Pattern,
		"visible":     true,
		"lightLevel":  block.LightLevel,
		"transparent": block.Transparent,
		"scale":       1.0,
		"animated":    false,
	}

	if len(block.TopColor) > 0 {
		renderComp["topColor"] = block.TopColor
	}
	if len(block.SideColor) > 0 {
		renderComp["sideColor"] = block.SideColor
	}
	if len(block.Colors) > 0 {
		renderComp["colors"] = block.Colors
	}
	if block.Texture != "" {
		renderComp["texture"] = block.Texture
	}

	template.Components["render"] = renderComp

	// Physics component
	physicsComp := map[string]interface{}{
		"type":       "physics",
		"hardness":   block.Hardness,
		"density":    1.0,
		"solid":      block.Solid,
		"gravity":    block.Gravity,
		"viscosity":  block.Viscosity,
		"friction":   0.6,
		"bounciness": 0.0,
		"collision":  block.Solid,
		"mass":       1.0,
	}

	template.Components["physics"] = physicsComp

	// Inventory component (for collectible blocks)
	if block.Collectible {
		inventoryComp := map[string]interface{}{
			"type":          "inventory",
			"stackSize":     64,
			"maxDurability": -1,
			"container":     false,
			"slots":         0,
			"weight":        1.0,
			"categories":    []string{"block", "building"},
		}
		template.Components["inventory"] = inventoryComp
		template.Tags = append(template.Tags, "collectible")
	}

	// Add flammable tag
	if block.Flammable {
		template.Tags = append(template.Tags, "flammable")
	}

	return template
}

// convertItemToTemplate converts an item config to an entity template
func (dl *DataLoader) convertItemToTemplate(itemID string, item *ItemConfig) *EntityTemplate {
	template := &EntityTemplate{
		ID:         itemID,
		Type:       "item",
		Name:       item.Name,
		Tags:       []string{"item"},
		Components: make(map[string]interface{}),
		Inherits:   []string{},
	}

	// Render component
	renderComp := map[string]interface{}{
		"type":        "render",
		"color":       item.IconColor,
		"pattern":     "solid",
		"visible":     true,
		"lightLevel":  0,
		"transparent": false,
		"scale":       0.8,
		"animated":    false,
	}
	template.Components["render"] = renderComp

	// Inventory component
	inventoryComp := map[string]interface{}{
		"type":              "inventory",
		"stackSize":         item.StackSize,
		"maxDurability":     item.Durability,
		"currentDurability": item.Durability,
		"container":         false,
		"slots":             0,
		"weight":            1.0,
		"categories":        []string{"item"},
	}

	// Add categories based on item properties
	if item.IsTool {
		inventoryComp["categories"] = append(inventoryComp["categories"].([]string), "tool")
		template.Tags = append(template.Tags, "tool")
	}
	if item.IsWeapon {
		inventoryComp["categories"] = append(inventoryComp["categories"].([]string), "weapon")
		template.Tags = append(template.Tags, "weapon")
	}
	if item.IsArmor {
		inventoryComp["categories"] = append(inventoryComp["categories"].([]string), "armor")
		template.Tags = append(template.Tags, "armor")
	}
	if item.IsPlaceable {
		inventoryComp["categories"] = append(inventoryComp["categories"].([]string), "placeable")
		template.Tags = append(template.Tags, "placeable")
	}

	template.Components["inventory"] = inventoryComp

	// Tool component
	if item.IsTool {
		toolComp := map[string]interface{}{
			"type":          "tool",
			"toolType":      dl.inferToolType(itemID),
			"power":         item.ToolPower,
			"efficiency":    item.ToolPower / 2.0,
			"durability":    item.Durability,
			"maxDurability": item.Durability,
			"effective":     dl.inferEffectiveBlocks(itemID),
		}
		template.Components["tool"] = toolComp
	}

	// Combat component
	if item.IsWeapon || item.IsArmor {
		combatComp := map[string]interface{}{
			"type": "combat",
		}

		if item.IsWeapon {
			combatComp["weaponType"] = item.WeaponType
			combatComp["damage"] = item.WeaponDamage
			combatComp["range"] = item.WeaponRange
			combatComp["speed"] = item.WeaponSpeed
			combatComp["health"] = 0
			combatComp["maxHealth"] = 0
		}

		if item.IsArmor {
			combatComp["armorType"] = item.ArmorType
			combatComp["armor"] = item.ArmorDefense
			combatComp["health"] = 0
			combatComp["maxHealth"] = 0
		}

		template.Components["combat"] = combatComp
	}

	return template
}

// convertOrganismToTemplate converts an organism config to an entity template
func (dl *DataLoader) convertOrganismToTemplate(organismID string, organism *OrganismConfig) *EntityTemplate {
	template := &EntityTemplate{
		ID:         organismID,
		Type:       "organism",
		Name:       organism.Name,
		Tags:       []string{"organism", "living"},
		Components: make(map[string]interface{}),
		Inherits:   []string{},
	}

	// Render component
	renderComp := map[string]interface{}{
		"type":        "render",
		"color":       []uint8{100, 200, 100, 255}, // Default green
		"pattern":     "solid",
		"visible":     true,
		"lightLevel":  0,
		"transparent": false,
		"scale":       1.0,
		"animated":    true,
	}

	// Extract color from appearance if available
	if color, ok := organism.Appearance["color"]; ok {
		if colorSlice, ok := color.([]uint8); ok {
			renderComp["color"] = colorSlice
		}
	}

	template.Components["render"] = renderComp

	// Physics component
	physicsComp := map[string]interface{}{
		"type":       "physics",
		"hardness":   1.0,
		"density":    1.0,
		"solid":      false,
		"gravity":    false,
		"viscosity":  0.0,
		"friction":   0.5,
		"bounciness": 0.0,
		"collision":  true,
		"mass":       1.0,
	}
	template.Components["physics"] = physicsComp

	// Behavior component
	behaviorComp := map[string]interface{}{
		"type":         "behavior",
		"aiType":       "passive",
		"passive":      true,
		"hostile":      false,
		"neutral":      false,
		"states":       []string{"idle", "active"},
		"currentState": "idle",
		"sightRange":   10.0,
		"hearingRange": 5.0,
		"speed":        1.0,
		"jumpHeight":   0.0,
		"abilities":    []string{},
	}

	// Extract behavior properties
	if isHostile, ok := organism.Properties["isHostile"]; ok {
		if hostile, ok := isHostile.(bool); ok && hostile {
			behaviorComp["hostile"] = true
			behaviorComp["passive"] = false
			behaviorComp["aiType"] = "hostile"
			template.Tags = append(template.Tags, "hostile")
		}
	}

	if maxHealth, ok := organism.Properties["maxHealth"]; ok {
		if _, ok := maxHealth.(float64); ok {
			// Will be used in combat component
		}
	}

	if damage, ok := organism.Properties["damage"]; ok {
		if _, ok := damage.(float64); ok {
			// Will be used in combat component
		}
	}

	if attackRange, ok := organism.Properties["attackRange"]; ok {
		if rng, ok := attackRange.(float64); ok {
			behaviorComp["sightRange"] = rng
		}
	}

	template.Components["behavior"] = behaviorComp

	// Combat component
	combatComp := map[string]interface{}{
		"type":         "combat",
		"weaponType":   "natural",
		"damage":       0.0,
		"range":        1.0,
		"speed":        1.0,
		"armor":        0.0,
		"health":       10.0,
		"maxHealth":    10.0,
		"regeneration": 0.1,
	}

	// Extract combat properties
	if maxHealth, ok := organism.Properties["maxHealth"]; ok {
		if health, ok := maxHealth.(float64); ok {
			combatComp["health"] = health
			combatComp["maxHealth"] = health
		}
	}

	if damage, ok := organism.Properties["damage"]; ok {
		if dmg, ok := damage.(float64); ok {
			combatComp["damage"] = dmg
		}
	}

	if attackRange, ok := organism.Properties["attackRange"]; ok {
		if rng, ok := attackRange.(float64); ok {
			combatComp["range"] = rng
		}
	}

	template.Components["combat"] = combatComp

	// Inventory component for drops
	if len(organism.Drops) > 0 {
		inventoryComp := map[string]interface{}{
			"type":          "inventory",
			"stackSize":     1,
			"maxDurability": -1,
			"container":     false,
			"slots":         0,
			"weight":        1.0,
			"categories":    []string{"organism"},
			"properties": map[string]interface{}{
				"drops": organism.Drops,
			},
		}
		template.Components["inventory"] = inventoryComp
	}

	return template
}

// convertCraftingRecipe converts a crafting recipe config to a crafting recipe
func (dl *DataLoader) convertCraftingRecipe(recipeID string, config *CraftingRecipeConfig) *CraftingRecipe {
	recipe := &CraftingRecipe{
		ID:           recipeID,
		Name:         config.Name,
		Inputs:       make(map[string]int),
		Outputs:      make(map[string]int),
		RequiredTool: config.RequiredTool,
		Category:     config.Category,
		Tier:         config.Tier,
	}

	// Convert inputs from list to map
	for _, input := range config.Inputs {
		if itemType, ok := input["item_type"].(int); ok {
			if quantity, ok := input["quantity"].(int); ok {
				// Convert int item_type to string using existing item mappings
				itemTypeStr := dl.getItemTypeString(itemType)
				recipe.Inputs[itemTypeStr] = quantity
			}
		}
	}

	// Convert outputs from list to map
	for _, output := range config.Outputs {
		if itemType, ok := output["item_type"].(int); ok {
			if quantity, ok := output["quantity"].(int); ok {
				// Convert int item_type to string using existing item mappings
				itemTypeStr := dl.getItemTypeString(itemType)
				recipe.Outputs[itemTypeStr] = quantity
			}
		}
	}

	// Parse crafting time
	if config.CraftingTime != "" {
		// Simple parsing - in production use proper duration parsing
		recipe.CraftingTime = 3000 * time.Millisecond // Default 3 seconds
	}

	return recipe
}

// ============================================================================
// Utility Functions
// ============================================================================

// inferToolType infers tool type from item ID
func (dl *DataLoader) inferToolType(itemID string) string {
	if strings.Contains(itemID, "pickaxe") {
		return "pickaxe"
	}
	if strings.Contains(itemID, "axe") {
		return "axe"
	}
	if strings.Contains(itemID, "shovel") {
		return "shovel"
	}
	if strings.Contains(itemID, "sword") {
		return "sword"
	}
	if strings.Contains(itemID, "bow") {
		return "bow"
	}
	return "tool"
}

// inferEffectiveBlocks infers which blocks a tool is effective against
func (dl *DataLoader) inferEffectiveBlocks(itemID string) []string {
	if strings.Contains(itemID, "pickaxe") {
		return []string{"stone", "coal_ore", "iron_ore", "gold_ore", "diamond_ore"}
	}
	if strings.Contains(itemID, "axe") {
		return []string{"log", "wood"}
	}
	if strings.Contains(itemID, "shovel") {
		return []string{"dirt", "sand", "gravel"}
	}
	return []string{"all"}
}

// getItemTypeString converts int item type to string using existing mappings
func (dl *DataLoader) getItemTypeString(itemType int) string {
	// Map item type IDs to strings based on existing items.go
	itemTypeMap := map[int]string{
		1:  "dirt_block",
		2:  "grass_block",
		3:  "stone_block",
		4:  "sand_block",
		5:  "log_block",
		6:  "coal",
		7:  "iron_ingot",
		8:  "gold_ingot",
		9:  "diamond",
		10: "wooden_pickaxe",
		11: "stone_pickaxe",
		12: "iron_pickaxe",
		13: "planks",
		14: "stick",
		15: "workbench",
		16: "furnace",
		17: "anvil",
	}

	if str, ok := itemTypeMap[itemType]; ok {
		return str
	}
	return "unknown"
}

// GetLoadedTemplates returns a list of loaded template IDs
func (dl *DataLoader) GetLoadedTemplates() []string {
	return dl.entityManager.ListTemplates()
}

// IsTemplateLoaded checks if a template is loaded
func (dl *DataLoader) IsTemplateLoaded(templateID string) bool {
	_, exists := dl.entityManager.GetTemplate(templateID)
	return exists
}

// Reload reloads all configurations
func (dl *DataLoader) Reload() error {
	// Clear loaded flags
	dl.loadedTemplates = make(map[string]bool)
	dl.loadedConfigs = make(map[string]bool)

	// Clear entity manager templates
	// This would need to be implemented in EntityManager

	return dl.LoadAll()
}
