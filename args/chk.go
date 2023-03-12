package args

import (
	"fmt"
	"github.com/timurkash/gek/utils"
)

func Check() error {
	for _, util := range Utils {
		if err := utils.IsUtilExists(util.Name); err != nil {
			return fmt.Errorf(`%s not installed
to install: %s
total list of utils run: gek -utl`, util.Name, util.Command)
		}
	}
	fmt.Println("all required utils is installed")
	return nil
}
