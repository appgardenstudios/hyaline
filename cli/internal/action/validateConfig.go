package action

import (
	"hyaline/internal/config"
	"hyaline/internal/io"
	"log/slog"
)

type ValidateConfigArgs struct {
	Config string
	Output string
}

type ValidateConfigOutput struct {
	Valid  bool                       `json:"valid"`
	Error  string                     `json:"error"`
	Detail ValidateConfigOutputDetail `json:"detail"`
}

type ValidateConfigOutputDetail struct {
	LLM     ValidateConfigOutputGlobal  `json:"llm"`
	GitHub  ValidateConfigOutputGlobal  `json:"github"`
	Extract ValidateConfigOutputCommand `json:"extract"`
	Check   ValidateConfigOutputCommand `json:"check"`
	Audit   ValidateConfigOutputCommand `json:"audit"`
}

type ValidateConfigOutputGlobal struct {
	Present bool   `json:"present"`
	Valid   bool   `json:"valid"`
	Error   string `json:"error"`
}

type ValidateConfigOutputCommand struct {
	Present  bool   `json:"present"`
	Disabled bool   `json:"disabled"`
	Valid    bool   `json:"valid"`
	Error    string `json:"error"`
}

func ValidateConfig(args *ValidateConfigArgs) error {
	slog.Info("Validating config", "config", args.Config, "output", args.Output)

	// Load Config
	cfg, err := config.Load(args.Config, false)
	if err != nil {
		slog.Debug("action.ValidateConfig could not load the config", "error", err)
		return err
	}

	// Ensure output JSON file does not exist
	outputFile, err := io.InitOutput(args.Output)
	if err != nil {
		slog.Debug("action.ValidateConfig could not initialize output file", "error", err)
		return err
	}
	defer outputFile.Close()

	// Initialize our result
	result := ValidateConfigOutput{}

	// Check overall config
	err = config.Validate(cfg)
	if err != nil {
		result.Error = err.Error()
	} else {
		result.Valid = true
	}

	// Check LLM
	if cfg.LLM.Provider != "" {
		result.Detail.LLM.Present = true
		err = config.ValidateLLM(cfg)
		if err != nil {
			result.Detail.LLM.Error = err.Error()
		} else {
			result.Detail.LLM.Valid = true
		}
	}

	// Check GitHub
	if cfg.GitHub.Token != "" {
		result.Detail.GitHub.Present = true
		result.Detail.GitHub.Valid = true
	}

	// Check extract
	if cfg.Extract != nil {
		result.Detail.Extract.Present = true
		err = config.ValidateExtract(cfg)
		if err != nil {
			result.Detail.Extract.Error = err.Error()
		} else {
			result.Detail.Extract.Valid = true
		}
		result.Detail.Extract.Disabled = cfg.Extract.Disabled
	}

	// Check check
	if cfg.Check != nil {
		result.Detail.Check.Present = true
		err = config.ValidateCheck(cfg)
		if err != nil {
			result.Detail.Check.Error = err.Error()
		} else {
			result.Detail.Check.Valid = true
		}
		result.Detail.Check.Disabled = cfg.Check.Disabled
	}

	// Check audit
	if cfg.Audit != nil {
		result.Detail.Audit.Present = true
		err = config.ValidateAudit(cfg)
		if err != nil {
			result.Detail.Audit.Error = err.Error()
		} else {
			result.Detail.Audit.Valid = true
		}
		result.Detail.Audit.Disabled = cfg.Audit.Disabled
	}

	// Output result
	err = io.WriteJSON(outputFile, result)
	if err != nil {
		slog.Debug("action.ValidateConfig could not write output file", "error", err)
		return err
	}

	slog.Info("Validate config completed successfully")

	return nil
}
