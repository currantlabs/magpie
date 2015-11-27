package magpie

import (
	"bytes"
	"compress/gzip"
	"io"
)

func Decompress(b []byte) (d []byte, err error) {
	r, err := gzip.NewReader(bytes.NewBuffer(b))
	if err != nil {
		return
	}
	var buf bytes.Buffer
	_, err = io.Copy(&buf, r)
	if err != nil {
		return
	}
	err = r.Close()
	if err != nil {
		return
	}
	d = buf.Bytes()
	return
}
