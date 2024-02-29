package args

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func JwtServer() error {
	rest, err := os.Getwd()
	if err != nil {
		return err
	}
	rest = strings.ReplaceAll(rest, fmt.Sprintf("%s/src/", fmt.Sprintf(os.Getenv("GOPATH"))), "")
	rest = strings.ReplaceAll(rest, "/proto", "")
	if err := filepath.Walk("gen/go", func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(path, "_http.pb.go") {
			if err := changeJwtHttpServer(path, info.Mode().Perm(), rest); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

var hCtx = []byte("h(ctx, &in)")

func changeJwtHttpServer(path string, perm os.FileMode, rest string) error {
	bytesIn, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if !bytes.Contains(bytesIn, context) {
		return nil
	}
	backCommonRest := []byte(fmt.Sprintf("\"%s/back/common/rest\"", rest))
	hasRest := bytes.Contains(bytesIn, backCommonRest)
	lines := bytes.Split(bytesIn, eof)
	for i, line := range lines {
		switch {
		case bytes.Contains(line, context) && !hasRest:
			lines[i] = append([]byte("\tcontext \"context\"\n\t"), backCommonRest...)
		case bytes.Contains(line, hCtx):
			lines[i] = []byte(`		bearerToken, err := rest.GetBearerToken(ctx)
		if err != nil {
			return err
		}
		out, err := h(rest.AppendTokenToContext(rest.AppendHeadersValues(ctx, ctx.Header()), bearerToken), &in)`)
		}
	}
	return os.WriteFile(path, bytes.Join(lines, eof), perm)
}
