package modifiers

import (
	"fmt"

	"github.com/awesome-goose/contracts"
)

type ColorTagsModifier struct{}

func (m *ColorTagsModifier) Modify(record contracts.Record) contracts.Record {
	colorStart := ""
	colorEnd := "\033[0m" // Reset code

	switch record.Level {
	case contracts.DebugLogLevel:
		colorStart = "\033[36m" // Cyan
	case contracts.InfoLogLevel, contracts.NoticeLogLevel:
		colorStart = "\033[32m" // Green
	case contracts.WarningLogLevel:
		colorStart = "\033[33m" // Yellow
	case contracts.ErrorLogLevel, contracts.CriticalLogLevel:
		colorStart = "\033[31m" // Red
	case contracts.AlertLogLevel, contracts.EmergencyLogLevel:
		colorStart = "\033[1;91m" // Bright red + bold
	default:
		colorStart = ""
		colorEnd = ""
	}

	record.Message = fmt.Sprintf("%s%s%s", colorStart, record.Message, colorEnd)
	return record
}
