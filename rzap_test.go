package rzap

import (
	"encoding/json"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"io/ioutil"
	"os"
	"testing"
)

type line struct {
	Status  int    `json:"status"`
	Level   string `json:"level"`
	TS      string `json:"ts"`
	Message string `json:"msg"`
}

func expectRemoveFile(filePath string, t *testing.T) {
	if _, err := os.Stat(filePath); err == nil {
		if err := os.Remove(filePath); err != nil {
			t.Errorf("Remove file: %s failed", filePath)
		}
	}
}

func expectExistFile(filePath string, t *testing.T) {
	if _, err := os.Stat(filePath); err != nil {
		t.Errorf("File: %s not exist", filePath)
	}
}

func expectReadFile(filePath string, t *testing.T) *line {
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.Errorf("Read file: %s failed", filePath)
		return nil
	}

	l := &line{}
	if err := json.Unmarshal(b, l); err != nil {
		t.Errorf("Unmarshal line: %s failed", string(b))
		return nil
	}
	return l
}

func TestNewLogger(t *testing.T) {
	logFile := "/tmp/rzap.log"
	expectRemoveFile(logFile, t)

	core := NewCore(&lumberjack.Logger{
		Filename: "/tmp/rzap.log",
	}, zap.InfoLevel)
	NewGlobalLogger([]zapcore.Core{core})

	zap.L().Info("info message", zap.Int("status", 0))
	expectExistFile(logFile, t)

	l := expectReadFile(logFile, t)
	if l == nil {
		return
	}

	if l.Status != 0 {
		t.Error("Log 'status' error")
	}
	if l.Level != "INFO" {
		t.Error("Log 'level' error")
	}
	if l.Message != "info message" {
		t.Error("Log 'msg' error")
	}
}

func TestNewTeeLogger(t *testing.T) {
	infoLogFile := "/tmp/rzap_info.log"
	errorLogFile := "/tmp/rzap_error.log"

	expectRemoveFile(infoLogFile, t)
	expectRemoveFile(errorLogFile, t)

	infoCore := NewCore(&lumberjack.Logger{
		Filename: infoLogFile,
	}, zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level <= zap.InfoLevel
	}))
	errorCore := NewCore(&lumberjack.Logger{
		Filename: errorLogFile,
	}, zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level > zap.InfoLevel
	}))
	NewGlobalLogger([]zapcore.Core{infoCore, errorCore})

	zap.L().Info("info message", zap.Int("status", 0))
	zap.L().Error("error message", zap.Int("status", 1))
	expectExistFile(infoLogFile, t)
	expectExistFile(errorLogFile, t)

	infoLine := expectReadFile(infoLogFile, t)
	errorLine := expectReadFile(errorLogFile, t)
	if infoLine == nil || errorLine == nil {
		return
	}

	if infoLine.Level != "INFO" || errorLine.Level != "ERROR" {
		t.Error("Log 'level' error")
	}
}
