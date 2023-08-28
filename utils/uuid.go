package utils

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetSessionId(c *fiber.Ctx) uuid.UUID {
	return c.Locals("userSessionId").(uuid.UUID)
}

func MustParseUUID(id string) uuid.UUID {
	return uuid.MustParse(id)
}

func ParseUUID(id string) (uuid.UUID, error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("The follow id '%s' is not a valid uuid!", id)
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
