package processors

import (
	"errors"
	"log/syslog"
	"runtime"
	"strings"
	"syscall"
)

type Syslog struct {
	writer *syslog.Writer
}

func NewSyslog(tag string) (*Syslog, error) {
	// 1. OS compatibility check
	if runtime.GOOS == "windows" {
		return nil, errors.New("syslog is not supported on Windows")
	}

	// 2. Permissions check: try writing a test line
	writer, err := syslog.New(syslog.LOG_INFO|syslog.LOG_USER, tag)
	if err != nil {
		// Try to distinguish permission errors
		if errno, ok := err.(syscall.Errno); ok && errno == syscall.EPERM {
			return nil, errors.New("permission denied: cannot write to syslog")
		}
		return nil, err
	}

	// 3. Test write
	testMsg := "[syslog-test] checking permissions"
	if testErr := writer.Info(testMsg); testErr != nil {
		writer.Close()
		return nil, errors.New("unable to write to syslog: " + testErr.Error())
	}

	return &Syslog{writer: writer}, nil
}

func (p *Syslog) Process(record []byte) {
	if p.writer == nil {
		return
	}

	msg := string(record)
	lower := strings.ToLower(msg)

	switch {
	case strings.Contains(lower, "debug"):
		_ = p.writer.Debug(msg)
	case strings.Contains(lower, "info"):
		_ = p.writer.Info(msg)
	case strings.Contains(lower, "notice"):
		_ = p.writer.Notice(msg)
	case strings.Contains(lower, "warning"):
		_ = p.writer.Warning(msg)
	case strings.Contains(lower, "error"):
		_ = p.writer.Err(msg)
	case strings.Contains(lower, "critical"):
		_ = p.writer.Crit(msg)
	case strings.Contains(lower, "alert"):
		_ = p.writer.Alert(msg)
	case strings.Contains(lower, "emergency"):
		_ = p.writer.Emerg(msg)
	default:
		_ = p.writer.Info(msg)
	}
}
