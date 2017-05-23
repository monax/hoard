package main

import "github.com/jawher/mow.cli"

func main() {
	hoard := cli.App("hoard",
		"A content-addressed deterministically encrypted blob storage system")

	hoard.Command("encrypt", "", func(cmd *cli.Cmd) {

	})

	hoard.Command("decrypt", "", func(cmd *cli.Cmd) {

	})

	hoard.Command("key",
		"Outputs the deterministically generated content key for the " +
				"plaintext provided on stdin", func(cmd *cli.Cmd) {

	})
}