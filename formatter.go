package formatter

import (
	"bytes"
	"fmt"
	"path"
	"runtime"
	"strings"
	"text/template"

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
	ModuleName       string
	NoColor          bool
}

func NewFormatter(modName string) *Formatter {
	return NewConsoleFormatter(modName)
}

func NewConsoleFormatter(modName string) *Formatter {
	return &Formatter{
		LogFormat:       "[{{.Time}}][{{.Level}}][{{.Module}}][{{.PathAndFunc}}] {{.Msg}}\n{{.YAML}}",
		TimestampFormat: "2006-01-02 15:04:05",
		ModuleName:      modName,
		CallerPrettyfier: func(f *runtime.Frame) string {
			if f != nil {
				filename := path.Base(f.File)
				fun := strings.Split(f.Function, "/")
				return fmt.Sprintf("%s:%d|%s()", filename, f.Line, fun[len(fun)-1])
			}
			return ""
		},
		LevelColor: map[logrus.Level]color.Color{
			logrus.TraceLevel: color.Gray,
			logrus.DebugLevel: color.Green,
			logrus.InfoLevel:  color.Blue,
			logrus.WarnLevel:  color.Yellow,
			logrus.ErrorLevel: color.Red,
			logrus.FatalLevel: color.Magenta,
			logrus.PanicLevel: color.Bold,
		},
		NoColor: false,
	}
}
func NewFileFormatter(modName string) *Formatter {
	return &Formatter{
		LogFormat:       "[{{.Time}}][{{.Level}}][{{.Module}}][{{.PathAndFunc}}] {{.Msg}}\n{{.YAML}}",
		TimestampFormat: "2006-01-02 15:04:05",
		ModuleName:      modName,
		CallerPrettyfier: func(f *runtime.Frame) string {
			if f != nil {
				filename := path.Base(f.File)
				fun := strings.Split(f.Function, "/")
				return fmt.Sprintf("%s:%d|%s()", filename, f.Line, fun[len(fun)-1])
			}
			return ""
		},
		NoColor: true,
	}
}

// Format building log message.
func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {

	timestampFormat := f.TimestampFormat
	t, _ := template.New("format").Parse(f.LogFormat)

	// 对error进行特殊处理
	if field, ok := entry.Data[logrus.ErrorKey]; ok {
		if err, ok := field.(error); ok {
			entry.Data[logrus.ErrorKey] = struct {
				Text  string
				Error error
			}{Text: err.Error(), Error: err}
		}
	}
	Sprint := func(a ...interface{}) string {
		if f.NoColor {
			return fmt.Sprint(a...)
		}

		col, has := f.LevelColor[entry.Level]
		if !has {
			col = color.White
		}
		return col.Sprint(a...)
	}

	log := struct {
		Time, Level, PathAndFunc, Msg, YAML, Module string
	}{
		Time: Sprint(entry.Time.Format(timestampFormat)),
		Level: func(lvl logrus.Level) string {
			level := strings.ToUpper(lvl.String())
			return Sprint((level + "     ")[:4])
		}(entry.Level),
		PathAndFunc: func() string {
			if entry.Caller != nil {
				return Sprint(f.CallerPrettyfier(entry.Caller))
			}
			return ""
		}(),
		YAML: func(data logrus.Fields) string {
			if len(data) > 0 {
				yml, err := yaml.Marshal(data)
				if err == nil {
					return string(yml)
				}
			}
			return ""
		}(entry.Data),
		Msg:    Sprint(entry.Message),
		Module: Sprint(f.ModuleName),
	}
	output := bytes.NewBuffer([]byte{})
	t.Execute(output, log)

	return output.Bytes(), nil
}
