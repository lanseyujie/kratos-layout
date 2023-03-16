package conf

import (
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
)

type SourceFile config.Source

func NewSourceFile(configDir string) SourceFile {
	return file.NewSource(configDir)
}
