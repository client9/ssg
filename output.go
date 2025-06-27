package ssg

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func WriteOutput(outdir string) Renderer {
	return func(wr io.Writer, source io.Reader, data any) (err error) {
		var f *os.File
		cs := data.(ContentSourceConfig)
		file := cs.OutputFile()
		fullpath := filepath.Join(outdir, file)
		dir := filepath.Dir(fullpath)

		if err = os.MkdirAll(dir, 0750); err != nil {
			return fmt.Errorf("output: unable to write dir %s: %v", dir, err)

		}
		if f, err = os.OpenFile(fullpath,
			os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666); err != nil {
			return err
		}
		defer func() {
			// don't overwrite the error if close doesn't work
			if tmperr := f.Close(); err == nil && tmperr != nil {
				err = tmperr
			}
		}()

		if _, err = io.Copy(f, source); err != nil {
			return err
		}
		return nil
	}
}
