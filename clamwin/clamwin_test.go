package clamwin

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
		BufferSize: 100,
	}
	buf := new(bytes.Buffer)
	err := tcwc.Download(buf)
	require.NoError(t, err)
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	require.Len(t, lines, 30)
	require.Contains(t, buf.String(), "45056:3ea7d00dedd30bcdf46191358c36ffa4:Win.Test.EICAR_MDB-1")
	require.Contains(t, buf.String(), "22016:9f8ab04e0302e3436f5b8ceb6d98abc8:Win.Spyware.846-2")
	require.Contains(t, buf.String(), "68:44d88612fea8a8f36de82e1278abb02f:Win.Test.EICAR_HDB-1")
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
