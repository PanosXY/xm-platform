package helpers

import (
	"github.com/google/uuid"
)

func IsValidUUID(id string) error {
	_, err := uuid.Parse(id)

	return err
}
