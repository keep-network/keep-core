package testutils

type MockLogger struct{}

func (ml *MockLogger) Debug(args ...interface{}) {}

func (ml *MockLogger) Debugf(format string, args ...interface{}) {}

func (ml *MockLogger) Error(args ...interface{}) {}

func (ml *MockLogger) Errorf(format string, args ...interface{}) {}

func (ml *MockLogger) Fatal(args ...interface{}) {}

func (ml *MockLogger) Fatalf(format string, args ...interface{}) {}

func (ml *MockLogger) Info(args ...interface{}) {}

func (ml *MockLogger) Infof(format string, args ...interface{}) {}

func (ml *MockLogger) Panic(args ...interface{}) {}

func (ml *MockLogger) Panicf(format string, args ...interface{}) {}

func (ml *MockLogger) Warn(args ...interface{}) {}

func (ml *MockLogger) Warnf(format string, args ...interface{}) {}

func (ml *MockLogger) Warning(args ...interface{}) {}

func (ml *MockLogger) Warningf(format string, args ...interface{}) {}
