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
	srcProtoRepo := fmt.Sprintf("%s%s/proto/", settings.GoPathSrc, settings.ProjectGroup)
	file, err := os.Open(fmt.Sprintf("%sinternal/service/%s.go", srcProtoRepo, settings.ServicePackage))
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Println(err)
		}
	}()
	scanner := bufio.NewScanner(file)
	fileWrite, err := os.Create(fmt.Sprintf("%s%s/internal/service/greeter.go", settings.GoPathSrc, settings.Repo))
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
		} else if strings.HasPrefix(line, fmt.Sprintf("func New%sService()", settings.Service)) {
			if _, err := writer.WriteString(fmt.Sprintf(`func New%sService(uc *biz.%sUsecase, logger log.Logger) *%sService {
`, settings.Service, settings.Service, settings.Service)); err != nil {
				return err
			}
			writeLine = false
		} else if line == fmt.Sprintf("\treturn &%sService{}", settings.Service) {
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
