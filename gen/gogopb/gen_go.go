package gogopb

import (
	"fmt"
	"github.com/davyxu/protoplus/codegen"
)

// 报错行号+7
const goCodeTemplate = `// Generated by github.com/davyxu/goobjfmt/objfmtgen
// DO NOT EDIT!

package {{.PackageName}}

import (
	"fmt"
	"reflect"
	_ "github.com/davyxu/cellnet/codec/binary"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
)
{{range $a, $enumobj := .Enums}}
type {{.Name}} int32
const (	{{range .Fields}}
	{{$enumobj.Name}}_{{.Name}} {{.Type}} = {{TagNumber $enumobj .}} {{end}}
)

var (
{{$enumobj.Name}}MapperValueByName = map[string]int32{ {{range .Fields}}
	"{{.Name}}": {{TagNumber $enumobj .}}, {{end}}
}

{{$enumobj.Name}}MapperNameByValue = map[int32]string{ {{range .Fields}}
	{{TagNumber $enumobj .}}: "{{.Name}}" , {{end}}
}
)

func (self {{$enumobj.Name}}) String() string {
	return {{$enumobj.Name}}MapperNameByValue[int32(self)]
}
{{end}}

{{range .Structs}}
{{ObjectLeadingComment .}}
type {{.Name}} struct{	{{range .Fields}}
	{{GoFieldName .}} {{GoTypeName .}} {{GoStructTag .}}{{FieldTrailingComment .}} {{end}}
}
{{end}}
{{range .Structs}}
func (self *{{.Name}}) String() string { return fmt.Sprintf("%+v",*self) } {{end}}

func init() {
	{{range .Structs}} {{ if IsMessage . }}
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("binary"),	
		Type:  reflect.TypeOf((*{{.Name}})(nil)).Elem(),
		ID:    {{StructMsgID .}},
	}) {{end}} {{end}}
}

`

func gen_go(ctx *Context) error {

	gen := codegen.NewCodeGen("go").
		RegisterTemplateFunc(codegen.UsefulFunc).
		ParseTemplate(goCodeTemplate, ctx).
		FormatGoCode()

	if gen.Error() != nil {
		fmt.Println(string(gen.Data()))
		return gen.Error()
	}

	return gen.WriteOutputFile(ctx.OutputFileName).Error()
}
