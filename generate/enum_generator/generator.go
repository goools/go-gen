package enum_generator

import (
	"fmt"
	"go/types"
	"os"
	"path/filepath"

	"github.com/dave/jennifer/jen"
	"github.com/goools/go-gen/generate"
	"github.com/goools/go-gen/packagex"
	"github.com/sirupsen/logrus"
)

type EnumOption struct {
	ConstValue int64  `json:"constValue"`
	Value      string `json:"value"`
	Doc        string `json:"doc"`
}

func (e EnumOption) MarshalJSON() ([]byte, error) {
	panic("implement me")
}

func (e EnumOption) UnmarshalJSON(bytes []byte) error {
	panic("implement me")
}

func NewEnum(pkgPath, enumName string, options []EnumOption) *Enum {
	return &Enum{
		PkgPath:      pkgPath,
		Name:         enumName,
		Options:      options,
		funcTypeName: fmt.Sprintf("*%s", enumName),
		funcThisName: "e",
	}
}

type Enum struct {
	PkgPath      string
	Name         string
	Options      []EnumOption
	funcTypeName string
	funcThisName string
}

func (enum *Enum) EnumCodeId(option EnumOption) string {
	return fmt.Sprintf("%s%s", enum.Name, option.Value)
}

func (enum *Enum) funcTypeParams() jen.Code {
	return jen.Id(enum.funcThisName).Id(enum.funcTypeName)
}

func (enum *Enum) writeString() jen.Code {
	funcName := "String"
	var cases []jen.Code
	for i := range enum.Options {
		enumValueName := fmt.Sprintf("%s%s", enum.Name, enum.Options[i].Value)
		cases = append(cases, jen.Case(jen.Id(enumValueName)).Block(
			jen.Return(jen.Lit(enum.Options[i].Value)),
		))
	}
	return jen.Func().Params(enum.funcTypeParams()).Id(funcName).Params().String().Block(
		jen.Switch(jen.Id(fmt.Sprintf("*%s", enum.funcThisName))).Block(cases...),
		jen.Return(jen.Lit("UNKNOWN")),
	)
}

func (enum *Enum) writeMarshalText() jen.Code {
	funcName := "MarshalText"
	var cases []jen.Code
	for i := range enum.Options {
		option := enum.Options[i]
		caseItem := jen.Case(jen.Id(enum.EnumCodeId(option))).Block(
			jen.Return(jen.List(jen.Index().Byte().Params(jen.Lit(enum.Options[i].Value)), jen.Nil())),
		)
		cases = append(cases, caseItem)
	}

	// func
	res := jen.Func()
	// func type params
	res = res.Params(enum.funcTypeParams())
	// func name
	res = res.Id(funcName)
	// func params
	res = res.Params().Params(jen.List(jen.Index().Byte(), jen.Error()))
	// func body
	unknownErr := jen.Qual("fmt", "Errorf").Call(
		jen.Lit(fmt.Sprintf("not found %s, value: %%v", enum.Name)),
		jen.Id(fmt.Sprintf("*%s", enum.funcThisName)),
	)
	res = res.Block(
		jen.Switch(jen.Id(fmt.Sprintf("*%s", enum.funcThisName))).Block(cases...),
		jen.Return(jen.List(jen.Index().Byte().Params(jen.Lit("UNKNOWN")), unknownErr)),
	)
	return res
}

func (enum *Enum) writeUnmarshalText() jen.Code {
	funcName := "UnmarshalText"
	// func
	res := jen.Func()
	// func type params
	res = res.Params(enum.funcTypeParams())
	// func name
	res = res.Id(funcName)
	// func params
	res = res.Params(jen.Id("enumBytes").Index().Byte()).Params(jen.Error())

	var cases []jen.Code
	for i := range enum.Options {
		option := enum.Options[i]
		caseItem := jen.Case(jen.Lit(option.Value)).Block(
			jen.Id(fmt.Sprintf("*%s", enum.funcThisName)).Op("=").Id(enum.EnumCodeId(option)),
			jen.Return(jen.Nil()),
		)
		cases = append(cases, caseItem)
	}

	enumStringId := "enumString"

	unknownErr := jen.Qual("fmt", "Errorf").Call(
		jen.Lit(fmt.Sprintf("not found %s, value: %%v", enum.Name)),
		jen.Id(enumStringId),
	)

	// func body
	res = res.Block(
		jen.Id(enumStringId).Op(":=").Id("string").Call(jen.Id("enumBytes")),
		jen.Switch(jen.Id(enumStringId)).Block(cases...),
		jen.Return(unknownErr),
	)
	return res
}

func NewEnumGenerator(pkg *packagex.Package) generate.Generator {
	return &EnumGenerator{
		pkg:   pkg,
		enums: map[*types.TypeName]*Enum{},
	}
}

type EnumGenerator struct {
	pkg   *packagex.Package
	enums map[*types.TypeName]*Enum
}

func (gen *EnumGenerator) Scan(enumNames ...string) {
	logrus.Debugf("enum_generator names: %v", enumNames)
	scanner := EnumScanner{
		pkg: gen.pkg,
	}
	for _, enumName := range enumNames {
		typeName := gen.pkg.TypeName(enumName)
		if typeName == nil {
			logrus.Fatalf("Not found enum_generator name %s", enumName)
		}
		gen.enums[typeName] = scanner.Scan(typeName)
	}
}

func (gen *EnumGenerator) WriteToFile() {
	packageFilePath, err := os.Getwd()
	if err != nil {
		logrus.Fatalf("get pwd have an err: %v", err)
	}
	for _, enum := range gen.enums {
		gen.writeToFile(packageFilePath, enum)
	}
}

func (gen *EnumGenerator) writeToFile(packageFilePath string, enum *Enum) {

	enumFileName := fmt.Sprintf("%s_generate.go", generate.ToSnakeCase(enum.Name))
	enumFilePath := filepath.Join(packageFilePath, enumFileName)
	logrus.Debugf("begin generate enum: %s, package path: %s, file path: %s",
		enum.Name, enum.PkgPath, enumFilePath)
	generateFile := jen.NewFilePath(enum.PkgPath)
	generateFile.HeaderComment(WriteDoNotEdit())
	generateFile.Add(enum.writeString())
	generateFile.Line()
	generateFile.Add(enum.writeMarshalText())
	generateFile.Line()
	generateFile.Add(enum.writeUnmarshalText())
	err := generateFile.Save(enumFilePath)
	if err != nil {
		logrus.Fatalf("save enum code to file have an err: %v, enum: %s, file: %s", err, enum.Name, enumFilePath)
	}
}

func WriteDoNotEdit() string {
	return "Code generated by go-gen DO NOT EDIT."
}
