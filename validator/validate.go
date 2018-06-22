package validator

import (
	. "github.com/matthewmueller/gqlparser"
	"github.com/matthewmueller/gqlparser/errors"
)

func Validate(schema *Schema, doc *QueryDocument) []errors.Validation {
	ctx := vctx{
		schema:   schema,
		document: doc,
	}

	ctx.walk()

	return ctx.errors
}
