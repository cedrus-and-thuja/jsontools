package generator

import (
	"io"
	"os"
	"path/filepath"
)

type Config struct {
	PackageName       string `json:"packageName"`
	OutputDirectory   string `json:"outputDirectory"`
	KotlinPackageName string `json:"kotlinPackageName"`
	KotlineFileName   string `json:"kotlinFileName"`
	GoFileName        string `json:"goFileName"`
}

func NewConfig() Config {
	return Config{
		PackageName:       EnvOrDefault("JANUS_GO_PACKAGE_NAME", "generated"),
		KotlinPackageName: EnvOrDefault("JANUS_KOTLIN_PACKAGE_NAME", "org.janus.generated"),
		OutputDirectory:   EnvOrDefault("JANUS_OUTPUT_DIRECTORY", "pkg/generated"),
		KotlineFileName:   EnvOrDefault("JANUS_KOTLIN_FILE_NAME", "schemas.kt"),
		GoFileName:        EnvOrDefault("JANUS_GO_FILE_NAME", "main.go"),
	}
}

func (c *Config) GetWriter(t OutputType) (io.Writer, error) {
	err := os.MkdirAll(c.OutputDirectory, 0750)
	if err != nil {
		return os.Stderr, err
	}
	switch t {
	case GO:
		return os.OpenFile(filepath.Join(c.OutputDirectory, c.GoFileName), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0750)
	case KOTLIN:
		return os.OpenFile(filepath.Join(c.OutputDirectory, c.KotlineFileName), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0750)
	}
	return os.Stdout, nil
}
