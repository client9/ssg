package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/client9/ssg"

	"encoding/json"
	"gopkg.in/yaml.v3"
)

func convert(headFrom, headTo ssg.HeadType, src []byte) ([]byte, error) {
	var meta []byte
	var body []byte
	head := make(map[string]any)
	var err error

	// If Splitter returns a head of nil, it means no metadata of the type
	// requested was found.  In this case, we skip
	if meta, body = ssg.Splitter(headFrom, src); meta == nil {
		return nil, nil
	}

	switch headFrom.Name {
	case "yaml":
		err = yaml.Unmarshal(meta, &head)
	case "json":
		err = json.Unmarshal(meta, &head)
	case "email":
		err = ssg.EmailUnmarshal(meta, head)
	default:
		panic("Unknown headFrom type of " + headFrom.Name)
	}
	if err != nil {
		return nil, err
	}

	switch headTo.Name {
	case "yaml":
		meta, err = yaml.Marshal(head)
	case "json":
		meta, err = json.MarshalIndent(head, "", "    ")
	case "email":
		meta, err = ssg.EmailMarshal(head)

	default:
		panic("Unknown headTo type of " + headTo.Name)
	}
	if err != nil {
		return nil, err
	}

	return ssg.Joiner(headTo, meta, body), nil
}

func convertFile(headFrom, headTo ssg.HeadType, file string, overwrite bool) (err error) {
	src, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	out, err := convert(headFrom, headTo, src)
	if err != nil {
		return err
	}

	// skip!
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

func strToHead(str string) (ssg.HeadType, error) {
	switch str {
	case "yaml":
		return ssg.HeadYaml, nil
	case "json":
		return ssg.HeadJson, nil
	case "email":
		return ssg.HeadEmail, nil
	}
	return ssg.HeadType{}, fmt.Errorf("unknown head type of %q", str)
}

func main() {
	var writeFile = flag.Bool("write", false, "over-write file")
	var frontFrom = flag.String("from", "yaml", "original front matter")
	var frontTo = flag.String("to", "json", "new front matter")
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
	files := flag.Args()
	for _, f := range files {
		if err := convertFile(headFrom, headTo, f, *writeFile); err != nil {
			log.Fatalf("Unable to convert %s: %v", f, err)
		}
	}
}
