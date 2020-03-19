package sources

import (
	"os"
	"strings"

	"testing"

	"github.com/stretchr/testify/require"
)

func Test_StdIn(t *testing.T) {

	testCases := []struct {
		input  string
		hasErr bool
		count  int
	}{
		{input: "[{},{}]", count: 2},
		{input: "[]", count: 0},
		{input: "", count: 0, hasErr: true},
		{input: "{}", count: 1},
		{input: "{", count: 0, hasErr: true},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.input, func(t *testing.T) {

			r := strings.NewReader(tc.input)
			f := StdIn(r)
			c, err := f()

			if tc.hasErr {
				require.NotNil(t, err)
			} else {
				require.Nil(t, err)
			}
			require.Equal(t, len(c), tc.count)

		})
	}
}

func Test_EnvVarValue(t *testing.T) {

	testCases := []struct {
		name   string
		input  string
		hasErr bool
		count  int
	}{
		{name: "one", input: " [{},{}]", count: 2},
		{name: "two", input: "[]", count: 0},
		{name: "three", input: "", count: 0, hasErr: true},
		{name: "four", input: "{  }", count: 1},
		{name: "five", input: "{", count: 0, hasErr: true},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {

			os.Setenv(tc.name, tc.input)
			c, err := EnvVarValue(tc.name)()

			if tc.hasErr {
				require.NotNil(t, err)
			} else {
				require.Nil(t, err)
			}
			require.Equal(t, len(c), tc.count)

		})
	}
}

func Test_Directory(t *testing.T) {

	path := "../../examples/"

	c, err := Directory(path)()

	require.Nil(t, err)
	require.Len(t, c, 4)
}
