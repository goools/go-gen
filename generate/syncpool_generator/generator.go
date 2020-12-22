package syncpool_generator

import (
	"fmt"
	"regexp"

	"github.com/dave/jennifer/jen"
	"github.com/goools/go-gen/generate"
	"github.com/goools/tools/strx"
	"github.com/sirupsen/logrus"
)

var (
	syncPoolDefRegexp = regexp.MustCompile(`([A-Za-z_][A-Za-z0-9]*)<([*A-Za-z][A-Za-z0-9/.]*)>`)
)

const (
	syncPoolDefRegexpTotal          = 3
	syncPoolDefRegexpNameIndex      = 1
	syncPoolDefRegexpValueTypeIndex = 2
	syncPoolFuncThisName            = "p"
)

type SyncPool struct {
	PkgPath      string
	Name         string
	Value        *generate.Type
	funcTypeName string
	funcThisName string
}

func NewSyncPool(pkgPath, syncPoolDef string) *SyncPool {
	syncPool := &SyncPool{
		PkgPath: pkgPath,
	}
	syncPool.parseDef(syncPoolDef)
	syncPool.funcThisName = syncPoolFuncThisName
	syncPool.funcTypeName = syncPool.Name
	return syncPool
}

func (syncPool *SyncPool) parseDef(syncPoolDef string) {
	res := syncPoolDefRegexp.FindStringSubmatch(syncPoolDef)
	if len(res) != syncPoolDefRegexpTotal {
		panic(fmt.Errorf("cannot parse sync pool def: %s", syncPoolDef))
	}
	syncPool.Name = res[syncPoolDefRegexpNameIndex]
	valueType := res[syncPoolDefRegexpValueTypeIndex]
	syncPool.Value = generate.NewType(valueType)
}

func (syncPool *SyncPool) WriteToFile() {
	logrus.Infof("begin generate sync pool: %s", syncPool.Name)
	syncPoolSnackName := strx.ToSnakeCase(syncPool.Name)
	syncPoolFileName := fmt.Sprintf("%s_syncpool_generate.go", syncPoolSnackName)

	generateFile := jen.NewFilePath(syncPool.PkgPath)
	generateFile.HeaderComment(generate.WriteDoNotEdit())

	generateFile.Add(syncPool.writeTypeDef())
	generateFile.Line()
	generateFile.Add(syncPool.writeFuncPut())
	generateFile.Line()
	generateFile.Add(syncPool.writeFuncGet())

	err := generateFile.Save(syncPoolFileName)
	if err != nil {
		logrus.Fatalf("save syncpool code to file have an err: %v, syncpool: %s, file: %s", err, syncPool.Name, syncPool)
	}
	logrus.Infof("complete generate syncpool: %s", syncPool.Name)
}

func (syncPool *SyncPool) writeTypeDef() jen.Code {
	res := jen.Type()
	res.Id(syncPool.Name).Qual("sync", "Pool")
	return res
}

func (syncPool *SyncPool) ValueType() jen.Code {
	return syncPool.Value.PtrCode()
}

func (syncPool *SyncPool) PreFuncCommon(funcName string) *jen.Statement {
	return jen.Func().
		Params(jen.Id(syncPool.funcThisName).Add(jen.Id(fmt.Sprintf("*%s", syncPool.Name)))).Id(funcName)
}

func (syncPool *SyncPool) ObjCode() *jen.Statement {
	ptrOp := jen.Op("*")
	objCode := jen.Params(ptrOp.Add(jen.Qual("sync", "Pool"))).Params(jen.Id(syncPool.funcThisName))
	return objCode
}

func (syncPool *SyncPool) writeFuncPut() jen.Code {
	res := syncPool.PreFuncCommon("Put")

	// func params
	res = res.Params(jen.Id("value").Add(syncPool.ValueType()))

	// func body
	objCode := syncPool.ObjCode()
	res = res.Block(
		objCode.Dot("Put").Call(jen.Id("value")),
	)

	return res
}

func (syncPool *SyncPool) writeFuncGet() jen.Code {
	res := syncPool.PreFuncCommon("Get")

	// func params
	res = res.Params()

	// func resp
	res = res.Params(syncPool.ValueType())

	// func body
	objCode := syncPool.ObjCode()
	res = res.Block(
		jen.Id("value").Op(":=").Add(objCode.Dot("Get").Call()),
		jen.Return(jen.Id("value").Assert(syncPool.ValueType())),
	)

	return res
}
