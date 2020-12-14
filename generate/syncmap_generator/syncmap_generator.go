package syncmap_generator

import (
	"github.com/goools/go-gen/generate"
	"github.com/goools/go-gen/packagex"
	"github.com/sirupsen/logrus"
)

type SyncMapGenerator struct {
	pkg      *packagex.Package
	syncMaps []*SyncMap
}

func NewSyncMapGenerator(pkg *packagex.Package) generate.Generator {
	return &SyncMapGenerator{
		pkg: pkg,
	}
}

func (s *SyncMapGenerator) WriteToFile() {
	for i := range s.syncMaps {
		syncMap := s.syncMaps[i]
		syncMap.WriteToFile()
	}
}

func (s *SyncMapGenerator) Scan(args ...string) {
	pkgPath := s.pkg.PkgPath
	typeDefNameSet := make(map[string]struct{})
	for i := range args {
		syncMapDef := args[i]
		syncMap := NewSyncMap(pkgPath, syncMapDef)
		s.syncMaps = append(s.syncMaps, syncMap)
		if _, exist := typeDefNameSet[syncMap.Name]; exist {
			logrus.Fatalf("duplicate syncmap name: %s", syncMap.Name)
		} else {
			typeDefNameSet[syncMap.Name] = struct{}{}
		}
	}
}
