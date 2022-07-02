package sql2struct

import (
	"fmt"
	"os"
	"text/template"

	"demo/ch01/internal/word"
)

// 预定义模板
/*
	此处的代码逻辑：
	首先是一个循环，{{range .Columns}} ... {{end}}
	在循环中包含了两个 if-else 语句：
	{{ if gt length 0}} ... {{end}} 注释长度大于0（即存在）则显示注释，否则显示字段名
	{{ if gt $type 0}} ... {{end}} 数据类型长度大于0（即存在）则显示类型和 JSON 标签，否则显示字段名
*/
const strcutTpl = `type {{.TableName | ToCamelCase}} struct {
{{range .Columns}}	{{ $length := len .Comment}} {{ if gt $length 0 }}// {{.Comment}} {{else}}// {{.Name}} {{ end }}
	{{ $typeLen := len .Type }} {{ if gt $typeLen 0 }}{{.Name | ToCamelCase}}	{{.Type}}	{{.Tag}}{{ else }}{{.Name}}{{ end }}
{{end}}}

func (model {{.TableName | ToCamelCase}}) TableName() string {
	return "{{.TableName}}"
}`

type StructTemplate struct {
	strcutTpl string
}

// 存储转换后 Go 结构体中的所有字段信息
type StructColumn struct {
	Name    string
	Type    string
	Tag     string
	Comment string
}

// 存储最终用于渲染的模板对象信息
type StructTemplateDB struct {
	TableName string
	Columns   []*StructColumn
}

func NewStructTemplate() *StructTemplate {
	return &StructTemplate{strcutTpl: strcutTpl}
}

// 对通过查询 COLUMNS 表所组装得到的 tbColumns 进行进一步的分解和转换
func (t *StructTemplate) AssemblyColumns(tbColumns []*TableColumn) []*StructColumn {
	tplColumns := make([]*StructColumn, 0, len(tbColumns))
	for _, column := range tbColumns {
		// 对 JSON Tag 的处理
		tag := fmt.Sprintf("`"+"json:"+"\"%s\""+"`", column.ColumnName)
		// 数据库类型到 Go 结构体的转换
		tplColumns = append(tplColumns, &StructColumn{
			Name:    column.ColumnName,
			Type:    DBTypeToStructType[column.DataType], // 进行简单的类型转换
			Tag:     tag,
			Comment: column.ColumnComment,
		})
	}

	return tplColumns
}

//
func (t *StructTemplate) Generate(tableName string, tplColumns []*StructColumn) error {
	// template.Must 包装对返回 (*Template, error) 的函数的调用，并在 error 为非 nil 时发生 panic
	// 声明了一个名为 sql2struct 的新模板对象
	// 定义了自定义函数 ToCamelCase，并与 word.UnderscoreToUpperCamelCase 方法进行绑定
	// 将文本解析为 t.strcutTpl 的模板主体
	tpl := template.Must(template.New("sql2struct").Funcs(template.FuncMap{
		"ToCamelCase": word.UnderscoreToUpperCamelCase,
	}).Parse(t.strcutTpl))

	// 组装符合预定义模板的模板对象
	tplDB := StructTemplateDB{
		TableName: tableName,
		Columns:   tplColumns,
	}
	// 进行渲染
	err := tpl.Execute(os.Stdout, tplDB)
	if err != nil {
		return err
	}

	return nil
}
