package files

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	t.Run("new filereader", func(t *testing.T) {
		result := New(func(string) ([]byte, error) { return []byte{}, nil })
		assert.NotNil(t, result)
	})
}

var testCases = []struct {
	description string
	readFile    func(string) ([]byte, error)
	path        string
	expect      *testStructure
	expectError bool
}{
	{
		description: "valid file",
		readFile: func(s string) ([]byte, error) {
			if s != "./test.json" {
				return []byte{}, errors.New("invalid path given")
			}
			return []byte("{ \"slug\": \"go\", \"name\": \"Go\" }"), nil
		},
		path:        "./test.json",
		expectError: false,
		expect:      &testStructure{Slug: "go", Name: "Go"},
	},
	{
		description: "invalid json",
		readFile: func(s string) ([]byte, error) {
			return []byte("{ \"slug\": 27, \"name\": false }"), nil
		},
		expectError: true,
		expect:      &testStructure{},
	},
	{
		description: "error when reading file",
		readFile: func(s string) ([]byte, error) {
			return nil, errors.New("testerror: error on readd")
		},
		expectError: true,
		expect:      &testStructure{},
	},
}

type testStructure struct {
	Slug string `json:"slug"`
	Name string `json:"name"`
}

func TestFileReader_JSONToStruct(t *testing.T) {
	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			reader := New(testCase.readFile)
			var cache = &testStructure{}
			err := reader.JSONToStruct(testCase.path, cache)
			assertError(t, err, testCase.expectError)
			assert.Equal(t, cache, testCase.expect)
		})
	}
}

func assertError(t *testing.T, err error, expectError bool) {
	if expectError {
		assert.NotNil(t, err)
	} else {
		assert.Nil(t, err)
	}
}
