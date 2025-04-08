package main

import (
	"fmt"
	"hyaline/internal/action"
	"log"
	"log/slog"
	"os"

	"github.com/urfave/cli/v2"
)

var Version = "unknown"

func main() {
	var logLevel = new(slog.LevelVar)
	h := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: logLevel})
	slog.SetDefault(slog.New(h))

	app := &cli.App{
		Name:  "hyaline",
		Usage: "Maintain Your Documentation - Find, Fix, and Prevent Documentation Issues",
		Action: func(*cli.Context) error {
			fmt.Println("hello world")
			return nil
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "debug",
				Usage: "Include debug output",
			},
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
				Name:  "check",
				Usage: "Check documentation for issues and errors",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "config",
						Required: true,
						Usage:    "Path to the config file",
					},
					&cli.StringFlag{
						Name:     "current",
						Required: true,
						Usage:    "Path to the current data set",
					},
					&cli.StringFlag{
						Name:     "change",
						Required: false,
						Usage:    "Path to the change data set",
					},
					&cli.StringFlag{
						Name:     "system",
						Required: true,
						Usage:    "ID of the system to extract",
					},
					&cli.BoolFlag{
						Name:  "recommend",
						Usage: "Include a recommended action when a check does not pass",
					},
				},
				Action: func(cCtx *cli.Context) error {
					// Set log level
					if cCtx.Bool("debug") {
						logLevel.Set(slog.LevelDebug)
					}

					// Execute action
					err := action.Check(&action.CheckArgs{
						Config:    cCtx.String("config"),
						Current:   cCtx.String("current"),
						Change:    cCtx.String("change"),
						System:    cCtx.String("system"),
						Recommend: cCtx.Bool("recommend"),
					})
					if err != nil {
						return cli.Exit(err.Error(), 1)
					}
					return nil
				},
			},
			{
				Name:  "extract",
				Usage: "Extract code, documentation, and other metadata",
				Subcommands: []*cli.Command{
					{
						Name:  "change",
						Usage: "Extract and create a change data set",
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
								Name:     "base",
								Required: true,
								Usage:    "Base branch (where changes will be applied)",
							},
							&cli.StringFlag{
								Name:     "head",
								Required: true,
								Usage:    "Head branch (which changes will be applied)",
							},
							&cli.StringFlag{
								Name:     "pull-request",
								Required: false,
								Usage:    "GitHub Pull Request to include in the change (OWNER/REPO/PR_NUMBER)",
							},
							&cli.StringSliceFlag{
								Name:     "issue",
								Required: false,
								Usage:    "GitHub Issue to include in the change (OWNER/REPO/PR_NUMBER). Accepts multiple issues by setting multiple times.",
							},
							&cli.StringFlag{
								Name:     "output",
								Required: true,
								Usage:    "Path of the sqlite database to create",
							},
						},
						Action: func(cCtx *cli.Context) error {
							// Set log level
							if cCtx.Bool("debug") {
								logLevel.Set(slog.LevelDebug)
							}

							// Execute action
							err := action.ExtractChange(&action.ExtractChangeArgs{
								Config:      cCtx.String("config"),
								System:      cCtx.String("system"),
								Base:        cCtx.String("base"),
								Head:        cCtx.String("head"),
								PullRequest: cCtx.String("pull-request"),
								Issues:      cCtx.StringSlice("issue"),
								Output:      cCtx.String("output"),
							})
							if err != nil {
								return cli.Exit(err.Error(), 1)
							}
							return nil
						},
					},
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
							// Set log level
							if cCtx.Bool("debug") {
								logLevel.Set(slog.LevelDebug)
							}

							// Execute action
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
