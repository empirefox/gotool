package crypt

import (
	"reflect"

	"github.com/mcuadros/go-defaults"
)

type Files struct {
	Xps    map[string][]byte
	filler *defaults.Filler
	tag    string
}

func NewFiles(xps map[string][]byte, tag string) *Files {
	if tag == "" {
		tag = "xps"
	}
	fs := &Files{
		Xps: xps,
		tag: tag,
	}
	fs.filler = fs.getDefaultFiller()
	return fs
}

func (fs Files) Equip(variable interface{}) {
	fs.getDefaultFiller().Fill(variable)
}

func (fs Files) getDefaultFiller() *defaults.Filler {
	if fs.filler == nil {
		fs.filler = fs.newDefaultFiller()
	}
	return fs.filler
}

func (fs Files) newDefaultFiller() *defaults.Filler {
	funcs := make(map[reflect.Kind]defaults.FillerFunc, 0)

	funcs[reflect.Slice] = func(field *defaults.FieldData) {
		k := field.Value.Type().Elem().Kind()
		switch k {
		case reflect.Uint8:
			if field.Value.Bytes() != nil {
				return
			}
			field.Value.SetBytes([]byte(fs.Xps[field.TagValue]))
		case reflect.Struct:
			count := field.Value.Len()
			for i := 0; i < count; i++ {
				fields := fs.getDefaultFiller().GetFieldsFromValue(field.Value.Index(i), nil)
				fs.getDefaultFiller().SetDefaultValues(fields)
			}
		}
	}

	funcs[reflect.Struct] = func(field *defaults.FieldData) {
		fields := fs.getDefaultFiller().GetFieldsFromValue(field.Value, nil)
		fs.getDefaultFiller().SetDefaultValues(fields)
	}

	return &defaults.Filler{FuncByKind: funcs, Tag: fs.tag}
}
