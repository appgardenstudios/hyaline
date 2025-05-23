package tool

import "github.com/invopop/jsonschema"

var Reflector jsonschema.Reflector

func init() {
	Reflector = jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
}
