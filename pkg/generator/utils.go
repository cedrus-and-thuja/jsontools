package generator

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v5"
)

func CapFirst(s string) string {
	if s == "" {
		return " "
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func LCapFirst(s string) string {
	if s == "" {
		return " "
	}
	return strings.ToLower(s[:1]) + s[1:]
}

func AllCap(s interface{}) string {
	if ss, ok := s.(string); ok {
		return strings.ToUpper(ss)
	} else if iv, ok := s.(int); ok {
		return fmt.Sprintf("%d", iv)
	} else if iv, ok := s.(json.Number); ok {
		return string(iv)
	} else {
		return "UNKNOWN"
	}
}

func Literalize(val interface{}) string {
	switch val.(type) {
	case string:
		return fmt.Sprintf("\"%s\"", val)
	default:
		return fmt.Sprintf("%v", val)
	}
}

func MapType(schema *jsonschema.Schema) string {
	ts := GetGoType(schema)
	if schema.Enum != nil && len(schema.Enum) > 0 {
		return CapFirst(schema.Title)
	}
	return ts
}

func MapTypeKotlin(schema *jsonschema.Schema) string {
	ts := GetKotlinType(schema)
	if schema.Enum != nil && len(schema.Enum) > 0 {
		return CapFirst(schema.Title)
	}
	switch ts {
	case "int":
		return "Int"
	case "float64":
		return "Double"
	case "bool":
		return "Boolean"
	default:
		return ts
	}
}

func GetGoType(schema *jsonschema.Schema) string {
	types := schema.Types
	if len(types) == 0 || len(types) > 1 {
		return "interface{}"
	}
	return goTypeForType(types[0], schema)
}

func GetKotlinType(schema *jsonschema.Schema) string {
	types := schema.Types
	if len(types) == 0 || len(types) > 1 {
		return "Any?"
	}
	return kotlinTypeForType(types[0], schema)
}

func goTypeForType(t string, schema *jsonschema.Schema) string {
	switch t {
	case "string":
		return "string"
	case "number":
		return "float64"
	case "integer":
		return "int"
	case "boolean":
		return "bool"
	case "object":
		return schema.Title
	case "array":
		subType := schema.Items2020
		return fmt.Sprintf("[]%s", MapType(subType))
	default:
		return fmt.Sprintf("/* unknown type: %s */ interface{}", t)
	}
}

func kotlinTypeForType(t string, schema *jsonschema.Schema) string {
	switch t {
	case "string":
		return "String"
	case "number":
		return "Double"
	case "integer":
		return "Int"
	case "boolean":
		return "Boolean"
	case "object":
		return schema.Title
	case "array":
		subType := schema.Items2020
		return fmt.Sprintf("List<%s>", MapTypeKotlin(subType))
	default:
		return fmt.Sprintf("/* unknown type: %s */ Any", t)
	}
}

func EnvOrDefault(key string, def string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return def
	}
	return val
}

func NameSafe(value interface{}) string {
	if v, ok := value.(string); ok {
		return strings.ReplaceAll(v, "-", "_")
	} else if iv, ok := value.(int); ok {
		return fmt.Sprintf("INT_%d", iv)
	} else if jv, ok := value.(json.Number); ok {
		return fmt.Sprintf("NUM_%s", string(jv))
	} else {
		return "UNKNOWN"
	}
}

func MarshalSchema(schema *jsonschema.Schema) string {
	b, err := json.Marshal(schema)
	if err != nil {
		return fmt.Sprintf("/* Error marshalling schema: %s */", err)
	}
	return string(b)
}
