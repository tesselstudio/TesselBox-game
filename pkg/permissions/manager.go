package permissions

import (
	"fmt"
	"strings"
	"sync"
)

// Manager handles all permission checking and role management
type Manager struct {
	roles    map[string]*Role
	registry *PlayerRegistry
	
	mu sync.RWMutex
}

// NewManager creates a new permission manager
func NewManager(registry *PlayerRegistry) *Manager {
	return &Manager{
		roles:    CreateDefaultRoles(),
		registry: registry,
	}
}

// HasPermission checks if a player has a permission
func (m *Manager) HasPermission(playerID string, node PermissionNode) bool {
	// Get player entry
	entry, exists := m.registry.GetByID(playerID)
	if !exists {
		return false // Unknown player has no permissions
	}
	
	return m.HasPermissionForEntry(entry, node)
}

// HasPermissionForEntry checks permission for a specific player entry
func (m *Manager) HasPermissionForEntry(entry *PlayerEntry, node PermissionNode) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	// 1. Check custom permission overrides (highest priority)
	if perm, exists := entry.CustomPerms[string(node)]; exists {
		return perm
	}
	
	// 2. Get role and check permissions
	role, exists := m.roles[entry.RoleID]
	if !exists {
		// Fall back to visitor role if role doesn't exist
		role, exists = m.roles["visitor"]
		if !exists {
			return false
		}
	}
	
	// 3. Get effective permissions (including inheritance)
	effectivePerms := role.GetEffectivePermissions(m.roles)
	
	// 4. Check explicit permission
	if granted, exists := effectivePerms[node]; exists {
		return granted
	}
	
	// 5. Check wildcards
	for perm := range effectivePerms {
		permStr := string(perm)
		if strings.HasSuffix(permStr, ".*") {
			prefix := strings.TrimSuffix(permStr, ".*")
			if strings.HasPrefix(string(node), prefix+".") {
				return true
			}
		}
	}
	
	return false
}

// HasAnyPermission checks if player has any of the given permissions
func (m *Manager) HasAnyPermission(playerID string, nodes ...PermissionNode) bool {
	for _, node := range nodes {
		if m.HasPermission(playerID, node) {
			return true
		}
	}
	return false
}

// HasAllPermissions checks if player has all of the given permissions
func (m *Manager) HasAllPermissions(playerID string, nodes ...PermissionNode) bool {
	for _, node := range nodes {
		if !m.HasPermission(playerID, node) {
			return false
		}
	}
	return true
}

// GetRole gets a role by ID
func (m *Manager) GetRole(roleID string) (*Role, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	role, exists := m.roles[roleID]
	return role, exists
}

// GetPlayerRole gets a player's role
func (m *Manager) GetPlayerRole(playerID string) (*Role, bool) {
	entry, exists := m.registry.GetByID(playerID)
	if !exists {
		return nil, false
	}
	
	return m.GetRole(entry.RoleID)
}

// SetPlayerRole changes a player's role
func (m *Manager) SetPlayerRole(playerID, roleID string) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	// Validate role exists
	if _, exists := m.roles[roleID]; !exists {
		return fmt.Errorf("role '%s' does not exist", roleID)
	}
	
	// Get player entry
	entry, exists := m.registry.GetByID(playerID)
	if !exists {
		return fmt.Errorf("player '%s' not found", playerID)
	}
	
	// Update role
	entry.SetRole(roleID)
	
	// Save registry
	return m.registry.Update(entry)
}

// GrantPlayerPermission grants a custom permission to a player
func (m *Manager) GrantPlayerPermission(playerID string, node PermissionNode) error {
	entry, exists := m.registry.GetByID(playerID)
	if !exists {
		return fmt.Errorf("player '%s' not found", playerID)
	}
	
	entry.GrantCustomPermission(node)
	return m.registry.Update(entry)
}

// DenyPlayerPermission denies a custom permission for a player
func (m *Manager) DenyPlayerPermission(playerID string, node PermissionNode) error {
	entry, exists := m.registry.GetByID(playerID)
	if !exists {
		return fmt.Errorf("player '%s' not found", playerID)
	}
	
	entry.DenyCustomPermission(node)
	return m.registry.Update(entry)
}

// RemovePlayerPermission removes a custom permission override
func (m *Manager) RemovePlayerPermission(playerID string, node PermissionNode) error {
	entry, exists := m.registry.GetByID(playerID)
	if !exists {
		return fmt.Errorf("player '%s' not found", playerID)
	}
	
	entry.RemoveCustomPermission(node)
	return m.registry.Update(entry)
}

// CreateRole creates a new custom role
func (m *Manager) CreateRole(id, name, description string, weight int) (*Role, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if _, exists := m.roles[id]; exists {
		return nil, fmt.Errorf("role '%s' already exists", id)
	}
	
	role := NewRole(id, name, description, weight)
	m.roles[id] = role
	
	return role, nil
}

