package cmd

import (
	"fmt"

	"github.com/jawher/mow.cli"
	"github.com/monax/hoard/release"
)

func AddVersionCommand(cmd *cli.Cli) {
	cmd.Command("version", "Get version number",
		func(versionCmd *cli.Cmd) {
			versionCmd.Action = func() {
				fmt.Println(release.Version())
			}

			versionCmd.Command("changes", "Get changes in this version",
				func(changesCmd *cli.Cmd) {
					changesCmd.Action = func() {
						fmt.Println(release.Changes())
					}
				})

			versionCmd.Command("changelog", "Get Hoard changelog",
				func(changesCmd *cli.Cmd) {
					changesCmd.Action = func() {
						fmt.Println(release.Changelog())
					}
				})

		})
}
