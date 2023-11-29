package args

import (
	"bytes"
	"os"
	"path/filepath"
)

var (
	bytesGoogleApiAnnotations = []byte(`/google/api/annotations";`)
	bytesComment              = []byte(`// `)
	bytesValidation           = []byte(`/validate/validate";`)
	bytesInterface            = []byte(`    interface `)
	bytesTsIgnore             = []byte(`    // @ts-ignore
`)
)

func TsIgnore() error {
	if err := filepath.Walk("gen/ts/api", func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
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
		case bytes.HasSuffix(line, bytesGoogleApiAnnotations):
			lines[i] = append(bytesComment, line...)
		case bytes.HasSuffix(line, bytesValidation):
			lines[i] = append(bytesComment, line...)
		case bytes.HasPrefix(line, bytesInterface):
			lines[i] = append(bytesTsIgnore, line...)
		}
	}
	return os.WriteFile(path, bytes.Join(lines, eof), perm)
}
