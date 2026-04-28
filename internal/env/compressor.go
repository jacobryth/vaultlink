package env

import (
	"fmt"
	"strings"
)

// CompressMode defines how blank/empty values are compressed in a secret map.
type CompressMode string

const (
	CompressModeNone   CompressMode = "none"
	CompressModeKeys   CompressMode = "keys"
	CompressModeValues CompressMode = "values"
	CompressModeBoth   CompressMode = "both"
)

var validCompressModes = map[CompressMode]bool{
	CompressModeNone:   true,
	CompressModeKeys:   true,
	CompressModeValues: true,
	CompressModeBoth:   true,
}

// EnvCompressor removes keys with blank names and/or empty values from a secret map.
type EnvCompressor struct {
	mode CompressMode
}

// NewEnvCompressor returns an EnvCompressor for the given mode.
func NewEnvCompressor(mode CompressMode) (*EnvCompressor, error) {
	if !validCompressModes[mode] {
		return nil, fmt.Errorf("env/compressor: unknown mode %q", mode)
	}
	return &EnvCompressor{mode: mode}, nil
}

// Compress filters the secrets map according to the configured mode.
func (c *EnvCompressor) Compress(secrets map[string]string) map[string]string {
	if secrets == nil {
		return nil
	}
	if c.mode == CompressModeNone {
		out := make(map[string]string, len(secrets))
		for k, v := range secrets {
			out[k] = v
		}
		return out
	}
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		blankKey := strings.TrimSpace(k) == ""
		emptyVal := strings.TrimSpace(v) == ""

		switch c.mode {
		case CompressModeKeys:
			if blankKey {
				continue
			}
		case CompressModeValues:
			if emptyVal {
				continue
			}
		case CompressModeBoth:
			if blankKey || emptyVal {
				continue
			}
		}
		out[k] = v
	}
	return out
}
