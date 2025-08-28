package formatters

import (
	"fmt"
	"strings"

	"github.com/awesome-goose/platform/contracts"
)

type Line struct{}

func (f *Line) Format(record contracts.Record) []byte {
	extra := ""
	if len(record.Extra) > 0 {
		extra = fmt.Sprintf(" | extra: %v", record.Extra)
	}

	line := fmt.Sprintf(
		"%s [%s] %s: %s%s",
		record.Datetime.Format("2006-01-02 15:04:05"),
		record.Channel,
		strings.ToUpper(string(record.Level)),
		record.Message,
		extra,
	)

	return []byte(line)
}
