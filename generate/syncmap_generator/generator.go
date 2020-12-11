package syncmap_generator

import (
	"regexp"

	"github.com/goools/go-gen/generate"
	"github.com/goools/go-gen/packagex"
	"github.com/sirupsen/logrus"
)

var (
	syncMapDefRegexp = regexp.MustCompile(`(^[A-Za-z_][A-Za-z0-9]*)<(^[A-Za-z][A-Za-z0-9]*),(^[A-Za-z][A-Za-z0-9]*)>`)
)

type SyncMap struct {
	PkgPath      string
	Name         string
	KeyType      string
	ValueType    string
	funcTypeName string
	funcThisName string
}

func NewSyncMap(pkgPath, name string) *SyncMap {
	logrus.Debugf("create new sync map, package path: %s, name: %s", pkgPath, name)
	res := &SyncMap{}
	return res
}

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
	logrus.Debugf("write to file")
}

func (s *SyncMapGenerator) syncMapDefParse(syncMapDef string) (typeName, keyType, valueType string) {
	return
}

func (s *SyncMapGenerator) checkSyncMaps(syncMapDefs []string) {
}

func (s *SyncMapGenerator) Scan(args ...string) {
	s.checkSyncMaps(args)
	pkgPath := s.pkg.PkgPath
	for i := range args {
		syncMapDef := args[i]
		syncMap := NewSyncMap(pkgPath, syncMapDef)
		s.syncMaps = append(s.syncMaps, syncMap)
	}
}
