package args

import (
	"bufio"
	"fmt"
	"gitlab.com/mcsolutions/tools/gek/utils/settings"
	"log"
	"os"
	"strings"
)

func processProto(settings *settings.Settings) error {
	srcProtoRepo := settings.GoPathSrc + settings.ProjectGroup + "/proto/"
	file, err := os.Open(srcProtoRepo + "internal/service/" + settings.ServicePackage + ".go")
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Println(err)
		}
	}()
	scanner := bufio.NewScanner(file)
	fileWrite, err := os.Create(settings.GoPathSrc + settings.Repo + "/internal/service/greeter.go")
	if err != nil {
		return err
	}
	writer := bufio.NewWriter(fileWrite)
	writeLine := true
	for scanner.Scan() {
		line := scanner.Text()
		writeLine = true
		if line == ")" {
			if _, err := writer.WriteString(fmt.Sprintf(`
	"github.com/go-kratos/kratos/v2/log"
	"%s/internal/biz"
`, settings.Repo)); err != nil {
				return err
			}
		} else if strings.HasPrefix(line, "\tpb.Unimplemented") {
			if _, err := writer.WriteString(fmt.Sprintf(`	uc  *biz.%sUsecase
	log *log.Helper

`, settings.Service)); err != nil {
				return err
			}
		} else if strings.HasPrefix(line, "func New"+settings.Service+"Service()") {
			if _, err := writer.WriteString(fmt.Sprintf(`func New%sService(uc *biz.%sUsecase, logger log.Logger) *%sService {
`, settings.Service, settings.Service, settings.Service)); err != nil {
				return err
			}
			writeLine = false
		} else if line == "\treturn &"+settings.Service+"Service{}" {
			if _, err := writer.WriteString(fmt.Sprintf(`	return &%sService{
		uc:  uc,
		log: log.NewHelper(logger),
	}
`, settings.Service)); err != nil {
				return err
			}
			writeLine = false
		}
		if writeLine {
			if _, err := writer.WriteString(fmt.Sprintf(`%s
`, line)); err != nil {
				return err
			}
		}
	}
	if err := writer.Flush(); err != nil {
		return err
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}
