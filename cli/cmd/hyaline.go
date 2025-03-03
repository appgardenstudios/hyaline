package main

import (
	"fmt"
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
				Usage: "print out the current version",
				Action: func(cCtx *cli.Context) error {
					fmt.Println(Version)
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
