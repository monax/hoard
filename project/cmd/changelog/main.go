package main

import (
	"fmt"

	"github.com/monax/hoard/v4/project"
)

func main() {
	fmt.Println(project.History.MustChangelog())
}
