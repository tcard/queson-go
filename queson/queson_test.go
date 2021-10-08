package queson

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestBackAndForth(t *testing.T) {
	t.Parallel()

	for _, c := range []struct {
		name string
		js   string
		qs   string
	}{
		{"true", `true`, `true`},
		{"false", `false`, `false`},
		{"null", `null`, `null`},

		{"123", `123`, `123`},
		{"-123", `-123`, `-123`},
		{"123.4", `123.4`, `123.4`},
		{"-123.4", `-123.4`, `-123.4`},
		{"123.4e+5", `123.4e+5`, `123.4e5`},
		{"123.4E+5", `123.4E+5`, `123.4e5`},
		{"-123.4E-5", `-123.4E-5`, `-123.4e-5`},

		{"string", `"hey"`, `w.hey.w`},
		{"spaces", `"hola qué tal"`, `w.hola_qué_tal.w`},
		{"escape chars", `"\"hola\t\bqué\r\ntal._\\"`, `w."hola.t.bqué.r.ntal..._\.w`},
		{"🌝", `"🌝"`, `w.🌝.w`},

		{"list", `[ 123, 456, [ "78", true ]]`, `I.123_456_I.w.78.w_true.I.I`},
		{"empty list", `[]`, `I..I`},

		{
			"object",
			`{ "foo": "bar", "qux-": {"list": [1,2,3]} }`,
			`X.w.foo.w-w.bar.w_w.qux-.w-X.w.list.w-I.1_2_3.I.X.X`,
		},
	} {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			gotQS, err := FromJSON(c.js)
			assertNoError(t, err)
			assertEqual(t, c.qs, gotQS)

			var expected, got interface{}
			assertNoError(t, json.Unmarshal([]byte(c.js), &expected))
			assertNoError(t, Unmarshal([]byte(gotQS), &got))
			assertEqual(t, expected, got)

			expectedJS, err := json.Marshal(expected)
			assertNoError(t, err)
			expectedMarshal, err := FromJSONBytes(expectedJS)
			assertNoError(t, err)
			gotMarshal, err := Marshal(got)
			assertNoError(t, err)
			assertEqual(t, string(expectedMarshal), string(gotMarshal))
		})
	}
}

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

func assertEqual(t *testing.T, expected, got interface{}) {
	t.Helper()
	if !reflect.DeepEqual(expected, got) {
		t.Fatalf("expected:\n\n%#v\n\ngot:\n\n%#v", expected, got)
	}
}
