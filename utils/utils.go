package utils

import (
	"fmt"
	"gitlab.com/mcsolutions/tools/gek/utils/commands"
)

type Util struct {
	Name    string
	Command string
}

var Utils = []Util{
	{"go", "https://go.dev/doc/install"},
	{"git", "sudo apt install git"},
	{"make", "sudo apt -y install make"},
	{"kratos", "go install  github.com/go-kratos/kratos/cmd/kratos/v2@latest"},
	{"protoc", "sudo apt install -y protobuf-compiler"},
	{"protoc-gen-go", "kratos upgrade"},
	{"protoc-gen-go-grpc", "kratos upgrade"},
	{"protoc-gen-go-http", "kratos upgrade"},
	{"protoc-gen-grpc-web", "sudo npm i -g protoc-gen-grpc-web"},
	{"protoc-gen-ts", "sudo npm i -g protoc-gen-ts"},
	{"protoc-gen-go-errors", "kratos upgrade"},
	//{"protoc-gen-validate", "kratos upgrade"},
	{"wire", "go install github.com/google/wire/cmd/wire@latest"},
	{"protoc-gen-openapiv2", "go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest"},
	{"protoc-gen-doc", "go install github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc@latest"},
	{"ent", "go install entgo.io/ent/cmd/ent@latest"},
	{"goimports", "go install golang.org/x/tools/cmd/goimports@latest"},
	{"bloomrpc", `download from https://github.com/bloomrpc/bloomrpc/releases
run
	sudo dpkg -i ~/Downloads/bloomrpc*amd64.deb`},
}

func IsExistsAll() error {
	for _, util := range Utils {
		if err := IsExists(util.Name); err != nil {
			return fmt.Errorf(`%s not installed
to install: %s
total list of utils run: gek -utl`, util.Name, util.Command)
		}
	}
	return nil
}

func IsExists(util string) error {
	_, err := commands.Exec("which", util)
	return err
}
