package clamwin

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

// Alternate format in hsb files
var hsbLine, _ = regexp.Compile(`\b([^:]+):([0-9]+):`)

func NewClamwinConnector() *ClamwinConnector {
	return &ClamwinConnector{
		GetClamwinStream: func() (io.ReadCloser, error) {
			res, err := http.Get("http://database.clamav.net/main.cvd")
			if err != nil {
				return nil, err
			}
			return res.Body, nil
		},
		BufferSize: 32 * 1024,
	}
}

type ClamwinConnectorImpl struct{}

type ClamwinConnector struct {
	GetClamwinStream func() (io.ReadCloser, error)
	BufferSize       int32
}

func (c *ClamwinConnector) Download(out io.Writer) error {
	in, err := c.GetClamwinStream()
	if err != nil {
		return err
	}
	defer in.Close()
	io.CopyN(ioutil.Discard, in, 512)
	out = NewClamwinNormalizer(out)
	return c.ExtractMainMDB(in, out)
}

func NewClamwinNormalizer(out io.Writer) io.Writer {
	return &ClamwinNormalizer{
		Out: out,
	}
}

type ClamwinNormalizer struct {
	Out         io.Writer
	DanglingBit string
}

func (w *ClamwinNormalizer) Write(p []byte) (n int, err error) {
	content := string(p)
	lastReturn := strings.LastIndex(content, "\n")
	lines := strings.Split(w.DanglingBit+content[0:lastReturn], "\n")
	newContent := ""
	for _, line := range lines {
		newContent += hsbLine.ReplaceAllString(line, "$2:$1:") + "\n"
	}
	w.DanglingBit = content[lastReturn+1:]
	_, err = w.Out.Write([]byte(newContent))
	return len([]byte(p)), err
}

func (c *ClamwinConnector) ExtractMainMDB(gzipStream io.Reader, out io.Writer) error {
	uncompressedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		return err
	}
	tarReader := tar.NewReader(uncompressedStream)
	count := 0
	buf := make([]byte, c.BufferSize)
	for true {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		switch header.Typeflag {
		case tar.TypeReg:
			if header.Name == "main.mdb" {
				fmt.Println("\n" + header.Name)
				if _, err := io.CopyBuffer(out, tarReader, buf); err != nil {
					return err
				}
				count++
			} else if header.Name == "main.hsb" {
				fmt.Println("\n" + header.Name)
				if _, err := io.CopyBuffer(out, tarReader, buf); err != nil {
					return err
				}
				count++
			}
		}
	}
	if count == 2 {
		return nil
	}
	return fmt.Errorf("didn't find main.mdb and/or main.hsb in this tar.gz")
}
