package codegen

import (
	"testing"

	"github.com/dave/jennifer/jen"
)

func TestXxx(t *testing.T) {
	f := jen.NewFile("estat")

	for _, name := range []string{"Class", "Note"} {
		f.Func().Params(jen.Id("c").Op("*").Id(name + "Helper")).Id("UnmarshalJSON").Call(jen.Id("d").Index().Byte()).Error().Block(
			jen.If(jen.Id("d").Index(jen.Lit(0)).Op("==").LitRune('{').Block(
				jen.Op("*").Id("c").Op("=").Make(jen.Index().Id(name), jen.Lit(1)),
				jen.Return().Qual("encoding/json", "Unmarshal").Call(jen.Id("d"), jen.Op("&").Call(jen.Op("*").Id("c")).Index(jen.Lit(0))),
			).Else().Block(
				jen.Return().Qual("encoding/json", "Unmarshal").Call(jen.Id("d"), jen.Call(jen.Op("*").Index().Id(name)).Call(jen.Id("c"))),
			)),
		).Line()
	}

	if err := f.Save("../estat.generated.go"); err != nil {
		t.Fatal(err)
	}
}

// func (c *NoteHelper) UnmarshalJSON(d []byte) error {
// 	if d[0] == '{' {
// 		*c = []Note{{}}
// 		return json.Unmarshal(d, &(*c)[0])
// 	} else {
// 		return json.Unmarshal(d, (*[]Note)(c))
// 	}
// }

// func (c NoteHelper) MarshalJSON() ([]byte, error) {
// 	if len(c) == 1 {
// 		return json.Marshal(c[0])
// 	} else {
// 		return json.Marshal([]Note(c))
// 	}
// }
