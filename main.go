package main

import (
	"flag"
	"github.com/timurkash/gek/args"
	"log"
	"os"
)

var (
	utl = flag.Bool("utl", false, "required utils")
	chk = flag.Bool("chk", false, "check if all utils is installed")
	gen = flag.Bool("gen", false, "generate the service project")
	mes = flag.Bool("mes", false, "adjust protobuf messages to json one")
	//ser = flag.Bool("ser", false, "some replaces in internal/service")
	crs = flag.Bool("crs", false, "adding cors to http")
	tsi = flag.Bool("tsi", false, "ts ignore unused interface")
)

var argFunc func() error

func main() {
	log.SetPrefix("[>error<] ")
	log.SetFlags(0)
	flag.Parse()
	if len(os.Args) == 1 {
		args.ShowDescription()
		return
	}
	switch {
	case *utl:
		argFunc = args.ShowUtils
	case *chk:
		argFunc = args.Check
	case *gen:
		argFunc = args.Generate
	case *mes:
		argFunc = args.MessagesServer
	//case *ser:
	//	argFunc = args.Services
	case *crs:
		argFunc = args.Cors
	case *tsi:
		argFunc = args.TsIgnore
	default:
		log.Fatalln("unknown option")
	}
	if err := argFunc(); err != nil {
		log.Fatalln(err)
	}
}
