package utils

import (
	"errors"

	"github.com/google/uuid"
)

func ParseUUIDs(ids []string, idList *[]uuid.UUID) error {
	for _, value := range ids {
		parsedId, err := uuid.Parse(value)
		if err != nil {
			return errors.New("You must provide a valid environment id!")
		} else {
			*idList = append(*idList, parsedId)
		}
	}

	return nil
}
