package helper

import (
	uuid "github.com/satori/go.uuid"
)

func GetUUID(input string) uuid.UUID {
	id, err := uuid.FromString(input)
	if err != nil {
		return id
	}
	return id
}

func IsValidUUID(u string) bool {
	_, err := uuid.FromString(u)
	return err == nil
}
