package build

import (
	"os"

	"github.com/wbuchwalter/azul/build/basic"
	"github.com/wbuchwalter/azul/build/golang"
	"github.com/wbuchwalter/azul/function"
)

type buildFn func(f *function.Function) (function.FilesMap, function.Config, error)

type builder struct {
	build buildFn
}

func Build(f *function.Function) (function.FilesMap, function.Config, error) {
	var b builder

	if _, err := os.Stat(f.Path + "main.go"); os.IsNotExist(err) {
		b.setBuilder(basic.Build)
	} else {
		b.setBuilder(golang.Build)
	}

	return b.build(f)
}

func (b *builder) setBuilder(bfn buildFn) {
	b.build = bfn
}
