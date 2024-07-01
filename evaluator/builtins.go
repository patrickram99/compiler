package evaluator

import (
	"fmt"
	"main/object"
)

var builtins = map[string]*object.Builtin{
	"len": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return createError("Numero equivocado de argumentos. Son: %d, deberian ser 1",
					len(args))
			}
			switch arg := args[0].(type) {
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			default:
				return createError("Tipo sin soporte para `len` no es string sino: %s",
					args[0].Type())
			}
		},
	},

	"debut": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return createError("Numero equivocado de argumentos. Son: %d, deberian ser 1",
					len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return createError("Tipo sin soporte para `debut` no es array sino: %s",
					args[0].Type())
			}
			arr := args[0].(*object.Array)
			if len(arr.Elements) > 0 {
				return arr.Elements[0]
			}
			return NULL
		},
	},

	"ttpd": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return createError("Numero equivocado de argumentos. Son: %d, deberian ser 1",
					len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return createError("Tipo sin soporte para `ttpd` no es array sino: %s",
					args[0].Type())
			}
			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			if length > 0 {
				return arr.Elements[length-1]
			}
			return NULL
		},
	},

	"rest": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return createError("Numero equivocado de argumentos. Son: %d, deberian ser 1",
					len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return createError("Tipo sin soporte para `ttpd` no es array sino: %s",
					args[0].Type())
			}
			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			if length > 0 {
				newElements := make([]object.Object, length-1, length-1)
				copy(newElements, arr.Elements[1:length])
				return &object.Array{Elements: newElements}
			}
			return NULL
		},
	},
	"billboard": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return createError("Numero equivocado de argumentos. Son: %d, deberian ser 1",
					len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return createError("Tipo sin soporte para `ttpd` no es array sino: %s",
					args[0].Type())
			}
			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			newElements := make([]object.Object, length+1, length+1)
			copy(newElements, arr.Elements)
			newElements[length] = args[1]
			return &object.Array{Elements: newElements}
		},
	},
	"SpeakNow": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}
			return NULL
		},
	},
}
