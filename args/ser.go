package args

import (
	"bytes"
	"os"
	"path/filepath"
)

const (
	internalService = "internal/service"
)

func Services() error {
	files, err := os.ReadDir(internalService)
	if err != nil {
		return err
	}
	for _, file := range files {
		if err := replace(file.Name()); err != nil {
			return err
		}
	}
	return nil
}

func replace(filename string) error {
	filename = filepath.Join(internalService, filename)
	dat, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	if bytes.Contains(dat, []byte("}\nfunc")) {
		dat = bytes.ReplaceAll(dat, []byte("}\nfunc"), []byte("}\n\nfunc"))
		if err := os.WriteFile(filename, dat, 0644); err != nil {
			return err
		}
	}
	return nil
}
