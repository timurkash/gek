package args

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/timurkash/back/files"
	"gitlab.com/mcsolutions/tools/gek/utils/commands"
	settingsPackage "gitlab.com/mcsolutions/tools/gek/utils/settings"
)

func downloadAndReplaceFromTemplateRepo(settings *settingsPackage.Settings) error {
	templateRepo := settings.TemplateRepo
	if templateRepo == "" {
		return nil
	}
	templateRepo = strings.ReplaceAll(templateRepo, "https://", "")
	templateRepo = strings.ReplaceAll(templateRepo, ".git", "")
	repoDir := fmt.Sprintf("%s%s", settings.GoPathSrc, templateRepo)
	if err := cloneOrPull(settings.GoPathSrc, templateRepo); err != nil {
		return err
	}
	projectDir := fmt.Sprintf("%s%s", settings.GoPathSrc, settings.Repo)
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

func modFile(path string, info os.FileInfo, settings *settingsPackage.Settings) error {
	filenameInRepoDir := strings.ReplaceAll(path,
		fmt.Sprintf("%s%s", settings.GoPathSrc, settings.TemplateRepo),
		"")
	if info.IsDir() {
		if filenameInRepoDir == "" {
			return nil
		}
		if _, err := commands.Exec(
			"mkdir",
			"-p",
			fmt.Sprintf("%s/%s%s", settings.Pwd, settings.NameVersion, filenameInRepoDir),
		); err != nil {
			return err
		}
		return nil
	}
	bytes, err := os.ReadFile(path)
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
	modFile, err := os.Create(fmt.Sprintf("%s%s%s", settings.GoPathSrc, settings.Repo, filenameInRepoDir))
	if err != nil {
		return err
	}
	defer modFile.Close()
	writer := bufio.NewWriter(modFile)
	defer writer.Flush()
	if err := temp.Execute(writer, *settings); err != nil {
		return err
	}
	return nil
}

func cloneOrPull(goPathSrc, repo string) error {
	repoDir := fmt.Sprintf("%s%s", goPathSrc, repo)
	if files.IsDirExists(repoDir) {
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
