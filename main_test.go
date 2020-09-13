package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDownload(t *testing.T) {
	tcwc := &ClamwinConnector{
		GetClamwinStream: func() (io.ReadCloser, error) {
			return os.Open("test.cvd")
		},
	}
	buf := new(bytes.Buffer)
	err := tcwc.Download(buf)
	require.NoError(t, err)
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	require.Len(t, lines, 10)
}

func TestBadStream(t *testing.T) {
	defer os.Remove("bad.cvd")
	tcwc := &ClamwinConnector{
		GetClamwinStream: func() (io.ReadCloser, error) {
			return os.Create("bad.cvd")
		},
	}
	err := tcwc.Download(new(bytes.Buffer))
	require.Error(t, err)
}

func TestCreate(t *testing.T) {
	_, err := NewClamwinConnector().GetClamwinStream()
	require.NoError(t, err)
}
