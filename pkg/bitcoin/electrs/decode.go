package electrs

import (
	"encoding/json"
	"fmt"
)

func decodeJSON[K any](input []byte) (K, error) {
	var result K

	if err := json.Unmarshal(input, &result); err != nil {
		return result, fmt.Errorf("failed to decode: [%w]", err)
	}

	return result, nil
}
