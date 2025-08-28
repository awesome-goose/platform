package modifiers

import (
	"github.com/awesome-goose/platform/contracts"
	"github.com/awesome-goose/platform/utils/rand"
)

type UUID struct{}

func (m *UUID) Modify(record contracts.Record) contracts.Record {
	record.Extra = append(record.Extra, map[string]any{
		"id": rand.UUID(),
	})

	return record
}
