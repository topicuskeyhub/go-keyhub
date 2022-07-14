package model

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

// DebugResponseDecoder Usage sling.ResponseDecoder(&model.DebugResponseDecoder{})......
type DebugResponseDecoder struct {
}

// Decode Print raw body while unmarshalling json Body
func (drd *DebugResponseDecoder) Decode(resp *http.Response, v interface{}) error {

	var err error
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Printf("Raw json response: %s \n", body)
	err = json.Unmarshal(body, &v)
	return err

}
