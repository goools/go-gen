package syncmap_generator

import (
	"github.com/dave/jennifer/jen"
)

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
	res.Block(
		jen.Id(syncMap.funcThisName).Dot("Store").Call(jen.Id("key"), jen.Id("value")),
	)

	return res
}

