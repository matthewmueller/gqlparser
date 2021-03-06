package gqlparser

import (
	"fmt"
)

type Schema struct {
	Query        *Definition
	Mutation     *Definition
	Subscription *Definition

	Types      map[string]*Definition
	Directives map[string]*DirectiveDefinition

	possibleTypes map[string][]*Definition
}

func (s *Schema) addPossibleType(name string, def *Definition) {
	s.possibleTypes[name] = append(s.possibleTypes[name], def)
}

func LoadSchema(input string) (*Schema, error) {
	ast, err := ParseSchema(input)
	if err != nil {
		return nil, err
	}

	schema := Schema{
		Types:         map[string]*Definition{},
		Directives:    map[string]*DirectiveDefinition{},
		possibleTypes: map[string][]*Definition{},
	}

	for i, def := range ast.Definitions {
		if schema.Types[def.Name] != nil {
			return nil, fmt.Errorf("Cannot redeclare type %s.", def.Name)
		}
		schema.Types[def.Name] = &ast.Definitions[i]

		if def.Kind != Interface {
			for _, intf := range def.Interfaces {
				schema.addPossibleType(intf.Name(), &ast.Definitions[i])
			}
			schema.addPossibleType(def.Name, &ast.Definitions[i])
		}
	}

	for _, ext := range ast.Extensions {
		def := schema.Types[ext.Name]
		if def == nil {
			return nil, fmt.Errorf("Cannot extend type %s because it does not exist.", ext.Name)
		}

		if def.Kind != ext.Kind {
			return nil, fmt.Errorf("Cannot extend type %s because the base type is a %s, not %s.", ext.Name, def.Kind, ext.Kind)
		}

		def.Directives = append(def.Directives, ext.Directives...)
		def.Interfaces = append(def.Interfaces, ext.Interfaces...)
		def.Fields = append(def.Fields, ext.Fields...)
		def.Types = append(def.Types, ext.Types...)
		def.Values = append(def.Values, ext.Values...)
	}

	for i, dir := range ast.Directives {
		if schema.Directives[dir.Name] != nil {
			return nil, fmt.Errorf("Cannot redeclare directive %s.", dir.Name)
		}
		schema.Directives[dir.Name] = &ast.Directives[i]
	}

	if len(ast.Schema) > 1 {
		return nil, fmt.Errorf("Cannot have multiple schema entry points, consider schema extensions instead.")
	}

	if len(ast.Schema) == 1 {
		for _, entrypoint := range ast.Schema[0].OperationTypes {
			def := schema.Types[entrypoint.Type.Name()]
			if def == nil {
				return nil, fmt.Errorf("Schema root %s refers to a type %s that does not exist.", entrypoint.Operation, entrypoint.Type)
			}
			switch entrypoint.Operation {
			case Query:
				schema.Query = def
			case Mutation:
				schema.Mutation = def
			case Subscription:
				schema.Subscription = def
			}
		}
	}

	for _, ext := range ast.SchemaExtension {
		for _, entrypoint := range ext.OperationTypes {
			def := schema.Types[entrypoint.Type.Name()]
			if def == nil {
				return nil, fmt.Errorf("Schema root %s refers to a type %s that does not exist.", entrypoint.Operation, entrypoint.Type)
			}
			switch entrypoint.Operation {
			case Query:
				schema.Query = def
			case Mutation:
				schema.Mutation = def
			case Subscription:
				schema.Subscription = def
			}
		}
	}

	return &schema, nil
}

// GetPossibleTypes will enumerate all the definitions for a given interface or union
func (s *Schema) GetPossibleTypes(def *Definition) []*Definition {
	if def.Kind == Union {
		var defs []*Definition
		for _, t := range def.Types {
			defs = append(defs, s.Types[t.Name()])
		}
		return defs
	}

	return s.possibleTypes[def.Name]
}
