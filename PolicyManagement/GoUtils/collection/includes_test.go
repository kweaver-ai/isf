package collection

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Arguments struct {
	collection []string
	value      string
}

type Case struct {
	args   Arguments
	expect bool
}

func TestIncludesString(t *testing.T) {
	var cases = []Case{
		{
			args:   Arguments{[]string{"Hello"}, "Hello"},
			expect: true,
		},
		{
			args:   Arguments{[]string{"Hello", "World"}, "World"},
			expect: true,
		},
		{
			args:   Arguments{[]string{"Hello"}, "World"},
			expect: false,
		},
		{
			args:   Arguments{[]string{}, "World"},
			expect: false,
		},
		{
			args:   Arguments{[]string{}, ""},
			expect: false,
		},
		{
			args:   Arguments{[]string{"", ""}, ""},
			expect: true,
		},
	}

	assert := assert.New(t)

	for _, c := range cases {
		get := IncludesString(c.args.collection, c.args.value)

		assert.Equal(get, c.expect)
	}
}
