/**
 * Copyright 2015 @ at3.net.
 * name : tool.go
 * author : jarryliu
 * date : 2016-11-11 12:19
 * description :
 * history :
 */
package generator

import (
	"bytes"
	"github.com/jsix/gof/db/orm"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
	"text/template"
	"unicode"
)

var (
	emptyReg             = regexp.MustCompile("\\s+\"\\s*\"\\s*\\n")
	emptyImportReg       = regexp.MustCompile("import\\s*\\(([\\n\\s\"]+)\\)")
	revertTemplateRegexp = regexp.MustCompile("([^\\$])\\${([^\\}]+)\\}")
)

const (
	//模型包名
	V_ModelPkgName = "ModelPkgName"
	//仓储结构包名
	V_RepoPkgName = "RepoPkgName"
	//仓储接口包名
	V_IRepoPkgName = "IRepoPkgName"
	//仓储结构引用模型包路径
	V_ModelPkg = "ModelPkg"
	//仓储接口引用模型包路径
	V_IRepoPkg = "IRepoPkg"
	// 仓储包路径
	V_RepoPkg = "RepoPkg"
)

type (
	// 表
	Table struct {
		// 表名
		Name string
		// 表前缀
		Prefix string
		// 表名单词首字大写
		Title string
		// 表注释
		Comment string
		Engine  string
		Charset string
		Raw     *orm.Table
		Columns []*Column
	}
	// 列
	Column struct {
		// 表名
		Name string
		// 表名首字大写
		Title   string
		Pk      bool
		Auto    bool
		NotNull bool
		Type    string
		Comment string
	}
)
type Session struct {
	//生成代码变量
	codeVars map[string]interface{}
	IdUpper  bool
}

// 数据库代码生成器
func DBCodeGenerator() *Session {
	return (&Session{
		codeVars: make(map[string]interface{}),
	}).init()
}

func (s *Session) init() *Session {
	s.Var(V_ModelPkgName, "model")
	s.Var(V_RepoPkgName, "repo")
	s.Var(V_IRepoPkgName, "repo")
	s.Var(V_ModelPkg, "")
	s.Var(V_IRepoPkg, "")
	s.Var(V_RepoPkg, "")
	return s
}

func (s *Session) title(str string) string {
	// 小于3且ID大写，则返回大写
	if s.IdUpper && len(str) < 3 {
		return strings.ToUpper(str)
	}
	arr := strings.Split(str, "_")
	for i, v := range arr {
		arr[i] = strings.Title(v)
	}
	return strings.Join(arr, "")
}

func (s *Session) prefix(str string) string {
	if i := strings.Index(str, "_"); i != -1 {
		return str[:i]
	}
	for i, l := 1, len(str); i < l-1; i++ {
		if unicode.IsUpper(rune(str[i])) {
			return strings.ToLower(str[:i])
		}
	}
	return ""
}

func (s *Session) goType(dbType string) string {
	l := len(dbType)
	switch true {
	case strings.HasPrefix(dbType, "tinyint"):
		return "int"
	case strings.HasPrefix(dbType, "bit"):
		return "bool"
	case strings.HasPrefix(dbType, "int("):
		if l == 6 {
			return "int"
		}
		return "int64"
	case strings.HasPrefix(dbType, "float"):
		return "float32"
	case strings.HasPrefix(dbType, "decimal"):
		return "float64"
	case dbType == "text", strings.HasPrefix(dbType, "varchar"):
		return "string"
	}
	return "interface{}"
}

// 获取所有的表
func (s *Session) ParseTables(tbs []*orm.Table, err error) ([]*Table, error) {
	n := make([]*Table, len(tbs))
	for i, tb := range tbs {
		n[i] = s.ParseTable(tb)
	}
	return n, err
}

// 获取表结构
func (s *Session) ParseTable(tb *orm.Table) *Table {
	n := &Table{
		Name:    tb.Name,
		Prefix:  s.prefix(tb.Name),
		Title:   s.title(tb.Name),
		Comment: tb.Comment,
		Engine:  tb.Engine,
		Charset: tb.Charset,
		Raw:     tb,
		Columns: make([]*Column, len(tb.Columns)),
	}
	for i, v := range tb.Columns {
		n.Columns[i] = &Column{
			Name:    v.Name,
			Title:   s.title(v.Name),
			Pk:      v.Pk,
			Auto:    v.Auto,
			NotNull: v.NotNull,
			Type:    v.Type,
			Comment: v.Comment,
		}
	}
	return n
}

