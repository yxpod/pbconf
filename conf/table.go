package conf

import (
	"bytes"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/tealeg/xlsx"
)

type Table struct {
	Title string
	Names []string
	Types []string
	Datas [][]interface{}
}

func LoadTable(path string) (*Table, error) {
	f, err := xlsx.OpenFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed open xlsx file, path:%v, err:%v", path, err)
	}

	if len(f.Sheet) == 0 {
		return nil, fmt.Errorf("no sheet found, path:%v", path)
	}

	sheet := f.Sheets[0]

	if len(sheet.Rows) < 2 {
		return nil, fmt.Errorf("invalid xlsx file, need at least 2 rows, path:%v", path)
	}

	var t Table

	t.Title = strings.Title(strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)))

	for col, cell := range sheet.Rows[0].Cells {
		typ := strings.ToLower(strings.TrimSpace(cell.String()))
		if typ != "string" && typ != "int" && typ != "float" {
			return nil, fmt.Errorf("invalid field type `%v`, path:%v, col:%v", typ, path, col+1)
		}
		t.Types = append(t.Types, typ)
	}

	fieldNum := len(t.Types)

	if len(sheet.Rows[1].Cells) != fieldNum {
		return nil, fmt.Errorf("row fieldNum not match, path:%v, row:%v", path, 2)
	}

	for _, cell := range sheet.Rows[1].Cells {
		t.Names = append(t.Names, strings.Title(strings.TrimSpace(cell.String())))
	}

	for row := 2; row < len(sheet.Rows); row++ {
		rd := sheet.Rows[row]

		if len(rd.Cells) > 0 && strings.HasPrefix(rd.Cells[0].String(), "#") {
			continue
		}

		if len(rd.Cells) > fieldNum {
			return nil, fmt.Errorf("row fieldNum not match, path:%v, row:%v", path, row+1)
		}

		data := make([]interface{}, fieldNum)

		for col := 0; col < fieldNum; col++ {
			switch t.Types[col] {
			case "int":
				if col >= len(rd.Cells) {
					data[col] = 0
				} else if n, err := rd.Cells[col].Int(); err == nil {
					data[col] = n
				} else {
					return nil, fmt.Errorf("invalid field, path:%v, row:%v, col:%v", path, row+1, col+1)
				}

			case "float":
				if col >= len(rd.Cells) {
					data[col] = float32(0)
				} else if n, err := rd.Cells[col].Float(); err == nil {
					data[col] = float32(n)
				} else {
					return nil, fmt.Errorf("invalid field, path:%v, row:%v, col:%v", path, row+1, col+1)
				}

			case "string":
				if col >= len(rd.Cells) {
					data[col] = ""
				} else {
					data[col] = rd.Cells[col].String()
				}
			}
		}
		t.Datas = append(t.Datas, data)
	}

	return &t, nil
}

func (t *Table) WriteProto(w io.Writer, num int, indent string) {
	var buf bytes.Buffer

	p := func(s string) {
		if s != "" {
			buf.WriteString(indent + s + "\n")
		} else {
			buf.WriteString("\n")
		}
	}

	pbt := func(s string) string {
		switch s {
		case "int":
			return "int32 "
		case "float":
			return "float "
		case "string":
			return "string"
		}
		return ""
	}

	lower := func(s string) string {
		if s == "" {
			return s
		}

		l := strings.ToLower(s)
		return string(l[0]) + s[1:]
	}

	p("")
	p(`message ` + t.Title + ` {`)
	for i := 0; i < len(t.Types); i++ {
		p(fmt.Sprintf(`%v%v %v = %v;`, indent, pbt(t.Types[i]), lower(t.Names[i]), i+1))
	}
	p(`}`)
	p(fmt.Sprintf("repeated %v %v = %v;", t.Title, lower(t.Title), num))

	w.Write(buf.Bytes())
}

func (t *Table) Encode(n int) []byte {
	var b []byte
	for _, d := range t.Datas {
		var row []byte
		for i, v := range d {
			switch v.(type) {
			case int:
				row = append(row, packFieldInt(v.(int), i+1)...)
			case float32:
				row = append(row, packFieldFloat(v.(float32), i+1)...)
			case string:
				row = append(row, packFieldBytes([]byte(v.(string)), i+1)...)
			}
		}
		b = append(b, packFieldBytes(row, n)...)
	}
	return b
}
