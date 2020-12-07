package enum_generator

import (
	"fmt"
	"github.com/goools/go-gen/generate"
	"github.com/goools/go-gen/packagex"
	"github.com/sirupsen/logrus"
	"go/types"
	"os"
	"path/filepath"
)

type EnumOption struct {
	ConstValue int64  `json:"constValue"`
	Value      string `json:"value"`
	Doc        string `json:"doc"`
}

func NewEnum(pkgPath, enumName string, options []EnumOption) *Enum {
	return &Enum{
		PkgPath: pkgPath,
		Name:    enumName,
		Options: options,
	}
}

type Enum struct {
	PkgPath string
	Name    string
	Options []EnumOption
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
		enumFileName := fmt.Sprintf("%s_generate.go", generate.ToSnakeCase(enum.Name))
		enumFilePath := filepath.Join(packageFilePath, enumFileName)
		logrus.Infof("begin generate enum: %s, package path: %s, file path: %s",
			enum.Name, enum.PkgPath, enumFilePath)
	}
}
