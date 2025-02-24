package identifier

import (
	"github.com/google/uuid"
)

type Identifier string

func (i Identifier) Validate() (string, error) {
	if err := uuid.Validate(string(i)); err != nil {
		return "", err
	}
	return string(i), nil
}
