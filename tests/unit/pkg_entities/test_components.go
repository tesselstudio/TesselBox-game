package pkg_entities

import (
	"image/color"
	"testing"

	"tesselbox/pkg/entities"
	"tesselbox/tests/internal"
)

func TestRenderComponent(t *testing.T) {
	helper := internal.NewTestHelper(t)

	// Test component creation
	component := &entities.RenderComponent{
		Type:  "render",
		Color: color.RGBA{255, 0, 0, 255},
		Scale: 1.5,
	}

	// Test component type
	helper.AssertEqual("render", component.GetType(), "Component type should be render")

	// Test component cloning
	clone := component.Clone()
	helper.AssertEqual(component.GetType(), clone.GetType(), "Cloned component should have same type")

	// Test component validation
	err := component.Validate()
	helper.AssertNoError(err, "Valid render component should not error")

	// Test invalid component
	invalidComponent := &entities.RenderComponent{
		Type:  "render",
		Color: color.RGBA{255, 0, 0, 255},
		Scale: -1.0, // Invalid scale
	}
	err = invalidComponent.Validate()
	helper.AssertError(err, "Invalid render component should error")
}

func TestPhysicsComponent(t *testing.T) {
	helper := internal.NewTestHelper(t)

	component := &entities.PhysicsComponent{
		Type:     "physics",
		Hardness: 5.0,
		Density:  2.0,
		Mass:     10.0,
	}

	// Test component type
	helper.AssertEqual("physics", component.GetType(), "Component type should be physics")

	// Test component cloning
	clone := component.Clone()
	helper.AssertEqual(component.GetType(), clone.GetType(), "Cloned component should have same type")

	// Test component validation
	err := component.Validate()
	helper.AssertNoError(err, "Valid physics component should not error")

	// Test invalid mass
	invalidComponent := &entities.PhysicsComponent{
		Type: "physics",
		Mass: -1.0, // Invalid mass
	}
	err = invalidComponent.Validate()
	helper.AssertError(err, "Invalid physics component should error")
}

func TestInventoryComponent(t *testing.T) {
	helper := internal.NewTestHelper(t)

	component := &entities.InventoryComponent{
		Type:     "inventory",
		Slots:    10,
		Contents: make(map[string]int),
	}

	// Test component type
	helper.AssertEqual("inventory", component.GetType(), "Component type should be inventory")

	// Test component cloning
	clone := component.Clone()
	helper.AssertEqual(component.GetType(), clone.GetType(), "Cloned component should have same type")

	// Test component validation
	err := component.Validate()
	helper.AssertNoError(err, "Valid inventory component should not error")

	// Test invalid slots
	invalidComponent := &entities.InventoryComponent{
		Type:  "inventory",
		Slots: -1, // Invalid slots
	}
	err = invalidComponent.Validate()
	helper.AssertError(err, "Invalid inventory component should error")
}

func TestComponentMerge(t *testing.T) {
	helper := internal.NewTestHelper(t)

	// Test merging render components
	base := &entities.RenderComponent{
		Type:  "render",
		Color: color.RGBA{255, 0, 0, 255},
		Scale: 1.0,
	}

	merge := &entities.RenderComponent{
		Type:  "render",
		Color: color.RGBA{0, 255, 0, 255},
		Scale: 2.0,
	}

	base.Merge(merge)
	helper.AssertEqual(2.0, base.Scale, "Scale should be merged")
	helper.AssertEqual(color.RGBA{0, 255, 0, 255}, base.Color, "Color should be merged")
}

func TestComponentRegistry(t *testing.T) {
	helper := internal.NewTestHelper(t)

	// Test component registration
	component := internal.CreateTestComponent("render")
	entities.RegisterComponent("test_render", component)

	// Test component retrieval
	retrieved, err := entities.CreateComponent("test_render")
	helper.AssertNotNil(retrieved, "Registered component should be retrievable")
	helper.AssertNoError(err, "Component creation should not error")

	// Test component creation from type
	instance := retrieved.Clone()
	helper.AssertEqual("render", instance.GetType(), "Created component should have correct type")
}

func TestComponentValidation(t *testing.T) {
	tests := []struct {
		name        string
		component   entities.Component
		shouldError bool
	}{
		{
			name: "valid render component",
			component: &entities.RenderComponent{
				Type:  "render",
				Color: color.RGBA{255, 255, 255, 255},
				Scale: 1.0,
			},
			shouldError: false,
		},
		{
			name: "invalid render scale",
			component: &entities.RenderComponent{
				Type:  "render",
				Color: color.RGBA{255, 255, 255, 255},
				Scale: -1.0,
			},
			shouldError: true,
		},
		{
			name: "valid physics component",
			component: &entities.PhysicsComponent{
				Type:     "physics",
				Hardness: 0,
				Density:  0,
				Mass:     1.0,
			},
			shouldError: false,
		},
		{
			name: "invalid physics mass",
			component: &entities.PhysicsComponent{
				Type:     "physics",
				Hardness: 0,
				Density:  0,
				Mass:     -1.0,
			},
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.component.Validate()
			if tt.shouldError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
			}
		})
	}
}
