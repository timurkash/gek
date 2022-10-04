package args

import (
	"bytes"
	"errors"
	"fmt"
	"gitlab.com/mcsolutions/tools/gek/utils/commands"
	"os"
	"strconv"
)

func ArgUp(level int) error {
	badTagError := errors.New("bad tag")
	version, err := commands.Exec("git", "describe", "--tags", "--abbrev=0")
	if err != nil {
		return err
	}
	version = version[:len(version)-1]
	if !bytes.HasPrefix(version, []byte("v")) {
		return errors.New("version has not prefix \"v\"")
	}
	parts := bytes.Split(version[1:], []byte("."))
	if len(parts) != 3 {
		return badTagError
	}
	major, err := getInt(parts[0])
	if err != nil {
		return err
	}
	minor, err := getInt(parts[1])
	if err != nil {
		return err
	}
	patch, err := getInt(parts[2])
	if err != nil {
		return err
	}
	switch level {
	case 0:
		major++
		minor = 0
		patch = 0
	case 1:
		minor++
		patch = 0
	case 2:
		patch++
	}
	versionNext := fmt.Sprintf("v%d.%d.%d", major, minor, patch)
	_, err = commands.Exec("git", "tag", "-a", versionNext, "-F", "tag_message")
	if err != nil {
		return err
	}
	valuesBytes, err := os.ReadFile("values.yaml")
	if err != nil {
		return err
	}
	valuesBytes = bytes.ReplaceAll(valuesBytes,
		[]byte(fmt.Sprintf("  tag: %s", string(version))),
		[]byte(fmt.Sprintf("  tag: %s", versionNext)))
	if err := os.WriteFile("values.yaml", valuesBytes, 0644); err != nil {
		return err
	}
	if err := os.WriteFile("tag", []byte(versionNext), 0644); err != nil {
		return err
	}
	return nil
}

func getInt(part []byte) (int, error) {
	return strconv.Atoi(string(part))
}
