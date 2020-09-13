package clamwin

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

func NewClamwinConnector() *ClamwinConnector {
	return &ClamwinConnector{
		GetClamwinStream: func() (io.ReadCloser, error) {
			res, err := http.Get("http://database.clamav.net/main.cvd")
			if err != nil {
				return nil, err
			}
			return res.Body, nil
		},
	}
}

type ClamwinConnectorImpl struct{}

type ClamwinConnector struct {
	GetClamwinStream func() (io.ReadCloser, error)
}

func (c *ClamwinConnector) Download(out io.Writer) error {
	in, err := c.GetClamwinStream()
	if err != nil {
		return err
	}
	defer in.Close()
	io.CopyN(ioutil.Discard, in, 512)
	return ExtractMainMDB(in, out)
}

func ExtractMainMDB(gzipStream io.Reader, out io.Writer) error {
	uncompressedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		return err
	}
	tarReader := tar.NewReader(uncompressedStream)
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
				if _, err := io.Copy(out, tarReader); err != nil {
					return err
				}
				return nil
			}
		}
	}
	return fmt.Errorf("didn't find main.mdb in this tar.gz")
}
