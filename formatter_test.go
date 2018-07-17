package easy

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"strings"
	"testing"
	"time"
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
		format string
		fields logrus.Fields
		result string
	}{
		{"[%lvl%]: %time% - %first%", map[string]interface{}{"first": "First Custom Param"}, "[PANIC]: 0001-01-01T00:00:00Z - First Custom Param"},
		{"[%lvl%]: %time% - %first% %second%", map[string]interface{}{"first": "First Custom Param", "second": "Second Custom Param"}, "[PANIC]: 0001-01-01T00:00:00Z - First Custom Param Second Custom Param"},
	}
	for _, tv := range testValues {
		f.LogFormat = tv.format
		b, _ := f.Format(logrus.WithFields(tv.fields))
		if string(b) != tv.result {
			t.Errorf("formatting expected result was %q instead of %q", string(b), tv.result)
		}
	}
}

func TestFormatterFormatWithCustomDateFormat(t *testing.T) {
	f := Formatter{}

	testValues := []struct {
		timestampFormat string
	}{
		{time.RFC822},
		{time.RFC850},
		{"2006-01-02 15:04:05"},
		{"2006-01-02"},
	}
	for _, tv := range testValues {
		f.TimestampFormat = tv.timestampFormat
		e := logrus.WithField("", "")
		e.Time = time.Now()

		b, _ := f.Format(e)
		if !bytes.Contains(b, []byte(e.Time.Format(tv.timestampFormat))) {
			t.Errorf("formatting expected format date was %q instead of %q", string(b), tv.timestampFormat)
		}
	}
}
