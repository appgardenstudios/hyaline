package main

import (
	"fmt"
	"hyaline/cmd/hyaline"
	"log"
	"log/slog"
	"os"

	"github.com/urfave/cli/v2"
)

var Version = "unknown"

var usage = "Maintain Your Documentation - Find, Fix, and Prevent Documentation Issues."

var betaNote = "Note: Hyaline is currently in an open beta. As such, this software is only licensed for evaluation and use until the open beta period ends."

func main() {
	var logLevel = new(slog.LevelVar)
	h := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: logLevel})
	slog.SetDefault(slog.New(h))

	app := &cli.App{
		Name:  "hyaline",
		Usage: fmt.Sprintf("%s\n%s", usage, betaNote),
		Action: func(*cli.Context) error {
			fmt.Printf("%s\n%s\n", usage, betaNote)
			return nil
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "debug",
				Usage: "Include debug output",
			},
		},
		Commands: []*cli.Command{
			hyaline.Version(Version, betaNote),
			hyaline.Check(logLevel),
			hyaline.Extract(logLevel),
			hyaline.Generate(logLevel),
			hyaline.Merge(logLevel),
			hyaline.Update(logLevel),
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
