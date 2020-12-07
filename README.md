[![Go Report Card](https://goreportcard.com/badge/github.com/klarkxy/logrus-formatter)](https://goreportcard.com/badge/github.com/klarkxy/logrus-formatter)
## My Logrus Formatter
A provided formatter allow to format [Logrus](https://github.com/sirupsen/logrus) log output

## How to use
#### Install
```shell
go get github.com/klarkxy/logrus-formatter
```
#### Use
```go
logrus.SetFormatter(formatter.NewFormatter("xxx"))
```
or
```go
log := logrus.New()
log.SetFormatter(formatter.NewFormatter("QQ-Bot"))
// log.SetReportCaller(true)   // Just do you like
```
