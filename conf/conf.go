package conf

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
)

type PBConf []*Table

func LoadPBConf(path string) (PBConf, error) {
	var conf PBConf

	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if f == nil || f.IsDir() || filepath.Ext(f.Name()) != ".xlsx" {
			return nil
		}

		t, err := LoadTable(path)
		if err != nil {
			return err
		}

		conf = append(conf, t)
		return nil
	})

	return conf, err
}

func (c PBConf) WriteProto(w io.Writer, pkg, topMsg string) {
	var buf bytes.Buffer

	p := func(s string) {
		buf.WriteString(s + "\n")
	}

	p(`syntax = "proto3";`)
	p("")
	p(`package ` + pkg + ";")
	p("")
	p(`message ` + topMsg + ` {`)

	for i, t := range c {
		t.WriteProto(&buf, i+1, "    ")
	}

	p(`}`)

	w.Write(buf.Bytes())
}

func (c PBConf) WriteData(w io.Writer) {
	for i, t := range c {
		w.Write(t.Encode(i + 1))
	}
}

func (c PBConf) WriteText(w io.Writer) {
	for _, t := range c {
		t.WriteText(w)
		w.Write([]byte("\n"))
	}
}
