package main

import (
	"github.com/satori/go.uuid"
)

// GenerateID by using UUID V1
func GenerateID() (string, error) {
	id, err := uuid.NewV1()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}
