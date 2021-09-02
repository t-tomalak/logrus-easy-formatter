package formatter

import (
	"io"

	"github.com/sirupsen/logrus"
)

type FileHook struct {
	Writer    io.Writer
	Level     logrus.Level
	Formatter logrus.Formatter
}

func (h *FileHook) Levels() (lvls []logrus.Level) {
	lvls = make([]logrus.Level, 0)
	for _, lvl := range logrus.AllLevels {
		if lvl <= h.Level {
			lvls = append(lvls, lvl)
		}
	}
	return
}

func (h *FileHook) Fire(entry *logrus.Entry) error {
	buf, err := h.Formatter.Format(entry)

	if err != nil {
		return err
	}

	h.Writer.Write(buf)
	return nil
}
