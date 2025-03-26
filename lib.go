package puregostruct

import (
	"fmt"
	"reflect"

	"github.com/ebitengine/purego"
)

func registerLibFunc(field string, val any, lib uintptr, symbol string) {
	defer func() {
		if err := recover(); err != nil {
			panic(fmt.Sprintf("failed to register field `%s` with purego tag `%s`: %v", field, symbol, err))
		}
	}()
	purego.RegisterLibFunc(val, lib, symbol)
}

func LoadLibrary(val any, names ...string) error {
	var err error
	var lib uintptr
	for _, name := range names {
		lib, err = openLibrary(name)
		if err == nil {
			break
		}
	}

	if err != nil {
		return err
	}

	v := reflect.ValueOf(val)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("val must be a pointer to a struct")
	}

	v = v.Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		symbol := field.Tag.Get("purego")
		if symbol != "" {
			if !v.Field(i).CanAddr() {
				return fmt.Errorf("field %s cannot be addressed", field.Name)
			}
			registerLibFunc(field.Name, v.Field(i).Addr().Interface(), lib, symbol)
		}
	}

	return nil
}
