package cmd

import (
	"log"

	"github.com/google/uuid"
)

const SlashSeparator = "/"

func mustGetUUID() string {
	u, err := uuid.NewRandom()
	if err != nil {
		log.Fatal(err)
	}
	return u.String()
}
