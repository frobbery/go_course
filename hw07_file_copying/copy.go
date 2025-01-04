package main

import (
	"errors"
	"io"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

type CopyBuffer struct {
	buf *[]byte
}

func (cb *CopyBuffer) initBuf(fileSize, limit, offset int64) {
	var buf []byte

	if fileSize-offset < limit || limit == 0 {
		buf = make([]byte, fileSize-offset)
	} else {
		buf = make([]byte, limit)
	}
	cb.buf = &buf
}

func (cb *CopyBuffer) readFrom(fromPath string, offset int64, limit int64) error {
	file, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	if cb.buf == nil {
		stats, err := file.Stat()
		if err != nil {
			return err
		}
		cb.initBuf(stats.Size(), limit, offset)
	}
	_, err = file.ReadAt(*cb.buf, offset)
	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}
	return nil
}

func (cb *CopyBuffer) writeTo(toPath string) error {
	file, err := os.Create(toPath)
	if err != nil {
		return err
	}
	_, err = file.Write(*cb.buf)
	if err != nil {
		return err
	}
	return nil
}

func Copy(fromPath, toPath string, offset, limit int64) error {
	cb := CopyBuffer{}
	var err error
	if err := cb.readFrom(fromPath, offset, limit); err == nil {
		return cb.writeTo(toPath)
	}
	return err
}
