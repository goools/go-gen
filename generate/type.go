package generate

import (
	"fmt"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/sirupsen/logrus"
)

type Type struct {
	packagePath string
	typeName    string
	code        jen.Code
}

func (t *Type) PtrCode() jen.Code {
	code := t.Code()
	return jen.Op("*").Add(code)
}

func (t *Type) Code() jen.Code {
	if t.code != nil {
		return t.code
	}
	if t.packagePath == "" {
		t.code = jen.Id(t.typeName)
	} else {
		t.code = jen.Qual(t.packagePath, t.typeName)
	}
	return t.code
}

func (t *Type) TypeName() string {
	if t.packagePath != "" {
		pkgPath := t.packagePath
		pkgPath = strings.ReplaceAll(pkgPath, ".", "_")
		pkgPath = strings.ReplaceAll(pkgPath, "/", "__")
		return fmt.Sprintf("%s_%s", pkgPath, t.typeName)
	}
	return t.typeName
}

func NewType(typeDef string) *Type {
	typeDef = strings.Trim(typeDef, "*")
	index := strings.LastIndex(typeDef, ".")
	var res *Type
	if index == -1 {
		res = &Type{
			typeName: typeDef,
		}
	} else {
		res = &Type{
			packagePath: typeDef[:index],
			typeName:    typeDef[index+1:],
		}
	}
	logrus.Debugf("type: %#v", res)
	return res
}
