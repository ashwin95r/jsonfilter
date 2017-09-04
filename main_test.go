package main

import (
	"bufio"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	f, err := os.Open("input.json")
	require.NoError(t, err)
	defer f.Close()

	p := bufio.NewReader(f)
	js, err := parse(p)
	require.NoError(t, err)

	resp, err := ioutil.ReadFile("output.json")
	require.NoError(t, err)

	require.JSONEq(t, string(resp), string(js))
}
