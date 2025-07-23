package logger

import (
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"strings"
	"time"
)

const (
	LevelDebug = "debug"
	LevelInfo  = "info"
	LevelWarn  = "warn"
	LevelError = "error"
	LevelFatal = "fatal"
)

func isPrintable(level string, printLevel string) bool {
	switch printLevel {
	case LevelDebug:
		return level == LevelDebug || level == LevelInfo || level == LevelWarn || level == LevelError || level == LevelFatal
	case LevelInfo:
		return level == LevelInfo || level == LevelWarn || level == LevelError || level == LevelFatal
	case LevelWarn:
		return level == LevelWarn || level == LevelError || level == LevelFatal
	case LevelError:
		return level == LevelError || level == LevelFatal
	case LevelFatal:
		return level == LevelFatal
	}
	return false
}

type F map[string]any

func (f F) Merge(others ...F) F {
	newF := make(F, len(f))
	maps.Copy(newF, f)
	for _, other := range others {
		for k, v := range other {
			if _, ok := newF[k]; !ok {
				newF[k] = v
			}
		}
	}
	return newF
}

type Logger struct {
	level  string
	format string
	output string
	kvs    F
}

func NewLogger(level string, format string, output string) *Logger {
	return &Logger{level: strings.ToLower(level), format: strings.ToLower(format), output: output}
}

func (l *Logger) With(kvs F) *Logger {
	return &Logger{level: l.level, format: l.format, output: l.output, kvs: l.kvs.Merge(kvs)}
}

func (l *Logger) Debug(msg string, kvs ...F) {
	l.print(LevelDebug, msg, kvs...)
}

func (l *Logger) Debugf(format string, args ...any) {
	l.print(LevelDebug, fmt.Sprintf(format, args...), l.kvs)
}

func (l *Logger) Info(msg string, kvs ...F) {
	l.print(LevelInfo, msg, kvs...)
}

func (l *Logger) Infof(format string, args ...any) {
	l.print(LevelInfo, fmt.Sprintf(format, args...), l.kvs)
}

func (l *Logger) Warn(msg string, kvs ...F) {
	l.print(LevelWarn, msg, kvs...)
}

func (l *Logger) Warnf(format string, args ...any) {
	l.print(LevelWarn, fmt.Sprintf(format, args...), l.kvs)
}

func (l *Logger) Error(msg string, kvs ...F) {
	l.print(LevelError, msg, kvs...)
}

func (l *Logger) Errorf(format string, args ...any) {
	l.print(LevelError, fmt.Sprintf(format, args...), l.kvs)
}

func (l *Logger) Fatal(msg string, kvs ...F) {
	l.print(LevelFatal, msg, kvs...)
}

func (l *Logger) Fatalf(format string, args ...any) {
	l.print(LevelFatal, fmt.Sprintf(format, args...), l.kvs)
}

func (l *Logger) print(level, msg string, kvs ...F) {
	if !isPrintable(level, l.level) {
		return
	}
	kv := make(F, 0)
	kv = kv.Merge(l.kvs)
	kv = kv.Merge(kvs...)
	kv["msg"] = msg
	kv["level"] = level
	kv["time"] = time.Now().Format(time.RFC3339)

	_msg := ""
	if l.format == "json" {
		b, _ := json.Marshal(kv)
		_msg = string(b)
	} else {
		_msg = fmt.Sprintf("%s %s %s", level, time.Now().Format(time.RFC3339), msg)
		for k, v := range kv {
			_msg += fmt.Sprintf(" %s=%v", k, v)
		}
		_msg += "\n"
	}
	if l.output == "stdout" {
		fmt.Println(_msg)
	} else {
		os.WriteFile(l.output, []byte(_msg), 0644)
	}
}
