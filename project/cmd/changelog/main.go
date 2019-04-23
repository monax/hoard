package main

import (
	"fmt"

	"github.com/monax/hoard/v3/project"
)

func main() {
	fmt.Println(project.History.MustChangelog())
}
