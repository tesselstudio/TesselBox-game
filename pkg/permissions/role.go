package permissions

import (
	"fmt"
	"strings"
)

// Role represents a permission role
type Role struct {
	ID          string
	Name        string
	Description string
	Prefix      string // Chat prefix [ADMIN]
	Weight      int    // Higher weight = more permissions (Owner = 100, Admin = 90, etc)

	// Permissions
	Permissions   map[PermissionNode]bool // true = explicit grant, false = explicit deny
	InheritsFrom  []string                // Role IDs to inherit from
	WildcardPerms []string                // Wildcard patterns like "build.*", "admin.*"
}

// NewRole creates a new role
func NewRole(id, name, description string, weight int) *Role {
	return &Role{
		ID:            id,
		Name:          name,
		Description:   description,
		Weight:        weight,
		Permissions:   make(map[PermissionNode]bool),
		InheritsFrom:  []string{},
		WildcardPerms: []string{},
	}
}

// SetPrefix sets the chat prefix
func (r *Role) SetPrefix(prefix string) *Role {
	r.Prefix = prefix
	return r
}

// Grant grants a permission
func (r *Role) Grant(node PermissionNode) *Role {
	r.Permissions[node] = true
	return r
}

// Deny explicitly denies a permission
func (r *Role) Deny(node PermissionNode) *Role {
	r.Permissions[node] = false
	return r
}

// GrantWildcard grants permissions matching a wildcard pattern
func (r *Role) GrantWildcard(pattern string) *Role {
	r.WildcardPerms = append(r.WildcardPerms, pattern)
	return r
}

// Inherit adds inheritance from another role
func (r *Role) Inherit(roleID string) *Role {
	r.InheritsFrom = append(r.InheritsFrom, roleID)
	return r
}

// HasPermission checks if this role has a permission (local only, no inheritance)
func (r *Role) HasPermission(node PermissionNode) (bool, bool) {
	// Check explicit permissions first
	if granted, exists := r.Permissions[node]; exists {
		return granted, true // explicit
	}

	// Check wildcards
	for _, pattern := range r.WildcardPerms {
		if matchesWildcard(string(node), pattern) {
			return true, true
		}
	}

	return false, false // not set
}

// matchesWildcard checks if a node matches a wildcard pattern
func matchesWildcard(node, pattern string) bool {
	if pattern == "*" {
		return true
	}

	// Handle patterns like "build.*" or "admin.*"
	if strings.HasSuffix(pattern, ".*") {
		prefix := strings.TrimSuffix(pattern, ".*")
		return strings.HasPrefix(node, prefix+".")
	}

	return node == pattern
}

// String returns a string representation
func (r *Role) String() string {
	return fmt.Sprintf("Role[%s:%s w=%d perms=%d]", r.ID, r.Name, r.Weight, len(r.Permissions))
}

// GetEffectivePermissions returns all permissions including inherited
func (r *Role) GetEffectivePermissions(roleRegistry map[string]*Role) map[PermissionNode]bool {
	effective := make(map[PermissionNode]bool)
	visited := make(map[string]bool)

	r.collectPermissions(roleRegistry, effective, visited)

	return effective
}

// collectPermissions recursively collects permissions from inheritance chain
func (r *Role) collectPermissions(registry map[string]*Role, effective map[PermissionNode]bool, visited map[string]bool) {
	if visited[r.ID] {
		return // Prevent circular inheritance
	}
	visited[r.ID] = true

	// First, apply inherited permissions (lower priority)
	for _, parentID := range r.InheritsFrom {
		if parent, exists := registry[parentID]; exists {
			parent.collectPermissions(registry, effective, visited)
		}
	}

	// Then, apply wildcards
	for _, pattern := range r.WildcardPerms {
		for _, node := range AllNodes() {
			if matchesWildcard(string(node), pattern) {
				effective[node] = true
			}
		}
	}

	// Finally, apply explicit permissions (highest priority, can override)
	for node, granted := range r.Permissions {
		effective[node] = granted
	}
}

// RoleBuilder helps construct roles fluently
type RoleBuilder struct {
	role *Role
}

// NewRoleBuilder starts building a new role
func NewRoleBuilder(id, name string) *RoleBuilder {
	return &RoleBuilder{
		role: NewRole(id, name, "", 0),
	}
}

// WithDescription sets description
func (rb *RoleBuilder) WithDescription(desc string) *RoleBuilder {
	rb.role.Description = desc
	return rb
}

// WithWeight sets weight
func (rb *RoleBuilder) WithWeight(w int) *RoleBuilder {
	rb.role.Weight = w
	return rb
}

// WithPrefix sets prefix
func (rb *RoleBuilder) WithPrefix(p string) *RoleBuilder {
	rb.role.Prefix = p
	return rb
}

// Grant grants permission
func (rb *RoleBuilder) Grant(node PermissionNode) *RoleBuilder {
	rb.role.Grant(node)
	return rb
}

// GrantWildcard grants a wildcard pattern
func (rb *RoleBuilder) GrantWildcard(pattern string) *RoleBuilder {
	rb.role.GrantWildcard(pattern)
	return rb
}

// Deny denies permission
func (rb *RoleBuilder) Deny(node PermissionNode) *RoleBuilder {
	rb.role.Deny(node)
	return rb
}

// Wildcard grants wildcard
func (rb *RoleBuilder) Wildcard(pattern string) *RoleBuilder {
	rb.role.GrantWildcard(pattern)
	return rb
}

// Inherit adds inheritance
func (rb *RoleBuilder) Inherit(roleID string) *RoleBuilder {
	rb.role.Inherit(roleID)
	return rb
}

// Build returns the role
func (rb *RoleBuilder) Build() *Role {
	return rb.role
}

// Copy creates a copy of the role
func (r *Role) Copy(newID string) *Role {
	roleCopy := NewRole(newID, r.Name, r.Description, r.Weight)
	roleCopy.Prefix = r.Prefix

	for k, v := range r.Permissions {
		roleCopy.Permissions[k] = v
	}

	roleCopy.InheritsFrom = make([]string, len(r.InheritsFrom))
	copy(roleCopy.InheritsFrom, r.InheritsFrom)

	roleCopy.WildcardPerms = make([]string, len(r.WildcardPerms))
	copy(roleCopy.WildcardPerms, r.WildcardPerms)

	return roleCopy
}
