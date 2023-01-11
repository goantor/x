package x

import (
	"github.com/fatih/color"
	_json "github.com/json-iterator/go"
	"os"
)

var (
	json = _json.Config{
		EscapeHTML:             true,
		SortMapKeys:            true,
		ValidateJsonRawMessage: true,
		UseNumber:              true,
	}.Froze()
)

type H map[string]interface{}

func DD(args ...interface{}) {
	js, _ := json.Marshal(args)
	color.Yellow("%s", js)
	os.Exit(0)
}
