// Package toml provides a frontmatter parser for TOML-formatted metadata.
package toml

import (
	"bytes"

	btoml "github.com/BurntSushi/toml"
	"github.com/client9/ssg"
)

// Parser returns a MetaParser that reads frontmatter as TOML.
func Parser() ssg.MetaParser {
	return func(s []byte) (ssg.ContentSourceConfig, error) {
		meta := ssg.ContentSourceConfig{}
		if _, err := btoml.NewDecoder(bytes.NewReader(s)).Decode(&meta); err != nil {
			return nil, err
		}
		return meta, nil
	}
}
