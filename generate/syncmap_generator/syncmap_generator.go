package syncmap_generator

import (
	"github.com/dave/jennifer/jen"
	"github.com/sirupsen/logrus"
	"runtime"
	"strconv"
	"strings"
)

func (syncMap *SyncMap) writeEmptyValue() jen.Code {
	emptyFunc := jen.Func().Params().Params(jen.Id("val").Add(jen.Id(syncMap.ValueType))).Block(jen.Return())
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
	res = res.Params(jen.List(jen.Id("key").Id(syncMap.KeyType), jen.Id("value").Id(syncMap.ValueType)))
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
	res = res.Params(jen.List(jen.Id("key").Id(syncMap.KeyType), jen.Id("value").Id(syncMap.ValueType))).
		Params(jen.List(jen.Id(syncMap.KeyType), jen.Bool()))

	// func body
	syncMapObjId := syncMap.syncMapObjId()
	result := jen.Id("res").Assert(jen.Id(syncMap.ValueType))
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
	res = res.Params(jen.List(jen.Id("key").Id(syncMap.KeyType))).Params(jen.List(jen.Id(syncMap.KeyType), jen.Bool()))

	// func body
	syncMapObjId := syncMap.syncMapObjId()
	result := jen.Id("res").Assert(jen.Id(syncMap.ValueType))
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
	res = res.Params(jen.List(jen.Id("key").Id(syncMap.KeyType)))

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
		jen.Id("key").Add(jen.Id(syncMap.KeyType)),
		jen.Id("value").Add(jen.Id(syncMap.ValueType)),
	)).Params(jen.Bool())
	res = res.Params(jen.Id("f").Add(rangeParamFunc))

	// func body
	convertRangeFunc := jen.Func().Params(jen.List(jen.Id("ikey"), jen.Id("ivalue").Add(jen.Interface()))).
		Params(jen.Bool()).Block(
		jen.Id("k").Op(":=").Add(jen.Id("ikey").Assert(jen.Id(syncMap.KeyType))),
		jen.Id("v").Op(":=").Add(jen.Id("ivalue").Assert(jen.Id(syncMap.ValueType))),
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
	res = res.Params(jen.Id("key").Add(jen.Id(syncMap.KeyType))).Params(jen.List(
		jen.Id(syncMap.ValueType),
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
		jen.Return(jen.Id("value").Assert(jen.Id(syncMap.ValueType)), jen.Id("ok")),
	)

	return res
}
