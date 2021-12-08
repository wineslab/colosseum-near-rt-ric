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

//  This source code is part of the near-RT RIC (RAN Intelligent Controller)
//  platform project (RICP).


package logger

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"testing"
)

func TestInitDebugLoggerSuccess(t *testing.T) {
	log, err := InitLogger(DebugLevel)
	assert.Nil(t, err)
	assert.NotNil(t, log)
	assert.True(t, log.Logger.Core().Enabled(zap.DebugLevel))
}

func TestInitInfoLoggerSuccess(t *testing.T) {
	log, err := InitLogger(InfoLevel)
	assert.Nil(t, err)
	assert.NotNil(t, log)
	assert.True(t, log.Logger.Core().Enabled(zap.InfoLevel))
}

func TestInitWarnLoggerSuccess(t *testing.T) {
	log, err := InitLogger(WarnLevel)
	assert.Nil(t, err)
	assert.NotNil(t, log)
	assert.True(t, log.Logger.Core().Enabled(zap.WarnLevel))
}

func TestInitErrorLoggerSuccess(t *testing.T) {
	log, err := InitLogger(ErrorLevel)
	assert.Nil(t, err)
	assert.NotNil(t, log)
	assert.True(t, log.Logger.Core().Enabled(zap.ErrorLevel))
}

func TestInitDPanicLoggerSuccess(t *testing.T) {
	log, err := InitLogger(DPanicLevel)
	assert.Nil(t, err)
	assert.NotNil(t, log)
	assert.True(t, log.Logger.Core().Enabled(zap.DPanicLevel))
}

func TestInitPanicLoggerSuccess(t *testing.T) {
	log, err := InitLogger(PanicLevel)
	assert.Nil(t, err)
	assert.NotNil(t, log)
	assert.True(t, log.Logger.Core().Enabled(zap.PanicLevel))
}

func TestInitInfoLoggerFailure(t *testing.T) {
	log, err := InitLogger(99)
	assert.NotNil(t, err)
	assert.Nil(t, log)
}

func TestSyncSuccess(t *testing.T){
	logFile, err := os.Create("./loggerTest.txt")
	if err != nil{
		t.Errorf("logger_test.TestSyncSuccess - failed to create file, error: %s", err)
	}
	old := os.Stdout
	os.Stdout = logFile
	log, err := InitLogger(DebugLevel)
	if err != nil {
		t.Errorf("logger_test.TestSyncSuccess - failed to initialize logger, error: %s", err)
	}
	err = log.Sync()
	assert.Nil(t, err)

	os.Stdout = old
	logFile, err = os.Open("./loggerTest.txt")
	if err != nil{
		t.Errorf("logger_test.TestSyncSuccess - failed to open file, error: %s", err)
	}
	var buf bytes.Buffer
	_, err = io.Copy(&buf, logFile)
	if err != nil {
		t.Errorf("logger_test.TestSyncSuccess - failed to copy bytes, error: %s", err)
	}
	debugRecord,_ :=buf.ReadString('\n')
	errorRecord,_ :=buf.ReadString('\n')

	assert.NotEmpty(t, debugRecord)
	assert.Empty(t, errorRecord)
	err = os.Remove("./loggerTest.txt")
	if err != nil {
		t.Errorf("logger_test.TestSyncSuccess - failed to remove file, error: %s", err)
	}

}

func TestSyncFailure(t *testing.T){
	log, err := InitLogger(DebugLevel)
	err = log.Sync()
	assert.NotNil(t, err)
}

func TestDebugEnabledFalse(t *testing.T){
	entryNum, log := countRecords(InfoLevel, t)
	assert.False(t, log.DebugEnabled())
	assert.Equal(t,3, entryNum)
}

func TestDebugEnabledTrue(t *testing.T){
	entryNum, log := countRecords(DebugLevel, t)
	assert.True(t, log.DebugEnabled())
	assert.Equal(t,4, entryNum)
}

func TestDPanicfDebugLevel(t *testing.T){
	assert.True(t,validateRecordExists(DebugLevel, zap.DPanicLevel, t))
}

func TestDPanicfInfoLevel(t *testing.T){
	assert.True(t,validateRecordExists(InfoLevel, zap.DPanicLevel, t))
}

func TestErrorfDebugLevel(t *testing.T)  {
	assert.True(t,validateRecordExists(DebugLevel, zap.ErrorLevel, t))
}

func TestErrorfInfoLevel(t *testing.T)  {
	assert.True(t,validateRecordExists(InfoLevel, zap.ErrorLevel, t))
}

func TestInfofDebugLevel(t *testing.T)  {
	assert.True(t,validateRecordExists(DebugLevel, zap.InfoLevel, t))
}

func TestInfofInfoLevel(t *testing.T)  {
	assert.True(t,validateRecordExists(InfoLevel, zap.InfoLevel, t))
}

func TestDebugfDebugLevel(t *testing.T)  {
	assert.True(t,validateRecordExists(DebugLevel, zap.DebugLevel, t))
}

