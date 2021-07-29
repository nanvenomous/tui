package main

import (
	"fmt"

	"github.com/mrgarelli/tui"
)

func main() {
	exampleChoices := []string{"one", "two", "three"}

	fc := tui.LaunchSelection(exampleChoices)
	fmt.Println(fc)
}
