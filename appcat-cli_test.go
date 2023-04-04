package main

import (
	"reflect"
	"testing"
)

func TestCleanInput(t *testing.T) {
	type cleanInputTest struct {
		arg      []string
		expected []string
		hasError bool
	}

	var cleanInputTests = []cleanInputTest{
		{[]string{"--foo=bar", "--baz"}, nil, true},
		{[]string{"--baz", "--foo=bar"}, nil, true},
		{[]string{"--foo", "--baz"}, nil, true},
		{[]string{"--foo=bar=baz", "--foo", "bar"}, []string{"--foo", "bar=baz", "--foo", "bar"}, false},
		{[]string{"--foo=bar", "--baz=2", "--fooby", "bary"}, []string{"--foo", "bar", "--baz", "2", "--fooby", "bary"}, false},
		{[]string{"--fooby", "bary", "--foo=bar", "--baz=2"}, []string{"--fooby", "bary", "--foo", "bar", "--baz", "2"}, false},
	}

	for _, test := range cleanInputTests {
		output, err := cleanInputArguments(test.arg)
		hasGenError := err != nil
		if !reflect.DeepEqual(output, test.expected) || hasGenError != test.hasError {
			if hasGenError != test.hasError {
				t.Errorf("Error should have been returned")
			} else {
				t.Errorf("%v\nOutput %q not equal to expected %q", err, output, test.expected)
			}
		}
	}
}

func TestSetFields(t *testing.T) {
	type testStruct struct {
		StringField  string
		Int8Field    int8
		Int16Field   int16
		Int32Field   int32
		Int64Field   int64
		UInt8Field   uint8
		UInt16Field  uint16
		UInt32Field  uint32
		UInt64Field  uint64
		Float32Field float32
		Float64Field float64
		Unsupported  bool
	}

	tests := []struct {
		fieldName string
		value     string
		expected  interface{}
		hasErr    bool
	}{
		{"StringField", "foo", "foo", false},
		{"Int8Field", "042", int64(42), false},
		{"Int16Field", "042", int64(42), false},
		{"Int32Field", "042", int64(42), false},
		{"Int64Field", "042", int64(42), false},
		{"UInt8Field", "0123", uint64(123), false},
		{"UInt16Field", "0123", uint64(123), false},
		{"UInt32Field", "0123", uint64(123), false},
		{"UInt64Field", "0123", uint64(123), false},
		{"Float32Field", "03.14", float64(3.14), false},
		{"Float64Field", "03.14", float64(3.14), false},
		{"Unsupported", "bar", nil, true},
	}

	for _, test := range tests {
		t.Run(test.fieldName, func(t *testing.T) {
			// Create a new instance of the test struct
			testFields := testStruct{}

			// Get the field we want to set
			field := reflect.ValueOf(&testFields).Elem().FieldByName(test.fieldName)

			// Call the setFields function with the field and value
			err := setFields(field, test.value)

			hasErrReturn := err != nil
			// Check that error return is correct
			if !hasErrReturn {
				if test.hasErr {
					t.Errorf("Error should have been returned")
				}
				// Check that the field was set as expected
				actual := field.Interface()
				expected := reflect.ValueOf(test.expected).Convert(field.Type()).Interface()
				if !reflect.DeepEqual(actual, expected) {
					t.Errorf("setFields(%s, %q) = %v; want %v", test.fieldName, test.value, actual, expected)
				}
			} else {
				if !test.hasErr {
					t.Errorf("setFields returned an unexpected error: %v", err)
				}
			}

		})
	}
}
