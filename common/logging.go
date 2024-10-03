package common

import (
	"fmt"
	"log"
	"os"
	"time"
)

const ERROR int = 40
const WARNING int = 30
const INFO int = 20
const DEBUG int = 10
const VERBOSE int = 0

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

func (c *Logger) formatLogMessage(level, msg string) string {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	return fmt.Sprintf("[%s] - <%s> %s: %s", timestamp, c.name, level, msg)
}

func concatMessage(msg ...any) string {
	// Concatenate all strings in msg
	log := ""

	for _, m := range msg {
		log += fmt.Sprintf("%v ", m)
	}

	return log
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

	c.logger.Printf(c.formatLogMessage("INFO", fmt.Sprintf(format, msg...)))
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

	c.logger.Printf(c.formatLogMessage("ERROR", fmt.Sprintf(format, msg...)))
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

	c.logger.Printf(c.formatLogMessage("WARNING", fmt.Sprintf(format, msg...)))
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

	c.logger.Printf(c.formatLogMessage("DEBUG", fmt.Sprintf(format, msg...)))
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

	c.logger.Printf(c.formatLogMessage("VERBOSE", fmt.Sprintf(format, msg...)))
}
