package generator

import (
	_ "embed"
	"slices"
	"text/template"

	"github.com/santhosh-tekuri/jsonschema/v5"
)

type OutputType string

const (
	GO     OutputType = "go"
	KOTLIN OutputType = "kotlin"
)

//go:embed templates/go.tmpl
var GO_TEMPLATE string

//go:embed templates/kotlin.tmpl
var KOTLIN_TEMPLATE string

type Context struct {
	Schemas []*jsonschema.Schema
	Enums   []*jsonschema.Schema
	Config  Config
}

func CollectSchemas(schema *jsonschema.Schema) []*jsonschema.Schema {
	var schemas []*jsonschema.Schema
	schemas = append(schemas, schema)
	for _, def := range schema.Dependencies {
		if dd := def.(*jsonschema.Schema); dd != nil {
			schemas = append(schemas, CollectSchemas(dd)...)
		}
	}
	for key, def := range schema.Properties {
		if def.Title == "" {
			def.Title = key
		}
		if slices.Contains(def.Types, "object") {
			schemas = append(schemas, CollectSchemas(def)...)
		} else if def.Items2020 != nil {
			schemas = append(schemas, CollectSchemas(def.Items2020)...)
		}
	}
	return schemas
}

func CollectEnums(schema *jsonschema.Schema) []*jsonschema.Schema {
	var schemas []*jsonschema.Schema
	for _, def := range schema.Dependencies {
		if dd := def.(*jsonschema.Schema); dd != nil {
			schemas = append(schemas, CollectEnums(dd)...)
		}
	}
	for key, def := range schema.Properties {
		if def.Title == "" {
			def.Title = key
		}
		if slices.Contains(def.Types, "object") {
			schemas = append(schemas, CollectEnums(def)...)
		} else if def.Items2020 != nil {
			schemas = append(schemas, CollectEnums(def.Items2020)...)
		} else if def.Enum != nil && len(def.Enum) > 0 {
			schemas = append(schemas, def)
		}
	}
	if schema.Enum != nil && len(schema.Enum) > 0 {
		schemas = append(schemas, schema)
	}
	return schemas
}

func Construct(schema *jsonschema.Schema, config Config, outType OutputType) error {
	var err error
	wr, err := config.GetWriter(outType)
	if err != nil {
		return err
	}
	funcMap := template.FuncMap{
		// The name "title" is what the function will be called in the template text.
		"MapType":       MapType,
		"GetGoType":     GetGoType,
		"CapFirst":      CapFirst,
		"LCapFirst":     LCapFirst,
		"AllCap":        AllCap,
		"Literalize":    Literalize,
		"MapTypeKotlin": MapTypeKotlin,
		"GetKotlinType": GetKotlinType,
		"NameSafe":      NameSafe,
		"MarshalSchema": MarshalSchema,
	}
	var tmpl *template.Template
	switch outType {
	case KOTLIN:
		tmpl, err = template.New("test").Funcs(funcMap).Parse(KOTLIN_TEMPLATE)
		break
	case GO:
		tmpl, err = template.New("test").Funcs(funcMap).Parse(GO_TEMPLATE)
		break
	}
	if err != nil {
		return err
	}
	context := Context{
		Schemas: CollectSchemas(schema),
		Config:  config,
		Enums:   CollectEnums(schema),
	}
	err = tmpl.Execute(wr, context)
	if err != nil {
		return err
	}
	return nil
}
