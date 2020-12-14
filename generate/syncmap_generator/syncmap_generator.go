package syncmap_generator

import (
	"fmt"

	"github.com/dave/jennifer/jen"
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

func (s *SyncMapGenerator) writeToFile(syncMap *SyncMap) {
	logrus.Infof("begin generate syncmap: %s", syncMap.Name)
	syncMapSnackName := generate.ToSnakeCase(syncMap.Name)
	syncMapFileName := fmt.Sprintf("%s_syncmap_generate.go", syncMapSnackName)
	generateFile := jen.NewFilePath(syncMap.PkgPath)

	generateFile.HeaderComment(generate.WriteDoNotEdit())
	generateFile.Add(syncMap.writeTypeDef())
	generateFile.Line()
	generateFile.Add(syncMap.writeEmptyValue())
	generateFile.Line()
	generateFile.Add(syncMap.writeFuncStore())
	generateFile.Line()
	generateFile.Add(syncMap.writeFuncLoadOrStore())
	generateFile.Line()
	generateFile.Add(syncMap.writeFuncLoad())
	generateFile.Line()
	generateFile.Add(syncMap.writeFuncDelete())
	generateFile.Line()
	generateFile.Add(syncMap.writeFuncRange())
	generateFile.Line()
	generateFile.Add(syncMap.writeFuncLoadAndDelete())

	err := generateFile.Save(syncMapFileName)
	if err != nil {
		logrus.Fatalf("save syncmap code to file have an err: %v, syncmap: %s, file: %s", err, syncMap.Name, syncMapFileName)
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
