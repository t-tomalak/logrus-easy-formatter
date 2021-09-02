package main

import (
	"errors"
	"os"

	formatter "github.com/klarkxy/logrus-formatter"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(formatter.NewFormatter("test-console"))
	logrus.SetReportCaller(true)
	logrus.AddHook(&formatter.FileHook{
		Formatter: formatter.NewFileFormatter("test-file"),
		Level:     logrus.ErrorLevel,
		Writer:    os.Stdout,
	})
	logrus.WithFields(logrus.Fields{
		"aaa": 1234,
		"bbb": map[string]interface{}{
			"ccc": 567.8,
			"ddd": "9",
		},
	}).Error("Hello World")
	logrus.WithFields(logrus.Fields{
		"aaa": 1234,
		"bbb": map[string]interface{}{
			"ccc": 567.8,
			"ddd": "9",
		},
	}).Debug("Hello World")
	logrus.WithFields(logrus.Fields{
		"aaa": 1234,
		"bbb": map[string]interface{}{
			"ccc": 567.8,
			"ddd": "9",
		},
	}).Info("Hello World")
	logrus.WithFields(logrus.Fields{
		"aaa": 1234,
		"bbb": map[string]interface{}{
			"ccc": 567.8,
			"ddd": "9",
		},
	}).Warn("Hello World")
	logrus.WithFields(logrus.Fields{
		"aaa": 1234,
		"bbb": map[string]interface{}{
			"ccc": 567.8,
			"ddd": "9",
		},
	}).Trace("Hello World")

	logrus.WithError(errors.New("hahaha")).Error("wahaha")

	logrus.WithFields(logrus.Fields{
		"aaa": 1234,
		"bbb": map[string]interface{}{
			"ccc": 567.8,
			"ddd": "9",
			"nil": nil,
		},
	}).Info("Hello World")

}
