package log

import "github.com/awesome-goose/contracts"

var (
	Log = defaultLog()
)

func defaultLog() *log {
	defaultChannel := "noop"
	defaultLoggers := []*logger{
		NewLogger(
			[]contracts.Modifier{
				NewNoopModifier(),
			},
			NewNoopFormatter(),
			NewNoopProcessor(),
		),
	}

	return newLog(
		defaultChannel,
		defaultLoggers...,
	)
}
