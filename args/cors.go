package args

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	r = []byte("r.")
)

func Cors() error {
	rest, err := os.Getwd()
	if err != nil {
		return err
	}
	rest = strings.ReplaceAll(rest, fmt.Sprintf("%s/src/", fmt.Sprintf(os.Getenv("GOPATH"))), "")
	rest = strings.ReplaceAll(rest, "/proto", "")
	if err := filepath.Walk("gen/go", func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(path, "_http.pb.go") {
			if err := changeHttpServer(path, info.Mode().Perm(), rest); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func changeHttpServer(path string, perm os.FileMode, rest string) error {
	bytesIn, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if !bytes.Contains(bytesIn, r) {
		return nil
	}
	lines := bytes.Split(bytesIn, eof)
	for i, line := range lines {
		switch {
		case bytes.Contains(line, r):
			lines[i] = changeRLine(line)
		case bytes.Contains(line, []byte(`"context"`)):
			lines[i] = []byte(fmt.Sprintf("\tcontext \"context\"\n\t\"%s/back/common/rest\"", rest))
		}
		if bytes.Contains(line, r) {
			lines[i] = changeRLine(line)
		}
	}
	return os.WriteFile(path, bytes.Join(lines, eof), perm)
}

var methods = []string{
	"GET",
	"HEAD",
	"POST",
	"PUT",
	"PATCH",
	"DELETE",
	"CONNECT",
	"OPTIONS",
	"TRACE",
}

func changeRLine(line []byte) []byte {
	line = bytes.ReplaceAll(line, r, []byte("rest.HandleRoute(r, "))
	for _, method := range methods {
		line = bytes.ReplaceAll(line,
			[]byte(fmt.Sprintf("%s(", method)),
			[]byte(fmt.Sprintf(`"%s", `, method)))
	}
	return line
}

//func rewriteFile(temp *template.Template, filename string, service *Service) error {
//	httpFile, err := os.Create(filename)
//	if err != nil {
//		return err
//	}
//	defer func() {
//		if err := httpFile.Close(); err != nil {
//			log.Fatalln(err)
//		}
//	}()
//	writer := bufio.NewWriter(httpFile)
//	defer func() {
//		if err := writer.Flush(); err != nil {
//			log.Fatalln(err)
//		}
//	}()
//	return temp.Execute(writer, *service)
//}
