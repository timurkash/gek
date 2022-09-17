package commands

import (
	"fmt"
	"io"
	"os/exec"
	"strings"
)

func ExecOnline(name string, args ...string) error {
	fmt.Println(name, strings.Join(args, " "))
	cmd := exec.Command(name, args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	oneByte := make([]byte, 1)
	for {
		_, err := stdout.Read(oneByte)
		if err != nil {
			if err == io.EOF {
				return nil
			} else {
				return err
			}
		}
		fmt.Print(string(oneByte))
	}
}

func Exec(name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	stdout, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return stdout, nil
}
