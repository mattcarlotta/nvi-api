package utils

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

func ParseUUIDs(ids []string) ([]uuid.UUID, error) {
	uuids := make([]uuid.UUID, len(ids))
	for _, value := range ids {
		parsedId, err := uuid.Parse(value)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("The follow id '%s' is not a valid uuid!", value))
		} else {
			uuids = append(uuids, parsedId)
		}
	}

	return uuids, nil
}
