package shared

import (
	"strings"

	"github.com/google/uuid"
)

func ParseToRawID[T ~[16]byte](ID string) (T, error) {
	var result T
	parts := strings.Split(ID, "_")
	idStr := parts[len(parts)-1]

	parsed, err := uuid.Parse(idStr)
	if err != nil {
		return result, err
	}

	return T(parsed), nil
}
