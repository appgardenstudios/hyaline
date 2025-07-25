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

var licenseNote = "Note: This software requires a license for any purpose other than evaluation or demonstration. By using this software you attest that you are either 1) using Hyaline for evaluation or demonstration purposes or 2) you have a current and valid license to use Hyaline."

func main() {
	// Configure logger
	var logLevel = new(slog.LevelVar)
	h := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: logLevel})
	slog.SetDefault(slog.New(h))

	// Output license information
	fmt.Fprintf(os.Stderr, "%s\n\n", licenseNote)

	// Configure app
	app := &cli.App{
		Name:  "hyaline",
		Usage: usage,
		Action: func(*cli.Context) error {
			fmt.Printf("%s\n", usage)
			return nil
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "debug",
				Usage: "Include debug output",
			},
		},
		Commands: []*cli.Command{
			hyaline.Version(Version),
			hyaline.License(),
			hyaline.Check(logLevel),
			hyaline.Extract(logLevel),
			hyaline.Generate(logLevel),
			hyaline.Merge(logLevel),
			hyaline.Update(logLevel),
			hyaline.MCP(logLevel, Version),
			hyaline.Serve(logLevel, Version),
		},
	}

	// Run the app
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
