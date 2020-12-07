package generate

import (
	"github.com/sirupsen/logrus"
	"regexp"
	"strings"
	"time"
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
