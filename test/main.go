package main

import (
	"errors"

	formatter "github.com/klarkxy/logrus-formatter"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(formatter.NewFormatter("test"))
	logrus.SetReportCaller(true)
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
}
