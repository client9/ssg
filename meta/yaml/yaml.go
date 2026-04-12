// Package yaml provides a frontmatter parser for YAML-formatted metadata.
package yaml

import (
	"github.com/client9/ssg"
	goyaml "gopkg.in/yaml.v3"
)

// Parser returns a MetaParser that reads frontmatter as YAML.
func Parser() ssg.MetaParser {
	return func(s []byte) (ssg.ContentSourceConfig, error) {
		meta := ssg.ContentSourceConfig{}
		if err := goyaml.Unmarshal(s, &meta); err != nil {
			return nil, err
		}
		return meta, nil
	}
}
