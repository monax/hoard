package main

import (
	"fmt"

	"github.com/monax/hoard/version"
)

// For use from tooling should do nothing but output version
func main() {
	fmt.Println(version.String())
}
