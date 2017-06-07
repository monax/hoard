package main

import (
	"fmt"

	"code.monax.io/platform/hoard/version"
)

// For use from tooling should do nothing but output version
func main() {
	fmt.Println(version.String())
}
