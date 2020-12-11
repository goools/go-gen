package enum_generator

import (
	"go/ast"
	"go/types"
	"strconv"
	"strings"

	"github.com/goools/go-gen/packagex"
	"github.com/sirupsen/logrus"
)

type EnumScanner struct {
	pkg *packagex.Package
}

func (scanner *EnumScanner) Scan(enum *types.TypeName) *Enum {
	pkgPath := enum.Pkg().Path()
	enumName := enum.Name()
	logrus.Debugf("[Scan] enum_generator name: %s, package path: %s", enumName, pkgPath)
	var options []EnumOption
	for ident, def := range scanner.pkg.TypesInfo.Defs {
		typeConst, ok := def.(*types.Const)
		if !ok {
			continue
		}
		if typeConst.Type() != enum.Type() {
			continue
		}
		constName := typeConst.Name()
		if !strings.HasPrefix(constName, enum.Name()) {
			continue
		}
		val := typeConst.Val()
		// label := ident.Obj.Decl.(*ast.ValueSpec).Comment.Text()
		doc := ident.Obj.Decl.(*ast.ValueSpec).Doc.Text()
		doc = strings.Trim(doc, " \n\t\r")
		valName := constName[len(enum.Name()):]
		logrus.Debugf("val: %d, doc: %s, valName: %s", val, doc, valName)
		intVal, err := strconv.ParseInt(val.String(), 10, 64)
		if err != nil {
			logrus.Fatalf("enum_generator Type: %s, val Name: %s, value not is a int, val: %s",
				enum.Name(), constName, val.String())
		}
		options = append(options, EnumOption{
			ConstValue: intVal,
			Value:      valName,
			Doc:        doc,
		})
	}
	return NewEnum(pkgPath, enumName, options)
}
