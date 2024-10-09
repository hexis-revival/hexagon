package common

import (
	"fmt"
	"log"
	"os"
	"time"
	"unicode"
)

const (
	ERROR   int = 40
	WARNING int = 30
	INFO    int = 20
	DEBUG   int = 10
	VERBOSE int = 0
)

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorWhite  = "\033[37m"
	ColorGrey   = "\033[90m"
)

type Logger struct {
	logger *log.Logger
	name   string
	level  int
}

func CreateLogger(name string, level int) *Logger {
	l := log.New(os.Stdout, "", 0)
	return &Logger{
		logger: l,
		name:   name,
		level:  level,
	}
}

func (c *Logger) SetLevel(level int) {
	c.level = level
}

func (c *Logger) GetLevel() int {
	return c.level
}

func (c *Logger) SetName(name string) {
	c.name = name
}

func (c *Logger) GetName() string {
	return c.name
}

func (c *Logger) formatLogMessage(level string, msg string) string {
	color := getColor(level)
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	return fmt.Sprintf(
		"[%s] - <%s> %s%s: %s%s",
		timestamp,
		c.name,
		color,
		level,
		msg,
		ColorReset,
	)
}

func getColor(level string) string {
	switch level {
	case "INFO":
		return ColorBlue
	case "ERROR":
		return ColorRed
	case "WARNING":
		return ColorYellow
	case "DEBUG":
		return ColorWhite
	case "VERBOSE":
		return ColorGrey
	default:
		return ColorReset
	}
}

func concatMessage(msg ...any) string {
	// Concatenate all strings in msg
	log := ""

	for _, m := range msg {
		log += fmt.Sprintf("%v ", m)
	}

	return log
}

func FormatBytes(data []byte) string {
	result := ""
	for _, b := range data {
		if unicode.IsPrint(rune(b)) {
			// Append ASCII character to result
			result += fmt.Sprintf("%c", b)
		} else {
			// Append \xhh for non-printable characters
			result += fmt.Sprintf("\\x%02x", b)
		}
	}
	return result
}

func (c *Logger) Info(msg ...any) {
	if c.level > INFO {
		return
	}

	c.logger.Println(c.formatLogMessage("INFO", concatMessage(msg...)))
}

func (c *Logger) Infof(format string, msg ...any) {
	if c.level > INFO {
		return
	}

	c.logger.Println(c.formatLogMessage("INFO", fmt.Sprintf(format, msg...)))
}

func (c *Logger) Error(msg ...any) {
	if c.level > ERROR {
		return
	}
	c.logger.Println(c.formatLogMessage("ERROR", concatMessage(msg...)))
}

func (c *Logger) Errorf(format string, msg ...any) {
	if c.level > ERROR {
		return
	}

	c.logger.Println(c.formatLogMessage("ERROR", fmt.Sprintf(format, msg...)))
}

func (c *Logger) Warning(msg ...any) {
	if c.level > WARNING {
		return
	}
	c.logger.Println(c.formatLogMessage("WARNING", concatMessage(msg...)))
}

func (c *Logger) Warningf(format string, msg ...any) {
	if c.level > WARNING {
		return
	}

	c.logger.Println(c.formatLogMessage("WARNING", fmt.Sprintf(format, msg...)))
}

func (c *Logger) Debug(msg ...any) {
	if c.level > DEBUG {
		return
	}
	c.logger.Println(c.formatLogMessage("DEBUG", concatMessage(msg...)))
}

func (c *Logger) Debugf(format string, msg ...any) {
	if c.level > DEBUG {
		return
	}

	c.logger.Println(c.formatLogMessage("DEBUG", fmt.Sprintf(format, msg...)))
}

func (c *Logger) Verbose(msg ...any) {
	if c.level > VERBOSE {
		return
	}
	c.logger.Println(c.formatLogMessage("VERBOSE", concatMessage(msg...)))
}

func (c *Logger) Verbosef(format string, msg ...any) {
	if c.level > VERBOSE {
		return
	}

	c.logger.Println(c.formatLogMessage("VERBOSE", fmt.Sprintf(format, msg...)))
}
