package log

import "github.com/awesome-goose/platform/contracts"

type NoopModifier struct{}

func NewNoopModifier() *NoopModifier {
	return &NoopModifier{}
}

func (m *NoopModifier) Modify(record contracts.Record) contracts.Record {
	return record
}

type NoopFormatter struct{}

func NewNoopFormatter() *NoopFormatter {
	return &NoopFormatter{}
}

func (f *NoopFormatter) Format(record contracts.Record) []byte {
	return []byte{}
}

type NoopProcessor struct{}

func NewNoopProcessor() *NoopProcessor {
	return &NoopProcessor{}
}

func (p *NoopProcessor) Process(record []byte) {}
