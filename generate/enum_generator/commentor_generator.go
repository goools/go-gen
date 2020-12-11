package enum_generator

import (
	"fmt"

	"github.com/dave/jennifer/jen"
)

func (enum *Enum) writeComment() jen.Code {
	funcName := "Comment"
	var cases []jen.Code
	for i := range enum.Options {
		option := enum.Options[i]
		comment := option.Doc
		if comment == "" {
			comment = option.Value
		}
		caseItem := jen.Case(jen.Id(enum.EnumCodeId(option))).Block(
			jen.Return(jen.Lit(comment)),
		)
		cases = append(cases, caseItem)
	}
	defaultCase := jen.Default().Block(
		jen.Return(jen.Lit("UNKNOWN")),
	)
	cases = append(cases, defaultCase)
	// func
	res := jen.Func()
	// func type params
	res = res.Params(enum.funcTypeParams())
	// func name
	res = res.Id(funcName)
	// func params
	res = res.Params().Params(jen.String())
	// func body
	res = res.Block(
		jen.Switch(jen.Id(fmt.Sprintf("*%s", enum.funcThisName))).Block(cases...),
	)
	return res
}
