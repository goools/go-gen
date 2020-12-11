package enum_generator

import (
	"fmt"

	"github.com/dave/jennifer/jen"
)

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
