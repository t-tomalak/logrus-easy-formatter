package main

import (
	formatter "github.com/klarkxy/logrus-formatter"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(formatter.NewFormatter("TEST"))
	logrus.SetReportCaller(true)
	logrus.Error("Hello World")
}
