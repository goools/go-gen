package enum_generator

import (
	"fmt"

	"github.com/dave/jennifer/jen"
)

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
