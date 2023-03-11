package utils

import (
	"github.com/timurkash/gek/utils/commands"
	"log"
	"os"
)

type Util struct {
	Name    string
	Command string
}

const (
	Api   = "api"
	Src   = "src"
	Back  = "back"
	Proto = "proto"
)

func IsUtilExists(util string) error {
	_, err := commands.Exec("which", util)
	return err
}

func IsFileExists(filename string) bool {
	fi, err := os.Stat(filename)
	if err == nil {
		if !fi.IsDir() {
			return true
		} else {
			log.Println(filename, "is directory")
			return false
		}
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func IsDirExists(dirname string) bool {
	fi, err := os.Stat(dirname)
	if err == nil {
		if fi.IsDir() {
			return true
		} else {
			log.Println(dirname, "is not directory")
			return false
		}
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func IsExists(filename string) bool {
	_, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
