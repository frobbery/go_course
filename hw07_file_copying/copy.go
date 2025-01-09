package main

import (
	"errors"
	"io"
	"os"
	"path/filepath"
)

var (
	ErrUnsupportedFile = errors.New("unsupported file")

	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")

	ErrSameInputOutput = errors.New("from path is equal to to path")
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

	defer file.Close()

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

	defer file.Close()

	_, err = file.Write(*cb.buf)
	if err != nil {
		return err
	}

	return nil
}

func Copy(fromPath, toPath string, offset, limit int64) error {
	err := validate(fromPath, toPath, offset)
	if err != nil {
		return err
	}

	cb := CopyBuffer{}

	if err := cb.readFrom(fromPath, offset, limit); err == nil {
		return cb.writeTo(toPath)
	}

	return err
}

func validate(fromPath, toPath string, offset int64) error {
	err := checkIfPathsEqual(fromPath, toPath)
	if err != nil {
		return err
	}

	err = checkIfInputSupported(fromPath)
	if err != nil {
		return err
	}

	err = checkIfOffsetGreaterThanFileSize(fromPath, offset)
	if err != nil {
		return err
	}

	return nil
}

func checkIfPathsEqual(fromPath, toPath string) error {
	fromAbsPath, err := filepath.Abs(fromPath)
	if err != nil {
		return err
	}

	toAbsPath, err := filepath.Abs(toPath)
	if err != nil {
		return err
	}

	if fromAbsPath == toAbsPath {
		return ErrSameInputOutput
	}

	return nil
}

func checkIfInputSupported(fromPath string) error {
	matchesUrandom, err := filepath.Match("*/dev/urandom/*", fromPath)
	if err != nil {
		return err
	}

	matchesProc, err := filepath.Match("*/proc*", fromPath)
	if err != nil {
		return err
	}

	if matchesUrandom || matchesProc {
		return ErrUnsupportedFile
	}

	return nil
}

func checkIfOffsetGreaterThanFileSize(fromPath string, offset int64) error {
	stats, err := os.Stat(fromPath)
	if err != nil {
		return err
	}

	if stats.IsDir() {
		return ErrUnsupportedFile
	}

	if stats.Size() < offset {
		return ErrOffsetExceedsFileSize
	}

	return nil
}
