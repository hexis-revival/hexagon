package common

import (
	"testing"
	"time"
)

type TestStruct struct {
	String  string
	Integer int
	Time    time.Time
	List    []string
}

func TestLogging(t *testing.T) {
	t.Parallel()

	logger := CreateLogger("test", QUIET)
	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warning("warn message")
	logger.Error("error message")

	logger.Debugf("debug message %s", "formatted")
	logger.Infof("info message %s", "formatted")
	logger.Warningf("warn message %s", "formatted")
	logger.Errorf("error message %s", "formatted")

	test := &TestStruct{
		String:  "test",
		Integer: 1,
		Time:    time.Now(),
		List:    []string{"a", "b", "c"},
	}

	formatted := FormatStruct(test)
	logger.Debug(formatted)

	bytesData := []byte{0x01, 0x02, 0x03, 0x04}
	formatted = FormatBytes(bytesData)
	logger.Debug(formatted)
}
