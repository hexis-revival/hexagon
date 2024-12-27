package common

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"time"
	"unicode"
)

const (
	ANOMALY int = 50
	ERROR   int = 40
	WARNING int = 30
	INFO    int = 20
	DEBUG   int = 10
	VERBOSE int = 0
)

const (
	ColorReset  = "\033[0m"
	ColorBlack  = "\033[30m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
	ColorGrey   = "\033[0;90m"
)

const (
	ColorBoldBlack  = "\033[1;30m"
	ColorBoldRed    = "\033[1;31m"
	ColorBoldGreen  = "\033[1;32m"
	ColorBoldYellow = "\033[1;33m"
	ColorBoldBlue   = "\033[1;34m"
	ColorBoldPurple = "\033[1;35m"
	ColorBoldCyan   = "\033[1;36m"
	ColorBoldWhite  = "\033[1;37m"
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

func (c *Logger) Anomaly(msg ...any) {
	if c.level > ANOMALY {
		return
	}
	c.logger.Println(c.formatLogMessage("ANOMALY", concatMessage(msg...)))
}

func (c *Logger) Anomalyf(format string, msg ...any) {
	if c.level > ANOMALY {
		return
	}
	c.logger.Println(c.formatLogMessage("ANOMALY", fmt.Sprintf(format, msg...)))
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
	case "ANOMALY":
		return ColorBoldYellow
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

// FormatBytes returns a string representation of a byte slice
// similar to how python handles byte strings
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

// FormatStruct returns a string representation of a struct
func FormatStruct(s interface{}) string {
	v := reflect.ValueOf(s)
	t := v.Type()

	// Check if the value is a pointer and dereference it
	if t.Kind() == reflect.Ptr {
		v = v.Elem()
		t = v.Type()
	}

	// Ensure we're dealing with a struct
	if t.Kind() != reflect.Struct {
		return fmt.Sprintf("%v", s)
	}

	var sb strings.Builder
	sb.WriteString(t.Name() + "{")

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		// Skip unexported fields
		if !field.IsExported() {
			continue
		}

		// Append field name and formatted value
		sb.WriteString(fmt.Sprintf("%s: %s", field.Name, FormatValue(value)))

		if i < v.NumField()-1 {
			sb.WriteString(", ")
		}
	}

	sb.WriteString("}")
	return sb.String()
}

// FormatValue handles different types and returns the formatted string
func FormatValue(v reflect.Value) string {
	defer recover()
	switch v.Kind() {
	case reflect.String:
		return fmt.Sprintf("\"%s\"", v.String())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return fmt.Sprintf("%d", v.Uint())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%d", v.Int())
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%f", v.Float())
	case reflect.Bool:
		return fmt.Sprintf("%t", v.Bool())
	case reflect.Slice:
		if v.Type().Elem().Kind() == reflect.Uint8 {
			return fmt.Sprintf("%v", v.Bytes())
		}
		return fmt.Sprintf("%v", v.Interface())
	case reflect.Struct:
		if v.Type() == reflect.TypeOf(time.Time{}) {
			return v.Interface().(time.Time).Format(time.RFC3339)
		}
		return FormatStruct(v.Interface())
	default:
		return fmt.Sprintf("%v", v.Interface())
	}
}
