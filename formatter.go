// Package easy allows to easily format output of Logrus logger
package easy

import (
	"fmt"
	"path"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/gookit/color"

	"github.com/sirupsen/logrus"
)

const (
	// Default log format will output [INFO]: 2006-01-02T15:04:05Z07:00 - Log message
	defaultLogFormat       = "[%lvl%]: %time% - %msg%"
	defaultTimestampFormat = time.RFC3339
)

// Formatter implements logrus.Formatter interface.
type Formatter struct {
	// Timestamp format
	TimestampFormat string
	// Available standard keys: time, msg, lvl, func, path,
	// Also can include custom fields but limited to strings.
	// All of fields need to be wrapped inside %% i.e %time% %msg%
	LogFormat string

	// CallerPrettyfier can be set by the user to modify the content
	// of the function and file keys in the data when ReportCaller is
	// activated. If any of the returned value is the empty string the
	// corresponding key will be removed from fields.
	CallerPrettyfier func(*runtime.Frame) (function string, file string)

	// UseColors
	// True: Use LevelColor's color.
	// False: No color
	UseColors bool

	// LevelColor is a map for tag level
	// The Color is just useful on the TTY
	LevelColor map[logrus.Level]color.Color
}

// Format building log message.
func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {

	output := f.LogFormat
	if output == "" {
		output = defaultLogFormat
	}

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = defaultTimestampFormat
	}

	output = strings.Replace(output, "%time%", entry.Time.Format(timestampFormat), 1)

	output = strings.Replace(output, "%msg%", entry.Message, 1)

	level := strings.ToUpper(entry.Level.String())
	if f.UseColors {
		col, has := f.getColor(entry.Level)
		if has {
			level = col.Sprint(level)
		}
	}
	output = strings.Replace(output, "%lvl%", level, 1)

	funcVal := ""
	fileVal := ""
	if entry.HasCaller() {
		fc := *entry.Caller
		filename := path.Base(fc.File)
		funcVal, fileVal = fmt.Sprintf("%s()", fc.Function), fmt.Sprintf("%s:%d", filename, fc.Line)

		if f.CallerPrettyfier != nil {
			funcVal, fileVal = f.CallerPrettyfier(entry.Caller)
		}
	} else {
		fmt.Println("No Caller.")
		fmt.Println(entry.Logger)
		fmt.Println(entry.Logger.ReportCaller)
		fmt.Println(entry.Caller)
	}
	output = strings.Replace(output, "%func%", funcVal, 1)
	output = strings.Replace(output, "%path%", fileVal, 1)

	for k, val := range entry.Data {
		switch v := val.(type) {
		case string:
			output = strings.Replace(output, "%"+k+"%", v, 1)
		case int:
			s := strconv.Itoa(v)
			output = strings.Replace(output, "%"+k+"%", s, 1)
		case bool:
			s := strconv.FormatBool(v)
			output = strings.Replace(output, "%"+k+"%", s, 1)
		}
	}

	return []byte(output), nil
}

var defaultColor = map[logrus.Level]color.Color{
	logrus.TraceLevel: color.Gray,
	logrus.DebugLevel: color.Cyan,
	logrus.InfoLevel:  color.Blue,
	logrus.WarnLevel:  color.Yellow,
	logrus.ErrorLevel: color.Red,
	logrus.FatalLevel: color.Magenta,
	logrus.PanicLevel: color.Magenta,
}

func (f *Formatter) getColor(lv logrus.Level) (color.Color, bool) {
	col, has := f.LevelColor[lv]
	if !has {
		col, has = defaultColor[lv]
	}
	if !has {
		return color.FgDefault, has
	}
	return col, has
}
