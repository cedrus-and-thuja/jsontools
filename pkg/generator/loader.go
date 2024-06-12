package generator

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v5"
)

func LoadSchemaFromFile(schemaFileLocation string) (string, *jsonschema.Schema, error) {
	schemaFile, err := os.Open(schemaFileLocation)
	if err != nil {
		fmt.Printf("error opening schema file: %s", err)
		return "", nil, err
	}
	schemaBytes, err := io.ReadAll(schemaFile)
	if err != nil {
		fmt.Printf("error loading schema: %s", err)
		return "", nil, err
	}
	schemaText := string(schemaBytes)
	schema, err := LoadSchemaFromString(schemaText)
	return schemaText, schema, err
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
