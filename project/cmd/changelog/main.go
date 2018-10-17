package main

import (
	"fmt"

	"github.com/monax/hoard/project"
)

func main() {
	fmt.Println(project.History.MustChangelog())
}
