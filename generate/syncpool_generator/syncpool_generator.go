package syncpool_generator

import (
	"github.com/goools/go-gen/generate"
	"github.com/goools/go-gen/packagex"
)

type SyncPoolGenerator struct {
	pkg       *packagex.Package
	syncPools []*SyncPool
}

func NewSyncPoolGenerator(pkg *packagex.Package) generate.Generator {
	return &SyncPoolGenerator{
		pkg: pkg,
	}
}

func (g *SyncPoolGenerator) WriteToFile() {
	for i := range g.syncPools {
		g.syncPools[i].WriteToFile()
	}
}

func (g *SyncPoolGenerator) Scan(args ...string) {
	for i := range args {
		syncPoolDef := args[i]
		syncPool := NewSyncPool(g.pkg.PkgPath, syncPoolDef)
		g.syncPools = append(g.syncPools, syncPool)
	}
}
