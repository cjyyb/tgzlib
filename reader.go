package tgzlib

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"fmt"
	"strings"
)

type reader struct {
	buf       *bytes.Buffer
	sr        io.Reader
	closeFlag bool
	*tar.Reader
	gr *gzip.Reader
}

func (r *reader) Close() error {
	if r.closeFlag {
		return fmt.Errorf("writer already close")
	}
	r.closeFlag = true
	if r.buf != nil {
		r.buf.Reset()
	}
	grerr := r.gr.Close()
	if grerr != nil {
		return grerr
	}
	return nil
}

type BufferFile struct {
	Name string
	Data []byte
}

func(r *reader) Read()([]*BufferFile, error) {
	if r.closeFlag {
		return nil, fmt.Errorf("writer already close")
	}
	defer r.Close()
	delimiter := "/"
	files := make([]*BufferFile, 0, 1)
	for {
		var buf bytes.Buffer
		header, err := r.Reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, err
			}
		}

		if header.FileInfo().IsDir() {
			continue
		}
		if strings.ContainsRune(header.Name, '\\') {
			delimiter = "\\"
		}

		parts := strings.Split(header.Name, delimiter)
		n := strings.Join(parts[1:], delimiter)
		n = strings.Replace(n, delimiter, "/", -1)
		if _, err := io.Copy(&buf, r.Reader); err != nil {
			return nil, err
		}
		files = append(files, &BufferFile{
			Name: n,
			Data: buf.Bytes(),
		})
		buf.Reset()
	}
	return files, nil
}

//NewReader ...
func NewReader(sr io.Reader) (*reader, error) {
	gzipReader, err := gzip.NewReader(sr)
	if err != nil {
		return nil, err
	}
	return &reader{
		sr:     sr,
		Reader: tar.NewReader(gzipReader),
		gr:     gzipReader,
	}, nil
}


