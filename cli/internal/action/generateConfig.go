package action

import (
	"hyaline/internal/config"
	"log/slog"
)

type GenerateConfigArgs struct {
	Config         string
	Current        string
	System         string
	Output         string
	IncludePurpose bool
}

func GenerateConfig(args *GenerateConfigArgs) error {
	slog.Info("Generating Config")
	slog.Debug("action.GenerateConfig Args", slog.Group("args",
		"config", args.Config,
		"current", args.Current,
		"system", args.System,
		"output", args.Output,
		"include-purpose", args.IncludePurpose,
	))

	// Load Config
	_, err := config.Load(args.Config)
	if err != nil {
		slog.Debug("action.GenerateConfig could not load the config", "error", err)
		return err
	}

	return nil
}
