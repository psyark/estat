package codegen

import (
	"fmt"
	"strings"
	"testing"

	"github.com/dave/jennifer/jen"
)

func TestXxx(t *testing.T) {
	f := jen.NewFile("estat")
	f.Comment("Code generated by github.com/psyark/estat/codegen; DO NOT EDIT.").Line()

	// TODO TableInfはファイルを分ける
	for _, name := range []string{"Class", "Note", "Value", "TableInf"} {
		helperName := name + "Helper"

		f.Comment(fmt.Sprintf("%s は、%s の単一の値または配列を透過的にUnmarshal/Marshalするスライスです", helperName, name))

		f.Type().Id(helperName).Index().Id(name)

		f.Func().Params(jen.Id("c").Op("*").Id(helperName)).Id("UnmarshalJSON").Call(jen.Id("d").Index().Byte()).Error().Block(
			jen.If(jen.Id("d").Index(jen.Lit(0)).Op("==").LitRune('{').Block(
				jen.Op("*").Id("c").Op("=").Make(jen.Index().Id(name), jen.Lit(1)),
				jen.Return().Qual("encoding/json", "Unmarshal").Call(jen.Id("d"), jen.Op("&").Call(jen.Op("*").Id("c")).Index(jen.Lit(0))),
			).Else().Block(
				jen.Return().Qual("encoding/json", "Unmarshal").Call(jen.Id("d"), jen.Call(jen.Op("*").Index().Id(name)).Call(jen.Id("c"))),
			)),
		).Line()

		f.Func().Params(jen.Id("c").Id(helperName)).Id("MarshalJSON").Call().Call(jen.Index().Byte(), jen.Error()).Block(
			jen.If(jen.Id("len").Call(jen.Id("c")).Op("==").Lit(1).Block(
				jen.Return().Qual("encoding/json", "Marshal").Call(jen.Id("c").Index(jen.Lit(0))),
			).Else().Block(
				jen.Return().Qual("encoding/json", "Marshal").Call(jen.Index().Id(name).Call(jen.Id("c"))),
			)),
		).Line()
	}

	for _, classId := range []string{"area", "time", "cat01", "cat02"} {
		upper := strings.ToUpper(classId[0:1]) + classId[1:]

		f.Func().Params(jen.Id("c").Id("ClassInf")).Id(upper).Params().Op("*").Id("ClassObj").Block(
			jen.Return().Id("c").Dot("GetClassObj").Call(jen.Lit(classId)),
		)
	}

	if err := f.Save("../getstatsdata.gen.go"); err != nil {
		t.Fatal(err)
	}
}
