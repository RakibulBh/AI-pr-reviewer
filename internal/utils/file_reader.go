package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ReadRepositoryRuleFile reads a file from docs/repository_rules directory
// and converts it to plain text. Supports .md and .json files.
func ReadRepositoryRuleFile(filename string) (string, error) {
	// Construct the full path
	fullPath := filepath.Join("docs", "repository_rules", filename)

	// Check if file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", fmt.Errorf("file not found: %s", fullPath)
	}

	// Read the file content
	content, err := os.ReadFile(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to read file %s: %w", fullPath, err)
	}

	// Get file extension
	ext := strings.ToLower(filepath.Ext(filename))

	switch ext {
	case ".md":
		// For markdown files, return content as-is (already plain text)
		return string(content), nil
	case ".json":
		// For JSON files, convert to formatted plain text
		return formatJSONToText(content)
	default:
		return "", fmt.Errorf("unsupported file type: %s. Only .md and .json files are supported", ext)
	}
}

// formatJSONToText converts JSON content to a readable plain text format
func formatJSONToText(jsonContent []byte) (string, error) {
	var data interface{}

	// Parse JSON
	if err := json.Unmarshal(jsonContent, &data); err != nil {
		return "", fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Convert to formatted text
	var result strings.Builder
	formatValue(&result, data, 0)

	return result.String(), nil
}

// formatValue recursively formats JSON values into readable text
func formatValue(builder *strings.Builder, value interface{}, indent int) {
	indentStr := strings.Repeat("  ", indent)

	switch v := value.(type) {
	case map[string]interface{}:
		for key, val := range v {
			builder.WriteString(fmt.Sprintf("%s%s:\n", indentStr, key))
			formatValue(builder, val, indent+1)
		}
	case []interface{}:
		for i, val := range v {
			builder.WriteString(fmt.Sprintf("%s[%d]:\n", indentStr, i))
			formatValue(builder, val, indent+1)
		}
	case string:
		builder.WriteString(fmt.Sprintf("%s%s\n", indentStr, v))
	case float64:
		builder.WriteString(fmt.Sprintf("%s%.2f\n", indentStr, v))
	case bool:
		builder.WriteString(fmt.Sprintf("%s%t\n", indentStr, v))
	case nil:
		builder.WriteString(fmt.Sprintf("%s<nil>\n", indentStr))
	default:
		builder.WriteString(fmt.Sprintf("%s%v\n", indentStr, v))
	}
}
