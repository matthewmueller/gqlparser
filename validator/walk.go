package validator

import (
	"fmt"

	"github.com/matthewmueller/gqlparser"
	"github.com/matthewmueller/gqlparser/errors"
)

var fieldVisitors []func(vctx *vctx, parentDef *gqlparser.Definition, fieldDef *gqlparser.FieldDefinition, field *gqlparser.Field)

type vctx struct {
	schema   *gqlparser.Schema
	document *gqlparser.QueryDocument
	errors   []errors.Validation
}

func (c *vctx) walk() {
	for _, child := range c.document.Operations {
		c.walkOperation(&child)
	}
	for _, child := range c.document.Fragments {
		c.walkFragment(&child)
	}
}

func (c *vctx) walkOperation(operation *gqlparser.OperationDefinition) {
	var def *gqlparser.Definition
	switch operation.Operation {
	case gqlparser.Query:
		def = c.schema.Query
	case gqlparser.Mutation:
		def = c.schema.Mutation
	case gqlparser.Subscription:
		def = c.schema.Subscription
	}

	for _, v := range operation.SelectionSet {
		c.walkSelection(def, v)
	}
}

func (c *vctx) walkFragment(it *gqlparser.FragmentDefinition) {
	parentDef := c.schema.Types[it.TypeCondition.Name()]
	if parentDef == nil {
		return
	}
	for _, child := range it.SelectionSet {
		c.walkSelection(parentDef, child)
	}
}

func (c *vctx) walkSelection(parentDef *gqlparser.Definition, it gqlparser.Selection) {
	switch it := it.(type) {
	case gqlparser.Field:
		var def *gqlparser.FieldDefinition
		if it.Name == "__typename" {
			def = &gqlparser.FieldDefinition{
				Name: "__typename",
				Type: gqlparser.NamedType("String"),
			}
		} else {
			def = parentDef.Field(it.Name)
		}
		for _, v := range fieldVisitors {
			v(c, parentDef, def, &it)
		}
		for _, sel := range it.SelectionSet {
			c.walkSelection(parentDef, sel)
		}

	case gqlparser.InlineFragment:
		if it.TypeCondition.Name() != "" {
			parentDef = c.schema.Types[it.TypeCondition.Name()]
		}

		for _, sel := range it.SelectionSet {
			c.walkSelection(parentDef, sel)
		}
	default:
		panic(fmt.Errorf("unsupported %T", it))

	}
}
