package enum_generator

import (
	"fmt"
	"go/types"
	"os"
	"path/filepath"

	"github.com/dave/jennifer/jen"
	"github.com/goools/go-gen/generate"
	"github.com/goools/go-gen/packagex"
	"github.com/goools/tools/strx"
	"github.com/sirupsen/logrus"
)

const (
	enumFuncThisName = "e"
)

type EnumOption struct {
	ConstValue int64  `json:"constValue"`
	Value      string `json:"value"`
	Doc        string `json:"doc"`
}

func NewEnum(pkgPath, enumName string, options []EnumOption) *Enum {
	return &Enum{
		PkgPath:      pkgPath,
		Name:         enumName,
		Options:      options,
		funcTypeName: fmt.Sprintf("*%s", enumName),
		funcThisName: enumFuncThisName,
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

	enumFileName := fmt.Sprintf("%s_enum_generate.go", strx.ToSnakeCase(enum.Name))
	enumFilePath := filepath.Join(packageFilePath, enumFileName)
	logrus.Debugf("begin generate enum: %s, package path: %s, file path: %s",
		enum.Name, enum.PkgPath, enumFilePath)
	generateFile := jen.NewFilePath(enum.PkgPath)
	generateFile.HeaderComment(generate.WriteDoNotEdit())
	generateFile.Add(enum.writeString())
	generateFile.Line()
	generateFile.Add(enum.writeMarshalText())
	generateFile.Line()
	generateFile.Add(enum.writeUnmarshalText())
	generateFile.Line()
	generateFile.Add(enum.writeComment())
	err := generateFile.Save(enumFilePath)
	if err != nil {
		logrus.Fatalf("save enum code to file have an err: %v, enum: %s, file: %s", err, enum.Name, enumFilePath)
	}
}
