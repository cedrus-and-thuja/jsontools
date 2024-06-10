package generator

import (
	"fmt"
	"os"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v5"
)

func LoadSchemaFromFile(schemaFileLocation string) (*jsonschema.Schema, error) {
	schemaFile, err := os.Open(schemaFileLocation)
	if err != nil {
		fmt.Printf("error loading schema: %s", err)
		return nil, err
	}
	compiler := jsonschema.NewCompiler()
	compiler.ExtractAnnotations = true
	if err := compiler.AddResource("schema.json", schemaFile); err != nil {
		return nil, err
	}
	schema, err := compiler.Compile("schema.json")
	if err != nil {
		return nil, err
	}
	return schema, nil
}

func LoadSchemaFromString(schemaText string) (*jsonschema.Schema, error) {
	compiler := jsonschema.NewCompiler()
	compiler.ExtractAnnotations = true
	if err := compiler.AddResource("schema.json", strings.NewReader(schemaText)); err != nil {
		return nil, err
	}
	schema, err := compiler.Compile("schema.json")
	if err != nil {
		return nil, err
	}
	return schema, nil
}
