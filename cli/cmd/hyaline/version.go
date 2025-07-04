package hyaline

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func Version(version string) *cli.Command {
	return &cli.Command{
		Name:  "version",
		Usage: "Print out the current version",
		Action: func(cCtx *cli.Context) error {
			fmt.Printf("%s\n", version)
			return nil
		},
	}
}
