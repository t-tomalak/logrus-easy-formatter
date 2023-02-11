// Package easy allows to easily format output of Logrus logger
package easy

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

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
	// Available standard keys: time, msg, lvl, caller
	// Also can include custom fields but limited to strings.
	// All of fields need to be wrapped inside %% i.e %time% %msg%
	LogFormat string
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

	caller := ""
	if entry.Caller != nil {
		// fileShortPath should contain at most the last 3 levels of the file path
		filePath := entry.Caller.File
		tokens := strings.Split(filePath, string(os.PathSeparator))
		startPos := len(tokens) - 3
		if startPos < 0 {
			startPos = 0
		}
		fileShortPath := strings.Join(tokens[startPos:], string(os.PathSeparator))

		// keep the function name only without the package prefix e.g. package.Func => Func
		funcName := entry.Caller.Func.Name()
		funcShortName := funcName[strings.LastIndex(funcName, ".")+1:]

		caller = fmt.Sprintf("%s:%d in %s", fileShortPath, entry.Caller.Line, funcShortName)
	}

	output = strings.Replace(output, "%time%", entry.Time.Format(timestampFormat), 1)
	output = strings.Replace(output, "%msg%", entry.Message, 1)
	output = strings.Replace(output, "%lvl%", strings.ToUpper(entry.Level.String()), 1)
	output = strings.Replace(output, "%caller%", caller, 1)

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
