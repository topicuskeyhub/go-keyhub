package model

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"
)

func additionalQueryParamsUrlEncoder(additionalQueryParams interface{}, key string, v *url.Values) error {
	val := reflect.ValueOf(additionalQueryParams)
	typ := val.Type()
	for i := 0; i < typ.NumField(); i++ {
		sf := typ.Field(i)
		if sf.PkgPath != "" && !sf.Anonymous { // unexported
			continue
		}

		sv := val.Field(i)
		tag := sf.Tag.Get("url")
		if tag == "-" {
			continue
		}

		if sv.Kind() == reflect.Bool {
			if sv.Bool() {
				parts := strings.Split(tag, ",")
				if parts[0] == "" {
					v.Add(key, sf.Name) // Use Field name
				} else {
					v.Add(key, parts[0]) // Use name from tag
				}
			}
		} else {
			return fmt.Errorf("additionalQueryParamsUrlEncoder expects struct input with only booleans, discovered '%v'", sv.Kind())
		}
	}

	return nil

}
