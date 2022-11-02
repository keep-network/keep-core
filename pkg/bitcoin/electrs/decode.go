package electrs

import (
	"encoding/json"
	"fmt"
	"io"
)

func decodeJSON[K any](r io.ReadCloser) (K, error) {
	var result K

	if err := json.NewDecoder(r).Decode(&result); err != nil {
		return result, fmt.Errorf("failed to decode: [%w]", err)
	}

	return result, nil
}
