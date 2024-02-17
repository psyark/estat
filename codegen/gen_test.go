package codegen

import (
	"fmt"
	"strings"
	"testing"

	"github.com/dave/jennifer/jen"
)

func TestXxx(t *testing.T) {
	gsdFile := jen.NewFile("estat")
	gsdFile.Comment("Code generated by github.com/psyark/estat/codegen; DO NOT EDIT.").Line()

	addClassAccesor(gsdFile, "area")
	addClassAccesor(gsdFile, "time")
	for i := 1; i <= 15; i++ {
		addClassAccesor(gsdFile, fmt.Sprintf("cat%02d", i))
	}

	addValueAccesor(gsdFile, "area")
	addValueAccesor(gsdFile, "time")
	for i := 1; i <= 15; i++ {
		addValueAccesor(gsdFile, fmt.Sprintf("cat%02d", i))
	}

	if err := gsdFile.Save("../getstatsdata.gen.go"); err != nil {
		t.Fatal(err)
	}
}

func TestList(t *testing.T) {
	listFile := jen.NewFile("estat")
	listFile.Comment("Code generated by github.com/psyark/estat/codegen; DO NOT EDIT.").Line()

	for _, name := range []string{"Class", "Note", "Value", "TableInf"} {
		addList(listFile, name)
	}

	if err := listFile.Save("../list.gen.go"); err != nil {
		t.Fatal(err)
	}
}

func addList(f *jen.File, name string) {
	listName := name + "List"

	f.Comment(fmt.Sprintf("%s は、%s または []%s を透過的にUnmarshalし、元通りにMarshalするスライスです", listName, name, name))

	f.Type().Id(listName).Index().Id(name)

	f.Func().Params(jen.Id("c").Op("*").Id(listName)).Id("UnmarshalJSON").Call(jen.Id("d").Index().Byte()).Error().Block(
		jen.Return().Id("unmarshalList").Call(
			jen.Call(jen.Op("*").Index().Id(name)).Call(jen.Id("c")),
			jen.Id("d"),
		),
	).Line()

	f.Func().Params(jen.Id("c").Id(listName)).Id("MarshalJSON").Call().Call(jen.Index().Byte(), jen.Error()).Block(
		jen.Return().Id("marshalList").Call(jen.Id("c")),
	).Line()
}

func addClassAccesor(f *jen.File, classID string) {
	funcName := strings.ToUpper(classID[0:1]) + classID[1:]

	f.Comment(fmt.Sprintf("%s は、IDが %q であるClassObjを返します", funcName, classID))
	f.Func().Params(jen.Id("c").Id("ClassInf")).Id(funcName).Params().Op("*").Id("ClassObj").Block(
		jen.Return().Id("c").Dot("GetClassObj").Call(jen.Lit(classID)),
	)
}

func addValueAccesor(f *jen.File, classID string) {
	propName := strings.ToUpper(classID[0:1]) + classID[1:]
	funcName := strings.ToUpper(classID[0:1]) + classID[1:] + "Class"

	f.Comment(fmt.Sprintf("%s は、このValueの %s に対応するClassを返します", funcName, propName))
	f.Func().Params(jen.Id("v").Id("Value")).Id(funcName).Params(jen.Id("c").Id("ClassInf")).Op("*").Id("Class").Block(
		jen.Return().Id("c").Dot("GetClassObj").Call(jen.Lit(classID)).Dot("GetClass").Call(jen.Id("v").Dot(propName)),
	)
}
