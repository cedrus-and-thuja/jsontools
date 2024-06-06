package main

import (
	_ "embed"
	"flag"
	"fmt"
	"os"

	"github.com/cedrus-and-thuja/jsontools/pkg/generator"
	"github.com/santhosh-tekuri/jsonschema/v5"
)

func main() {
	var schemaFileLocation string
	var generateGo bool
	var generateKotlin bool

	flag.StringVar(&schemaFileLocation, "s", "", "location of schema file, required")
	flag.BoolVar(&generateGo, "go", true, "generate go structs, defaults to true")
	flag.BoolVar(&generateKotlin, "kotlin", false, "generate kotlin classes, defaults to false")
	flag.Parse()

	schemaFile, err := os.Open(schemaFileLocation)
	if err != nil {
		fmt.Printf("error loading schema: %s", err)
		os.Exit(8)
	}
	compiler := jsonschema.NewCompiler()
	compiler.ExtractAnnotations = true
	if err := compiler.AddResource("schema.json", schemaFile); err != nil {
		panic(err)
	}
	schema, err := compiler.Compile("schema.json")
	if err != nil {
		fmt.Printf("error parsing schema: %s", err)
		os.Exit(1)
	}
	if generateGo {
		err = generator.Construct(schema, generator.NewConfig(), generator.GO)
		if err != nil {
			fmt.Printf("error generating go structs: %s", err)
			os.Exit(2)
		}
	}
	if generateKotlin {
		err = generator.Construct(schema, generator.NewConfig(), generator.KOTLIN)
		if err != nil {
			fmt.Printf("error generating kotlin classes: %s", err)
			os.Exit(2)
		}
	}
}
