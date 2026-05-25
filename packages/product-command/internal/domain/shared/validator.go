package shared

import "strings"

type Validate interface {
	~string
	IsValid() bool
}

func ValidateEnum[T Validate](input string, err error) (T, error) {
	value := T(strings.ToUpper(strings.TrimSpace(input)))
	if !value.IsValid() {
		return "", err
	}

	return value, nil
}
