package args

import (
	"fmt"
	"github.com/timurkash/gek/utils/commands"
	settingsPackage "github.com/timurkash/gek/utils/settings"
	"os"
	"strings"
)

func Generate() error {
	settings := new(settingsPackage.Settings)
	if err := settings.LoadAndCheck(); err != nil {
		return err
	}
	if err := settings.CheckEnv(true); err != nil {
		return err
	}
	kratosNew := fmt.Sprintf("kratos new %s", settings.Repo)
	if settings.KratosLayout != "" {
		kratosNew = fmt.Sprintf("%s -r %s", kratosNew, settings.KratosLayout)
	}
	scriptCommands := []string{
		fmt.Sprintf("cd %s", settings.GoPathSrc),
		kratosNew,
		fmt.Sprintf("cd %s", settings.BackDir),
		"kratos -v",
		fmt.Sprintf("cp .settings %s", settings.NameVersion),
		fmt.Sprintf("cd %s", settings.NameVersion),
		"git init",
		"git add .",
		"git commit -m dummy -q",
	}
	if err := executeScriptCommands(scriptCommands); err != nil {
		return err
	}
	if err := downloadAndReplaceFromTemplateRepo(settings); err != nil {
		return err
	}
	if err := processProto(settings); err != nil {
		return err
	}
	printDetails(settings.NameVersion)
	return nil
}

func executeScriptCommands(scriptCommands []string) error {
	for _, script := range scriptCommands {
		words := strings.Split(script, " ")
		command := words[0]
		if command == "cd" {
			if err := os.Chdir(words[1]); err != nil {
				return err
			}
			wd, err := os.Getwd()
			if err != nil {
				return err
			}
			fmt.Println("current dir is", wd)
		} else {
			if err := commands.ExecOnline(command, words[1:]...); err != nil {
				return err
			}
		}
	}
	return nil
}

const tsn = "\t%s\n"

func printDetails(nameVersion string) {
	fmt.Printf("\nMake\n")
	fmt.Printf(tsn, blue(fmt.Sprintf("cd %s", nameVersion)))
	fmt.Printf(tsn, blue("make init-all"))
	fmt.Printf("To use ent.ORM you can\n")
	fmt.Printf(tsn, blue("make init-ent"))
	fmt.Printf("To build\n")
	fmt.Printf(tsn, blue("make build"))
	fmt.Printf("To run\n")
	fmt.Printf(tsn, blue("make run"))
	fmt.Printf("To init push\n")
	fmt.Printf(tsn, blue("make init-push"))
	fmt.Printf("To upgrade kratos\n")
	fmt.Printf(tsn, blue("kratos upgrade"))
}
