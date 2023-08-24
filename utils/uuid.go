package utils

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetSessionId(c *fiber.Ctx) uuid.UUID {
	return c.Locals("userSessionId").(uuid.UUID)
}

func ParseUUID(id string) (uuid.UUID, error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return uuid.UUID{}, errors.New(fmt.Sprintf("The follow id '%s' is not a valid uuid!", id))
	}
	return parsedUUID, nil
}

func ParseUUIDs(ids []string) ([]uuid.UUID, error) {
	var UUIDS []uuid.UUID
	for _, value := range ids {
		parsedUUID, err := ParseUUID(value)
		if err != nil {
			return nil, err
		}
		UUIDS = append(UUIDS, parsedUUID)
	}
	return UUIDS, nil
}
