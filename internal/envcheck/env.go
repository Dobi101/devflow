package envcheck

import (
	"fmt"
	"os"
	"strings"
)

func Parse(env string) map[string]string {
	lines := strings.Split(env, "\n")
	result := make(map[string]string)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		result[key] = value
	}
	return result
}

func Load(path string) (map[string]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read env file: %w", err)
	}

	return Parse(string(data)), nil
}

func MissingRequired(env map[string]string, required []string) []string {
	var missing []string

	for _, key := range required {
		key = strings.TrimSpace(key)
		if key == "" {
			continue
		}

		value, ok := env[key]
		if !ok || strings.TrimSpace(value) == "" {
			missing = append(missing, key)
		}
	}

	return missing
}
