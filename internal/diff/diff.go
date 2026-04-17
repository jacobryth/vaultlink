package diff

import "fmt"

// Change represents a single secret change between two snapshots.
type Change struct {
	Key    string
	OldVal string
	NewVal string
	Type   ChangeType
}

// ChangeType describes the kind of change detected.
type ChangeType string

const (
	Added   ChangeType = "added"
	Removed ChangeType = "removed"
	Updated ChangeType = "updated"
)

// Compare returns the list of changes between previous and current secret maps.
func Compare(previous, current map[string]string) []Change {
	var changes []Change

	for k, newVal := range current {
		oldVal, exists := previous[k]
		if !exists {
			changes = append(changes, Change{Key: k, NewVal: newVal, Type: Added})
		} else if oldVal != newVal {
			changes = append(changes, Change{Key: k, OldVal: oldVal, NewVal: newVal, Type: Updated})
		}
	}

	for k, oldVal := range previous {
		if _, exists := current[k]; !exists {
			changes = append(changes, Change{Key: k, OldVal: oldVal, Type: Removed})
		}
	}

	return changes
}

// Summary returns a human-readable summary of the changes.
func Summary(changes []Change) string {
	if len(changes) == 0 {
		return "no changes detected"
	}
	added, removed, updated := 0, 0, 0
	for _, c := range changes {
		switch c.Type {
		case Added:
			added++
		case Removed:
			removed++
		case Updated:
			updated++
		}
	}
	return fmt.Sprintf("%d added, %d updated, %d removed", added, updated, removed)
}
