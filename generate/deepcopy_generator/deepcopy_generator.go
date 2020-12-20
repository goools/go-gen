package deepcopy_generator

import (
	"github.com/goools/go-gen/generate"
	"github.com/goools/go-gen/packagex"
	"github.com/sirupsen/logrus"
)

type DeepCopyGenerator struct {
	pkg       *packagex.Package
	deepCopys []*DeepCopy
}

func NewDeepCopyGenerator(pkg *packagex.Package) generate.Generator {
	return &DeepCopyGenerator{
		pkg: pkg,
	}
}

func (g *DeepCopyGenerator) WriteToFile() {
	for i := range g.deepCopys {
		g.deepCopys[i].WriteToFile()
	}
}

func (g *DeepCopyGenerator) Scan(args ...string) {
	for i := range args {
		structNameDef := args[i]
		logrus.Infof("struct name: %s", structNameDef)
		deepCopy := NewDeepCopy(g.pkg, structNameDef)
		g.deepCopys = append(g.deepCopys, deepCopy)
	}
}
