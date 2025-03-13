package code

type Preset struct {
	Glob  string
	Files []string
}

var presets = map[string]Preset{
	"js": {
		Glob:  "./**/*.js",
		Files: []string{"package.json", "Makefile"},
	},
}
