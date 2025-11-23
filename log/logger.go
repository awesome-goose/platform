package log

import "github.com/awesome-goose/contracts"

type Logger struct {
	modifiers []contracts.Modifier
	formatter contracts.Formatter
	processor contracts.Processor
}

func NewLogger(modifiers []contracts.Modifier, formatter contracts.Formatter, processor contracts.Processor) *Logger {
	return &Logger{modifiers, formatter, processor}
}

func (c *Logger) Write(records ...contracts.Record) {
	for _, record := range records {
		for _, modifier := range c.modifiers {
			modifier.Modify(record)
		}

		c.processor.Process(c.formatter.Format(record))
	}
}
