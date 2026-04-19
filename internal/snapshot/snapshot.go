package snapshot

import (
	"encoding/json"
	"os"
	"time"
)

// Snapshot represents a point-in-time record of synced secrets metadata.
type Snapshot struct {
	Timestamp  time.Time         `json:"timestamp"`
	SecretPath string            `json:"secret_path"`
	Keys       []string          `json:"keys"`
	Checksum   map[string]string `json:"checksum"`
}

// Manager handles reading and writing snapshots to disk.
type Manager struct {
	filePath string
}

// NewManager creates a new snapshot Manager.
func NewManager(filePath string) *Manager {
	return &Manager{filePath: filePath}
}

// Save writes the snapshot to disk as JSON.
func (m *Manager) Save(s *Snapshot) error {
	s.Timestamp = time.Now().UTC()
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(m.filePath, data, 0600)
}

// Load reads and parses the snapshot from disk.
// Returns nil, nil if no snapshot file exists yet.
func (m *Manager) Load() (*Snapshot, error) {
	data, err := os.ReadFile(m.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var s Snapshot
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}
	return &s, nil
}

// HasChanged returns true if the provided keys differ from the snapshot.
func (m *Manager) HasChanged(current map[string]string) (bool, error) {
	prev, err := m.Load()
	if err != nil {
		return false, err
	}
	if prev == nil {
		return true, nil
	}
	for k, v := range current {
		if prev.Checksum[k] != v {
			return true, nil
		}
	}
	return len(current) != len(prev.Checksum), nil
}

// Delete removes the snapshot file from disk.
// Returns nil if the file does not exist.
func (m *Manager) Delete() error {
	err := os.Remove(m.filePath)
	if err != nil && os.IsNotExist(err) {
		return nil
	}
	return err
}
