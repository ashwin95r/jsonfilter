package main

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	f, err := ioutil.ReadFile("input.json")
	require.NoError(t, err)

	fmt.Println(f)
	js, err := parse(f)
	require.NoError(t, err)

	resp, err := ioutil.ReadFile("output.json")
	require.NoError(t, err)

	require.JSONEq(t, string(resp), string(js))
}
