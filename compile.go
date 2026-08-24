package client

import (
	"fmt"
	"github.com/evanw/esbuild/pkg/api"
)

type Format uint8

const (
	FormatESM Format = iota
	FormatIIFE
)

type CompileOptions struct {
	Format    Format
	Target    string
	Minify    bool
	SourceMap bool
}

func Compile(p ProgramNode) (string, error) {
	return CompileWithOptions(p, CompileOptions{Format: FormatESM, Target: "es2020", Minify: true})
}
func CompileWithOptions(p ProgramNode, o CompileOptions) (string, error) {
	ts, err := TypeScript(p)
	if err != nil {
		return "", err
	}
	f := api.FormatESModule
	if o.Format == FormatIIFE {
		f = api.FormatIIFE
	}
	target := api.ES2020
	if o.Target != "" {
		target = api.ESNext
	}
	r := api.Transform(ts, api.TransformOptions{Loader: api.LoaderTS, Format: f, Target: target, MinifySyntax: o.Minify, MinifyWhitespace: o.Minify, MinifyIdentifiers: o.Minify, Sourcemap: func() api.SourceMap {
		if o.SourceMap {
			return api.SourceMapInline
		}
		return api.SourceMapNone
	}()})
	if len(r.Errors) > 0 {
		return "", fmt.Errorf("esbuild: %s", r.Errors[0].Text)
	}
	return string(r.Code), nil
}
