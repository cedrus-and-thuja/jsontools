package main

import (
	_ "embed"
	"flag"
	"fmt"
	"os"

	"github.com/cedrus-and-thuja/jsontools/pkg/generator"
)

func main() {
	var schemaFileLocation string
	var generateGo bool
	var generateKotlin bool

	flag.StringVar(&schemaFileLocation, "s", "", "location of schema file, required")
	flag.BoolVar(&generateGo, "go", false, "generate go structs, defaults to false")
	flag.BoolVar(&generateKotlin, "kotlin", false, "generate kotlin classes, defaults to false")
	flag.Parse()

	schema, err := generator.LoadSchemaFromFile(schemaFileLocation)
	if err != nil {
		fmt.Printf("error parsing schema: %s\n", err)
		os.Exit(1)
	}
	if !generateGo && !generateKotlin {
		fmt.Printf("nothing generated, schema was parsed successfully\n")
		os.Exit(4)
	}
	if generateGo {
		err = generator.Construct(schema, generator.NewConfig(), generator.GO)
		if err != nil {
			fmt.Printf("error generating go structs: %s\n", err)
			os.Exit(2)
		}
	}
	if generateKotlin {
		err = generator.Construct(schema, generator.NewConfig(), generator.KOTLIN)
		if err != nil {
			fmt.Printf("error generating kotlin classes: %s\n", err)
			os.Exit(3)
		}
	}
}
