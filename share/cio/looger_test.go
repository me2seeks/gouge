package cio

import (
	"strings"
	"testing"
)

func TestFork(t *testing.T) {
	// 创建一个 Logger 实例
	originalLogger := &Logger{
		prefix: "original",
		Info:   true,
		Debug:  true,
	}

	// 调用 Fork 方法
	newLogger := originalLogger.Fork("new %s", "logger")

	// 检查新 Logger 实例的前缀是否正确
	expectedPrefix := "original: new logger"
	if !strings.Contains(newLogger.prefix, expectedPrefix) {
		t.Errorf("Expected prefix '%s', but got '%s'", expectedPrefix, newLogger.prefix)
	}

	// 检查 Info 和 Debug 字段是否与原始 Logger 实例相同
	if newLogger.Info != originalLogger.Info {
		t.Errorf("Expected Info to be '%v', but got '%v'", originalLogger.Info, newLogger.Info)
	}

	if newLogger.Debug != originalLogger.Debug {
		t.Errorf("Expected Debug to be '%v', but got '%v'", originalLogger.Debug, newLogger.Debug)
	}
}
