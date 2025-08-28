package log

import "github.com/awesome-goose/platform/contracts"

type logger struct {
	modifiers []contracts.Modifier
	formatter contracts.Formatter
	processor contracts.Processor
}

func NewLogger(modifiers []contracts.Modifier, formatter contracts.Formatter, processor contracts.Processor) *logger {
	return &logger{modifiers, formatter, processor}
}

func (c *logger) Write(records ...contracts.Record) {
	for _, record := range records {
		for _, modifier := range c.modifiers {
			modifier.Modify(record)
		}

		c.processor.Process(c.formatter.Format(record))
	}
}
