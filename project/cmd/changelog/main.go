package main

import (
	"fmt"

	"github.com/monax/hoard/v8/project"
)

func main() {
	fmt.Println(project.History.MustChangelog())
}