func TestDebugfInfoLevel(t *testing.T)  {
	assert.False(t,validateRecordExists(InfoLevel, zap.DebugLevel, t))
}

func TestInfofFatalLevel(t *testing.T)  {
	assert.False(t,validateRecordExists(FatalLevel, zap.InfoLevel, t))
}

func TestDebugfFatalLevel(t *testing.T)  {
	assert.False(t,validateRecordExists(FatalLevel, zap.DebugLevel, t))
}

func TestWarnfWarnLevel(t *testing.T)  {
	assert.True(t,validateRecordExists(WarnLevel, zap.WarnLevel, t))
}

func TestWarnfDebugLevel(t *testing.T)  {
	assert.True(t,validateRecordExists(DebugLevel, zap.WarnLevel, t))
}

func TestWarnfInfoLevel(t *testing.T)  {
	assert.True(t,validateRecordExists(InfoLevel, zap.WarnLevel, t))
}

func TestWarnfFatalLevel(t *testing.T)  {
	assert.False(t,validateRecordExists(FatalLevel, zap.WarnLevel, t))
}

func TestLogLevelTokenToLevel(t *testing.T) {
	level, ok := LogLevelTokenToLevel("deBug")
	assert.True(t, ok)
	assert.True(t, level == DebugLevel)

	level, ok = LogLevelTokenToLevel("infO")
	assert.True(t, ok)
	assert.True(t, level == InfoLevel)

	level, ok = LogLevelTokenToLevel("Warn")
	assert.True(t, ok)
	assert.True(t, level == WarnLevel)

	level, ok = LogLevelTokenToLevel("eRror")
	assert.True(t, ok)
	assert.True(t, level == ErrorLevel)

	level, ok = LogLevelTokenToLevel("Dpanic ")
	assert.True(t, ok)
	assert.True(t, level == DPanicLevel)

	level, ok = LogLevelTokenToLevel("    panic ")
	assert.True(t, ok)
	assert.True(t, level == PanicLevel)

	level, ok = LogLevelTokenToLevel("fatal")
	assert.True(t, ok)
	assert.True(t, level == FatalLevel)

	level, ok = LogLevelTokenToLevel("zzz")
	assert.False(t, ok)
	assert.True(t, level > FatalLevel)

}
func countRecords(logLevel LogLevel, t *testing.T) (int, *Logger){
	old := os.Stdout
	r, w, _ :=os.Pipe()
	os.Stdout = w
	log, err := InitLogger(logLevel)
	if err != nil {
		t.Errorf("logger_test.TestSyncFailure - failed to initialize logger, error: %s", err)
	}
	log.Infof("%v, %v, %v", 1, "abc", 0.1)
	log.Debugf("%v, %v, %v", 1, "abc", 0.1)
	log.Errorf("%v, %v, %v", 1, "abc", 0.1)
	log.DPanicf("%v, %v, %v", 1, "abc", 0.1)
	err = w.Close()
	if err != nil {
		t.Errorf("logger_test.TestSyncFailure - failed to close writer, error: %s", err)
	}
	os.Stdout = old
	var buf bytes.Buffer
	_, err = io.Copy(&buf, r)
	if err != nil {
		t.Errorf("logger_test.TestSyncFailure - failed to copy bytes, error: %s", err)
	}
	entryNum := 0
	s,_:= buf.ReadString('\n')
	for len(s) > 0{
		entryNum +=1
		s,_= buf.ReadString('\n')
	}
	return entryNum, log
}

func validateRecordExists(logLevel LogLevel, recordLevel zapcore.Level, t *testing.T) bool {
	old := os.Stdout
	r, w, _ :=os.Pipe()
	os.Stdout = w
	log, err := InitLogger(logLevel)
	if err != nil {
		t.Errorf("logger_test.TestSyncFailure - failed to initialize logger, error: %s", err)
	}
	switch recordLevel{
	case  zap.DebugLevel:
		log.Debugf("%v, %v, %v", 1, "abc", 0.1)
	case zap.InfoLevel:
		log.Infof("%v, %v, %v", 1, "abc", 0.1)
	case zap.WarnLevel:
		log.Warnf("%v, %v, %v", 1, "abc", 0.1)
	case zap.ErrorLevel:
		log.Errorf("%v, %v, %v", 1, "abc", 0.1)
	case zap.DPanicLevel:
		log.DPanicf("%v, %v, %v", 1, "abc", 0.1)
	}
	err = w.Close()
	if err != nil {
		t.Errorf("logger_test.TestSyncFailure - failed to close writer, error: %s", err)
	}
	os.Stdout = old
	var buf bytes.Buffer
	_, err = io.Copy(&buf, r)
	if err != nil {
		t.Errorf("logger_test.TestSyncFailure - failed to copy bytes, error: %s", err)
	}
	entryNum := 0
	s,_:= buf.ReadString('\n')
	for len(s) > 0{
		entryNum +=1
		s,_= buf.ReadString('\n')
	}
	return entryNum == 1
}