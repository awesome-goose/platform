package contracts

import "time"

type LogLevel string

var (
	DebugLogLevel     LogLevel = "debug"
	InfoLogLevel      LogLevel = "info"
	NoticeLogLevel    LogLevel = "notice"
	WarningLogLevel   LogLevel = "warning"
	ErrorLogLevel     LogLevel = "error"
	CriticalLogLevel  LogLevel = "critical"
	AlertLogLevel     LogLevel = "alert"
	EmergencyLogLevel LogLevel = "emergency"
)

type Record struct {
	Datetime time.Time `json:"time"`
	Channel  string    `json:"channel"`
	Level    LogLevel  `json:"level"`
	Message  string    `json:"message"`
	Extra    []any     `json:"extra"`
}

type Modifier interface {
	Modify(record Record) Record
}

type Formatter interface {
	Format(record Record) []byte
}

type Processor interface {
	Process(record []byte)
}
