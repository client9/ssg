package ssg

import (
	"fmt"
	"os"
	"path/filepath"
)

// WriteOutput is a Stage that writes pipeline output to disk.
// The output directory is taken from ctx.OutputDir (default "public"); the
// file path within it comes from cfg.OutputFile(). Parent directories are
// created as needed.
var WriteOutput = Step("write-output", writeOutput)

func writeOutput(ctx *Context, cfg ContentSourceConfig, in []byte) (struct{}, error) {
	outdir := "public"
	if ctx != nil && ctx.OutputDir != "" {
		outdir = ctx.OutputDir
	}

	file := cfg.OutputFile()
	if file == "" {
		return struct{}{}, fmt.Errorf("WriteOutput: OutputFile is empty")
	}
	fullpath := filepath.Join(outdir, file)
	dir := filepath.Dir(fullpath)

	if err := os.MkdirAll(dir, 0750); err != nil {
		return struct{}{}, fmt.Errorf("WriteOutput: mkdir %s: %w", dir, err)
	}

	f, err := os.OpenFile(fullpath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return struct{}{}, err
	}
	_, werr := f.Write(in)
	cerr := f.Close()
	if werr != nil {
		return struct{}{}, werr
	}
	return struct{}{}, cerr
}
