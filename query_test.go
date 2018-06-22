package gqlparser

import (
	"testing"

	"github.com/matthewmueller/gqlparser/spec"
)

func TestQueryDocument(t *testing.T) {
	spec.Test(t, "spec/query.yml", func(t *testing.T, input string) spec.Spec {
		doc, err := ParseQuery(input)
		return spec.Spec{
			Error: err,
			AST:   spec.DumpAST(doc),
		}
	})
}
