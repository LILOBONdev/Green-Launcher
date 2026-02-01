package utils

import (
	"fmt"
	//"github.com/inancgumus/screen"
)

func ClearScreen() {
	fmt.Print("\033[H\033[2J")
}
