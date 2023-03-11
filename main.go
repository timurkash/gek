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
	htp = flag.Bool("htp", false, "generate empty http-server")
	mes = flag.Bool("mes", false, "adjust protobuf messages to json one")
)

var argFunc func() error

func main() {
	log.SetPrefix("[>error<] ")
	log.SetFlags(0)
	flag.Parse()
	if len(os.Args) != 2 {
		args.ShowDescription()
	} else {
		if *utl {
			argFunc = args.ShowUtils
		} else if *chk {
			argFunc = args.Check
		} else if *gen {
			argFunc = args.Generate
		} else if *htp {
			argFunc = args.HttpServer
		} else if *mes {
			argFunc = args.MessagesServer
		} else {
			log.Fatalln("unknown option")
		}
		if err := argFunc(); err != nil {
			log.Fatalln(err)
		}
	}
}
