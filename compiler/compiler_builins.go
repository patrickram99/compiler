package compiler

import (
	"main/object"
)

// This file is no longer needed for MIPS generation, but we'll keep it for reference
var builtins = map[string]*object.Builtin{
	"SpeakNow": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			// The actual implementation is now in the MIPS generator
			return nil
		},
	},
	// Other builtins can be added here
}
