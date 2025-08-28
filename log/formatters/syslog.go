package formatters

import (
	"fmt"
	"os"
	"strings"

	"github.com/awesome-goose/platform/contracts"
)

type Syslog struct {
	AppName string
	PID     int
}

func (f *Syslog) Format(record contracts.Record) []byte {
	if f.PID == 0 {
		f.PID = os.Getpid()
	}

	if f.AppName == "" {
		f.AppName = os.Getenv("APP_NAME")
	}

	timestamp := record.Datetime.Format("Jan 2 15:04:05") // RFC3164 style
	hostname, _ := os.Hostname()
	msg := fmt.Sprintf("%s %s %s[%d]: %s",
		timestamp,
		hostname,
		f.AppName,
		f.PID,
		strings.TrimSpace(record.Message),
	)

	return []byte(msg)
}
