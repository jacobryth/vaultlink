package filter

import "strings"

// Role represents a named role with a set of allowed secret key prefixes.
type Role struct {
	Name     string
	Prefixes []string
}

// Filter holds a collection of roles used to filter secrets.
type Filter struct {
	roles map[string]Role
}

// NewFilter creates a Filter from a slice of Role definitions.
func NewFilter(roles []Role) *Filter {
	m := make(map[string]Role, len(roles))
	for _, r := range roles {
		m[r.Name] = r
	}
	return &Filter{roles: m}
}

// Apply returns only the key/value pairs from secrets whose keys match
// any of the prefixes defined for the given role. If the role is not
// found, an empty map is returned.
func (f *Filter) Apply(role string, secrets map[string]string) map[string]string {
	r, ok := f.roles[role]
	if !ok {
		return map[string]string{}
	}

	// Empty prefix list means allow all keys.
	if len(r.Prefixes) == 0 {
		result := make(map[string]string, len(secrets))
		for k, v := range secrets {
			result[k] = v
		}
		return result
	}

	result := make(map[string]string)
	for k, v := range secrets {
		for _, prefix := range r.Prefixes {
			if strings.HasPrefix(k, prefix) {
				result[k] = v
				break
			}
		}
	}
	return result
}

// HasRole reports whether the given role name is registered.
func (f *Filter) HasRole(role string) bool {
	_, ok := f.roles[role]
	return ok
}
