package syncmap_generator

import (
	"fmt"
	"regexp"

	"github.com/dave/jennifer/jen"
	"github.com/goools/go-gen/generate"
	"github.com/goools/go-gen/packagex"
	"github.com/sirupsen/logrus"
)

var (
	syncMapDefRegexp = regexp.MustCompile(`([A-Za-z_][A-Za-z0-9]*)<([A-Za-z][A-Za-z0-9]*),([A-Za-z][A-Za-z0-9]*)>`)
)

const (
	syncMapDefRegexpTotal          = 4
	syncMapDefRegexpNameIndex      = 1
	syncMapDefRegexpKeyTypeIndex   = 2
	syncMapDefRegexpValueTypeIndex = 3
	syncMapFuncThisName            = "s"
)

type SyncMap struct {
	PkgPath      string
	Name         string
	KeyType      string
	ValueType    string
	funcTypeName string
	funcThisName string
}

func (syncMap *SyncMap) syncMapDefParse(syncMapDef string) {
	var res []string
	res = syncMapDefRegexp.FindStringSubmatch("Pill<int,int>")
	if len(res) != syncMapDefRegexpTotal {
		panic("cannot find syncmap name, key type and value type")
	}
	syncMap.Name = res[syncMapDefRegexpNameIndex]
	syncMap.KeyType = res[syncMapDefRegexpKeyTypeIndex]
	syncMap.ValueType = res[syncMapDefRegexpValueTypeIndex]
	return
}

func (syncMap *SyncMap) funcTypeParams() jen.Code {
	return jen.Id(syncMap.funcThisName).Id(fmt.Sprintf("*%s", syncMap.funcTypeName))
}

func NewSyncMap(pkgPath, syncMapDef string) *SyncMap {
	logrus.Debugf("create new sync map, package path: %s, syncMapDef: %s", pkgPath, syncMapDef)
	res := &SyncMap{
		PkgPath: pkgPath,
	}
	res.syncMapDefParse(syncMapDef)
	res.funcTypeName = res.Name
	res.funcThisName = syncMapFuncThisName
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

func (s *SyncMapGenerator) writeToFile(syncMap *SyncMap) {
	logrus.Infof("begin generate syncmap: %s", syncMap.Name)
	syncMapSnackName := generate.ToSnakeCase(syncMap.Name)
	syncMapFileName := fmt.Sprintf("%s_syncmap_generate.go", syncMapSnackName)
	generateFile := jen.NewFilePath(syncMap.PkgPath)

	generateFile.HeaderComment(generate.WriteDoNotEdit())
	generateFile.Add(syncMap.writeTypeDef())
	generateFile.Add(syncMap.writeFuncStore())

	err := generateFile.Save(syncMapFileName)
	if err != nil {
		logrus.Fatalf("save enum code to file have an err: %v, syncmap: %s, file: %s", err, syncMap.Name, syncMapFileName)
	}
	logrus.Infof("complete generate syncmap: %s", syncMap.Name)
}

func (s *SyncMapGenerator) WriteToFile() {
	for i := range s.syncMaps {
		syncMap := s.syncMaps[i]
		s.writeToFile(syncMap)
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
