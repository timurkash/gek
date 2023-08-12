package args

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
)

var (
	bytesJson = []byte(",json=")
	eof       = []byte("\n")
	comma     = []byte(",")
)

func MessagesServer() error {
	err := filepath.Walk("gen/go", func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(path, ".pb.go") {
			if err := changeMessages(path, info.Mode().Perm()); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func changeMessages(path string, perm os.FileMode) error {
	bytesIn, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if !bytes.Contains(bytesIn, bytesJson) {
		return nil
	}
	lines := bytes.Split(bytesIn, eof)
	for i, line := range lines {
		if bytes.Contains(line, bytesJson) {
			lines[i] = changeLine(line)
		}
	}
	return os.WriteFile(path, bytes.Join(lines, eof), perm)
}

func changeLine(line []byte) []byte {
	lineItems := bytes.Split(line, comma)
	var json []byte
	for i, item := range lineItems {
		if bytes.HasPrefix(item, []byte("json=")) {
			json = item[5:]
		}
		if p := bytes.Index(item, []byte(`json:"`)); p > 0 {
			item = append(item[:p+6], json...)
			lineItems[i] = item
		}
	}
	return bytes.Join(lineItems, comma)
}
