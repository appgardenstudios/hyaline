package code

import (
	"database/sql"
	"hyaline/internal/config"
)

func ExtractChange(system *config.System, head string, base string, db *sql.DB) (err error) {
	// // Process each code source
	// for _, c := range system.Code {
	// 	slog.Debug("code.ExtractChange extracting code", "system", system, "code", c.ID)

	// 	// Insert Code
	// 	codeId := system.ID + "-" + c.ID
	// 	err = sqlite.InsertChangeCode(sqlite.ChangeCode{
	// 		ID:       codeId,
	// 		SystemID: system.ID,
	// 		Path:     c.Path,
	// 	}, db)

	// 	// Make sure we have a valid preset. If not, skip
	// 	preset, ok := presets[c.Preset]
	// 	if !ok {
	// 		slog.Info("Code Preset Not Found. Skipping...", "system", system, "code", c.ID, "preset", c.Preset)
	// 		continue
	// 	}
	// 	glob, err := zglob.New(preset.Glob)
	// 	if err != nil {
	// 		slog.Debug("code.ExtractChange could not instantiate preset glob", "error", err, "glob", preset.Glob)
	// 		return err
	// 	}
	// 	slog.Debug("code.ExtractChange extracting code using preset", "presetID", c.Preset, "preset", preset)

	// 	// Get our diffs
	// 	diffs, err := repo.GetDiff(c.Path, head, base)
	// 	if err != nil {
	// 		slog.Debug("code.ExtractChange could not get diff", "error", err)
	// 		return err
	// 	}

	// 	// Load any files in the diff that match our preset
	// 	for _, diff := range diffs {
	// 		slog.Debug("code.ExtractChange processing diff", "diff", diff.String())
	// 		action, err := diff.Action()
	// 		if err != nil {
	// 			slog.Debug("code.ExtractChange could not retrieve action for diff", "error", err, "diff", diff)
	// 			return err
	// 		}
	// 		_, to, err := diff.Files()
	// 		if err != nil {
	// 			slog.Debug("code.ExtractChange could not retrieve files for diff", "error", err, "diff", diff)
	// 			return err
	// 		}
	// 		switch action {
	// 		case merkletrie.Insert:
	// 			fallthrough
	// 		case merkletrie.Modify:
	// 			if glob.Match(diff.To.Name) || slices.Contains(preset.Files, diff.To.Name) {
	// 				slog.Debug("code.ExtractChange inserting file", "file", diff.To.Name, "action", action)
	// 				bytes, err := repo.GetBlobBytes(to.Blob)
	// 				if err != nil {
	// 					slog.Debug("code.ExtractChange could not retrieve blob from diff", "error", err)
	// 					return err
	// 				}
	// 				err = sqlite.InsertChangeFile(sqlite.ChangeFile{
	// 					ID:           diff.To.Name,
	// 					CodeID:       codeId,
	// 					SystemID:     system.ID,
	// 					RelativePath: diff.To.Name,
	// 					Action:       action.String(),
	// 					RawData:      string(bytes),
	// 				}, db)
	// 				if err != nil {
	// 					slog.Debug("code.ExtractChange could not insert file", "error", err)
	// 					return err
	// 				}
	// 			}
	// 		case merkletrie.Delete:
	// 			if glob.Match(diff.From.Name) || slices.Contains(preset.Files, diff.From.Name) {
	// 				slog.Debug("code.ExtractChange inserting file", "file", diff.From.Name, "action", action)
	// 				err = sqlite.InsertChangeFile(sqlite.ChangeFile{
	// 					ID:           diff.From.Name,
	// 					CodeID:       codeId,
	// 					SystemID:     system.ID,
	// 					RelativePath: diff.From.Name,
	// 					Action:       action.String(),
	// 					RawData:      "",
	// 				}, db)
	// 				if err != nil {
	// 					slog.Debug("code.ExtractChange could not insert file", "error", err)
	// 					return err
	// 				}
	// 			}
	// 		}
	// 	}
	// }

	return
}
