package tgzlib

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"time"
)

type writer struct {
	buf       *bytes.Buffer
	sw        io.Writer
	closeFlag bool
	*tar.Writer
	gw *gzip.Writer
}

const (
	//DefaultMode ...
	DefaultMode = 0666
	//DefaultCompressLevel ...
	DefaultCompressLevel = gzip.DefaultCompression
)

//Write write file or directory to .tgz file or buffer
func (w *writer) Write(name string) error {
	if w.closeFlag {
		return fmt.Errorf("writer already close")
	}
	walkFunc := func(name string, data []byte) error {
		header := tar.Header{
			Name:    name,
			Mode:    DefaultMode,
			ModTime: time.Now().Local(),
			Size:    int64(len(data)),
		}
		if err := w.Writer.WriteHeader(&header); err != nil {
			return err
		}
		if _, err := w.Writer.Write(data); err != nil {
			return err
		}
		return nil
	}
	if err := WalkWriteFile(name, "", walkFunc); err != nil {
		return err
	}
	return nil
}

func (w *writer) WriteBody(data []byte) (int, error) {
	if w.closeFlag {
		return 0, fmt.Errorf("writer already close")
	}
	return w.Writer.Write(data)
}

func (w *writer) Body() []byte {
	return w.buf.Bytes()
}

func (w *writer) Close() error {
	if w.closeFlag {
		return fmt.Errorf("writer already close")
	}
	w.closeFlag = true
	if w.buf != nil {
		w.buf.Reset()
	}
	twerr := w.Writer.Close()
	gwerr := w.gw.Close()
	if gwerr != nil {
		return gwerr
	}
	if twerr != nil {
		return twerr
	}
	return nil
}

//NewDefaultWriter ...
func NewDefaultWriter() *writer {
	var buf bytes.Buffer
	w, _ := NewWriter(&buf, gzip.DefaultCompression)
	w.buf = &buf
	return w
}

//NewWriter ...
func NewWriter(sw io.Writer, level int) (*writer, error) {
	gzipWriter, err := gzip.NewWriterLevel(sw, level)
	if err != nil {
		return nil, err
	}
	return &writer{
		sw:     sw,
		Writer: tar.NewWriter(gzipWriter),
		gw:     gzipWriter,
	}, nil
}

//TODO add filter
func ignore(name string) bool {
	return true
}
