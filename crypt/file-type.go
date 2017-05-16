package crypt

import (
	"encoding/json"
	"errors"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/rolldever/go-json5"
	"gopkg.in/yaml.v2"
)

var ErrUnsupportedFile = errors.New("File type is not supported")

func DetectFileType(filetype, filename string) string {
	if filetype != "" {
		return filetype
	}
	ext := filepath.Ext(filename)
	if ext == "" {
		return ext
	}
	return ext[1:]
}

func UnmarshalFormat(data []byte, v interface{}, filetype string) (err error) {
	switch filetype {
	case "json":
		err = json.Unmarshal(data, v)
	case "yaml":
		err = yaml.Unmarshal(data, v)
	case "toml":
		err = toml.Unmarshal(data, v)
	case "json5":
		err = json5.Unmarshal(data, v)
	default:
		err = ErrUnsupportedFile
	}
	return
}
