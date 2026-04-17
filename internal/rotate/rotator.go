package rotate

import (
	"fmt"
	"time"
)

// Record holds metadata about a secret rotation event.
type Record struct {
	Key       string
	RotatedAt time.Time
	PrevHash  string
	NewHash   string
}

// Rotator detects which secrets have changed and records rotation events.
type Rotator struct {
	hasher func(string) string
}

// New returns a new Rotator.
func New() *Rotator {
	return &Rotator{hasher: defaultHash}
}

// Detect compares previous and current secret maps and returns rotation records
// for any keys whose values have changed or are newly added.
func (r *Rotator) Detect(prev, curr map[string]string) []Record {
	var records []Record
	for k, newVal := range curr {
		newHash := r.hasher(newVal)
		if oldVal, ok := prev[k]; ok {
			oldHash := r.hasher(oldVal)
			if oldHash != newHash {
				records = append(records, Record{
					Key:       k,
					RotatedAt: time.Now().UTC(),
					PrevHash:  oldHash,
					NewHash:   newHash,
				})
			}
		} else {
			records = append(records, Record{
				Key:       k,
				RotatedAt: time.Now().UTC(),
				PrevHash:  "",
				NewHash:   newHash,
			})
		}
	}
	return records
}

// Summary returns a human-readable summary of rotation records.
func Summary(records []Record) string {
	if len(records) == 0 {
		return "no rotations detected"
	}
	return fmt.Sprintf("%d secret(s) rotated", len(records))
}

func defaultHash(val string) string {
	h := uint32(2166136261)
	for i := 0; i < len(val); i++ {
		h ^= uint32(val[i])
		h *= 16777619
	}
	return fmt.Sprintf("%08x", h)
}
