package modifiers

import (
	"os"
	"runtime"

	"github.com/awesome-goose/contracts"
)

type SystemInfo struct{}

func (m *SystemInfo) Modify(record contracts.Record) contracts.Record {
	hostname, _ := os.Hostname()

	info := map[string]any{
		"go_version": runtime.Version(),
		"go_os":      runtime.GOOS,
		"go_arch":    runtime.GOARCH,
		"num_cpu":    runtime.NumCPU(),
		"hostname":   hostname,
		"pid":        os.Getpid(),
	}

	record.Extra = append(record.Extra, map[string]any{
		"system-info": info,
	})

	return record
}
