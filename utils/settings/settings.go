package settings

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/timurkash/gek/utils"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

const FileName = ".settings"

type (
	Settings struct {
		Name           string `yaml:"name" valid:"required,lowercase"`
		Version        string `yaml:"version" valid:"required,lowercase"`
		Service        string `yaml:"service" valid:"required"`
		ServicePackage string `yaml:"servicePackage"`
		KratosLayout   string `yaml:"kratosLayout"`
		TemplateRepo   string `yaml:"templateRepo" valid:"required"`

		Project     string `yaml:"project"`
		Email       string `yaml:"email" valid:"email"`
		Description string `yaml:"description"`

		ConfigVolume string
		Port         *Port `yaml:"port"`
		ServiceLower string
		Repo         string
		GitRepo      string
		ProtoRepo    string
		ProjectGroup string
		GoPathSrc    string
		BackDir      string
		NameVersion  string
	}
	Port struct {
		Grpc int `yaml:"grpc" valid:"int"`
		Http int `yaml:"http" valid:"int"`
	}
)

var cas = cases.Title(language.English)

func (s *Settings) LoadAndCheck() error {
	if !utils.IsFileExists(FileName) {
		return fmt.Errorf("go settings file %s is not exists in current directory", FileName)
	}
	file, err := os.Open(FileName)
	if err != nil {
		return err
	}
	if err := yaml.NewDecoder(bufio.NewReader(file)).Decode(s); err != nil {
		return err
	}
	if _, err := govalidator.ValidateStruct(s); err != nil {
		return err
	}
	//if s.ConfigVolume != "" && !strings.HasPrefix(s.ConfigVolume, "/") {
	//	return errors.New(".settings.configVolume has not prefix /")
	//}
	firstLetter := s.Service[:1]
	if firstLetter != cas.String(firstLetter) {
		return errors.New(".settings.Service has to be in title case")
	}
	s.ServiceLower = fmt.Sprintf("%s%s", strings.ToLower(s.Service[:1]), s.Service[1:])
	if s.ServicePackage == "" {
		s.ServicePackage = strings.ToLower(s.Service)
	}
	s.ConfigVolume = "/data/conf"
	s.Port = &Port{
		Grpc: 9000,
		Http: 8000,
	}

	//if s.ConfigVolume == "" {
	//	s.ConfigVolume = "/data/conf"
	//}
	//if s.Port != nil {
	//	if s.Port.Grpc == 0 {
	//		s.Port.Grpc = 9000
	//	}
	//	if s.Port.Http == 0 {
	//		s.Port.Http = 8000
	//	}
	//} else {
	//	s.Port = &Port{
	//		Grpc: 9000,
	//		Http: 8000,
	//	}
	//}
	return nil
}

func (s *Settings) CheckEnv(gen bool) error {
	nameVersion := fmt.Sprintf("%s-%s", s.Name, s.Version)
	s.NameVersion = nameVersion
	dirExists := utils.IsDirExists(nameVersion)
	if gen && dirExists {
		return fmt.Errorf("directory %s already exists", nameVersion)
	}
	if !gen && !dirExists {
		return fmt.Errorf("directory %s not exists", nameVersion)
	}
	var err error
	if s.BackDir, err = os.Getwd(); err != nil {
		return err
	}
	goPathSrc := filepath.Join(os.Getenv("GOPATH"), utils.Src)
	s.GoPathSrc = goPathSrc
	if !strings.HasPrefix(s.BackDir, goPathSrc) {
		return errors.New("you are not in GOPATH")
	}
	projectGroupBack := s.BackDir[len(goPathSrc)+1:]
	i := strings.LastIndex(projectGroupBack, "/")
	if projectGroupBack[i:] != "/"+utils.Back {
		return errors.New("you are not in /back directory")
	}
	projectGroup := projectGroupBack[:i]
	s.ProtoRepo = filepath.Join(projectGroup, utils.Proto)
	if utils.IsDirExists(s.ProtoRepo) {
		return errors.New("no /proto dir in projectGroup")
	}
	s.ProjectGroup = projectGroup
	s.Repo = filepath.Join(projectGroup, utils.Back, nameVersion)
	s.GitRepo = fmt.Sprintf("git@%s.git", strings.Replace(s.Repo, "/", ":", 1))
	srcProtoRepo := filepath.Join(s.GoPathSrc, s.ProjectGroup, utils.Proto)
	if !utils.IsDirExists(srcProtoRepo) {
		return fmt.Errorf("%s not exists", srcProtoRepo)
	}
	srcProtoRepoService := filepath.Join(srcProtoRepo, utils.Api, s.ServicePackage)
	if !utils.IsDirExists(srcProtoRepoService) {
		return fmt.Errorf("%s not exists", srcProtoRepoService)
	}
	srcProtoRepoServiceProto := filepath.Join(srcProtoRepoService, fmt.Sprintf("%s.proto", s.ServicePackage))
	fileBytes, err := os.ReadFile(srcProtoRepoServiceProto)
	if err != nil {
		return err
	}
	if !bytes.Contains(fileBytes, []byte(fmt.Sprintf("service %s {", s.Service))) {
		return fmt.Errorf("%s not contains service %s", srcProtoRepoServiceProto, s.Service)
	}
	if strings.ToLower(s.ServicePackage) != s.ServicePackage {
		return fmt.Errorf("servicePackage not in lower case")
	}
	packageService := fmt.Sprintf("package api.%s;", s.ServicePackage)
	if !bytes.Contains(fileBytes, []byte(packageService)) {
		return fmt.Errorf("%s not contains package %s", srcProtoRepoServiceProto, packageService)
	}
	goOption := fmt.Sprintf(`option go_package = "%s/proto/gen/go/api/%s;%s";`,
		s.ProjectGroup, s.ServicePackage, s.ServicePackage)
	if !bytes.Contains(fileBytes, []byte(goOption)) {
		return fmt.Errorf("%s not contains go_option %s", srcProtoRepoServiceProto, goOption)
	}
	return nil
}
