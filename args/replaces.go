package args

import (
	"bufio"
	"fmt"
	"github.com/timurkash/gek/utils"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/timurkash/gek/utils/commands"
	"github.com/timurkash/gek/utils/settings"
)

func downloadAndReplaceFromTemplateRepo(settings *settings.Settings) error {
	templateRepo := settings.TemplateRepo
	if templateRepo == "" {
		return nil
	}
	templateRepo = strings.ReplaceAll(templateRepo, "https://", "")
	templateRepo = strings.ReplaceAll(templateRepo, ".git", "")
	repoDir := path.Join(settings.GoPathSrc, templateRepo)
	if err := cloneOrPull(settings.GoPathSrc, templateRepo); err != nil {
		return err
	}
	projectDir := path.Join(settings.GoPathSrc, settings.Repo)
	if err := os.Chdir(projectDir); err != nil {
		return err
	}
	if err := filepath.Walk(repoDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.Contains(path, "/.git/") ||
			info.Name() == ".git" && info.IsDir() ||
			strings.Contains(path, "/.idea/") ||
			info.Name() == ".idea" && info.IsDir() {
			return nil
		}
		fmt.Println(path)
		if err := modFile(path, info, settings); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func modFile(pathString string, info os.FileInfo, settings *settings.Settings) error {
	filenameInRepoDir := strings.ReplaceAll(pathString,
		path.Join(settings.GoPathSrc, settings.TemplateRepo),
		"")
	if info.IsDir() {
		if filenameInRepoDir == "" {
			return nil
		}
		if _, err := commands.Exec(
			"mkdir",
			"-p",
			path.Join(settings.Pwd, settings.NameVersion, filenameInRepoDir),
		); err != nil {
			return err
		}
		return nil
	}
	bytes, err := os.ReadFile(pathString)
	if err != nil {
		return err
	}
	temp := template.New(info.Name()).
		Funcs(map[string]interface{}{
			"lower": strings.ToLower,
		})
	if _, err := temp.Parse(string(bytes)); err != nil {
		return err
	}
	if filenameInRepoDir == "/gitlab-ci.yml" {
		filenameInRepoDir = "/.gitlab-ci.yml"
	}
	modFile, err := os.Create(path.Join(settings.GoPathSrc, settings.Repo, filenameInRepoDir))
	if err != nil {
		return err
	}
	defer func() {
		if err := modFile.Close(); err != nil {
			log.Fatalln(err)
		}
	}()
	writer := bufio.NewWriter(modFile)
	defer func() {
		if err := writer.Flush(); err != nil {
			log.Fatalln(err)
		}
	}()
	if err := temp.Execute(writer, *settings); err != nil {
		return err
	}
	return nil
}

func cloneOrPull(goPathSrc, repo string) error {
	repoDir := path.Join(goPathSrc, repo)
	if utils.IsDirExists(repoDir) {
		if err := os.Chdir(repoDir); err != nil {
			return err
		}
		if err := commands.ExecOnline("git", "pull"); err != nil {
			return err
		}
	} else {
		if err := commands.ExecOnline("git", "clone", fmt.Sprintf("https://%s.git", repo), repoDir); err != nil {
			return err
		}
	}
	return nil
}
