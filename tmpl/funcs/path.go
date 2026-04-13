package funcs

import (
	"path"
	"text/template"
)

func pathFuncMap() template.FuncMap {
	return template.FuncMap{
		"pathBase":  path.Base,
		"pathDir":   path.Dir,
		"pathExt":   path.Ext,
		"pathJoin":  path.Join,
		"pathClean": path.Clean,
	}
}
