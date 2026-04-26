package internal

import (
	"fmt"
	"os"
	"strconv"
)

func stringToFileMode(permStr string) (os.FileMode, error) {
	// Parse the string as Base-8 (octal), with a 32-bit size limit
	parsedInt, err := strconv.ParseUint(permStr, 8, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid permission format: %w", err)
	}

	// Cast the resulting integer to os.FileMode
	return os.FileMode(parsedInt), nil
}
