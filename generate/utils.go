package generate

import (
	"regexp"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

func costTime() func() {
	startTime := time.Now()

	return func() {
		logrus.Infof("costs %d ms", time.Now().Sub(startTime).Milliseconds())
	}
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func WriteDoNotEdit() string {
	return "Code generated by go-gen DO NOT EDIT."
}
