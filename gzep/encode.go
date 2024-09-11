//go:build !solution

package gzep

import (
	"compress/gzip"
	"io"
	"sync"
)

var Pool = sync.Pool{
	New: func() interface{} {
		return gzip.NewWriter(nil)
	},
}

func Encode(data []byte, w io.Writer) error {
	opt := Pool.Get().(*gzip.Writer)
	defer func() {
		opt.Close()
		Pool.Put(opt)
	}()

	opt.Reset(w)
	if _, err := opt.Write(data); err != nil {
		return err
	}
	return opt.Flush()
}
