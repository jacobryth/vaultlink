package checksum

import (
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"sort"
	"strings"
)

// Level controls which checksum algorithm is used.
type Level string

const (
	None   Level = "none"
	MD5    Level = "md5"
	SHA256 Level = "sha256"
)

var validLevels = map[Level]bool{
	None:   true,
	MD5:    true,
	SHA256: true,
}

// Checksummer computes a deterministic checksum over a secrets map.
type Checksummer struct {
	level Level
}

// New returns a Checksummer for the given level.
// Returns an error if the level is unknown.
func New(level Level) (*Checksummer, error) {
	if !validLevels[level] {
		return nil, fmt.Errorf("checksum: unknown level %q (want none|md5|sha256)", level)
	}
	return &Checksummer{level: level}, nil
}

// Compute returns a hex checksum string for the provided secrets map.
// Returns an empty string when level is None.
func (c *Checksummer) Compute(secrets map[string]string) string {
	if c.level == None || len(secrets) == 0 {
		return ""
	}

	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for _, k := range keys {
		sb.WriteString(k)
		sb.WriteByte('=')
		sb.WriteString(secrets[k])
		sb.WriteByte('\n')
	}
	raw := []byte(sb.String())

	switch c.level {
	case MD5:
		sum := md5.Sum(raw)
		return fmt.Sprintf("%x", sum)
	case SHA256:
		sum := sha256.Sum256(raw)
		return fmt.Sprintf("%x", sum)
	default:
		return ""
	}
}
