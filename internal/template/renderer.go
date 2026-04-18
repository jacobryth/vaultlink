package template

import (
	"fmt"
	"os"
	"strings"
)

// Renderer renders secrets into a user-provided template string,
// replacing {{KEY}} placeholders with their corresponding values.
type Renderer struct {
	strict bool
}

// New creates a Renderer. When strict is true, missing keys return an error.
func New(strict bool) *Renderer {
	return &Renderer{strict: strict}
}

// Render replaces all {{KEY}} placeholders in tmpl with values from secrets.
func (r *Renderer) Render(tmpl string, secrets map[string]string) (string, error) {
	result := tmpl
	for key, val := range secrets {
		placeholder := "{{" + key + "}}"
		result = strings.ReplaceAll(result, placeholder, val)
	}
	if r.strict {
		if start := strings.Index(result, "{{"); start != -1 {
			end := strings.Index(result[start:], "}}")
			if end != -1 {
				missing := result[start+2 : start+end]
				return "", fmt.Errorf("template: unresolved placeholder: %s", missing)
			}
		}
	}
	return result, nil
}

// RenderFile reads a template file and renders it with the given secrets.
func (r *Renderer) RenderFile(path string, secrets map[string]string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("template: read file: %w", err)
	}
	return r.Render(string(data), secrets)
}