// DeleteRole deletes a custom role (cannot delete built-in roles)
func (m *Manager) DeleteRole(roleID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	// Prevent deletion of built-in roles
	builtins := []string{"owner", "co-owner", "admin", "moderator", "helper", "player", "visitor", "shopkeeper", "builder"}
	for _, builtin := range builtins {
		if roleID == builtin {
			return fmt.Errorf("cannot delete built-in role '%s'", roleID)
		}
	}
	
	if _, exists := m.roles[roleID]; !exists {
		return fmt.Errorf("role '%s' does not exist", roleID)
	}
	
	// Reassign players with this role to default
	players := m.registry.GetByRole(roleID)
	for _, player := range players {
		player.SetRole(GetDefaultRole())
		m.registry.Update(player)
	}
	
	delete(m.roles, roleID)
	return nil
}

// GetAllRoles returns all roles
func (m *Manager) GetAllRoles() map[string]*Role {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	// Return a copy
	result := make(map[string]*Role)
	for k, v := range m.roles {
		result[k] = v
	}
	return result
}

// GetRoleHierarchy returns roles sorted by weight (highest first)
func (m *Manager) GetRoleHierarchy() []*Role {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	return RoleHierarchy(m.roles)
}

// CanPromoteTo checks if one role can promote to another
func (m *Manager) CanPromoteTo(promoterRoleID, targetRoleID string) bool {
	promoterWeight := GetRoleWeight(promoterRoleID, m.roles)
	targetWeight := GetRoleWeight(targetRoleID, m.roles)
	
	// Can only promote to roles of lower or equal weight
	return promoterWeight >= targetWeight
}

// CanDemoteFrom checks if one role can demote from another  
func (m *Manager) CanDemoteFrom(demoterRoleID, targetRoleID string) bool {
	demoterWeight := GetRoleWeight(demoterRoleID, m.roles)
	targetWeight := GetRoleWeight(targetRoleID, m.roles)
	
	// Can only demote roles of lower or equal weight
	return demoterWeight >= targetWeight
}

// IsAtLeastRole checks if player has at least the specified role level
func (m *Manager) IsAtLeastRole(playerID, minRoleID string) bool {
	entry, exists := m.registry.GetByID(playerID)
	if !exists {
		return false
	}
	
	playerWeight := GetRoleWeight(entry.RoleID, m.roles)
	minWeight := GetRoleWeight(minRoleID, m.roles)
	
	return playerWeight >= minWeight
}

// GetPlayerPrefix gets the chat prefix for a player
func (m *Manager) GetPlayerPrefix(playerID string) string {
	entry, exists := m.registry.GetByID(playerID)
	if !exists {
		return ""
	}
	
	role, exists := m.roles[entry.RoleID]
	if !exists {
		return ""
	}
	
	return role.Prefix
}

// CheckCommandPermission checks if a player can use a command
func (m *Manager) CheckCommandPermission(playerID, command string) bool {
	// Map command to permission node
	node := PermissionNode("commands." + command)
	return m.HasPermission(playerID, node)
}

// FilterAllowedCommands returns commands a player is allowed to use
func (m *Manager) FilterAllowedCommands(playerID string, allCommands []string) []string {
	allowed := make([]string, 0)
	
	for _, cmd := range allCommands {
		if m.CheckCommandPermission(playerID, cmd) {
			allowed = append(allowed, cmd)
		}
	}
	
	return allowed
}

// ReloadRoles reloads default roles (doesn't affect custom roles)
func (m *Manager) ReloadRoles() {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	defaults := CreateDefaultRoles()
	
	// Only reload built-in roles, preserve custom ones
	builtins := []string{"owner", "co-owner", "admin", "moderator", "helper", "player", "visitor", "shopkeeper", "builder"}
	
	for _, id := range builtins {
		if role, exists := defaults[id]; exists {
			m.roles[id] = role
		}
	}
}

// GetPlayerEffectivePermissions returns all effective permissions for a player
func (m *Manager) GetPlayerEffectivePermissions(playerID string) map[PermissionNode]bool {
	entry, exists := m.registry.GetByID(playerID)
	if !exists {
		return make(map[PermissionNode]bool)
	}
	
	return m.GetPlayerEffectivePermissionsForEntry(entry)
}

// GetPlayerEffectivePermissionsForEntry returns effective permissions for a player entry
func (m *Manager) GetPlayerEffectivePermissionsForEntry(entry *PlayerEntry) map[PermissionNode]bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	role, exists := m.roles[entry.RoleID]
	if !exists {
		return make(map[PermissionNode]bool)
	}
	
	// Get role effective permissions
	effective := role.GetEffectivePermissions(m.roles)
	
	// Apply custom overrides
	for node, granted := range entry.CustomPerms {
		effective[PermissionNode(node)] = granted
	}
	
	return effective
}
