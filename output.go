package ssg

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func WriteOutput(outdir string, file string, data []byte) error {
	log.Printf("WriteOutput: %s %s len(%d)", outdir, file, len(data))
	// make directory
	fullpath := filepath.Join(outdir, file)

	dir := filepath.Dir(fullpath)
	log.Printf("full = %s, dir=%s", fullpath, dir)
	if err := os.MkdirAll(dir, 0750); err != nil {
		return fmt.Errorf("output: unable to write dir %s: %v", dir, err)

	}

	// open file
	f, err := os.OpenFile(fullpath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}

	_, err = f.Write(data)

	if err != nil {
		return err
	}

	err = f.Close()

	if err != nil {
		return err
	}
	return err
}
