package hyaline

import (
	"fmt"
	"time"

	"github.com/urfave/cli/v2"
)

func License() *cli.Command {
	return &cli.Command{
		Name:  "license",
		Usage: "Print out license information",
		Action: func(cCtx *cli.Context) error {
			fmt.Printf("Copyright %d App Garden Studios, LLC\n\n", time.Now().Year())
			fmt.Print("All Rights Reserved\n\n")
			fmt.Print("THE CONTENTS OF THIS PROJECT ARE PROPRIETARY AND CONFIDENTIAL. UNAUTHORIZED COPYING, TRANSFER, OR REPRODUCTION OF THIS PROJECT VIA ANY MEDIUM IS STRICTLY PROHIBITED.\n")
			return nil
		},
	}
}
