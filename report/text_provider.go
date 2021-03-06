package report

import (
	"bytes"
	"strings"
)

var _ IExportProvider = new(TextProvider)

type TextProvider struct {
	delimer string
}

func NewTextProvider() IExportProvider {
	return &TextProvider{
		delimer: ",",
	}
}

func (t *TextProvider) Export(rows []map[string]interface{},
	fields []string, names []string, formatter []IExportFormatter) (binary []byte) {
	buf := bytes.NewBufferString("")
	// 显示表头
	showHeader := fields != nil && len(fields) > 0
	if showHeader {
		for i, k := range names {
			if i > 0 {
				buf.WriteString(t.delimer)
			}
			buf.WriteString(k)
		}
	}
	l := len(rows)
	for i, row := range rows {
		if i < l {
			buf.WriteString("\n")
		}
		for ki, field := range fields {

			data := row[field]
			if formatter != nil {
				for _, f := range formatter {
					data = f.Format(field, names[ki], i, data)
				}
			}
			t.appendField(buf, ki, data)
		}
	}
	return buf.Bytes()
}

func (t *TextProvider) appendField(buf *bytes.Buffer, ki int, data interface{}) {
	if ki > 0 {
		buf.WriteString(t.delimer)
	}
	dataStr := data.(string)
	specData := strings.Index(dataStr, " ") != -1 ||
		strings.Index(dataStr, "-") != -1 ||
		strings.Index(dataStr, "'") != -1

	if strings.Index(dataStr, "\"") != -1 {
		dataStr = strings.Replace(dataStr, "\"", "\"\"", -1)
		specData = true
	}
	//防止里面含有特殊符号
	if specData {
		buf.WriteString("\"")
		buf.WriteString(dataStr)
		buf.WriteString("\"")
	} else {
		buf.WriteString(dataStr)
	}
}
