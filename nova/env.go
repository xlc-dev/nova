package nova

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

// LoadDotenv loads variables from a .env file, expands them, and sets them
// as environment variables. If no path is provided, ".env" is used.
// If the specified file doesn't exist, it is silently ignored.
func LoadDotenv(paths ...string) error {
	path := ".env"
	if len(paths) > 0 {
		path = paths[0]
	}

	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Silently ignore if the file doesn't exist.
			fmt.Printf("%s not found, ignoring....\n", path)
			return nil
		}
		return fmt.Errorf("error opening file %s: %w", path, err)
	}
	defer file.Close()

	envMap := make(map[string]string)
	scanner := bufio.NewScanner(file)
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "export ") {
			line = strings.TrimSpace(strings.TrimPrefix(line, "export "))
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		if key == "" || strings.ContainsAny(key, " #'\"") {
			continue
		}

		rawValue := strings.TrimSpace(parts[1])
		var valueBuilder strings.Builder
		var state struct {
			inSingleQuotes   bool
			inDoubleQuotes   bool
			isEscaped        bool
			potentialComment bool
		}

		for i, r := range rawValue {
			if state.potentialComment && !unicode.IsSpace(r) {
				state.potentialComment = false
				valueBuilder.WriteRune('#')
			}
			if state.potentialComment && unicode.IsSpace(r) {
				break
			}
			if state.isEscaped {
				state.isEscaped = false
				switch r {
				case 'n':
					valueBuilder.WriteRune('\n')
				case 'r':
					valueBuilder.WriteRune('\r')
				case 't':
					valueBuilder.WriteRune('\t')
				case '\\', '"', '\'':
					if state.inDoubleQuotes || (!state.inSingleQuotes && !state.inDoubleQuotes) {
						valueBuilder.WriteRune(r)
					} else {
						valueBuilder.WriteRune('\\')
						valueBuilder.WriteRune(r)
					}
				case '$':
					if state.inDoubleQuotes {
						valueBuilder.WriteRune(r)
					} else {
						valueBuilder.WriteRune('\\')
						valueBuilder.WriteRune(r)
					}
				default:
					valueBuilder.WriteRune('\\')
					valueBuilder.WriteRune(r)
				}
				continue
			}
			if r == '\\' {
				state.isEscaped = true
				continue
			}
			if r == '\'' {
				if state.inDoubleQuotes {
					valueBuilder.WriteRune(r)
				} else {
					state.inSingleQuotes = !state.inSingleQuotes
				}
				continue
			}
			if r == '"' {
				if state.inSingleQuotes {
					valueBuilder.WriteRune(r)
				} else {
					state.inDoubleQuotes = !state.inDoubleQuotes
				}
				continue
			}
			if r == '#' && !state.inSingleQuotes && !state.inDoubleQuotes {
				isStart := i == 0
				var prevChar rune = ' '
				if i > 0 {
					prevChar = []rune(rawValue)[i-1]
				}
				if isStart || unicode.IsSpace(prevChar) {
					state.potentialComment = true
					continue
				}
			}
			valueBuilder.WriteRune(r)
		}
		if state.inSingleQuotes || state.inDoubleQuotes {
			// Skip lines with unclosed quotes.
			continue
		}
		if state.isEscaped {
			valueBuilder.WriteRune('\\')
		}
		envMap[key] = valueBuilder.String()
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error scanning file %s: %w", path, err)
	}

	expandedMap := make(map[string]string, len(envMap))
	expanding := make(map[string]bool)

	var expandValue func(string) string
	expandValue = func(key string) string {
		value, ok := envMap[key]
		if !ok {
			return ""
		}
		if expanding[key] {
			return ""
		}
		expanding[key] = true
		expanded := os.Expand(value, func(varName string) string {
			if _, exists := envMap[varName]; exists {
				return expandValue(varName)
			}
			return os.Getenv(varName)
		})
		delete(expanding, key)
		expandedMap[key] = expanded
		return expanded
	}

	for key := range envMap {
		if _, ok := expandedMap[key]; !ok {
			expandValue(key)
		}
	}

	for key, value := range expandedMap {
		if _, exists := os.LookupEnv(key); !exists {
			if err := os.Setenv(key, value); err != nil {
				return fmt.Errorf("failed to set environment variable %s: %w", key, err)
			}
		}
	}

	return nil
}
