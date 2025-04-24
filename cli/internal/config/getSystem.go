package config

import (
	"errors"
	"log/slog"
)

func GetSystem(system string, cfg *Config) (targetSystem *System, err error) {
	for _, s := range cfg.Systems {
		if s.ID == system {
			targetSystem = &s
		}
	}
	if targetSystem == nil {
		slog.Debug("config.GetSystem target system not found", "system", system)
		err = errors.New("system not found: " + system)
	}

	return
}
