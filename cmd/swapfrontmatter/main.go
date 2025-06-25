package main

import (
	"flag"
	"log"
	"os"

	"github.com/client9/ssg"

	"encoding/json"
	"gopkg.in/yaml.v3"
)

func convert(fromhead, tohead string, src []byte) ([]byte, error) {
	meta := []byte{}
	body := []byte{}
	head := make(map[string]any)
	var err error

	// If Splitter returns a head of nil, it means no metadata of the type
	// requested was found.  In this case, we skip
	switch fromhead {
	case "yaml":
		if meta, body = ssg.Splitter(ssg.HeadYaml, src); meta == nil {
			return nil, nil
		}
		err = yaml.Unmarshal(meta, &head)
	case "json":
		if meta, body = ssg.Splitter(ssg.HeadJson, src); meta == nil {
			return nil, nil
		}
		err = json.Unmarshal(meta, &head)
	}

	if err != nil {
		return nil, err
	}

	switch tohead {
	case "yaml":
		meta, err = yaml.Marshal(head)
	case "json":
		meta, err = json.MarshalIndent(head, "", "    ")
	}

	out := []byte{}
	out = append(out, meta...)
	out = append(out, byte('\n'))
	out = append(out, body...)
	return out, nil
}

func convertFile(fromhead string, tohead string, file string) (err error) {
	src, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	out, err := convert(fromhead, tohead, src)
	if err != nil {
		return err
	}

	// skip!
	if out == nil {
		log.Printf("skipping %q", file)
		return nil
	}
	return os.WriteFile(file, out, 0644)
}

func main() {
	var frontFrom = flag.String("from", "yaml", "original front matter")
	var frontTo = flag.String("to", "json", "new front matter")
	flag.Parse()

	files := flag.Args()
	for _, f := range files {
		if err := convertFile(*frontFrom, *frontTo, f); err != nil {
			log.Fatalf("Unable to convert %s: %v", f, err)
		}
	}
}
