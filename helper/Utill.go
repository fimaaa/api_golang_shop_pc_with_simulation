package helper

import (
	"fmt"

	"other/simulasi_pc/conf"
)

func PrintCommand(command string) {
	if conf.Configuration().Log.Verbose {
		fmt.Println(command)
	}
}
