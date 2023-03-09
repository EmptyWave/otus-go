package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/cheggaaa/pb"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	if len(fromPath) == 0 {
		return fmt.Errorf("%w. Wrong key 'from': %q", ErrUnsupportedFile, fromPath)
	}
	if len(toPath) == 0 {
		return fmt.Errorf("%w. Wrong key 'to': %q", ErrUnsupportedFile, toPath)
	}

	from, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer from.Close()

	info, err := from.Stat()
	if err != nil {
		return err
	}
	sizeFile := info.Size()

	if sizeFile == 0 {
		return fmt.Errorf("wrong key 'from': %q, unknown file length. %w ", fromPath, ErrUnsupportedFile)
	}
	if offset > sizeFile {
		return fmt.Errorf("wrong key 'offset' = %d, file size = %d. %w ", offset, sizeFile, ErrOffsetExceedsFileSize)
	}
	if offset > 0 {
		_, err = from.Seek(offset, 0)
		if err != nil {
			return err
		}
	}
	if limit == 0 {
		limit = sizeFile
	}

	replica, err := os.OpenFile(toPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	defer replica.Close()

	realLimit := realLimit(offset, limit, sizeFile)

	bar := pb.New64(realLimit)
	barReader := bar.NewProxyReader(from)

	_, err = io.CopyN(replica, barReader, realLimit)

	bar.Finish()

	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}
	return nil
}

func realLimit(offset, limit, fileSize int64) int64 {
	if limit == 0 || offset+limit > fileSize {
		return fileSize - offset
	}
	return limit
}
