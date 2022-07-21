package model

import (
	"net/url"
	"testing"
)

func init() {

}

func TestAdditionalQueryParamsUrlEncoder(t *testing.T) {

	v := &url.Values{}
	key := "test"
	goodQueryParams := struct {
		ParamA bool `url:"paramA"`
		ParamB bool `url:"-"`
		ParamC bool
		paramD bool `url:"paramD"`
		ParamE bool `url:"paramE"`
	}{ParamA: true, ParamB: true, ParamC: true, paramD: true, ParamE: false}

	err := additionalQueryParamsUrlEncoder(goodQueryParams, key, v)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err.Error())
	}

	gotOutput := v.Encode()
	expectedOutput := "test=paramA&test=ParamC"
	if gotOutput != expectedOutput {
		t.Errorf("Result differs, want `%v`, got `%v`", gotOutput, expectedOutput)
	}

	faultyQueryParams := struct {
		ParamA bool `url:"paramA"`
		ParamB string
	}{ParamA: true, ParamB: "test"}

	err = additionalQueryParamsUrlEncoder(faultyQueryParams, key, v)
	if err == nil {
		t.Fatalf("Unexpected error: %s", err.Error())
	} else {
		gotOutput := err.Error()
		expectedOutput := "additionalQueryParamsUrlEncoder expects struct input with only booleans, discovered 'string'"
		if gotOutput != expectedOutput {
			t.Errorf("Error differs, want `%v`, got `%v`", gotOutput, expectedOutput)
		}
	}

}
