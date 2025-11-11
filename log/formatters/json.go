package formatters

import (
	"encoding/json"
	"time"

	"github.com/awesome-goose/contracts"
)

type JSON struct{}

func (j *JSON) Format(record contracts.Record) []byte {
	b, err := json.Marshal(record)
	if err != nil {
		fallback := map[string]any{
			"time":    time.Now(),
			"channel": "system",
			"level":   contracts.ErrorLogLevel,
			"message": "log/formatter/json: failed to marshal log record",
		}

		b, _ = json.Marshal(fallback)
	}

	return b
}
