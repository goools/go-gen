package syncmap_generator

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/goools/go-gen/generate"
	"github.com/goools/go-gen/packagex"
	"github.com/sirupsen/logrus"
)

var (
	syncMapDefRegexp = regexp.MustCompile(`([A-Za-z_][A-Za-z0-9]*)<([A-Za-z][A-Za-z0-9/.]*),([A-Za-z][A-Za-z0-9/.]*)>`)
)

const (
	syncMapDefRegexpTotal          = 4
	syncMapDefRegexpNameIndex      = 1
	syncMapDefRegexpKeyTypeIndex   = 2
	syncMapDefRegexpValueTypeIndex = 3
	syncMapFuncThisName            = "s"
)

type Type struct {
	packagePath string
	typeName    string
	code        jen.Code
}

func (t *Type) Code() jen.Code {
	if t.code != nil {
		return t.code
	}
	if t.packagePath == "" {
		t.code = jen.Id(t.typeName)
	} else {
		t.code = jen.Qual(t.packagePath, t.typeName)
	}
	return t.code
}

func (t *Type) TypeName() string {
	if t.packagePath != "" {
		pkgPath := t.packagePath
		pkgPath = strings.ReplaceAll(pkgPath, ".", "_")
		pkgPath = strings.ReplaceAll(pkgPath, "/", "__")
		return fmt.Sprintf("%s_%s", pkgPath, t.typeName)
	}
	return t.typeName
}

func NewType(typeDef string) *Type {
	index := strings.LastIndex(typeDef, ".")
	var res *Type
	if index == -1 {
		res = &Type{
			typeName: typeDef,
		}
	} else {
		res = &Type{
			packagePath: typeDef[:index],
			typeName:    typeDef[index+1:],
		}
	}
	logrus.Debugf("type: %#v", res)
	return res
}

type SyncMap struct {
	PkgPath        string
	Name           string
	Key            *Type
	Value          *Type
	funcTypeName   string
	funcThisName   string
	emptyValueName string
}

func (syncMap *SyncMap) KeyType() jen.Code {
	return syncMap.Key.Code()
}

func (syncMap *SyncMap) ValueType() jen.Code {
	return syncMap.Value.Code()
}

func (syncMap *SyncMap) syncMapDefParse(syncMapDef string) {
	var res []string
	res = syncMapDefRegexp.FindStringSubmatch(syncMapDef)
	if len(res) != syncMapDefRegexpTotal {
		panic("cannot find syncmap name, key type and value type")
	}
	syncMap.Name = res[syncMapDefRegexpNameIndex]
	keyType := res[syncMapDefRegexpKeyTypeIndex]
	syncMap.Key = NewType(keyType)
	valueType := res[syncMapDefRegexpValueTypeIndex]
	syncMap.Value = NewType(valueType)
	syncMap.emptyValueName = fmt.Sprintf("_%s_%s_empty_value", syncMap.Value.TypeName(), syncMap.Name)
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
