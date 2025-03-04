package main

import (
	"fmt"
	"hyaline/internal/action"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

var Version = "unknown"

func main() {
	app := &cli.App{
		Name:  "hyaline",
		Usage: "Maintain Your Documentation - Find, Fix, and Prevent Documentation Issues",
		Action: func(*cli.Context) error {
			fmt.Println("hello world")
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:  "version",
				Usage: "Print out the current version",
				Action: func(cCtx *cli.Context) error {
					fmt.Println(Version)
					return nil
				},
			},
			{
				Name:  "extract",
				Usage: "Extract code, documentation, and other metadata",
				Subcommands: []*cli.Command{
					{
						Name:  "current",
						Usage: "Extract and create a current data set",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "config",
								Required: true,
								Usage:    "Path to the config file",
							},
							&cli.StringFlag{
								Name:     "system",
								Required: true,
								Usage:    "ID of the system to extract",
							},
							&cli.StringFlag{
								Name:     "output",
								Required: true,
								Usage:    "Path of the sqlite database to create",
							},
						},
						Action: func(cCtx *cli.Context) error {
							err := action.ExtractCurrent(&action.ExtractCurrentArgs{
								Config: cCtx.String("config"),
								System: cCtx.String("system"),
								Output: cCtx.String("output"),
							})
							if err != nil {
								return cli.Exit(err.Error(), 1)
							}
							return nil
						},
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
