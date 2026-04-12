package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/client9/ssg"
	"github.com/client9/ssg/meta/email"
	goyaml "gopkg.in/yaml.v3"
)

func copyComments(src []byte) []byte {
	out := []byte{}
	lines := bytes.Split(src, []byte{'\n'})
	for _, line := range lines {
		if len(line) > 0 && line[0] == '#' {
			out = append(out, line...)
			out = append(out, '\n')
		}
	}
	return out
}

func convert(headFrom, headTo ssg.MetaHeadType, src []byte) ([]byte, error) {
	head := make(map[string]any)
	var err error

	comments := copyComments(src)

	// If Splitter returns a nil head, no matching frontmatter was found; skip.
	meta, body := ssg.Splitter(headFrom, src)
	if meta == nil {
		return nil, nil
	}

	switch headFrom.Name {
	case "yaml":
		err = goyaml.Unmarshal(meta, &head)
	case "json":
		err = json.Unmarshal(meta, &head)
	case "email":
		err = email.Unmarshal(meta, head)
	default:
		panic("unknown headFrom type: " + headFrom.Name)
	}
	if err != nil {
		return nil, err
	}

	switch headTo.Name {
	case "yaml":
		meta, err = goyaml.Marshal(head)
	case "json":
		meta, err = json.MarshalIndent(head, "", "    ")
	case "email":
		meta, err = email.Marshal(head)
	default:
		panic("unknown headTo type: " + headTo.Name)
	}
	if err != nil {
		return nil, err
	}

	comments = append(comments, '\n')
	meta = append(comments, meta...)

	return ssg.Joiner(headTo, meta, body), nil
}

func convertFile(headFrom, headTo ssg.MetaHeadType, file string, overwrite bool) (err error) {
	src, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	out, err := convert(headFrom, headTo, src)
	if err != nil {
		return err
	}
	if out == nil {
		log.Printf("skipping %q", file)
		return nil
	}
	if overwrite {
		return os.WriteFile(file, out, 0644)
	}
	fmt.Println(string(out))
	return nil
}

func strToHead(str string) (ssg.MetaHeadType, error) {
	switch str {
	case "yaml":
		return ssg.MetaHeadYaml, nil
	case "json":
		return ssg.MetaHeadJson, nil
	case "email":
		return ssg.MetaHeadEmail, nil
	case "toml":
		return ssg.MetaHeadToml, nil
	}
	return ssg.MetaHeadType{}, fmt.Errorf("unknown head type %q", str)
}

func main() {
	var writeFile = flag.Bool("write", false, "overwrite file in place")
	var frontFrom = flag.String("from", "yaml", "original frontmatter format")
	var frontTo = flag.String("to", "json", "target frontmatter format")
	flag.Parse()

	headFrom, err := strToHead(*frontFrom)
	if err != nil {
		fmt.Println(err)
		return
	}
	headTo, err := strToHead(*frontTo)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, f := range flag.Args() {
		if err := convertFile(headFrom, headTo, f, *writeFile); err != nil {
			log.Fatalf("unable to convert %s: %v", f, err)
		}
	}
}
