package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yxpod/pbconf/conf"
)

var (
	msg = flag.String("msg", "PBConf", "pb top message name")
	pkg = flag.String("pkg", "pbconf", "pb package name")
	f   = flag.String("f", "pbconf.proto", ".proto output file")
	o   = flag.String("o", "pbconf.pb", "encoded tables output file")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: pbconf [-msg msg] [-pkg pkg] [-f file] [-o file] table-path")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	conf, err := conf.LoadPBConf(flag.Args()[0])
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	ff, err := os.OpenFile(*f, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	defer ff.Close()
	conf.WriteProto(ff, *pkg, *msg)

	of, err := os.OpenFile(*o, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	defer of.Close()
	conf.WriteData(of)
}
