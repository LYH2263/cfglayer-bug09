package audit

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type Logger struct {
	mu   sync.Mutex
	path string
	f    *os.File
}

func Open(dir string) (*Logger, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	path := filepath.Join(dir, "audit.log")
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return &Logger{path: path, f: f}, nil
}

func (l *Logger) Printf(format string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.f != nil {
		fmt.Fprintf(l.f, format+"\n", args...)
	}
}

func (l *Logger) Rotate() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.f == nil {
		return nil
	}
	if err := l.f.Close(); err != nil {
		return err
	}
	l.f = nil
	backup := l.path + ".1"
	if err := os.Rename(l.path, backup); err != nil && !os.IsNotExist(err) {
		return err
	}
	f, err := os.OpenFile(l.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	l.f = f
	return nil
}

func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.f != nil {
		err := l.f.Close()
		l.f = nil
		return err
	}
	return nil
}
