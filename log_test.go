package log

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"

	. "gopkg.in/check.v1"
)

func TestLog(t *testing.T) { TestingT(t) }

type LogSuite struct{}

var _ = Suite(&LogSuite{})

func (s *LogSuite) SetUpTest(c *C) {
	// reset global loggers chain before every test
	loggers = []Logger{}
}

func (s *LogSuite) TestInit(c *C) {
	Init(newTestLogger("log1"), newTestLogger("log2"))
	c.Assert(loggers[0].String(), Equals, "testLogger(log1)")
	c.Assert(loggers[1].String(), Equals, "testLogger(log2)")
}

func (s *LogSuite) TestInitWithConfig(c *C) {
	InitWithConfig(LogConfig{ConsoleLoggerName, "info"}, LogConfig{SysLoggerName, "info"})
	c.Assert(loggers[0].String(), Equals, "consoleLogger(INFO)")
	c.Assert(loggers[1].String(), Equals, "sysLogger(INFO)")
}

func (s *LogSuite) TestNewLogger(c *C) {
	l, err := NewLogger(LogConfig{ConsoleLoggerName, "info"})
	c.Assert(err, IsNil)
	c.Assert(l.String(), Equals, "consoleLogger(INFO)")

	l, err = NewLogger(LogConfig{SysLoggerName, "warn"})
	c.Assert(err, IsNil)
	c.Assert(l.String(), Equals, "sysLogger(WARN)")

	l, err = NewLogger(LogConfig{UDPLoggerName, "error"})
	c.Assert(err, IsNil)
	c.Assert(l.String(), Equals, "udpLogger(ERROR)")

	l, err = NewLogger(LogConfig{"SuperDuperLogger", "info"})
	c.Assert(err, NotNil)
	c.Assert(l, IsNil)
}

func (s *LogSuite) TestInfof(c *C) {
	logger1 := newTestLogger("log1")
	logger2 := newTestLogger("log2")
	Init(logger1, logger2)

	Infof("hello %s", "world")
	c.Assert(logger1.b.String(), Equals, "INFO hello world\n")
	c.Assert(logger2.b.String(), Equals, "INFO hello world\n")
}

func (s *LogSuite) TestWarnf(c *C) {
	logger1 := newTestLogger("log1")
	logger2 := newTestLogger("log2")
	Init(logger1, logger2)

	Warnf("hello %s", "world")
	c.Assert(logger1.b.String(), Equals, "WARN hello world\n")
	c.Assert(logger2.b.String(), Equals, "WARN hello world\n")
}

func (s *LogSuite) TestErrorf(c *C) {
	logger1 := newTestLogger("log1")
	logger2 := newTestLogger("log2")
	Init(logger1, logger2)

	Errorf("hello %s", "world")
	c.Assert(logger1.b.String(), Equals, "ERROR hello world\n")
	c.Assert(logger2.b.String(), Equals, "ERROR hello world\n")
}

func (s *LogSuite) TestFatalf(c *C) {
	logger1 := newTestLogger("log1")
	logger2 := newTestLogger("log2")
	Init(logger1, logger2)

	Fatalf("hello %s", "world")
	c.Assert(strings.Split(logger1.b.String(), "\n")[0], Equals, "FATAL hello world")
	c.Assert(strings.Split(logger2.b.String(), "\n")[0], Equals, "FATAL hello world")
}

// testLogger helps in tests.
type testLogger struct {
	id string
	b  *bytes.Buffer
}

func newTestLogger(id string) *testLogger {
	return &testLogger{id, &bytes.Buffer{}}
}

func (l *testLogger) Infof(format string, args ...interface{}) {
}

func (l *testLogger) Warnf(format string, args ...interface{}) {
}

func (l *testLogger) Errorf(format string, args ...interface{}) {
}

func (l *testLogger) Fatalf(format string, args ...interface{}) {
}

func (l *testLogger) Writer(sev Severity) io.Writer {
	return l.b
}

func (l *testLogger) FormatMessage(sev Severity, caller *callerInfo, format string, args ...interface{}) string {
	return fmt.Sprintf("%s %s\n", sev, fmt.Sprintf(format, args...))
}

func (l testLogger) String() string {
	return fmt.Sprintf("testLogger(%s)", l.id)
}
