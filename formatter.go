package formatter

import (
	"bytes"
	"fmt"
	"html/template"
	"path"
	"runtime"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/gookit/color"

	"github.com/sirupsen/logrus"
)

// Formatter implements logrus.Formatter interface.
type Formatter struct {
	TimestampFormat  string
	LogFormat        string
	CallerPrettyfier func(*runtime.Frame) (ret string)
	LevelColor       map[logrus.Level]color.Color
}

// NewFormatter 新增一个Formatter
func NewFormatter() *Formatter {
	// 默认输出
	// [2006-01-02 15:04:05][INFO]main.go:123|main.main()
	// Log message
	return &Formatter{
		LogFormat:       "[{{.Time}}][{{.Level}}][{{.PathAndFunc}}] {{.Msg}}\n{{.YAML}}",
		TimestampFormat: "2006-01-02 15:04:05",
		CallerPrettyfier: func(f *runtime.Frame) string {
			if f != nil {
				filename := path.Base(f.File)
				return fmt.Sprintf("%s:%d|%s()", filename, f.Line, f.Function)
			} else {
				return ""
			}
		},
		LevelColor: map[logrus.Level]color.Color{
			logrus.TraceLevel: color.Gray,
			logrus.DebugLevel: color.Green,
			logrus.InfoLevel:  color.Blue,
			logrus.WarnLevel:  color.Yellow,
			logrus.ErrorLevel: color.Red,
			logrus.FatalLevel: color.Magenta,
			logrus.PanicLevel: color.Magenta,
		},
	}
}

// Format building log message.
func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {

	timestampFormat := f.TimestampFormat
	t, _ := template.New("format").Parse(f.LogFormat)

	log := struct {
		Time, Level, PathAndFunc, Msg, YAML string
	}{
		Time: entry.Time.Format(timestampFormat),
		Level: func(lvl logrus.Level) string {
			level := strings.ToUpper(lvl.String())
			col, has := f.LevelColor[lvl]
			if has {
				level = col.Sprint(level)
			}
			return level
		}(entry.Level),
		PathAndFunc: f.CallerPrettyfier(entry.Caller),
		YAML: func(data logrus.Fields) string {
			if len(data) > 0 {
				yml, _ := yaml.Marshal(data)
				return string(yml)
			}
			return ""
		}(entry.Data),
		Msg: entry.Message,
	}

	output := bytes.NewBuffer([]byte{})
	t.Execute(output, log)

	return []byte(output.String()), nil
}
