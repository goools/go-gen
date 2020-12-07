package generate

import (
	"github.com/goools/go-gen/packagex"
	"github.com/sirupsen/logrus"
	"os"
)

type Generator interface {
	WriteToFile()
	Scan(args ...string)
}

type CreateGenerator func(pkg *packagex.Package) Generator

func RunGenerator(createGenerator CreateGenerator, args []string) {
	defer costTime()()
	pwd, err := os.Getwd()
	if err != nil {
		logrus.Fatalf("Getwd have an err: %v", err)
	}
	logrus.Debugf("pwd: %v", pwd)
	pkg, err := packagex.Load(pwd)
	generator := createGenerator(pkg)
	if err != nil {
		logrus.Fatalf("packages load %s have an err: %v", pwd, err)
	}
	generator.Scan(args...)
	generator.WriteToFile()
}
