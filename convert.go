package main

import (
	"encoding/json"
	"os"

	"gopkg.in/yaml.v3"
)

// BlockJSON represents the YAML structure for block definitions
type BlockJSON struct {
	ID          string                 `json:"id" yaml:"id"`
	Name        string                 `json:"name" yaml:"name"`
	Color       []uint8                `json:"color" yaml:"color"`
	TopColor    []uint8                `json:"topColor,omitempty" yaml:"topColor,omitempty"`
	SideColor   []uint8                `json:"sideColor,omitempty" yaml:"sideColor,omitempty"`
	Hardness    float64                `json:"hardness" yaml:"hardness"`
	Transparent bool                   `json:"transparent" yaml:"transparent"`
	Solid       bool                   `json:"solid" yaml:"solid"`
	Collectible bool                   `json:"collectible" yaml:"collectible"`
	Flammable   bool                   `json:"flammable" yaml:"flammable"`
	LightLevel  int                    `json:"lightLevel" yaml:"lightLevel"`
	Gravity     bool                   `json:"gravity" yaml:"gravity"`
	Viscosity   float64                `json:"viscosity" yaml:"viscosity"`
	Pattern     string                 `json:"pattern" yaml:"pattern"`
	UI          map[string]interface{} `json:"ui" yaml:"ui"`
	Function    map[string]interface{} `json:"function" yaml:"function"`
}

// OrganismJSON represents the YAML structure for organism definitions
type OrganismJSON struct {
	ID         string                 `json:"id" yaml:"id"`
	Name       string                 `json:"name" yaml:"name"`
	Type       string                 `json:"type" yaml:"type"`
	Appearance map[string]interface{} `json:"appearance" yaml:"appearance"`
	Properties map[string]interface{} `json:"properties" yaml:"properties"`
	Behavior   map[string]interface{} `json:"behavior" yaml:"behavior"`
	Function   map[string]interface{} `json:"function" yaml:"function"`
	Drops      []string               `json:"drops" yaml:"drops"`
}

// ItemQuantity represents an item and its quantity in crafting
type ItemQuantity struct {
	ItemType int `json:"item_type" yaml:"item_type"`
	Quantity int `json:"quantity" yaml:"quantity"`
}

// CraftingJSON represents the YAML structure for crafting recipes
type CraftingJSON struct {
	ID           string         `json:"id" yaml:"id"`
	Name         string         `json:"name" yaml:"name"`
	Description  string         `json:"description" yaml:"description"`
	Inputs       []ItemQuantity `json:"inputs" yaml:"inputs"`
	Outputs      []ItemQuantity `json:"outputs" yaml:"outputs"`
	CraftingTime float64        `json:"crafting_time" yaml:"crafting_time"`
	RequiredTool int            `json:"required_tool" yaml:"required_tool"`
}

// ItemJSON represents the YAML structure for item definitions
type ItemJSON struct {
	ID           string  `json:"id" yaml:"id"`
	Name         string  `json:"name" yaml:"name"`
	IconColor    []uint8 `json:"iconColor" yaml:"iconColor"`
	Description  string  `json:"description" yaml:"description"`
	StackSize    int     `json:"stackSize" yaml:"stackSize"`
	Durability   int     `json:"durability" yaml:"durability"`
	IsTool       bool    `json:"isTool" yaml:"isTool"`
	ToolPower    float64 `json:"toolPower" yaml:"toolPower"`
	IsPlaceable  bool    `json:"isPlaceable" yaml:"isPlaceable"`
	BlockType    string  `json:"blockType" yaml:"blockType"`
	IsWeapon     bool    `json:"isWeapon" yaml:"isWeapon"`
	WeaponDamage float64 `json:"weaponDamage" yaml:"weaponDamage"`
	WeaponRange  float64 `json:"weaponRange" yaml:"weaponRange"`
	WeaponSpeed  float64 `json:"weaponSpeed" yaml:"weaponSpeed"`
	WeaponType   string  `json:"weaponType" yaml:"weaponType"`
	IsArmor      bool    `json:"isArmor" yaml:"isArmor"`
	ArmorType    string  `json:"armorType" yaml:"armorType"`
	ArmorDefense float64 `json:"armorDefense" yaml:"armorDefense"`
}

func convertBlocks() {
	data, err := os.ReadFile("config/blocks.json")
	if err != nil {
		return
	}
	var blocks map[string]*BlockJSON
	json.Unmarshal(data, &blocks)
	data2, err := yaml.Marshal(blocks)
	if err != nil {
		return
	}
	os.WriteFile("config/blocks.yaml", data2, 0644)
}

func convertOrganisms() {
	data, err := os.ReadFile("config/organisms.json")
	if err != nil {
		return
	}
	var organisms map[string]*OrganismJSON
	json.Unmarshal(data, &organisms)
	data2, err := yaml.Marshal(organisms)
	if err != nil {
		return
	}
	os.WriteFile("config/organisms.yaml", data2, 0644)
}

func convertCrafting() {
	data, err := os.ReadFile("config/crafting_recipes.json")
	if err != nil {
		return
	}
	var recipes []CraftingJSON
	json.Unmarshal(data, &recipes)
	data2, err := yaml.Marshal(recipes)
	if err != nil {
		return
	}
	os.WriteFile("config/crafting_recipes.yaml", data2, 0644)
}

func convertItems() {
	data, err := os.ReadFile("config/items.json")
	if err != nil {
		return
	}
	var items map[string]*ItemJSON
	json.Unmarshal(data, &items)
	data2, err := yaml.Marshal(items)
	if err != nil {
		return
	}
	os.WriteFile("config/items.yaml", data2, 0644)
}

func main() {
	os.MkdirAll("config", 0755)
	convertBlocks()
	convertOrganisms()
	convertCrafting()
	convertItems()
}
