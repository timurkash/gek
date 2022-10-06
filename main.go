package main

import (
	"fmt"
	"github.com/timurkash/gek/args"
	"github.com/timurkash/gek/utils"
	"log"
	"os"
)

func main() {
	log.SetPrefix("[>error<] ")
	log.SetFlags(0)
	argsStrings := os.Args
	if len(argsStrings) == 1 {
		args.ShowDescription()
		return
	}
	arg := argsStrings[1]
	if len(arg) != 4 {
		log.Fatalln("argument must be 4 characters")
	}
	if arg[0] != '-' {
		log.Fatalln("argument must begin with dash")
	}
	switch arg {
	case "-utl":
		fmt.Println("required utils:")
		for _, util := range utils.Utils {
			if err := utils.IsUtilExists(util.Name); err != nil {
				fmt.Printf(" - %s: To install run '%s'\n", util.Name, util.Command)
			} else {
				fmt.Printf(" - %s: installed\n", util.Name)
			}
		}
	case "-gen":
		if err := args.ArgGen(); err != nil {
			log.Fatalln(err)
		}
	case "-htp":
		if err := args.ArgHttpServer(); err != nil {
			log.Fatalln(err)
		}
	case "-mes":
		if err := args.ArgMessagesServer(); err != nil {
			log.Fatalln(err)
		}
	default:
		log.Fatalf("option %s not defined\n", arg)
	}
}
