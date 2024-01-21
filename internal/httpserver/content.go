package httpserver

import (
	"errors"
	"io"
)

const contentCharset = "-ABCDEFGHIJKLMNOPQRSTUVWXYZ"

type contentReader struct {
	size    int64
	current int64
}

// Read implements the io.Read interface.
func (r *contentReader) Read(p []byte) (int, error) {
	length := r.size - 1

	if r.current >= length {
		return 0, io.EOF
	}
	if len(p) == 0 {
		return 0, nil
	}

	var n int
	if r.current == 0 {
		p[n] = '|'
		r.current++
		n++
	}

	for n < len(p) && r.current <= length {
		p[n] = contentCharset[int(r.current)%len(contentCharset)]
		r.current++
		n++
	}

	if r.current >= length {
		p[n-1] = '|'
	}

	return n, nil
}

// Seek implements the io.Seek interface.
func (r *contentReader) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	default:
		return 0, errors.New("seek: invalid whence")
	case io.SeekStart:
		offset += 0
	case io.SeekCurrent:
		offset += r.current
	case io.SeekEnd:
		offset += r.size - 1
	}

	if offset < 0 {
		return 0, errors.New("seek: invalid offset")
	}

	r.current = offset

	return offset, nil
}
