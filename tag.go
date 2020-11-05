package multiconfig

import (
	"reflect"

	"github.com/fatih/structs"
)

// TagLoader satisfies the loader interface. It parses a struct's field tags
// and populates the each field with that given tag.
type TagLoader struct {
	// DefaultTagName is the default tag name for struct fields to define
	// default values for a field. Example:
	//
	//   // Field's default value is "koding".
	//   Name string `default:"koding"`
	//
	// The default value is "default" if it's not set explicitly.
	DefaultTagName string
}

func (t *TagLoader) Load(s interface{}) error {
	if t.DefaultTagName == "" {
		t.DefaultTagName = "default"
	}

	for _, field := range structs.Fields(s) {

		if err := t.processField(t.DefaultTagName, field); err != nil {
			return err
		}
	}

	return nil
}

// processField gets tagName and the field, recursively checks if the field has the given
// tag, if yes, sets it otherwise ignores
func (t *TagLoader) processField(tagName string, field *structs.Field) error {
	switch field.Kind() {
	case reflect.Struct:
		for _, f := range field.Fields() {
			if err := t.processField(tagName, f); err != nil {
				return err
			}
		}

	case reflect.Slice:
		val := reflect.ValueOf(field.Value())
		for i := 0; i < val.Len(); i++ {
			switch {
			case val.Index(i).Kind() == reflect.Ptr && val.Index(i).Elem().Type().Kind() == reflect.Struct: // []*Server
				for _, field := range structs.Fields(val.Index(i).Interface()) {
					if err := t.processField(tagName, field); err != nil {
						return err
					}
				}

			case val.Index(i).Kind() == reflect.Struct: // []Server
				vp := reflect.New(val.Index(i).Type())
				vp.Elem().Set(val.Index(i))
				for _, field := range structs.Fields(vp.Interface()) {
					if err := t.processField(tagName, field); err != nil {
						return err
					}
				}
				val.Index(i).Set(vp.Elem())
			}
		}

	default:
		defaultVal := field.Tag(t.DefaultTagName)
		if defaultVal == "" {
			return nil
		}

		err := fieldSet(field, defaultVal)
		if err != nil {
			return err
		}
	}

	return nil
}
