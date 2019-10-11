package main

import (
	"fmt"

	"github.com/monax/hoard/v6/project"
)

func main() {
	fmt.Println(project.History.CurrentVersion().String())
}