// 表生成结构
func (s *Session) TableToGoStruct(tb *Table) string {
	if tb == nil {
		return ""
	}
	pkgName := ""
	if p, ok := s.codeVars[V_ModelPkgName]; ok {
		pkgName = p.(string)
	} else {
		pkgName = "model"
	}

	//log.Println(fmt.Sprintf("%#v", tb))
	buf := bytes.NewBufferString("")
	buf.WriteString("package ")
	buf.WriteString(pkgName)

	buf.WriteString("\n// ")
	buf.WriteString(tb.Comment)
	buf.WriteString("\ntype ")
	buf.WriteString(s.title(tb.Name))
	buf.WriteString(" struct{\n")

	for _, col := range tb.Columns {
		if col.Comment != "" {
			buf.WriteString("    // ")
			buf.WriteString(col.Comment)
			buf.WriteString("\n")
		}
		buf.WriteString("    ")
		buf.WriteString(s.title(col.Name))
		buf.WriteString(" ")
		buf.WriteString(s.goType(col.Type))
		buf.WriteString(" `")
		buf.WriteString("db:\"")
		buf.WriteString(col.Name)
		buf.WriteString("\"")
		if col.Pk {
			buf.WriteString(" pk:\"yes\"")
		}
		if col.Auto {
			buf.WriteString(" auto:\"yes\"")
		}
		buf.WriteString("`")
		buf.WriteString("\n")
	}

	buf.WriteString("}")
	return buf.String()
}

// 解析模板
func (s *Session) Resolve(t CodeTemplate) CodeTemplate {
	t = resolveRepTag(t)
	return t
}

// 定义变量或修改变量
func (s *Session) Var(key string, v interface{}) {
	if v == nil {
		delete(s.codeVars, key)
		return
	}
	//if strings.HasSuffix(key, "PkgName") {
	//	if s := v.(string); s != "" && s[len(s)-1] != '.' {
	//		v = s + "."
	//	}
	//}
	s.codeVars[key] = v
}

// 还原模板的标签: ${...} -> {{...}}, $$ -> $
func (s *Session) revertTemplateVariable(str string) string {
	str = revertTemplateRegexp.ReplaceAllString(str, "$1{{$2}}")
	return strings.Replace(str, "$$", "$", -1)
}

// 转换成为模板
func (s *Session) ParseTemplate(file string) (CodeTemplate, error) {
	data, err := ioutil.ReadFile(file)
	if err == nil {
		return CodeTemplate(string(data)), nil
	}
	return CodeTemplate(""), err
}

// 生成代码
func (s *Session) GenerateCode(tb *Table, tpl CodeTemplate,
	structSuffix string, sign bool, ePrefix string) string {
	if tb == nil {
		return ""
	}

	var err error
	t := &template.Template{}
	t, err = t.Parse(string(tpl))
	if err != nil {
		panic(err)
	}

	pk := "<PK>"
	for i, v := range tb.Columns {
		if i == 0 {
			pk = v.Name
		}
		if v.Pk {
			pk = v.Name
			break
		}
	}
	n := s.title(tb.Name)
	r2 := ""
	if sign {
		r2 = n
	}
	mp := map[string]interface{}{
		"VAR": s.codeVars,
		"T":   tb,
		"R":   n + structSuffix,
		"R2":  r2,
		"E":   n,
		"E2":  ePrefix + n,
		"Ptr": strings.ToLower(tb.Name[:1]),
		"PK":  s.title(pk),
	}
	buf := bytes.NewBuffer(nil)
	err = t.Execute(buf, mp)
	if err == nil {
		code := buf.String()
		//去除空引用
		code = emptyImportReg.ReplaceAllString(code, "")
		//如果不包含模型，则可能为引用空的包
		code = emptyReg.ReplaceAllString(code, "")
		return s.revertTemplateVariable(code)
	}
	log.Println("execute template error:", err.Error())
	return ""
}

func (s *Session) GenerateTablesCode(tables []*Table, tpl CodeTemplate) string {
	if tables == nil || len(tables) == 0 {
		return ""
	}

	var err error
	t := &template.Template{}
	t, err = t.Parse(string(tpl))
	if err != nil {
		panic(err)
	}
	mp := map[string]interface{}{
		"VAR":    s.codeVars,
		"Tables": tables,
	}
	buf := bytes.NewBuffer(nil)
	err = t.Execute(buf, mp)
	if err == nil {
		code := buf.String()
		//去除空引用
		code = emptyImportReg.ReplaceAllString(code, "")
		//如果不包含模型，则可能为引用空的包
		code = emptyReg.ReplaceAllString(code, "")
		return s.revertTemplateVariable(code)
	}
	log.Println("execute template error:", err.Error())
	return ""
}

// 表生成仓储结构,sign:函数后是否带签名，ePrefix:实体是否带前缀
func (s *Session) TableToGoRepo(tb *Table,
	sign bool, ePrefix string) string {
	return s.GenerateCode(tb, TPL_ENTITY_REP,
		"Repo", sign, ePrefix)
}

// 表生成仓库仓储接口
func (s *Session) TableToGoIRepo(tb *Table,
	sign bool, ePrefix string) string {
	return s.GenerateCode(tb, TPL_ENTITY_REP_INTERFACE,
		"Repo", sign, ePrefix)
}