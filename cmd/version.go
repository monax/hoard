package cmd

import (
	"fmt"

	cli "github.com/jawher/mow.cli"
	"github.com/monax/hoard/project"
)

func AddVersionCommand(cmd *cli.Cli) {
	cmd.Command("version", "Get version number",
		func(versionCmd *cli.Cmd) {
			versionCmd.Action = func() {
				fmt.Println(project.FullVersion())
			}

			versionCmd.Command("notes", "Get release notes for this version",
				func(changesCmd *cli.Cmd) {
					changesCmd.Action = func() {
						fmt.Println(project.History.CurrentNotes())
					}
				})

			versionCmd.Command("changelog", "Get Hoard changelog",
				func(changesCmd *cli.Cmd) {
					changesCmd.Action = func() {
						fmt.Println(project.History.MustChangelog())
					}
				})

		})
}
