package args

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
)

var (
	bytesGoogleApiAnnotations  = []byte(`import * as dependency_1 from "./../../google/api/annotations";`)
	bytesGoogleApiAnnotations_ = append([]byte(`//`), bytesGoogleApiAnnotations...)
	bytesInterface             = []byte(`    interface `)
	bytesTsIgnore              = []byte(`    // @ts-ignore
`)
)

func TsIgnore() error {
	if err := filepath.Walk("gen/ts/api", func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && !strings.HasSuffix(path, "/messages.ts") {
			if err := changeTs(path, info.Mode()); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func changeTs(path string, perm os.FileMode) error {
	bytesIn, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	lines := bytes.Split(bytesIn, eof)
	for i, line := range lines {
		switch {
		case bytes.Equal(line, bytesGoogleApiAnnotations):
			lines[i] = bytesGoogleApiAnnotations_
		case bytes.HasPrefix(line, bytesInterface):
			lines[i] = append(bytesTsIgnore, line...)
		}
	}
	return os.WriteFile(path, bytes.Join(lines, eof), perm)
}
