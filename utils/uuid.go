package utils

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetAPIKey(c *fiber.Ctx) string {
	return c.Locals("apiKey").(string)
}

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

// strip "-" from UUID
var normalizeUUID = strings.NewReplacer("-", "")

// replaces 9 with 99, "-" with 90, and "_" with 91
var normalizer = strings.NewReplacer("9", "99", "-", "90", "_", "91")

func CreateBase64EncodedUUID() string {
	var encodedUUIDs string
	for i := 0; i < 2; i++ {
		normalizedID := normalizeUUID.Replace(uuid.NewString())
		hexID, _ := hex.DecodeString(normalizedID)
		encodedUUIDs += normalizer.Replace(base64.RawURLEncoding.EncodeToString(hexID))
	}
	return encodedUUIDs
}

// var denormalizer = strings.NewReplacer("99", "9", "90", "-", "91", "_")
// func DecodeUUID(encodedID string) string {
//  decodedID, _ := base64.RawURLEncoding.DecodeString(denormalizer.Replace(encodedID))
// 	return string(decodedID)
// }
