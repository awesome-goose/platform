package modifiers

import (
	"runtime/debug"

	"github.com/awesome-goose/contracts"
)

type StackTrace struct{}

func (m *StackTrace) Modify(record contracts.Record) contracts.Record {
	trace := debug.Stack()
	record.Extra = append(record.Extra, map[string]any{
		"stack": string(trace),
	})
	return record
}
