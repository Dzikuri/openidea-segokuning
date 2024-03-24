package helper

import (
	"fmt"
	"strings"

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

func PrepareQueryToString(query string, args ...interface{}) string {
	for i, arg := range args {
		placeholder := fmt.Sprintf("$%d", i+1)
		query = strings.Replace(query, placeholder, fmt.Sprintf("%v", arg), -1)
	}
	return query
}
