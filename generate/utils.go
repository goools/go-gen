package generate

import (
	"github.com/sirupsen/logrus"
	"time"
)

func costTime() func() {
	startTime := time.Now()

	return func() {
		logrus.Infof("costs %d ms", time.Now().Sub(startTime).Milliseconds())
	}
}
