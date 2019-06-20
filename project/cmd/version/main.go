package main

import (
	"fmt"

	"github.com/monax/hoard/v5/project"
)

func main() {
	fmt.Println(project.History.CurrentVersion().String())
}
