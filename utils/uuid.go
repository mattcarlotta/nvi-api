package utils

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetSessionID(c *fiber.Ctx) uuid.UUID {
	return c.Locals("userSessionID").(uuid.UUID)
}

func MustParseUUID(id string) uuid.UUID {
	return uuid.MustParse(id)
}

func ParseUUID(id string) (uuid.UUID, error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("the follow id '%s' is not a valid uuid", id)
	}
	return parsedUUID, nil
}

func ParseUUIDs(ids []string) ([]uuid.UUID, error) {
	UUIDS := make([]uuid.UUID, len(ids))
	for _, value := range ids {
		parsedUUID, err := ParseUUID(value)
		if err != nil {
			return nil, err
		}
		UUIDS = append(UUIDS, parsedUUID)
	}
	return UUIDS, nil
}
