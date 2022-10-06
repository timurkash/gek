package args

import "fmt"

const description = `gek (generate kratos) is add-on util under go-kratos framework to create the microservices using http and grpc.

kratos is described in https://go-kratos.dev/

To install kratos cli use 
	go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
To upgrade us
	kratos upgrade

kratos uses
Ent:  https://entgo.io/
Wire: https://github.com/google/wire
gRPC: https://grpc.io/docs/languages/go/quickstart/

Options:
	-gen - project generation. .settings file required
	-utl - lists utils to be installed

.settings file described in`

func ShowDescription() {
	fmt.Println()
	fmt.Println("gek compiled from github.com/timurkash/gek")
	fmt.Println()
	fmt.Println(description)
	fmt.Println(`
type Settings struct {
	Name         string yaml:"name" valid:"required,lowercase"
	Version      string yaml:"version" valid:"required,lowercase"
	Service      string yaml:"service" valid:"required,lowercase"
	ConfigVolume string yaml:"configVolume" default:"/data/config"
	TemplateRepo string yaml:"templateRepo" valid:"required"
	ProtoService string yaml:"protoService" valid:"required"
	Port         *struct {
		Grpc int yaml:"grpc"
		Http int yaml:"http"
	} yaml:"port"

	Project      string yaml:"project" valid:"required,lowercase"
	Email        string yaml:"email" valid:"email"
	Description  string yaml:"description"

	CiCd *struct {
		Registry string yaml:"registry"
		Staging  string yaml:"staging"
		Helm     *struct {
			Repo   string yaml:"repo"
			Museum string yaml:"museum"
			Chart  string yaml:"chart"
		} yaml:"helm"
	} yaml:"cicd"
}

For example

name: service-name
version: v1
service: Service
kratosLayout: https://github.com/timurkash/kratos-layout.git
templateRepo: your template repo

project: the-map
email: t.kashaev@mapcard.pro
description: BackOffice Payment Router

The cimis of this bootstrapper is to rewrite files set in template repo settings.TempRepo
And the template project is ready to modify. The files to be rewrote is `)
	files := []string{
		".dockerignore",
		".gitignore",
		"Dockerfile",
		"Makefile",
		"README.md",
	}
	for _, file := range files {
		fmt.Printf(" - %s\n", file)
	}
}
