package atomicvalue_generator

import (
	"fmt"
	"regexp"

	"github.com/dave/jennifer/jen"
	"github.com/goools/go-gen/generate"
	"github.com/goools/tools/strx"
	"github.com/sirupsen/logrus"
)

var (
	atomicValueDefRegexp = regexp.MustCompile(`([A-Za-z_][A-Za-z0-9]*)<([*A-Za-z][A-Za-z0-9/.]*)>`)
)

const (
	atomicValueDefRegexpTotal          = 3
	atomicValueDefRegexpNameIndex      = 1
	atomicValueDefRegexpValueTypeIndex = 2
	atomicValueFuncThisName            = "a"
)

type AtomicValue struct {
	PkgPath      string
	Name         string
	Value        *generate.Type
	funcTypeName string
	funcThisName string
}

func NewAtomicValue(pkgPath, typeDef string) *AtomicValue {
	atomicValue := &AtomicValue{
		PkgPath: pkgPath,
	}
	atomicValue.parseDef(typeDef)
	atomicValue.funcThisName = atomicValueFuncThisName
	atomicValue.funcTypeName = atomicValue.Name
	return atomicValue
}

func (atomicValue *AtomicValue) parseDef(typeDef string) {
	res := atomicValueDefRegexp.FindStringSubmatch(typeDef)
	if len(res) != atomicValueDefRegexpTotal {
		panic(fmt.Errorf("cannot parse atomic value def: %s", typeDef))
	}
	atomicValue.Name = res[atomicValueDefRegexpNameIndex]
	valueType := res[atomicValueDefRegexpValueTypeIndex]
	atomicValue.Value = generate.NewType(valueType)
}

func (atomicValue *AtomicValue) WriteToFile() {
	logrus.Infof("begin generate atomic value: %s", atomicValue.Name)
	atomicValueSnackName := strx.ToSnakeCase(atomicValue.Name)
	atomicValueFileName := fmt.Sprintf("%s_atomicvalue_generate.go", atomicValueSnackName)

	generateFile := jen.NewFilePath(atomicValue.PkgPath)
	generateFile.HeaderComment(generate.WriteDoNotEdit())

	generateFile.Add(atomicValue.writeTypeDef())
	generateFile.Line()
	generateFile.Add(atomicValue.writeStore())
	generateFile.Line()
	generateFile.Add(atomicValue.writeLoad())

	err := generateFile.Save(atomicValueFileName)
	if err != nil {
		logrus.Fatalf("save atomic value code to file have an err: %v, atomic value: %s, file: %s", err, atomicValue.Name, atomicValue)
	}
	logrus.Infof("complete generate atomic value: %s", atomicValue.Name)
}

func (atomicValue *AtomicValue) ValueType() jen.Code {
	return atomicValue.Value.PtrCode()
}

func (atomicValue *AtomicValue) PreFuncCommon(funcName string) *jen.Statement {
	return jen.Func().
		Params(jen.Id(atomicValue.funcThisName).Add(jen.Id(fmt.Sprintf("*%s", atomicValue.Name)))).Id(funcName)
}

func (atomicValue *AtomicValue) ObjCode() *jen.Statement {
	ptrOp := jen.Op("*")
	objCode := jen.Params(ptrOp.Add(jen.Qual("sync/atomic", "Value"))).Params(jen.Id(atomicValue.funcThisName))
	return objCode
}

func (atomicValue *AtomicValue) writeTypeDef() jen.Code {
	res := jen.Type()
	res.Id(atomicValue.Name).Qual("sync/atomic", "Value")
	return res
}

func (atomicValue *AtomicValue) writeStore() jen.Code {
	res := atomicValue.PreFuncCommon("Store")

	// func params
	res = res.Params(jen.Id("value").Add(atomicValue.ValueType()))

	// func body
	res = res.Block(
		atomicValue.ObjCode().Dot("Store").Call(jen.Id("value")),
	)
	return res
}

func (atomicValue *AtomicValue) writeLoad() jen.Code {
	res := atomicValue.PreFuncCommon("Load")

	// func params
	res = res.Params()

	// result params
	res = res.Params(atomicValue.ValueType())

	// func body
	objCode := atomicValue.ObjCode()

	res = res.Block(
		jen.Id("value").Op(":=").Add(objCode.Dot("Load")).Call(),
		jen.Return(jen.Id("value").Assert(atomicValue.ValueType())),
	)
	return res
}
