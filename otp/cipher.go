//go:build !solution

package otp

import (
	"io"
)

type streamCipherReader struct {
	r    io.Reader
	prng io.Reader
}

type streamCipherWriter struct {
	w    io.Writer
	prng io.Reader
}

func NewReader(r io.Reader, prng io.Reader) io.Reader {
	return &streamCipherReader{r: r, prng: prng}
}

func NewWriter(w io.Writer, prng io.Reader) io.Writer {
	return &streamCipherWriter{w: w, prng: prng}
}

func (scr *streamCipherReader) Read(p []byte) (n int, err error) {
	n, err = scr.r.Read(p)

	prngByte := make([]byte, 1)
	for i := 0; i < n; i++ {
		_, _ = scr.prng.Read(prngByte)

		p[i] ^= prngByte[0]
	}

	return n, err
}

func (scw *streamCipherWriter) Write(p []byte) (n int, err error) {
	prngByte := make([]byte, 1)
	data := make([]byte, len(p))
	for i := 0; i < len(p); i++ {
		_, _ = scw.prng.Read(prngByte)
		data[i] = p[i] ^ prngByte[0]
	}

	n, err = scw.w.Write(data)
	return n, err
}
