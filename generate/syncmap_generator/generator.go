package syncmap_generator

import (
	"fmt"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/goools/go-gen/generate"
	"github.com/goools/tools/strx"
	"github.com/sirupsen/logrus"
)

var (
	syncMapDefRegexp = regexp.MustCompile(`([A-Za-z_][A-Za-z0-9]*)<([*A-Za-z][A-Za-z0-9/.]*),([*A-Za-z][A-Za-z0-9/.]*)>`)
)

const (
	syncMapDefRegexpTotal          = 4
	syncMapDefRegexpNameIndex      = 1
	syncMapDefRegexpKeyTypeIndex   = 2
	syncMapDefRegexpValueTypeIndex = 3
	syncMapFuncThisName            = "s"
)

type SyncMap struct {
	PkgPath        string
	Name           string
	Key            *generate.Type
	Value          *generate.Type
	funcTypeName   string
	funcThisName   string
	emptyValueName string
}

func (syncMap *SyncMap) WriteToFile() {
	logrus.Infof("begin generate syncmap: %s", syncMap.Name)
	syncMapSnackName := strx.ToSnakeCase(syncMap.Name)
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
	syncMap.Key = generate.NewType(keyType)
	valueType := res[syncMapDefRegexpValueTypeIndex]
	syncMap.Value = generate.NewType(valueType)
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

func (syncMap *SyncMap) writeEmptyValue() jen.Code {
	emptyFunc := jen.Func().Params().Params(jen.Id("val").Add(syncMap.ValueType())).Block(jen.Return())
	return jen.Var().Id(syncMap.emptyValueName).Op("=").Add(emptyFunc.Call())
}

func (syncMap *SyncMap) syncMapObjId() *jen.Statement {
	ptrOp := jen.Op("*")
	syncMapObjId := jen.Params(ptrOp.Add(jen.Qual("sync", "Map"))).Params(jen.Id(syncMap.funcThisName))
	return syncMapObjId
}

func (syncMap *SyncMap) writeTypeDef() jen.Code {
	res := jen.Type()
	res.Id(syncMap.Name).Qual("sync", "Map")
	return res
}

func (syncMap *SyncMap) writeFuncStore() jen.Code {
	funcName := "Store"
	// func
	res := jen.Func()
	// func type params
	res = res.Params(syncMap.funcTypeParams())
	// func name
	res = res.Id(funcName)
	// func params
	res = res.Params(jen.List(
		jen.Id("key").Add(syncMap.KeyType()),
		jen.Id("value").Add(syncMap.ValueType()),
	))
	// func body
	syncMapObjId := syncMap.syncMapObjId()

	res.Block(
		syncMapObjId.Dot("Store").Call(jen.Id("key"), jen.Id("value")),
	)

	return res
}

func (syncMap *SyncMap) writeFuncLoadOrStore() jen.Code {
	funcName := "LoadOrStore"
	// func
	res := jen.Func()
	// func type params
	res = res.Params(syncMap.funcTypeParams())
	// func name
	res = res.Id(funcName)
	// func params
	res = res.Params(jen.List(
		jen.Id("key").Add(syncMap.KeyType()),
		jen.Id("value").Add(syncMap.ValueType()),
	)).Params(jen.List(syncMap.ValueType(), jen.Bool()))

	// func body
	syncMapObjId := syncMap.syncMapObjId()
	result := jen.Id("res").Assert(syncMap.ValueType())
	res = res.Block(
		jen.List(jen.Id("res"), jen.Id("ok")).Op(":=").Add(syncMapObjId.Dot("LoadOrStore").Call(jen.Id("key"), jen.Id("value"))),
		jen.If(jen.Op("!").Add().Id("ok")).Block(
			jen.Return(jen.List(jen.Id(syncMap.emptyValueName), jen.Id("ok"))),
		),
		jen.Return(jen.List(result, jen.Id("ok"))),
	)

	return res
}

func (syncMap *SyncMap) writeFuncLoad() jen.Code {
	funcName := "Load"
	// func
	res := jen.Func()
	// func type params
	res = res.Params(syncMap.funcTypeParams())
	// func name
	res = res.Id(funcName)
	// func params
	res = res.Params(jen.List(jen.Id("key").Add(syncMap.KeyType()))).Params(jen.List(syncMap.ValueType(), jen.Bool()))

	// func body
	syncMapObjId := syncMap.syncMapObjId()
	result := jen.Id("res").Assert(syncMap.ValueType())
	res = res.Block(
		jen.List(jen.Id("res"), jen.Id("ok")).Op(":=").Add(syncMapObjId.Dot("Load").Call(jen.Id("key"))),
		jen.If(jen.Op("!").Add().Id("ok")).Block(
			jen.Return(jen.List(jen.Id(syncMap.emptyValueName), jen.Id("ok"))),
		),
		jen.Return(jen.List(result, jen.Id("ok"))),
	)

	return res
}

func (syncMap *SyncMap) writeFuncDelete() jen.Code {
	funcName := "Delete"
	// func
	res := jen.Func()
	// func type params
	res = res.Params(syncMap.funcTypeParams())
	// func name
	res = res.Id(funcName)
	// func params
	res = res.Params(jen.List(jen.Id("key").Add(syncMap.KeyType())))

	// func body
	syncMapObjId := syncMap.syncMapObjId()
	res = res.Block(
		syncMapObjId.Dot("Delete").Call(jen.Id("key")),
	)
	return res
}

func (syncMap *SyncMap) writeFuncRange() jen.Code {
	funcName := "Range"
	// func
	res := jen.Func()
	// func type params
	res = res.Params(syncMap.funcTypeParams())
	// func name
	res = res.Id(funcName)
	// func params
	rangeParamFunc := jen.Func().Params(jen.List(
		jen.Id("key").Add(syncMap.KeyType()),
		jen.Id("value").Add(syncMap.ValueType()),
	)).Params(jen.Bool())
	res = res.Params(jen.Id("f").Add(rangeParamFunc))

	// func body
	convertRangeFunc := jen.Func().Params(jen.List(jen.Id("ikey"), jen.Id("ivalue").Add(jen.Interface()))).
		Params(jen.Bool()).Block(
		jen.Id("k").Op(":=").Add(jen.Id("ikey").Assert(syncMap.KeyType())),
		jen.Id("v").Op(":=").Add(jen.Id("ivalue").Assert(syncMap.ValueType())),
		jen.Return(jen.Id("f").Call(jen.Id("k"), jen.Id("v"))),
	)
	convertRangeFuncObj := jen.Id("rangeF").Op(":=").Add(convertRangeFunc)
	syncMapObjId := syncMap.syncMapObjId()
	res = res.Block(
		convertRangeFuncObj,
		syncMapObjId.Dot("Range").Call(jen.Id("rangeF")),
	)

	return res
}

func (syncMap *SyncMap) writeFuncLoadAndDelete() jen.Code {
	version := runtime.Version()
	version = version[2:]
	resVersion := strings.Split(version, ".")
	secondVersionStr := resVersion[1]
	secondVersion, err := strconv.ParseInt(secondVersionStr, 10, 64)
	if err != nil {
		logrus.Fatalf("cannot parse int, from %s to int", secondVersionStr)
	}
	if secondVersion < 15 {
		return jen.Empty()
	}

	funcName := "LoadAndDelete"
	// func
	res := jen.Func()
	// func type params
	res = res.Params(syncMap.funcTypeParams())
	// func name
	res = res.Id(funcName)
	// func params
	res = res.Params(jen.Id("key").Add(syncMap.KeyType())).Params(jen.List(
		syncMap.ValueType(),
		jen.Bool(),
	))

	// func body
	syncMapObjId := syncMap.syncMapObjId()
	loadAndDelete := jen.List(jen.Id("value"), jen.Id("ok")).Op(":=").Add(
		syncMapObjId.Dot("LoadAndDelete").Call(jen.Id("key")),
	)

	judgeOk := jen.If(jen.Op("!").Add(jen.Id("ok"))).Block(
		jen.Return(jen.Id(syncMap.emptyValueName), jen.Id("ok")),
	)

	res = res.Block(
		loadAndDelete,
		judgeOk,
		jen.Return(jen.Id("value").Assert(syncMap.ValueType()), jen.Id("ok")),
	)

	return res
}
