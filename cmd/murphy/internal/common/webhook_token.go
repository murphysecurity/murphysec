package common

import (
	"fmt"
	"strings"

	"github.com/murphysecurity/murphysec/api"
)

// ParseWebhookToken parses webhook token strings in key=value format
// and returns a slice of NoticeApiHeadersArray
func ParseWebhookToken(tokens []string) ([]api.NoticeApiHeadersArray, error) {
	var result []api.NoticeApiHeadersArray
	for _, token := range tokens {
		parts := strings.SplitN(token, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid webhook token format: %s (expected key=value)", token)
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		if key == "" {
			return nil, fmt.Errorf("webhook token key cannot be empty: %s", token)
		}
		if value == "" {
			return nil, fmt.Errorf("webhook token value cannot be empty: %s", token)
		}
		result = append(result, api.NoticeApiHeadersArray{
			Key:   key,
			Value: value,
		})
	}
	return result, nil
}
