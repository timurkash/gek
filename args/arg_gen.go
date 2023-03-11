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
		kratosNew,
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
		} else {
			if err := commands.ExecOnline(command, words[1:]...); err != nil {
				return err
			}
		}
	}
	return nil
}

func printDetails(nameVersion string) {
	fmt.Println("\nMake")
	fmt.Println("\tcd", nameVersion)
	fmt.Println("\tmake init-all")
	fmt.Println("To use ent.ORM you can")
	fmt.Println("\tmake init-ent")
	fmt.Println("To build")
	fmt.Println("\tmake build")
	fmt.Println("To run")
	fmt.Println("\tmake run")
	fmt.Println("To init push")
	fmt.Println("\tmake init-push")
	fmt.Println("To upgrade kratos")
	fmt.Println("\tkratos upgrade")
}
