package env

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// BackupMode controls when backups are created.
type BackupMode string

const (
	BackupNone    BackupMode = "none"
	BackupAlways  BackupMode = "always"
	BackupOnWrite BackupMode = "on_write"
)

var validBackupModes = map[BackupMode]bool{
	BackupNone:    true,
	BackupAlways:  true,
	BackupOnWrite: true,
}

// Backup holds configuration for the env file backup manager.
type Backup struct {
	mode   BackupMode
	suffix string
}

// NewBackup creates a Backup manager with the given mode.
// suffix is appended to the backup filename (e.g. ".bak", ".backup").
func NewBackup(mode BackupMode, suffix string) (*Backup, error) {
	if !validBackupModes[mode] {
		return nil, fmt.Errorf("env/backup: unknown mode %q", mode)
	}
	if suffix == "" {
		suffix = ".bak"
	}
	return &Backup{mode: mode, suffix: suffix}, nil
}

// ShouldBackup reports whether a backup should be created given the mode.
func (b *Backup) ShouldBackup() bool {
	return b.mode != BackupNone
}

// Create copies src to a timestamped backup file next to src.
// Returns the backup path or an error.
func (b *Backup) Create(src string) (string, error) {
	if b.mode == BackupNone {
		return "", nil
	}

	data, err := os.ReadFile(src)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", fmt.Errorf("env/backup: read %s: %w", src, err)
	}

	timestamp := time.Now().UTC().Format("20060102T150405")
	dir := filepath.Dir(src)
	base := filepath.Base(src)
	dest := filepath.Join(dir, fmt.Sprintf("%s.%s%s", base, timestamp, b.suffix))

	if err := os.WriteFile(dest, data, 0o600); err != nil {
		return "", fmt.Errorf("env/backup: write %s: %w", dest, err)
	}
	return dest, nil
}
