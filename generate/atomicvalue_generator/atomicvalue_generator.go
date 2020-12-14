package atomicvalue_generator

import (
	"github.com/goools/go-gen/generate"
	"github.com/goools/go-gen/packagex"
)

type AtomicValueGenerator struct {
	pkg          *packagex.Package
	atomicValues []*AtomicValue
}

func NewAtomicValueGenerator(pkg *packagex.Package) generate.Generator {
	return &AtomicValueGenerator{
		pkg: pkg,
	}
}

func (g *AtomicValueGenerator) WriteToFile() {
	for i := range g.atomicValues {
		g.atomicValues[i].WriteToFile()
	}
}

func (g *AtomicValueGenerator) Scan(args ...string) {
	for i := range args {
		syncPoolDef := args[i]
		syncPool := NewAtomicValue(g.pkg.PkgPath, syncPoolDef)
		g.atomicValues = append(g.atomicValues, syncPool)
	}
}
