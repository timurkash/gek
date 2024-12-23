package args

import (
	"fmt"

	"github.com/timurkash/gek/utils"

	"github.com/fatih/color"
)

type Util struct {
	Name    string
	Command string
}

var green = color.New(color.FgGreen).SprintFunc()
var blue = color.New(color.BgBlue).SprintFunc()

var Utils = []Util{
	{"go", "https://go.dev/doc/install"},
	{"git", "sudo apt install git"},
	{"make", "sudo apt -y install make"},
	{"kratos", "go install github.com/go-kratos/kratos/cmd/kratos/v2@latest"},
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
	{"golangci-lint", "go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.50.1"},
	{"bloomrpc", `download from https://github.com/bloomrpc/bloomrpc/releases
run
	sudo dpkg -i ~/Downloads/bloomrpc*.deb`},
}

func ShowUtils() error {
	fmt.Println("required utils:")
	for _, util := range Utils {
		if err := utils.IsUtilExists(util.Name); err != nil {
			fmt.Printf(" - %s: To install run '%s'\n", util.Name, blue(util.Command))
		} else {
			fmt.Printf(" - %s: %s\n", green("installed"), util.Name)
		}
	}
	return nil
}
