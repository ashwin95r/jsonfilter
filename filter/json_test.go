package filter

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	f, err := ioutil.ReadFile("input.json")
	require.NoError(t, err)

	js, err := Parse(f)
	require.NoError(t, err)

	resp, err := ioutil.ReadFile("output.json")
	require.NoError(t, err)

	require.JSONEq(t, string(resp), string(js))
}

func TestParseNull(t *testing.T) {
	req := ``
	_, err := Parse([]byte(req))
	require.Error(t, err)
}

func TestParseInvalid(t *testing.T) {
	req := `{asdasd}`
	_, err := Parse([]byte(req))
	require.Error(t, err)
}

func TestParseInvalid1(t *testing.T) {
	req := `{{{}}`
	_, err := Parse([]byte(req))
	require.Error(t, err)
}
