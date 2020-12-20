package deepcopy_generator

import (
	"fmt"

	"github.com/goools/go-gen/packagex"
)

type DeepCopy struct {
	pkg          *packagex.Package
	Name         string
	funcTypeName string
	funcThisName string
}

func NewDeepCopy(pkg *packagex.Package, structName string) *DeepCopy {
	deepCopy := &DeepCopy{
		pkg:  pkg,
		Name: structName,
	}
	deepCopy.Init()
	return deepCopy
}

func (d *DeepCopy) Init() {
	typeName := d.pkg.TypeName(d.Name)
	if typeName == nil {
		panic(fmt.Errorf("cannot find struct %s", d.Name))
	}

	packagex.Struct(d.pkg.PkgPath)

	// structType := typeName.Type().(*types.Struct)
	// fieldNum := structType.NumFields()
	// for i := 0; i < fieldNum; i++ {
	// 	field := structType.Field(i)
	// 	logrus.Infof("field index: %d, field name: %s", i, field.Name())
	// }
	// logrus.Infof("typeName: %#v", typeName)
}

func (d *DeepCopy) WriteToFile() {

}
