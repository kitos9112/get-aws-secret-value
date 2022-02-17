package main

import (
	"flag"
	"testing"
)

var loc = flag.String("location", "World", "Name of location to greet")

func TestInitialise(t *testing.T) {

	cases := map[string]struct{ A, B, Expected string }{
		"noSecretName": {
			A:        "",
			B:        "",
			Expected: "",
		},
		"noSourceProfile": {
			A:        "",
			B:        "",
			Expected: "",
		},
	}
	res := Greet(*loc)

	if res != "Hello, San Francisco!" {
		t.Errorf("String mismatch on test")
	}
}
