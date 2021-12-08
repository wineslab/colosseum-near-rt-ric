//
// Copyright 2019 AT&T Intellectual Property
// Copyright 2019 Nokia
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

//  This source code is part of the near-RT RIC (RAN Intelligent Controller)
//  platform project (RICP).

package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strings"
	"time"
)

type Logger struct {
	Logger     *zap.Logger
}

// Copied from zap logger
//
// A Level is a logging priority. Higher levels are more important.
type LogLevel int8

const (
	// DebugLevel logs are typically voluminous, and are usually disabled in
	// production.
	DebugLevel LogLevel = iota - 1
	// InfoLevel is the default logging priority.
	InfoLevel
	// WarnLevel logs are more important than Info, but don't need individual
	// human review.
	WarnLevel
	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	ErrorLevel
	// DPanicLevel logs are particularly important errors. In development the
	// logger panics after writing the message.
	DPanicLevel
	// PanicLevel logs a message, then panics.
	PanicLevel
	// FatalLevel logs a message, then calls os.Exit(1).
	FatalLevel

	_minLevel = DebugLevel
	_maxLevel = FatalLevel
)

var logLevelTokenToLevel = map[string] LogLevel {
	"debug" : DebugLevel,
	"info": InfoLevel,
	"warn": WarnLevel,
	"error": ErrorLevel,
	"dpanic": DPanicLevel,
	"panic": PanicLevel,
	"fatal": FatalLevel,
}

func LogLevelTokenToLevel(level string) (LogLevel, bool) {
	if level, ok := logLevelTokenToLevel[strings.TrimSpace(strings.ToLower(level))];ok {
		return level, true
	}
	return _maxLevel+1, false
}

func InitLogger(requested LogLevel) (*Logger, error) {
	var logger *zap.Logger
	var err error
	switch requested {
	case DebugLevel:
		logger, err = initLoggerByLevel(zapcore.DebugLevel)
	case InfoLevel:
		logger, err = initLoggerByLevel(zapcore.InfoLevel)
	case WarnLevel:
		logger, err = initLoggerByLevel(zapcore.WarnLevel)
	case ErrorLevel:
		logger, err = initLoggerByLevel(zapcore.ErrorLevel)
	case DPanicLevel:
		logger, err = initLoggerByLevel(zapcore.DPanicLevel)
	case PanicLevel:
		logger, err = initLoggerByLevel(zapcore.PanicLevel)
	case FatalLevel:
		logger, err = initLoggerByLevel(zapcore.FatalLevel)
	default:
		err = fmt.Errorf("Invalid logging Level :%d",requested)
	}
	if err != nil {
		return nil, err
	}
	return &Logger{Logger:logger}, nil

}
func(l *Logger)Sync() error {
	l.Debugf("#logger.Sync - Going to flush buffered log")
	return l.Logger.Sync()
}

func (l *Logger)Infof(formatMsg string, a ...interface{})  {
	if l.InfoEnabled() {
		msg := fmt.Sprintf(formatMsg, a...)
		l.Logger.Info(msg, zap.Any("mdc", l.getTimeStampMdc()))
	}
}

func (l *Logger)Debugf(formatMsg string, a ...interface{})  {
	if l.DebugEnabled(){
		msg := fmt.Sprintf(formatMsg, a...)
		l.Logger.Debug(msg, zap.Any("mdc", l.getTimeStampMdc()))
	}
}

func (l *Logger)Errorf(formatMsg string, a ...interface{})  {
	msg := fmt.Sprintf(formatMsg, a...)
	l.Logger.Error(msg, zap.Any("mdc", l.getTimeStampMdc()))
}

func (l *Logger)Warnf(formatMsg string, a ...interface{})  {
	msg := fmt.Sprintf(formatMsg, a...)
	l.Logger.Warn(msg, zap.Any("mdc", l.getTimeStampMdc()))
}

func (l *Logger) getTimeStampMdc() map[string]string {
	timeStr := time.Now().Format("2006-01-02 15:04:05.000")
	mdc := map[string]string{"time": timeStr}
	return mdc
}

func (l *Logger)InfoEnabled()bool{
	return l.Logger.Core().Enabled(zap.InfoLevel)
}

func (l *Logger)DebugEnabled()bool{
	return l.Logger.Core().Enabled(zap.DebugLevel)
}

func (l *Logger)DPanicf(formatMsg string, a ...interface{})  {
	msg := fmt.Sprintf(formatMsg, a...)
	l.Logger.DPanic(msg, zap.Any("mdc", l.getTimeStampMdc()))
}

func initLoggerByLevel(l zapcore.Level) (*zap.Logger, error) {
	cfg := zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(l),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "msg",

			LevelKey:    "crit",
			EncodeLevel: zapcore.CapitalLevelEncoder,

			TimeKey:    "ts",
			EncodeTime: epochMillisIntegerTimeEncoder,

			CallerKey: "id",
			EncodeCaller: xAppMockCallerEncoder,
		},
	}
	return cfg.Build()
}

func xAppMockCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("xAppMock")
}

func epochMillisIntegerTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	nanos := t.UnixNano()
	millis := int64(nanos) / int64(time.Millisecond)
	enc.AppendInt64(millis)
}

