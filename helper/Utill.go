package helper

import (
	"fmt"

	"other/simulasi_pc/conf"
)

func PrintCommand(command ...any) {
	if conf.Configuration().Log.Verbose {
		fmt.Println(command...)
	}
}
