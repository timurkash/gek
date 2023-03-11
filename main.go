package main

import (
	"errors"
	"flag"
	"github.com/timurkash/gek/args"
	"log"
	"os"
)

var (
	utl = flag.Bool("utl", false, "required utils")
	chk = flag.Bool("chk", false, "check if all utils is installed")
	gen = flag.Bool("gen", false, "generate the service project")
	htp = flag.Bool("htp", false, "generate empty http-server")
	mes = flag.Bool("mes", false, "adjust protobuf messages to json one")
)

func main() {
	log.SetPrefix("[>error<] ")
	log.SetFlags(0)
	flag.Parse()
	argsStrings := os.Args
	if len(argsStrings) != 2 {
		args.ShowDescription()
		return
	}
	var err error
	if *utl {
		args.ShowUtils()
	} else if *chk {
		err = args.Check()
	} else if *gen {
		err = args.Generate()
	} else if *htp {
		err = args.HttpServer()
	} else if *mes {
		err = args.MessagesServer()
	} else {
		err = errors.New("unknown option")
	}
	if err != nil {
		log.Fatalln(err)
	}
}
