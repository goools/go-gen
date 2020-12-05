package generate

import (
	"github.com/goools/go-gen/packagex"
	"github.com/sirupsen/logrus"
	"go/ast"
	"go/types"
	"strings"
)

func NewEnumGenerator(pkg *packagex.Package) Generator {
	return &EnumGenerator{pkg: pkg}
}

type EnumGenerator struct {
	pkg   *packagex.Package
	enums map[*types.TypeName]*Enum
}

func (gen *EnumGenerator) parseEnum(enum *types.TypeName) {
	pkg := enum.Pkg()
	pkgPath := pkg.Path()
	logrus.Debugf("parse enum, enum name: %s, package path: %s", enum.Name(), pkgPath)
	for ident, def := range gen.pkg.TypesInfo.Defs {
		typeConst, ok := def.(*types.Const)
		if !ok {
			continue
		}
		if typeConst.Type() != enum.Type() {
			continue
		}
		name := typeConst.Name()
		if !strings.HasPrefix(name, enum.Name()) {
			continue
		}
		val := typeConst.Val()
		label := ident.Obj.Decl.(*ast.ValueSpec).Comment.Text()
		doc := ident.Obj.Decl.(*ast.ValueSpec).Doc.Text()
		logrus.Debugf("val: %s, label: %s, doc: %s", val, label, doc)
	}
}

func (gen *EnumGenerator) Scan(enumNames ...string) {
	logrus.Debugf("enum names: %v", enumNames)
	for _, enumName := range enumNames {
		typeName := gen.pkg.TypeName(enumName)
		if typeName == nil {
			logrus.Fatalf("Not found enum name %s", enumName)
		}
		gen.parseEnum(typeName)
	}
}

func (gen *EnumGenerator) WriteToFile(filePath string) {
	logrus.Infof("Write to file, file path: %s", filePath)
}

type EnumOption struct {
	ConstValue int    `json:"constValue"`
	Value      string `json:"value"`
	Label      string `json:"label"`
}

func NewEnum(pkgTypeOrName string, options []EnumOption) *Enum {
	parts := strings.Split(pkgTypeOrName, ".")
	pkgPath, name := "", ""

	switch len(parts) {
	case 1:
		name = parts[0]
	default:
		pkgPath = strings.Join(parts[0:len(parts)-1], ".")
		name = parts[len(parts)-1]
	}

	return &Enum{
		PkgPath: pkgPath,
		Name:    name,
		Options: options,
	}
}

type Enum struct {
	PkgPath string
	Name    string
	Options []EnumOption
}
