package modifiers

import (
	"github.com/awesome-goose/contracts"
	"github.com/awesome-goose/utils/rand"
)

type UUID struct{}

func (m *UUID) Modify(record contracts.Record) contracts.Record {
	record.Extra = append(record.Extra, map[string]any{
		"id": rand.UUID(),
	})

	return record
}
