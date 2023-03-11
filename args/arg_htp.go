package args

import (
	"bufio"
	"fmt"
	"github.com/timurkash/gek/utils"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type Service struct {
	Name    string
	Service string
}

func HttpServer() error {
	dirRead, err := os.Open("api")
	if err != nil {
		return err
	}
	dirFiles, err := dirRead.ReadDir(0)
	if err != nil {
		return err
	}
	temp := template.New("http_server")
	if _, err := temp.Parse(httpServer); err != nil {
		return err
	}
	for _, file := range dirFiles {
		if file.IsDir() {
			name := file.Name()
			service, err := getService(filepath.Join(utils.Api, name, fmt.Sprintf("%s.proto", name)))
			if err != nil {
				return err
			}
			filename := filepath.Join("gen/go/api", name, fmt.Sprintf("%s_http.pb.go", name))
			if !utils.IsExists(filename) {
				if err := rewriteFile(temp, filename, &Service{
					Name:    name,
					Service: service,
				}); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func rewriteFile(temp *template.Template, filename string, service *Service) error {
	httpFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func() {
		if err := httpFile.Close(); err != nil {
			log.Fatalln(err)
		}
	}()
	writer := bufio.NewWriter(httpFile)
	defer func() {
		if err := writer.Flush(); err != nil {
			log.Fatalln(err)
		}
	}()
	return temp.Execute(writer, *service)
}

const (
	httpServer = `// Code generated by gek -htp. DO NOT EDIT.

package {{ .Name }}

import (
	http "github.com/go-kratos/kratos/v2/transport/http"
)

type {{ .Service }}HTTPServer interface {
}

func Register{{ .Service }}HTTPServer(s *http.Server, srv {{ .Service }}HTTPServer) {
}

type {{ .Service }}HTTPClient interface {
}

type {{ .Service }}HTTPClientImpl struct {
	cc *http.Client
}

func New{{ .Service }}HTTPClient(client *http.Client) {{ .Service }}HTTPClient {
	return &{{ .Service }}HTTPClientImpl{client}
}
`
)

func getService(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Fatalln(err)
		}
	}()
	errOut := fmt.Errorf("bad %s file", filename)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "service ") {
			lexemes := strings.Split(line, " ")
			if len(lexemes) < 3 {
				return "", errOut
			}
			return lexemes[1], nil
		}
	}
	return "", fmt.Errorf("bad %s file", filename)
}
