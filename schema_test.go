package gqlparser

import (
	"testing"

	"github.com/matthewmueller/gqlparser/spec"
)

func TestSchemaDocument(t *testing.T) {
	spec.Test(t, "spec/schema.yml", func(t *testing.T, input string) spec.Spec {
		doc, err := ParseSchema(input)
		return spec.Spec{
			Error: err,
			AST:   spec.DumpAST(doc),
		}
	})
}
