package log

import (
	"time"

	"github.com/awesome-goose/platform/contracts"
)

type log struct {
	channels map[string][]*logger
	channel  string
}

func newLog(channel string, loggers ...*logger) *log {
	l := &log{}
	l.channels = make(map[string][]*logger)
	l.channels[channel] = loggers

	return l.Use(channel)
}

func (l *log) Use(channel string) *log {
	log := *l
	log.channel = channel
	return &log
}

func (l *log) Add(channel string, loggers ...*logger) *log {
	l.channels[channel] = append(l.channels[channel], loggers...)
	return l.Use(channel)
}

func (l *log) Debug(message string, extra ...any) {
	l.log(contracts.Record{
		Datetime: time.Now(),
		Channel:  l.channel,
		Level:    contracts.DebugLogLevel,
		Message:  message,
		Extra:    extra,
	})
}

func (l *log) Info(message string, extra ...any) {
	l.log(contracts.Record{
		Datetime: time.Now(),
		Channel:  l.channel,
		Level:    contracts.InfoLogLevel,
		Message:  message,
		Extra:    extra,
	})
}

func (l *log) Notice(message string, extra ...any) {
	l.log(contracts.Record{
		Datetime: time.Now(),
		Channel:  l.channel,
		Level:    contracts.NoticeLogLevel,
		Message:  message,
		Extra:    extra,
	})
}

func (l *log) Warning(message string, extra ...any) {
	l.log(contracts.Record{
		Datetime: time.Now(),
		Channel:  l.channel,
		Level:    contracts.WarningLogLevel,
		Message:  message,
		Extra:    extra,
	})
}

func (l *log) Error(message string, extra ...any) {
	l.log(contracts.Record{
		Datetime: time.Now(),
		Channel:  l.channel,
		Level:    contracts.ErrorLogLevel,
		Message:  message,
		Extra:    extra,
	})
}

func (l *log) Critical(message string, extra ...any) {
	l.log(contracts.Record{
		Datetime: time.Now(),
		Channel:  l.channel,
		Level:    contracts.CriticalLogLevel,
		Message:  message,
		Extra:    extra,
	})
}

func (l *log) Alert(message string, extra ...any) {
	l.log(contracts.Record{
		Datetime: time.Now(),
		Channel:  l.channel,
		Level:    contracts.AlertLogLevel,
		Message:  message,
		Extra:    extra,
	})
}

func (l *log) Emergency(message string, extra ...any) {
	l.log(contracts.Record{
		Datetime: time.Now(),
		Channel:  l.channel,
		Level:    contracts.EmergencyLogLevel,
		Message:  message,
		Extra:    extra,
	})
}

func (l *log) log(records ...contracts.Record) {
	loggers, ok := l.channels[l.channel]
	if !ok {
		panic("No loggers found for channel: " + l.channel)
	}

	for _, logger := range loggers {
		logger.Write(records...)
	}
}
