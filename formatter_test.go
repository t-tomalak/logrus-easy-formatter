package easy

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

func TestFormatterDefaultFormat(t *testing.T) {
	f := Formatter{}

	e := logrus.WithField("", "")
	e.Message = "Test Message"
	e.Level = logrus.WarnLevel
	e.Time = time.Now()

	b, _ := f.Format(e)

	expected := strings.Join([]string{"[WARNING]:", e.Time.Format(time.RFC3339), "- Test Message"}, " ")
	if string(b) != expected {
		t.Errorf("formatting expected result was %q instead of %q", string(b), expected)
	}

}

func TestFormatterFormatWithCustomData(t *testing.T) {
	f := Formatter{}

	testValues := []struct {
		name   string
		format string
		fields logrus.Fields
		result string
	}{
		{
			"Single custom param",
			"[%lvl%]: %time% - %first%",
			map[string]interface{}{"first": "First Custom Param"},
			"[PANIC]: 0001-01-01T00:00:00Z - First Custom Param",
		},
		{
			"Multiple custom params of type string",
			"[%lvl%]: %time% - %first% %second%",
			map[string]interface{}{"first": "First Custom Param", "second": "Second Custom Param"},
			"[PANIC]: 0001-01-01T00:00:00Z - First Custom Param Second Custom Param",
		},
		{
			"Multiple custom params of different type",
			"[%lvl%]: %time% - %string%, %bool%, %int%",
			map[string]interface{}{"string": "String param", "bool": true, "int": 42},
			"[PANIC]: 0001-01-01T00:00:00Z - String param, true, 42",
		},
		{
			"Omits fields not included in format",
			"[%lvl%]: %time% - %first% %random%",
			map[string]interface{}{"first": "String param", "not_included": "random string"},
			"[PANIC]: 0001-01-01T00:00:00Z - String param %random%",
		},
	}
	for _, tv := range testValues {
		t.Run(tv.name, func(t *testing.T) {
			f.LogFormat = tv.format
			b, _ := f.Format(logrus.WithFields(tv.fields))
			if string(b) != tv.result {
				t.Errorf("formatting expected result was %q instead of %q", string(b), tv.result)
			}
		})
	}
}

func TestFormatterFormatWithCustomDateFormat(t *testing.T) {
	f := Formatter{}

	testValues := []struct {
		name            string
		timestampFormat string
	}{
		{
			"Timestamp with RFC822 format",
			time.RFC822,
		},
		{

			"Timestamp with RFC850 format",
			time.RFC850,
		},
		{
			"Timestamp with `yyyy-mm-dd hh:mm:ss` format",
			"2006-01-02 15:04:05",
		},
		{
			"Timestamp with `yyyy-mm-dd` format",
			"2006-01-02",
		},
	}

	for _, tv := range testValues {
		t.Run(tv.name, func(t *testing.T) {
			f.TimestampFormat = tv.timestampFormat
			e := logrus.WithField("", "")
			e.Time = time.Now()

			b, _ := f.Format(e)
			if !bytes.Contains(b, []byte(e.Time.Format(tv.timestampFormat))) {
				t.Errorf("formatting expected format date was %q instead of %q", string(b), tv.timestampFormat)
			}
		})

	}
}
